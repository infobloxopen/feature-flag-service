package storage

import (
	// "github.com/Infoblox-CTO/atlas.feature.flag/pkg/crd"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"

	ffv1 "github.com/Infoblox-CTO/atlas.feature.flag/api/v1"
)

type (
	// Storage ...
	Storage interface {
		Define(*ffv1.FeatureFlag)
		Override(*ffv1.FeatureFlagOverride)
		Find(featureName string, labels map[string]string) *pb.FeatureFlag
		RemoveDefinition(*ffv1.FeatureFlag)
		RemoveOverride(*ffv1.FeatureFlagOverride)
		FindAll(labels map[string]string) []*pb.FeatureFlag
	}
)
