package controllers

import (
	"github.com/infobloxopen/feature-flag-service/signals"
	featureflagv1 "github.com/infobloxopen/feature-flag-service/api/v1"
	"github.com/infobloxopen/feature-flag-service/controllers"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var (
	Scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(Scheme)

	_ = featureflagv1.AddToScheme(Scheme)
	// +kubebuilder:scaffold:scheme
}

func StartKubeController(exitSignals signals.ExitSignals, kubeConfig *rest.Config, scheme *runtime.Scheme, logger logr.Logger) (mgr ctrl.Manager, err error) {
	// Set the logger used by controller-runtime
	ctrl.SetLogger(logger)
	setupLog := ctrl.Log.WithName("kubecontroller")

	mgr, err = controllers.NewDefaultManager(exitSignals, kubeConfig, scheme, false, logger)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return nil, err
	}

	return mgr, nil
}
