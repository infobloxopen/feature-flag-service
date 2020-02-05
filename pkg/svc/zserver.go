package svc

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
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
func (server) List(context.Context, *pb.ListFeatureFlagRequest) (*pb.ListFeatureFlagResponse, error) {
	// TODO: business logic
	return &pb.ListFeatureFlagResponse{}, nil
}

// Read will return a particular value for the requested feature flag
func (server) Read(context.Context, *pb.ReadFeatureFlagRequest) (*pb.ReadFeatureFlagResponse, error) {
	// TODO: business logic
	return &pb.ReadFeatureFlagResponse{}, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer() (pb.AtlasFeatureFlagServer, error) {
	return &server{}, nil
}
