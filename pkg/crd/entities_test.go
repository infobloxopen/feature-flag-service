package crd

import "testing"

//due to unsafe type casting we have to make sure that all functions was implemented for concrete struct
func TestFeatureFlag_CRBaseImplemented(t *testing.T) {
	var ff CRBase = &FeatureFlag{}
	t.Log(ff)
}

func TestFeatureFlagOverride_CRBaseImplemented(t *testing.T) {
	var ff CRBase = &FeatureFlagOverride{}
	t.Log(ff)
}
