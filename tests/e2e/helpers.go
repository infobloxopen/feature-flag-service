package e2e

import (
	"fmt"
	"os"

	options "github.com/Infoblox-CTO/atlas-app-definition-controller/pkg/featureflag/client/options"
	"github.com/Infoblox-CTO/atlas-app-definition-controller/pkg/util/net/addr"
	"github.com/Infoblox-CTO/atlas-app-definition-controller/pkg/util/signals"
	. "github.com/Infoblox-CTO/atlas-app-definition-controller/tests/helpers"
	featureflagv1 "github.com/Infoblox-CTO/atlas.feature.flag/api/v1"
	"github.com/Infoblox-CTO/atlas.feature.flag/cmd/server/controllers"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func init() {
	featureflagv1.AddToScheme(Scheme)
}

func CreateFeatureFlagController(exitSignals signals.ExitSignals, restConfig *rest.Config, logger logr.Logger) (ctrl.Manager, error) {
	mgr, err := controllers.StartKubeController(exitSignals, restConfig, Scheme, logger)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

func CreateFeatureFlagService(exitSignals signals.ExitSignals, kubeConfigPath string, client client.Client) (*Server, error) {
	os.Setenv("KUBECONFIG", kubeConfigPath)
	viper.Set("kubeconfig", kubeConfigPath)

	Info("starting feature-flag service")

	ffPort, ffHost, err := addr.Suggest()
	ffHost = ""
	if err != nil {
		return nil, err
	}
	ffBindAddress := fmt.Sprintf("%s:%d", ffHost, ffPort)

	ffGatewayPort, ffGatewayHost, err := addr.Suggest()
	if err != nil {
		return nil, err
	}
	ffGatewayBindAddress := fmt.Sprintf("%s:%d", ffGatewayHost, ffGatewayPort)

	fflogger := NewLogger()
	fflogger.SetOutput(GinkgoWriter)

	jwtTokenFields := []string{"account_id", "user_id"}
	jwtTokenFields = []string{}
	server, err := NewFeatureFlagServer(client, ffBindAddress, ffGatewayBindAddress, "/terminus/v1", jwtTokenFields, fflogger)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func CreateFeatureFlagClient(bindAddress string, insecure bool) (pb.AtlasFeatureFlagClient, error) {
	o := &options.Options{
		Address:  bindAddress,
		Insecure: insecure,
	}
	c, err := o.Config()
	if err != nil {
		return nil, err
	}
	return c.FeatureFlagClient()
}

func StartKubernetesEnvironment(crdPaths []string) (testEnv *envtest.Environment, restConfig *rest.Config, k8sClient client.Client, kubeConfigPath string, err error) {
	logrus.SetOutput(GinkgoWriter)
	logrus.SetLevel(logrus.DebugLevel)

	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	t := true
	if os.Getenv("TEST_USE_EXISTING_CLUSTER") == "true" {
		testEnv = &envtest.Environment{
			UseExistingCluster: &t,
		}
	} else {
		err = os.Unsetenv("KUBECONFIG")
		if err != nil {
			return nil, nil, nil, "", err
		}
		testEnv = &envtest.Environment{
			AttachControlPlaneOutput: false,
			CRDDirectoryPaths:        crdPaths,
		}
	}

	restConfig, err = testEnv.Start()
	if err != nil {
		return nil, nil, nil, "", err
	}

	kubeConfigPath, err = WriteKubeConfig(restConfig, os.Getenv("BUILD_TEMP"))
	if err != nil {
		return nil, nil, nil, "", err
	}

	k8sClient, err = CreateKubeClient(restConfig, Scheme)
	if err != nil {
		return nil, nil, nil, "", err
	}
	return testEnv, restConfig, k8sClient, kubeConfigPath, err
}
