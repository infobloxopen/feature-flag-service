package svc

import (
	"context"
	"errors"
	"sort"

	"github.com/go-logr/logr"
	"github.com/golang/protobuf/ptypes/empty"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/infobloxopen/atlas-app-toolkit/auth"

	featureflagv1 "github.com/infobloxopen/feature-flag-service/api/v1"
	ctrls "github.com/infobloxopen/feature-flag-service/controllers"
	"github.com/infobloxopen/feature-flag-service/pkg/pb"
)

// Default implementation of the AtlasFeatureFlag server interface
type (
	server struct {
		client    ctrlclient.Client
		jwtLabels []string
		logr.Logger
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
func (s *server) getLabels(ctx context.Context, request labelsProvider) (labels.Labels, error) {
	resultLabels := labels.Set{}
	//get labels from labels
	for labelName, labelValue := range request.GetLabels() {
		resultLabels[labelName] = labelValue
	}
	//get labels from jwt
	for _, labelName := range s.jwtLabels {
		labelValue, err := auth.GetJWTField(ctx, labelName, nil)
		if err != nil {
			continue
		}
		resultLabels[labelName] = labelValue
	}
	return resultLabels, nil
}

func (s *server) getFeatureFlags(ctx context.Context, featureName string) ([]*featureflagv1.FeatureFlag, error) {
	// List all features with matching FeatureName (indexed query from local cache)
	// reset the listOptions (won't make memory of slice elements GC eligible until function exits)
	var listOptions []ctrlclient.ListOption
	if featureName != "" {
		listOptions = append(listOptions, ctrlclient.MatchingFields{ctrls.FeatureFlagNameKey: featureName})
	}

	ffList := &featureflagv1.FeatureFlagList{}
	err := s.client.List(ctx, ffList, listOptions...)
	if err != nil {
		return nil, err
	}

	matchedFlags := make([]*featureflagv1.FeatureFlag, 0, len(ffList.Items))
	for _, ff := range ffList.Items {
		matchedFlags = append(matchedFlags, ff.DeepCopy())
	}

	return matchedFlags, nil
}

func (s *server) getFeatureFlagOverrides(ctx context.Context, featureName string, labelSet labels.Labels) ([]*featureflagv1.FeatureFlagOverride, error) {
	var err error
	s.V(1).Info("getting FeatureFlagOverrides", "featureName", featureName, "labels", labelSet)
	// List all overrides with matching FeatureName (indexed query from local cache)
	var listOptions []ctrlclient.ListOption
	if featureName != "" {
		listOptions = append(listOptions, ctrlclient.MatchingFields{ctrls.FeatureFlagOverrideFeatureNameKey: featureName})
	}
	ffoList := &featureflagv1.FeatureFlagOverrideList{}
	err = s.client.List(ctx, ffoList, listOptions...)
	if err != nil {
		return nil, err
	}
	// Iterate and test each override's labelSelector against the label query
	var matchedOverrides []featureflagv1.FeatureFlagOverride
	for _, ffo := range ffoList.Items {
		selector, err := metav1.LabelSelectorAsSelector(ffo.Spec.LabelSelector)
		if err != nil {
			return nil, err
		}
		if !selector.Matches(labelSet) {
			continue
		}
		matchedOverrides = append(matchedOverrides, ffo)
		s.V(1).Info("matched FeatureFlagOverride", "featureName", featureName, "labels", labelSet,
			"selector", ffo.Spec.LabelSelector.MatchLabels, "name", ffo.Name, "priority", ffo.Spec.Priority, "value", ffo.Spec.Value)

	}
	// Sort matching overrides by higher priority
	sort.Slice(matchedOverrides, func(i, j int) bool {
		return matchedOverrides[i].Spec.Priority < matchedOverrides[j].Spec.Priority
	})
	matchedOverridesP := make([]*featureflagv1.FeatureFlagOverride, 0, len(matchedOverrides))
	for _, ffo := range matchedOverrides {
		s.V(1).Info("sorted FeatureFlagOverride", "featureName", featureName, "labels", labelSet,
			"selector", ffo.Spec.LabelSelector.MatchLabels, "name", ffo.Name, "priority", ffo.Spec.Priority, "value", ffo.Spec.Value)
		matchedOverridesP = append(matchedOverridesP, &ffo)
	}
	s.V(1).Info("matched FeatureFlagOverrides", "featureName", featureName, "labels", labelSet, "count", len(matchedOverrides))
	return matchedOverridesP, nil
}

// List will return a list of all feature flags
func (s *server) List(ctx context.Context, req *pb.ListFeatureFlagsRequest) (*pb.ListFeatureFlagsResponse, error) {
	s.V(1).Info("listing FeatureFlags/FeatureFlagOverrides", "labels", req.GetLabels())
	labels, err := s.getLabels(ctx, req)
	if err != nil {
		return nil, err
	}

	matchedFlags, err := s.getFeatureFlags(ctx, "")
	if err != nil {
		return nil, err
	}

	featureFlags := make([]*pb.FeatureFlag, 0, len(matchedFlags))
	for _, ff := range matchedFlags {
		matchedOverrides, err := s.getFeatureFlagOverrides(ctx, ff.Name, labels)
		if err != nil {
			return nil, err
		}

		ffPB := &pb.FeatureFlag{
			FeatureName: ff.Name,
		}

		// return highest priority override that matched
		if len(matchedOverrides) > 0 {
			ffo := matchedOverrides[len(matchedOverrides)-1]
			ffPB.Value = ffo.Spec.Value
			objectKey, err := ctrlclient.ObjectKeyFromObject(ffo)
			if err != nil {
				return nil, err
			}
			ffPB.Origin = "FeatureFlagOverride:" + objectKey.String()
			featureFlags = append(featureFlags, ffPB)
			continue
		}

		ffPB.Value = ff.Spec.Value
		objectKey, err := ctrlclient.ObjectKeyFromObject(ff)
		if err != nil {
			return nil, err
		}

		ffPB.Origin = "FeatureFlag:" + objectKey.String()
		featureFlags = append(featureFlags, ffPB)
	}

	s.V(1).Info("found FeatureFlags/FeatureFlagOverrides", "labels", req.GetLabels(), "count", len(featureFlags))
	return &pb.ListFeatureFlagsResponse{
		Results: featureFlags,
	}, nil
}

// Read will return a particular value for the requested feature flag
func (s *server) Read(ctx context.Context, req *pb.ReadFeatureFlagRequest) (*pb.ReadFeatureFlagResponse, error) {
	s.V(1).Info("finding value for feature request", "featureName", req.GetFeatureName(), "labels", req.GetLabels())
	labels, err := s.getLabels(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := &pb.ReadFeatureFlagResponse{
		Result: &pb.FeatureFlag{
			FeatureName: req.GetFeatureName(),
		},
	}
	matchedOverrides, err := s.getFeatureFlagOverrides(ctx, req.GetFeatureName(), labels)
	if err != nil {
		return nil, err
	}
	// return highest priority override that matched
	if len(matchedOverrides) > 0 {
		ffo := matchedOverrides[len(matchedOverrides)-1]
		s.V(1).Info("FeatureFlagOverride selected", "featureName", ffo.Spec.FeatureName, "selector", ffo.Spec.LabelSelector.MatchLabels, "name", ffo.Name, "priority", ffo.Spec.Priority, "value", ffo.Spec.Value)
		resp.Result.Value = ffo.Spec.Value
		objectKey, err := ctrlclient.ObjectKeyFromObject(ffo)
		if err != nil {
			return nil, err
		}
		resp.Result.Origin = "FeatureFlagOverride:" + objectKey.String()
		return resp, nil
	}
	matchedFlags, err := s.getFeatureFlags(ctx, req.GetFeatureName())
	if err != nil {
		return nil, err
	}
	// Return an error if more than one Feature is found referencing the same FeatureName
	// this case can be avoided in future with validating webhook
	if len(matchedFlags) > 1 {
		names := make([]string, 0, len(matchedFlags))
		for _, ff := range matchedFlags {
			objectKey, err := ctrlclient.ObjectKeyFromObject(ff)
			if err != nil {
				return nil, err
			}
			names = append(names, objectKey.String())
		}
		return nil, errors.New("multiple Feature resources exist with the same FeatureName")
	}
	// Return an error if no Feature exists with the FeatureName
	if len(matchedFlags) == 0 {
		return nil, errors.New("FeatureFlag or FeatureFlagOverride not found for FeatureName")
	}
	ff := matchedFlags[0]
	resp.Result.Value = ff.Spec.Value
	objectKey, err := ctrlclient.ObjectKeyFromObject(ff)
	if err != nil {
		return nil, err
	}
	resp.Result.Origin = "FeatureFlag:" + objectKey.String()

	return resp, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer(client client.Client, jwtLabels []string) (pb.AtlasFeatureFlagServer, error) {
	return &server{client: client, jwtLabels: jwtLabels, Logger: ctrl.Log.WithName("grpc")}, nil
}
