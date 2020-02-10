package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"

	"k8s.io/client-go/tools/cache"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/crd"
)

const (
	kubeconfig string = "default"
	// LWFeatureFlag ...
	LWFeatureFlag = "Feature Flag List Watcher"

	// LWFeatureFlagOverride ...
	LWFeatureFlagOverride = "FeatureFlagOverride List Watcher"
)

var (
	// Watchers ...
	Watchers = map[string]*cache.ListWatch{}
)

// boilerplate for controller generation
func getCRDControllers() map[string]cache.Controller {
	refresh := time.Minute * 5

	crdChangeHandler := cache.FilteringResourceEventHandler{
		FilterFunc: crd.ShouldProcessCR,
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc:    crd.HandleCRAdd,
			DeleteFunc: crd.HandleCRDelete,
			UpdateFunc: crd.HandleCRUpdate},
	}

	// Watch for changes in PARGs objects and fire Add, Delete, Update callbacks
	_, featureflagcontroller := cache.NewInformer(
		Watchers[LWFeatureFlag],
		&crd.FeatureFlag{},
		refresh,
		crdChangeHandler,
	)

	// Watch for changes in PARGs objects and fire Add, Delete, Update callbacks
	_, featureflagoverridecontroller := cache.NewInformer(
		Watchers[LWFeatureFlagOverride],
		&crd.FeatureFlagOverride{},
		refresh,
		crdChangeHandler,
	)

	return map[string]cache.Controller{
		crd.FeatureFlagCrdDefinition.Kind:         featureflagcontroller,
		crd.FeatureFlagOverrideCrdDefinition.Kind: featureflagoverridecontroller,
	}
}

// creates a list of feature flag controllers for Kubernetes CRD watching
// and connects them to the cluster;
// as of Kubernetes 1.9, these controllers will not return to the calling
// code on error, so we don't have a way of reacting to CRD deletion
func watchCRs() {
	featureflags := ConnectToCluster(kubeconfig, crd.FeatureFlagCrdDefinition)
	featureflagoverrides := ConnectToCluster(kubeconfig, crd.FeatureFlagOverrideCrdDefinition)

	Watchers[LWFeatureFlag] = featureflags.NewListWatch()
	Watchers[LWFeatureFlagOverride] = featureflagoverrides.NewListWatch()

	controllers := getCRDControllers()

	stopfeatureflag := make(chan struct{})
	go controllers[crd.FeatureFlagCrdDefinition.Kind].Run(stopfeatureflag)

	stopfeatureflagoverride := make(chan struct{})
	go controllers[crd.FeatureFlagOverrideCrdDefinition.Kind].Run(stopfeatureflagoverride)
}

// addJwtToContext adds a JWT token to context as a Bearer auth token
func addJwtToContext(ctx context.Context, jwt string) context.Context {
	// AppendToOutgoingContext creates a clone of context plus the new auth header
	newCtx := metadata.AppendToOutgoingContext(ctx, "Terminus", fmt.Sprintf("Bearer %s", jwt))

	return newCtx
}
