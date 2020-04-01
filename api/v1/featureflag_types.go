/*
Copyright 2020 Infoblox.

*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FeatureFlagSpec defines the desired state of FeatureFlag
type FeatureFlagSpec struct {
	// Value is the opaque data for the feature
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}

// +kubebuilder:object:root=true

// FeatureFlag is the Schema for the applications API
// +kubebuilder:printcolumn:name="FeatureName",type=string,JSONPath=`.spec.featureName`
// +kubebuilder:printcolumn:name="Value",type=string,JSONPath=`.spec.value`
type FeatureFlag struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FeatureFlagSpec   `json:"spec,omitempty"`
	Status FeatureFlagStatus `json:"status,omitempty"`
}

// FeatureFlagStatus defines the observed state of FeatureFlag
type FeatureFlagStatus struct {
}

// +kubebuilder:object:root=true

// FeatureFlagList contains a list of FeatureFlag
type FeatureFlagList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FeatureFlag `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FeatureFlag{}, &FeatureFlagList{})
}

// IsBeingDeleted returns true if a deletion timestamp is set
func (r *FeatureFlag) IsBeingDeleted() bool {
	return !r.ObjectMeta.DeletionTimestamp.IsZero()
}
