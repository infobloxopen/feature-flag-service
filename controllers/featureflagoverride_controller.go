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
)

const (
	FeatureFlagOverrideFeatureNameKey = ".spec.featureName"
)

// FeatureFlagOverrideReconciler reconciles a FeatureFlagOverride object
type FeatureFlagOverrideReconciler struct {
	client.Client
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

	if err := mgr.GetFieldIndexer().IndexField(&ffv1.FeatureFlagOverride{}, FeatureFlagOverrideFeatureNameKey, func(rawObj runtime.Object) []string {
		obj := rawObj.(*ffv1.FeatureFlagOverride)
		if obj.Spec.FeatureName == "" {
			return nil
		}
		return []string{obj.Spec.FeatureName}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&ffv1.FeatureFlagOverride{}).
		Complete(r)
}
