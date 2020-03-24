package client

import (
	"github.com/sirupsen/logrus"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/crd"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage"
)

var (
	Cache storage.Storage

	// CRCacheFn ...
	CRCacheFn = map[string]func(oldObj interface{}, newObj interface{}){
		crd.FeatureFlagCrdDefinition.Kind:         cacheFeatureFlag,
		crd.FeatureFlagOverrideCrdDefinition.Kind: cacheFeatureFlagOverride,
	}
)

// ShouldProcessCR ...
func ShouldProcessCR(obj interface{}) bool {
	cr := obj.(crd.CRBase)
	label := cr.GetObjectMeta().GetLabels()["unique_label_identifier"] // TODO insert label KEY to check for unique receiver

	return label == label // TODO insert label VALUE to be matched
}

// HandleCRDelete ...
func HandleCRDelete(oldCR interface{}) {
	cr := oldCR.(crd.CRBase)
	logrus.Info("CR delete "+cr.GetKind()+" discovered: ", cr.GetAppName()+":"+cr.GetName())
	logrus.Debug("Creation: ", cr.GetCreationTime(), " Version: ", cr.GetVersion())
	logrus.Debug("CR object: ", cr)

	HandleCRUpdate(cr, nil)
}

// HandleCRAdd ...
func HandleCRAdd(newCR interface{}) {
	cr := newCR.(crd.CRBase)
	logrus.Info("CR add "+cr.GetKind()+" discovered: ", cr.GetAppName()+":"+cr.GetName())
	logrus.Debug("Creation: ", cr.GetCreationTime(), " Version: ", cr.GetVersion())
	logrus.Debug("CR object: ", cr)

	HandleCRUpdate(nil, cr)
}

// HandleCRUpdate ...
func HandleCRUpdate(oldCR, newCR interface{}) {
	// "fake" updates is when k8s api sends periodic list as update (though no real updates happened)
	if oldCR == newCR {
		cr := newCR.(crd.CRBase)
		logrus.Debug("Fake CR update "+cr.GetKind()+" discovered, cr: ", cr.GetAppName()+":"+cr.GetName())
		return
	}

	if oldCR != nil && newCR != nil {
		cr := newCR.(crd.CRBase)
		logrus.Info("CR update "+cr.GetKind()+" discovered, new cr: ", cr.GetAppName()+":"+cr.GetName())
		logrus.Debug("Creation: ", cr.GetCreationTime(), " Version: ", cr.GetVersion())
		logrus.Debug("New CR object: ", newCR)
	}

	cr := newCR
	if cr == nil {
		cr = oldCR
	}

	customResource := cr.(crd.CRBase)
	CRCacheFn[customResource.GetKind()](oldCR, newCR)
}

func cacheFeatureFlag(oldCR interface{}, newCR interface{}) {
	if oldCR != nil {
		oldFeatureFlag := oldCR.(*crd.FeatureFlag)
		Cache.RemoveDefinition(oldFeatureFlag.FeatureID)
	}

	if newCR != nil {
		newFeatureFlag := newCR.(*crd.FeatureFlag)
		featureFlag := &crd.FeatureFlag{
			newFeatureFlag.CRBaseImpl,
			newFeatureFlag.FeatureID,
			newFeatureFlag.Value,
		}

		Cache.Define(*featureFlag)
	}
}

func cacheFeatureFlagOverride(oldCR interface{}, newCR interface{}) {
	if oldCR != nil {
		oldFeatureFlagOverride := oldCR.(*crd.FeatureFlagOverride)
		Cache.RemoveOverride(oldFeatureFlagOverride.FeatureID, oldFeatureFlagOverride.GetLabels())
	}

	if newCR != nil {
		newFeatureFlagOverride := newCR.(*crd.FeatureFlagOverride)
		featureFlagOverride := &crd.FeatureFlagOverride{
			newFeatureFlagOverride.CRBaseImpl,
			newFeatureFlagOverride.FeatureID,
			newFeatureFlagOverride.OverrideName,
			newFeatureFlagOverride.Priority,
			newFeatureFlagOverride.Value,
			newFeatureFlagOverride.CRBaseImpl.Labels,
		}

		Cache.Override(*featureFlagOverride)
	}
}
