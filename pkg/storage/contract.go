package storage

import "github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"

type (
	Storage interface {
		Define(FeatureFlagDefinition)
		Override(FeatureFlagOverride)
		Find(featureName string, labels map[string]string) *pb.FeatureFlag
		RemoveDefinition(featureName string)
		RemoveOverride(featureName string, labels map[string]string)
	}
	FeatureFlagDefinition struct {
		FeatureName  string
		DefaultValue string
	}
	FeatureFlagOverride struct {
		FeatureName string
		Value       string
		Origin      string
		Priority    int
		Labels      map[string]string
	}
)
