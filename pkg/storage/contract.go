package storage

import (
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/crd"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
)

type (
	// Storage ...
	Storage interface {
		Define(crd.FeatureFlag)
		Override(crd.FeatureFlagOverride)
		Find(featureName string, labels map[string]string) *pb.FeatureFlag
		RemoveDefinition(featureName string)
		RemoveOverride(featureName string, labels map[string]string)
		FindAll(labels map[string]string) []*pb.FeatureFlag
	}
)
