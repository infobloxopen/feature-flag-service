package main

import (
	"github.com/spf13/pflag"
)

const (
	// configuration defaults support local development (i.e. "go run ...")

	//
	defaultKubeConfig = ""

	// Server
	defaultServerAddress = "0.0.0.0"
	defaultServerPort    = "9090"

	// Gateway
	defaultGatewayEnable      = true
	defaultGatewayAddress     = "0.0.0.0"
	defaultGatewayPort        = "8080"
	defaultGatewayURL         = "/terminus/v1/"
	defaultGatewaySwaggerFile = "pkg/pb/service.swagger.json"

	// Authz
	defaultAuthzEnable  = false
	defaultAuthzAddress = "authz.atlas"
	defaultAuthzPort    = "5555"

	// Audit Logging
	defaultAuditEnable  = false
	defaultAuditAddress = "audit.atlas"
	defaultAuditPort    = "5555"

	// Tagging
	defaultTaggingEnable  = false
	defaultTaggingAddress = "tagging.atlas"
	defaultTaggingPort    = "5555"

	// Health
	defaultInternalEnable    = true
	defaultInternalAddress   = "0.0.0.0"
	defaultInternalPort      = "8081"
	defaultInternalHealth    = "/healthz"
	defaultInternalReadiness = "/ready"

	defaultConfigDirectory = "deploy/"
	defaultConfigFile      = ""
	defaultSecretFile      = ""
	defaultApplicationID   = "atlas.feature.flag"

	// Logging
	defaultLoggingLevel = "debug"
)

var (
	defaultJwtLabels = []string{"account_id", "user_id"}
)

var (
	// define flag overrides
	flagKubeConfig = pflag.String("kubeconfig", defaultKubeConfig, "path to kubernetes cluster config file")

	flagServerAddress = pflag.String("server.address", defaultServerAddress, "adress of gRPC server")
	flagServerPort    = pflag.String("server.port", defaultServerPort, "port of gRPC server")

	flagGatewayEnable      = pflag.Bool("gateway.enable", defaultGatewayEnable, "enable gatway")
	flagGatewayAddress     = pflag.String("gateway.address", defaultGatewayAddress, "address of gateway server")
	flagGatewayPort        = pflag.String("gateway.port", defaultGatewayPort, "port of gateway server")
	flagGatewayURL         = pflag.String("gateway.endpoint", defaultGatewayURL, "endpoint of gateway server")
	flagGatewaySwaggerFile = pflag.String("gateway.swaggerFile", defaultGatewaySwaggerFile, "directory of swagger.json file")

	flagAuthzEnable  = pflag.Bool("atlas.authz.enable", defaultAuthzEnable, "enable application with authorization")
	flagAuthzAddress = pflag.String("atlas.authz.address", defaultAuthzAddress, "address or FQDN of the authorization service")
	flagAuthzPort    = pflag.String("atlas.authz.port", defaultAuthzPort, "port of the authorization service")

	flagAuditEnable  = pflag.Bool("atlas.audit.enable", defaultAuditEnable, "enable logging of gRPC requests on Atlas audit service")
	flagAuditAddress = pflag.String("atlas.audit.address", defaultAuditAddress, "address or FQDN of Atlas audit log service")
	flagAuditPort    = pflag.String("atlas.audit.port", defaultAuditPort, "port of Atlas audit log service")

	flagTaggingEnable  = pflag.Bool("atlas.tagging.enable", defaultTaggingEnable, "enable tagging")
	flagTaggingAddress = pflag.String("atlas.tagging.address", defaultTaggingAddress, "address or FQDN of Atlas tagging service")
	flagTaggingPort    = pflag.String("atlas.tagging.port", defaultTaggingPort, "port of Atlas tagging service")

	flagInternalEnable    = pflag.Bool("internal.enable", defaultInternalEnable, "enable internal http server")
	flagInternalAddress   = pflag.String("internal.address", defaultInternalAddress, "address of internal http server")
	flagInternalPort      = pflag.String("internal.port", defaultInternalPort, "port of internal http server")
	flagInternalHealth    = pflag.String("internal.health", defaultInternalHealth, "endpoint for health checks")
	flagInternalReadiness = pflag.String("internal.readiness", defaultInternalReadiness, "endpoint for readiness checks")

	flagConfigDirectory = pflag.String("config.source", defaultConfigDirectory, "directory of the configuration file")
	flagConfigFile      = pflag.String("config.file", defaultConfigFile, "directory of the configuration file")
	flagSecretFile      = pflag.String("config.secret.file", defaultSecretFile, "directory of the secrets configuration file")
	flagApplicationID   = pflag.String("app.id", defaultApplicationID, "identifier for the application")

	flagLoggingLevel = pflag.String("logging.level", defaultLoggingLevel, "log level of application")

	jwtToLabels = pflag.StringSlice("jwt.labels", defaultJwtLabels, "jwt field matched to labels")
)
