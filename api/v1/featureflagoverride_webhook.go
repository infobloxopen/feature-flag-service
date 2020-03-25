/*
Copyright 2020 Infoblox.

*/

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var featureflagoverridelog = logf.Log.WithName("featureflagoverride-resource")

func (r *FeatureFlagOverride) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-terminus-infoblox-com-v1-featureflagoverride,mutating=true,failurePolicy=fail,groups=terminus.infoblox.com,resources=featureflagoverrides,verbs=create;update,versions=v1,name=mfeatureflagoverride.kb.io

var _ webhook.Defaulter = &FeatureFlagOverride{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *FeatureFlagOverride) Default() {
	featureflagoverridelog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-terminus-infoblox-com-v1-featureflagoverride,mutating=false,failurePolicy=fail,groups=terminus.infoblox.com,resources=featureflagoverrides,versions=v1,name=vfeatureflagoverride.kb.io

var _ webhook.Validator = &FeatureFlagOverride{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *FeatureFlagOverride) ValidateCreate() error {
	featureflagoverridelog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *FeatureFlagOverride) ValidateUpdate(old runtime.Object) error {
	featureflagoverridelog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *FeatureFlagOverride) ValidateDelete() error {
	featureflagoverridelog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
