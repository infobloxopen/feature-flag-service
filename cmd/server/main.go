package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-logr/zapr"
	"github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/go-grpc-middleware"
	_ "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lyft/protoc-gen-validate/validate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/infobloxopen/atlas-app-toolkit/gorm/resource"
	"github.com/infobloxopen/atlas-app-toolkit/health"
	"github.com/infobloxopen/atlas-app-toolkit/requestid"
	"github.com/infobloxopen/atlas-app-toolkit/server"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/svc"

	"github.com/Infoblox-CTO/atlas-app-definition-controller/pkg/util/signals"
	"github.com/Infoblox-CTO/atlas.feature.flag/cmd/server/controllers"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	svc.DumpBuildManifest("Terminus / Atlas Feature Flag")

	logger := NewLogger()

	z, _ := zap.NewProduction()
	opts := []zap.Option{
		zap.AddCallerSkip(1),
		//zap.AddStacktrace(zap.DebugLevel),
	}
	z = z.WithOptions(opts...)
	zlogger := zapr.NewLogger(z)

	exitSignals := signals.SetupExitHandlers(zlogger)

	var kubeConfig *rest.Config
	var err error
	incluster, ierr := ctrl.GetConfig()
	if ierr == nil {
		kubeConfig = incluster
	} else {
		kubeConfig, _ = clientcmd.BuildConfigFromFlags("", viper.GetString("kubeconfig"))
	}
	if kubeConfig == nil {
		logger.Errorf("creating rest.Config failed for: %s", viper.GetString("kubeconfig"))
		os.Exit(1)
	}

	mgr, err := controllers.StartKubeController(exitSignals, kubeConfig, controllers.Scheme, zlogger)
	if err != nil {
		logger.Fatal(err)
	}

	if viper.GetBool("internal.enable") {
		go func() { exitSignals.DoneCh() <- ServeInternal(logger) }()
	}

	go func() { exitSignals.DoneCh() <- ServeExternal(mgr.GetClient(), logger) }()

	<-exitSignals.StopCh()
}

func NewLogger() *logrus.Logger {
	logger := logrus.StandardLogger()
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Set the log level on the default logger based on command line flag
	logLevels := map[string]logrus.Level{
		"debug":   logrus.DebugLevel,
		"info":    logrus.InfoLevel,
		"warning": logrus.WarnLevel,
		"error":   logrus.ErrorLevel,
		"fatal":   logrus.FatalLevel,
		"panic":   logrus.PanicLevel,
	}
	if level, ok := logLevels[viper.GetString("logging.level")]; !ok {
		logger.Errorf("Invalid %q provided for log level", viper.GetString("logging.level"))
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}

	return logger
}

// ServeInternal builds and runs the server that listens on InternalAddress
func ServeInternal(logger *logrus.Logger) error {
	healthChecker := health.NewChecksHandler(
		viper.GetString("internal.health"),
		viper.GetString("internal.readiness"),
	)

	healthChecker.AddLiveness("ping", health.HTTPGetCheck(
		fmt.Sprint("http://", viper.GetString("internal.address"), ":", viper.GetString("internal.port"), "/ping"), time.Minute),
	)

	s, err := server.NewServer(
		// register our health checks
		server.WithHealthChecks(healthChecker),
		// this endpoint will be used for our health checks
		server.WithHandler("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("pong"))
		})),
	)
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("internal.address"), viper.GetString("internal.port")))
	if err != nil {
		return err
	}

	logger.Debugf("serving internal http at %q", fmt.Sprintf("%s:%s", viper.GetString("internal.address"), viper.GetString("internal.port")))
	return s.Serve(nil, l)
}

// ServeExternal builds and runs the server that listens on ServerAddress and GatewayAddress
func ServeExternal(client client.Client, logger *logrus.Logger) error {
	grpcServer, err := NewGRPCServer(client, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	s, err := server.NewServer(
		server.WithGrpcServer(grpcServer),
		server.WithGateway(
			gateway.WithGatewayOptions(
				runtime.WithForwardResponseOption(forwardResponseOption),
				runtime.WithIncomingHeaderMatcher(gateway.ExtendedDefaultHeaderMatcher(
					requestid.DefaultRequestIDKey)),
				runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
			),
			gateway.WithServerAddress(fmt.Sprintf("%s:%s", viper.GetString("server.address"), viper.GetString("server.port"))),
			gateway.WithEndpointRegistration(viper.GetString("gateway.endpoint"), pb.RegisterAtlasFeatureFlagHandlerFromEndpoint),
		),
		server.WithHandler("/swagger/", NewSwaggerHandler(viper.GetString("gateway.swaggerFile"))),
	)
	if err != nil {
		logger.Fatalln(err)
	}

	grpcL, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("server.address"), viper.GetString("server.port")))
	if err != nil {
		logger.Fatalln(err)
	}

	httpL, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("gateway.address"), viper.GetString("gateway.port")))
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Printf("serving gRPC at %s:%s", viper.GetString("server.address"), viper.GetString("server.port"))
	logger.Printf("serving http at %s:%s", viper.GetString("gateway.address"), viper.GetString("gateway.port"))

	return s.Serve(grpcL, httpL)
}

func init() {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(viper.GetString("config.source"))
	if viper.GetString("config.file") != "" {
		log.Printf("Serving from configuration file: %s", viper.GetString("config.file"))
		viper.SetConfigName(viper.GetString("config.file"))
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("cannot load configuration: %v", err)
		}
	} else {
		log.Printf("Serving from default values, environment variables, and/or flags")
	}
	resource.RegisterApplication(viper.GetString("app.id"))
	resource.SetPlural()
}

func forwardResponseOption(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	return nil
}
