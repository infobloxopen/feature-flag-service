package main

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/infobloxopen/atlas-app-toolkit/requestid"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/svc"
)

func NewGRPCServer(logger *logrus.Logger) (*grpc.Server, error){
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
	
	// register service implementation with the grpcServer
	s, err := svc.NewBasicServer()
	if err != nil {
		return nil, err
	}
	pb.RegisterAtlasFeatureFlagServer(grpcServer, s)

	return grpcServer, nil
}
