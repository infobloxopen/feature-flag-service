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
	FeatureFlagNameKey = ".metadata.name"
)

// FeatureFlagReconciler reconciles a FeatureFlag object
type FeatureFlagReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=terminus.infoblox.com,resources=featureflags,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=terminus.infoblox.com,resources=featureflags/status,verbs=get;update;patch

func (r *FeatureFlagReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("featureflag", req.NamespacedName)

	return ctrl.Result{}, nil
}

func (r *FeatureFlagReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(&ffv1.FeatureFlag{}, FeatureFlagNameKey, func(rawObj runtime.Object) []string {
		obj := rawObj.(*ffv1.FeatureFlag)
		if obj.ObjectMeta.Name == "" {
			return nil
		}
		return []string{obj.ObjectMeta.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&ffv1.FeatureFlag{}).
		Complete(r)
}
