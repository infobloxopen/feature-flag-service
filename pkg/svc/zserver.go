package svc

import (
	"context"
	"errors"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/client"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage/tree"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/infobloxopen/atlas-app-toolkit/auth"
	"github.com/sirupsen/logrus"
)

// Default implementation of the AtlasFeatureFlag server interface
type (
	server struct {
		jwtLabels []string
	}
	labelsProvider interface {
		GetLabels() map[string]string
	}
)

// GetVersion returns the current version of the service
func (*server) GetVersion(context.Context, *empty.Empty) (*pb.VersionResponse, error) {
	if Manifest == nil {
		return &pb.VersionResponse{Version: "manifest not found"}, nil
	}
	return &pb.VersionResponse{Version: Manifest.Version}, nil
}

// get map of labels
func (s *server) getLabels(ctx context.Context, request labelsProvider) (map[string]string, error) {
	resultLabels := map[string]string{}
	//get labels from labels
	for labelName, labelValue := range request.GetLabels() {
		resultLabels[labelName] = labelValue
	}
	//get labels from jwt
	for _, labelName := range s.jwtLabels {
		labelValue, err := auth.GetJWTField(ctx, labelName, nil)
		if err != nil {
			return nil, err
		}
		resultLabels[labelName] = labelValue
	}
	return resultLabels, nil
}

// List will return a list of all feature flags
func (s *server) List(ctx context.Context, request *pb.ListFeatureFlagRequest) (*pb.ListFeatureFlagResponse, error) {
	labels, err := s.getLabels(ctx, request)
	if err != nil {
		return nil, err
	}
	featureFlags := client.Cache.FindAll(labels)
	return &pb.ListFeatureFlagResponse{
		Results: featureFlags,
	}, nil
}

// Read will return a particular value for the requested feature flag
func (s *server) Read(ctx context.Context, request *pb.ReadFeatureFlagRequest) (*pb.ReadFeatureFlagResponse, error) {
	labels, err := s.getLabels(ctx, request)
	if err != nil {
		return nil, err
	}
	featureFlag := client.Cache.Find(request.FeatureName, labels)
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
func NewBasicServer(useKCRDs bool, jwtLabels []string) (pb.AtlasFeatureFlagServer, error) {
	client.Cache = tree.NewInMemoryStorage()
	logrus.Debug(client.Cache)
	if useKCRDs {
		client.WatchCRs()
	}
	return &server{jwtLabels: jwtLabels}, nil
}
