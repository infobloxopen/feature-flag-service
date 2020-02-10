package crd

import (
	"github.com/sirupsen/logrus"
)

var (
	// CRCacheFn ...
	CRCacheFn = map[string]func(oldObj interface{}, newObj interface{}){
		FeatureFlagCrdDefinition.Kind:         cacheFeatureFlag,
		FeatureFlagOverrideCrdDefinition.Kind: cacheFeatureFlagOverride,
	}
)

// ShouldProcessCR ...
func ShouldProcessCR(obj interface{}) bool {
	cr := obj.(CRBase)
	label := cr.GetObjectMeta().GetLabels()["unique_label_identifier"] // TODO insert label KEY to check for unique receiver

	return label == label // TODO insert label VALUE to be matched
}

// HandleCRDelete ...
func HandleCRDelete(oldCR interface{}) {
	cr := oldCR.(CRBase)
	logrus.Info("CR delete "+cr.GetKind()+" discovered: ", cr.GetAppName()+":"+cr.GetName())
	logrus.Debug("Creation: ", cr.GetCreationTime(), " Version: ", cr.GetVersion())
	logrus.Debug("CR object: ", cr)

	HandleCRUpdate(cr, nil)
}

// HandleCRAdd ...
func HandleCRAdd(newCR interface{}) {
	cr := newCR.(CRBase)
	logrus.Info("CR add "+cr.GetKind()+" discovered: ", cr.GetAppName()+":"+cr.GetName())
	logrus.Debug("Creation: ", cr.GetCreationTime(), " Version: ", cr.GetVersion())
	logrus.Debug("CR object: ", cr)

	HandleCRUpdate(nil, cr)
}

// HandleCRUpdate ...
func HandleCRUpdate(oldCR, newCR interface{}) {
	// "fake" updates is when k8s api sends periodic list as update (though no real updates happened)
	if oldCR == newCR {
		cr := newCR.(CRBase)
		logrus.Debug("Fake CR update "+cr.GetKind()+" discovered, cr: ", cr.GetAppName()+":"+cr.GetName())
		return
	}

	if oldCR != nil && newCR != nil {
		cr := newCR.(CRBase)
		logrus.Info("CR update "+cr.GetKind()+" discovered, new cr: ", cr.GetAppName()+":"+cr.GetName())
		logrus.Debug("Creation: ", cr.GetCreationTime(), " Version: ", cr.GetVersion())
		logrus.Debug("New CR object: ", newCR)
	}

	cr := newCR
	if cr == nil {
		cr = oldCR
	}

	customResource := cr.(CRBase)
	CRCacheFn[customResource.GetKind()](oldCR, newCR)
}

func cacheFeatureFlag(oldCR interface{}, newCR interface{}) {
	if newCR != nil {
		newFeatureFlag := newCR.(FeatureFlag)
		featureFlag := &FeatureFlag{
			newFeatureFlag.CRBaseImpl,
			newFeatureFlag.FeatureID,
			newFeatureFlag.Value,
		}

		logrus.Info(featureFlag) // delete me
		// TODO HERE IS WHERE WE ADD TO OR UPDATE THE CACHE
	} else if oldCR != nil {
		// TODO HERE IS WHERE WE DELETE FROM THE CACHE
	}
}

func cacheFeatureFlagOverride(oldCR interface{}, newCR interface{}) {
	if newCR != nil {
		newFeatureFlagOverride := newCR.(FeatureFlagOverride)
		featureFlagOverride := &FeatureFlagOverride{
			newFeatureFlagOverride.CRBaseImpl,
			newFeatureFlagOverride.FeatureID,
			newFeatureFlagOverride.OverrideName,
			newFeatureFlagOverride.Priority,
			newFeatureFlagOverride.Value,
		}

		logrus.Info(featureFlagOverride) // delete me
		// TODO HERE IS WHERE WE ADD TO OR UPDATE THE CACHE
	} else if oldCR != nil {
		// TODO HERE IS WHERE WE DELETE FROM THE CACHE
	}
}
