package controllers

import (
	"github.com/infobloxopen/feature-flag-service/signals"
	ffv1 "github.com/infobloxopen/feature-flag-service/api/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

func NewDefaultManager(exitSignals signals.ExitSignals, restConfig *rest.Config, scheme *runtime.Scheme, webhookEnabled bool, logger logr.Logger) (ctrl.Manager, error) {
	// Set the logger used by controller-runtime
	ctrl.SetLogger(logger)
	setupLog := ctrl.Log.WithName("kubecontroller")

	// Initialize the controller-runtime.Manager
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:             scheme,
		LeaderElection:     false,
		MetricsBindAddress: ":9000",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return nil, err
	}

	// Setup FeatureFlag and FeatureFlagOverride CRD controllers
	if err = (&FeatureFlagReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("FeatureFlag"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("featureflag-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FeatureFlag")
		return nil, err
	}

	if err = (&FeatureFlagOverrideReconciler{
		Client:   mgr.GetClient(),
		Log:      ctrl.Log.WithName("controllers").WithName("FeatureFlagOverride"),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("featureflagoverride-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "FeatureFlagOverride")
		return nil, err
	}

	// Setup webhooks for FeatureFlag and FeatureFlagOverride CRDs
	if webhookEnabled {
		setupLog.Info("setting up webhooks with manager")
		if err = (&ffv1.FeatureFlag{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "FeatureFlag")
			return nil, err
		}

		if err = (&ffv1.FeatureFlagOverride{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "FeatureFlagOverride")
			return nil, err
		}
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	go func() {
		if err := mgr.Start(exitSignals.StopCh()); err != nil {
			setupLog.Error(err, "problem running manager")
			exitSignals.DoneCh() <- err
		}
	}()
	return mgr, nil
}
