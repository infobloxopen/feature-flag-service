package crd

import (
	"github.com/sirupsen/logrus"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CRBase ...
type CRBase interface {
	GetTypeMeta() *meta_v1.TypeMeta
	GetObjectMeta() *meta_v1.ObjectMeta
	GetAppName() string
	GetKind() string
	GetName() string
	GetVersion() string
	GetCreationTime() string
}

// CRBaseImpl ...
type CRBaseImpl struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
}

// GetTypeMeta ...
func (pbi *CRBaseImpl) GetTypeMeta() *meta_v1.TypeMeta {
	return &pbi.TypeMeta
}

// GetObjectMeta ...
func (pbi *CRBaseImpl) GetObjectMeta() *meta_v1.ObjectMeta {
	return &pbi.ObjectMeta
}

// GetKind ...
func (pbi *CRBaseImpl) GetKind() string {
	return pbi.TypeMeta.Kind
}

// GetName ...
func (pbi *CRBaseImpl) GetName() string {
	return pbi.Name
}

// GetVersion ...
func (pbi *CRBaseImpl) GetVersion() string {
	return pbi.ResourceVersion
}

// GetAppName ...
func (pbi *CRBaseImpl) GetAppName() string {
	return pbi.ObjectMeta.Name
}

// GetCreationTime ...
func (pbi *CRBaseImpl) GetCreationTime() string {
	return pbi.CreationTimestamp.String()
}

// FeatureFlag identified with an id that contains feature flag value
// this struct is also a pb.imitation so we can have lowercase manifest files
type FeatureFlag struct {
	CRBaseImpl
	FeatureID string `json:"feature_id,omitempty"`
	Value     string `json:"value,omitempty"`
}

// GetKind ...
func (pbi *FeatureFlag) GetKind() string {
	return FeatureFlagCrdDefinition.Kind
}

// FeatureFlagList is a list of feature flags
type FeatureFlagList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []FeatureFlag `json:"items"`
}

// DeepCopyObject necessary for implementation
func (in *FeatureFlag) DeepCopyObject() runtime.Object {
	logrus.Debug("FeatureFlag DEEPCOPY")
	return &FeatureFlag{
		in.CRBaseImpl,
		in.FeatureID,
		in.Value,
	}
}

// DeepCopyObject necessary for implementation
func (in *FeatureFlagList) DeepCopyObject() runtime.Object {
	logrus.Debug("FeatureFlagList DEEPCOPY")
	return &FeatureFlagList{
		in.TypeMeta,
		in.ListMeta,
		in.Items,
	}
}

// FeatureFlagOverride identified with an id that contains feature flag value and override details
// this struct is also a pb.imitation so we can have lowercase manifest files
type FeatureFlagOverride struct {
	CRBaseImpl
	FeatureID    string `json:"feature_id,omitempty"`
	OverrideName string `json:"override_name,omitempty"`
	Priority     int    `json:"priority,omitempty"`
	Value        string `json:"value,omitempty"`
	Labels       map[string]string
}

// GetKind ...
func (pbi *FeatureFlagOverride) GetKind() string {
	return FeatureFlagOverrideCrdDefinition.Kind
}

// FeatureFlagOverrideList is a list of feature flag overrides
type FeatureFlagOverrideList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []FeatureFlagOverride `json:"items"`
}

// DeepCopyObject necessary for implementation
func (in *FeatureFlagOverride) DeepCopyObject() runtime.Object {
	logrus.Debug("FeatureFlagOverride DEEPCOPY")
	return &FeatureFlagOverride{
		in.CRBaseImpl,
		in.FeatureID,
		in.OverrideName,
		in.Priority,
		in.Value,
		in.CRBaseImpl.Labels,
	}
}

// DeepCopyObject necessary for implementation
func (in *FeatureFlagOverrideList) DeepCopyObject() runtime.Object {
	logrus.Debug("FeatureFlagOverrideList DEEPCOPY")
	return &FeatureFlagOverrideList{
		in.TypeMeta,
		in.ListMeta,
		in.Items,
	}
}
