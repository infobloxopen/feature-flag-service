package crd

import (
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// K8sRegex (DNS-1123)
	K8sRegex = `^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	// AllRegex ...
	AllRegex = `.*`
)

// CrdDefinition ...
type CrdDefinition struct {
	Singular   string
	Plural     string
	Group      string
	Kind       string
	Version    string
	Validation apiextv1beta1.CustomResourceValidation // This field si k8s-specific, it should probably be more generic.
	GoTypes    []runtime.Object
}

// AddKnownTypes ...
func (cd CrdDefinition) AddKnownTypes(scheme *runtime.Scheme) error {
	schemeGroupVersion := schema.GroupVersion{Group: cd.Group, Version: cd.Version}
	scheme.AddKnownTypes(schemeGroupVersion, cd.GoTypes...)
	meta_v1.AddToGroupVersion(scheme, schemeGroupVersion)
	return nil
}

// FeatureFlagCrdDefinition ...
var FeatureFlagCrdDefinition = CrdDefinition{
	Singular: "featureflag",
	Plural:   "featureflags",
	Group:    "terminus.infoblox.com",
	Kind:     "FeatureFlag",
	Version:  "v1",
	Validation: apiextv1beta1.CustomResourceValidation{
		OpenAPIV3Schema: &apiextv1beta1.JSONSchemaProps{
			Properties: map[string]apiextv1beta1.JSONSchemaProps{
				"feature_id": {
					Type:    "string",
					Pattern: K8sRegex,
				},
				"value": {
					Type:    "string",
					Pattern: K8sRegex,
				},
			},
		},
	},
	GoTypes: []runtime.Object{&FeatureFlag{}, &FeatureFlagList{}},
}

// FeatureFlagOverrideCrdDefinition ...
var FeatureFlagOverrideCrdDefinition = CrdDefinition{
	Singular: "featureflagoverride",
	Plural:   "featureflagoverrides",
	Group:    "terminus.infoblox.com",
	Kind:     "FeatureFlagOverride",
	Version:  "v1",
	Validation: apiextv1beta1.CustomResourceValidation{
		OpenAPIV3Schema: &apiextv1beta1.JSONSchemaProps{
			Properties: map[string]apiextv1beta1.JSONSchemaProps{
				"feature_id": {
					Type:    "string",
					Pattern: K8sRegex,
				},
				"override_name": {
					Type:    "string",
					Pattern: AllRegex,
				},
				"priority": {
					Type: "integer",
				},
				"value": {
					Type:    "string",
					Pattern: AllRegex,
				},
			},
		},
	},
	GoTypes: []runtime.Object{&FeatureFlagOverride{}, &FeatureFlagOverrideList{}},
}
