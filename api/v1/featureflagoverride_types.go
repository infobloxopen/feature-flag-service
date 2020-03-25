/*
Copyright 2020 Infoblox.

*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FeatureFlagOverrideSpec defines the desired state of FeatureFlagOverride
type FeatureFlagOverrideSpec struct {
	// FeatureID is the unique identifier of the feature
	// +kubebuilder:validation:Required
	FeatureID string `json:"featureID"`

	// Value is the opaque data for the feature
	// +kubebuilder:validation:Required
	Value string `json:"value"`

	// OverrideName is the opaque data for the feature
	// +kubebuilder:validation:Required
	OverrideName string `json:"overrideName"`

	// Priority is the ordering of
	// +kubebuilder:validation:Required
	Priority int `json:"priority"`

	// LabelSelector is a metav1.LabelSelector which matches against labels sent in host requests by config.generator
	// +kubebuilder:validation:Required
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty"`
}

// +kubebuilder:object:root=true

// FeatureFlagOverride is the Schema for the applications API
// +kubebuilder:printcolumn:name="FeatureID",type=string,JSONPath=`.feature_id`
// +kubebuilder:printcolumn:name="Value",type=string,JSONPath=`.value`
// +kubebuilder:printcolumn:name="OverrideName",type=string,JSONPath=`.override_name`
// +kubebuilder:printcolumn:name="Priority",type=integer,JSONPath=`.priority`
type FeatureFlagOverride struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeatureFlagOverrideSpec   `json:"spec,omitempty"`
	Status FeatureFlagOverrideStatus `json:"status,omitempty"`
}

// FeatureFlagOverrideStatus defines the observed state of FeatureFlag
type FeatureFlagOverrideStatus struct {
}

// +kubebuilder:object:root=true

// FeatureFlagOverrideList contains a list of FeatureFlagOverride
type FeatureFlagOverrideList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FeatureFlagOverride `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FeatureFlagOverride{}, &FeatureFlagOverrideList{})
}

// IsBeingDeleted returns true if a deletion timestamp is set
func (r *FeatureFlagOverride) IsBeingDeleted() bool {
	return !r.ObjectMeta.DeletionTimestamp.IsZero()
}

// Labels returns true if a deletion timestamp is set
func (r *FeatureFlagOverride) Labels() map[string]string {
	if r.Spec.LabelSelector == nil {
		return map[string]string{}
	}
	return r.Spec.LabelSelector.MatchLabels
}
