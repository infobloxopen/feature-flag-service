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
var featureflaglog = logf.Log.WithName("featureflag-resource")

func (r *FeatureFlag) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-terminus-infoblox-com-v1-featureflag,mutating=true,failurePolicy=fail,groups=terminus.infoblox.com,resources=featureflags,verbs=create;update,versions=v1,name=mfeatureflag.kb.io

var _ webhook.Defaulter = &FeatureFlag{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *FeatureFlag) Default() {
	featureflaglog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-terminus-infoblox-com-v1-featureflag,mutating=false,failurePolicy=fail,groups=terminus.infoblox.com,resources=featureflags,versions=v1,name=vfeatureflag.kb.io

var _ webhook.Validator = &FeatureFlag{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *FeatureFlag) ValidateCreate() error {
	featureflaglog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *FeatureFlag) ValidateUpdate(old runtime.Object) error {
	featureflaglog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *FeatureFlag) ValidateDelete() error {
	featureflaglog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
