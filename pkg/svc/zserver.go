package svc

import (
	"context"
	"errors"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/client"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage/tree"
)

// Default implementation of the AtlasFeatureFlag server interface
type server struct{}

// GetVersion returns the current version of the service
func (server) GetVersion(context.Context, *empty.Empty) (*pb.VersionResponse, error) {
	if Manifest == nil {
		return &pb.VersionResponse{Version: "manifest not found"}, nil
	}
	return &pb.VersionResponse{Version: Manifest.Version}, nil
}

// List will return a list of all feature flags
func (server) List(ctx context.Context, request *pb.ListFeatureFlagRequest) (*pb.ListFeatureFlagResponse, error) {
	featureFlags := client.Cache.FindAll(request.Labels)
	return &pb.ListFeatureFlagResponse{
		Results: featureFlags,
	}, nil
}

// Read will return a particular value for the requested feature flag
func (server) Read(ctx context.Context, request *pb.ReadFeatureFlagRequest) (*pb.ReadFeatureFlagResponse, error) {
	featureFlag := client.Cache.Find(request.FeatureName, request.Labels)
	if featureFlag == nil {
		return nil, errors.New("feature flag not found")
	}
	return &pb.ReadFeatureFlagResponse{
		Result: &pb.FeatureFlag{
			FeatureName: featureFlag.FeatureName,
			Value:       featureFlag.Value,
			Origin:      featureFlag.Origin,
		},
	}, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer(useKCRDs bool) (pb.AtlasFeatureFlagServer, error) {
	client.Cache = tree.NewInMemoryStorage()
	logrus.Debug(client.Cache)
	if useKCRDs {
		client.WatchCRs()
	}
	return &server{}, nil
}
