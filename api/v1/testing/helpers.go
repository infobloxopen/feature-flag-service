package featureflag

import (
	"context"

	. "github.com/Infoblox-CTO/atlas-app-definition-controller/tests/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	featureflagv1 "github.com/Infoblox-CTO/atlas.feature.flag/api/v1"
)

func init() {
	// Adding to the runtime Scheme in helpers
	featureflagv1.AddToScheme(Scheme)
}

func CreateOrUpdateFeatureFlag(k8sClient client.Client, name string, namespace string, featureID string, value string) error {
	obj := &featureflagv1.FeatureFlag{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.TODO(), k8sClient, obj, func() error {
		obj.Spec.Value = value
		obj.Spec.FeatureID = featureID

		return nil
	})
	return err
}

func DeleteFeatureFlag(k8sClient client.Client, name string, namespace string) error {
	obj := &featureflagv1.FeatureFlag{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return k8sClient.Delete(context.TODO(), obj)
}

func CreateOrUpdateFeatureFlagOverride(k8sClient client.Client, name string, namespace string, featureID string, value string, priority int, labels map[string]string) error {
	obj := &featureflagv1.FeatureFlagOverride{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.TODO(), k8sClient, obj, func() error {
		if obj.Spec.LabelSelector == nil {
			obj.Spec.LabelSelector = &metav1.LabelSelector{
				MatchLabels: map[string]string{},
			}
		}
		obj.Spec.LabelSelector.MatchLabels = labels
		obj.Spec.Value = value
		obj.Spec.FeatureID = featureID
		obj.Spec.Priority = priority
		obj.Spec.OverrideName = name

		return nil
	})
	return err
}

func DeleteFeatureFlagOverride(k8sClient client.Client, name string, namespace string) error {
	obj := &featureflagv1.FeatureFlagOverride{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return k8sClient.Delete(context.TODO(), obj)
}
