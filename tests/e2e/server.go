package e2e

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	_ "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	_ "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lyft/protoc-gen-validate/validate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/infobloxopen/atlas-app-toolkit/requestid"
	"github.com/infobloxopen/atlas-app-toolkit/server"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/svc"
)

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

type Server struct {
	*server.Server
	GRPCListener       net.Listener
	HTTPListener       net.Listener
	BindAddress        string
	GatewayBindAddress string
}

func (s *Server) Serve() error {
	return s.Server.Serve(s.GRPCListener, s.HTTPListener)
}

func NewFeatureFlagServer(client client.Client, bindAddress string, gatewayBindAddress string, gatewayEndpoint string, jwtLabels []string, logger *logrus.Logger) (*Server, error) {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				Time:    time.Duration(viper.GetInt("config.keepalive.time")) * time.Second,
				Timeout: time.Duration(viper.GetInt("config.keepalive.timeout")) * time.Second,
			},
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// logging middleware
				grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)),

				// Request-Id interceptor
				requestid.UnaryServerInterceptor(),

				// validation middleware
				grpc_validator.UnaryServerInterceptor(),

				// collection operators middleware
				gateway.UnaryServerInterceptor(),
			),
		),
	)
	//
	// register service implementation with the grpcServer
	s, err := svc.NewBasicServer(client, jwtLabels)
	if err != nil {
		return nil, err
	}
	pb.RegisterAtlasFeatureFlagServer(grpcServer, s)

	ffs, err := server.NewServer(
		server.WithGrpcServer(grpcServer),
		server.WithGateway(
			gateway.WithGatewayOptions(
				runtime.WithForwardResponseOption(forwardResponseOption),
				runtime.WithIncomingHeaderMatcher(gateway.ExtendedDefaultHeaderMatcher(
					requestid.DefaultRequestIDKey)),
			),
			gateway.WithServerAddress(bindAddress),
			gateway.WithEndpointRegistration(gatewayEndpoint, pb.RegisterAtlasFeatureFlagHandlerFromEndpoint),
		),
	)

	grpcL, err := net.Listen("tcp", bindAddress)
	if err != nil {
		logger.Fatalln(err)
	}

	httpL, err := net.Listen("tcp", gatewayBindAddress)
	if err != nil {
		logger.Fatalln(err)
	}

	return &Server{ffs, grpcL, httpL, bindAddress, gatewayBindAddress}, err
}

func forwardResponseOption(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	return nil
}
