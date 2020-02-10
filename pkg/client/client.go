package client

import (
	"time"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/crd"
)

type CrdClient struct {
	cl     *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

func getClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

// CreateCRD Create the CRD resource, ignore error if it already exists
func CreateCRD(clientset clientset.Interface, crdDef crd.CrdDefinition) error {
	FullCRDName := crdDef.Plural + "." + crdDef.Group

	crd := &apiextv1beta1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{Name: FullCRDName},
		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
			Group:   crdDef.Group,
			Version: crdDef.Version,
			Scope:   apiextv1beta1.NamespaceScoped,
			Names: apiextv1beta1.CustomResourceDefinitionNames{
				Singular: crdDef.Singular,
				Plural:   crdDef.Plural,
				Kind:     crdDef.Kind,
			},
			Validation: &crdDef.Validation,
		},
	}

	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	// Wait for the CRD to be created before we use it (only needed if it's a new one)
	time.Sleep(3 * time.Second)
	return err // note that wait time is done in calling function
}

func ConnectToCluster(kube string, crdDef crd.CrdDefinition) CrdClient {
	config, err := getClientConfig(kube)
	if err != nil {
		panic(err.Error())
	}

	// create clientset and create our CRD, this only need to run once
	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// note: if the CRD exist our CreateCRD function is set to exit without an error
	err = CreateCRD(clientset, crdDef)
	if err != nil {
		panic(err)
	}

	// Create a new clientset which include our CRD schema
	crdcs, scheme, err := NewClient(config, crdDef)
	if err != nil {
		panic(err)
	}

	crdclient := NewCrdClient(crdcs, scheme, "", crdDef)

	return *crdclient
}

func (f *CrdClient) NewListWatch() *cache.ListWatch {
	return cache.NewListWatchFromClient(f.cl, f.plural, f.ns, fields.Everything())
}

// NewClient ...
func NewClient(cfg *rest.Config, crdDef crd.CrdDefinition) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(crdDef.AddKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	schemeGroupVersion := schema.GroupVersion{Group: crdDef.Group, Version: crdDef.Version}
	config := *cfg
	config.GroupVersion = &schemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{ // DirectCodecFactory deprecated
		CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}
	return client, scheme, nil
}

func NewCrdClient(cl *rest.RESTClient, scheme *runtime.Scheme, namespace string, crdDef crd.CrdDefinition) *CrdClient {
	return &CrdClient{cl: cl, ns: namespace, plural: crdDef.Plural,
		codec: runtime.NewParameterCodec(scheme)}
}
