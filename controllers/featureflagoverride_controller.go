/*
Copyright 2020 Infoblox.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ffv1 "github.com/Infoblox-CTO/atlas.feature.flag/api/v1"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage"
)

const (
	FeatureFlagOverrideNameKey      = ".spec.overrideName"
	FeatureFlagOverrideFeatureIDKey = ".spec.featureID"
)

// FeatureFlagOverrideReconciler reconciles a FeatureFlagOverride object
type FeatureFlagOverrideReconciler struct {
	client.Client
	storage.Storage
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=terminus.infoblox.com,resources=featureflagoverrides,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=terminus.infoblox.com,resources=featureflagoverrides/status,verbs=get;update;patch

func (r *FeatureFlagOverrideReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("featureflagoverride", req.NamespacedName)

	return ctrl.Result{}, nil
}

func (r *FeatureFlagOverrideReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(&ffv1.FeatureFlagOverride{}, FeatureFlagOverrideNameKey, func(rawObj runtime.Object) []string {
		obj := rawObj.(*ffv1.FeatureFlagOverride)
		if obj.Spec.OverrideName == "" {
			return nil
		}
		return []string{obj.Spec.OverrideName}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(&ffv1.FeatureFlagOverride{}, FeatureFlagOverrideFeatureIDKey, func(rawObj runtime.Object) []string {
		obj := rawObj.(*ffv1.FeatureFlagOverride)
		if obj.Spec.FeatureID == "" {
			return nil
		}
		return []string{obj.Spec.FeatureID}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&ffv1.FeatureFlagOverride{}).
		Complete(r)
}
