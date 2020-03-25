package e2e

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/Infoblox-CTO/atlas-app-definition-controller/pkg/util/signals"
	// "github.com/Infoblox-CTO/atlas-app-definition-controller/tests/featureflag"

	. "github.com/Infoblox-CTO/atlas-app-definition-controller/tests/helpers"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var restConfig *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

var (
	kubeConfigPath   string
	retainKubeConfig = true

	exitSignals signals.TestExitSignals
	fakeSignal  chan struct{}

	featureFlagGRPCServer *Server
	featureFlagController ctrl.Manager
	featureFlagClient     pb.AtlasFeatureFlagClient

	namespace = "default"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Feature Flag Suite")

	//RunSpecsWithDefaultAndCustomReporters(t,
	//	"Controller Suite",
	//	[]Reporter{envtest.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	var err error

	logrus.SetOutput(GinkgoWriter)
	logrus.SetLevel(logrus.TraceLevel)
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	logger := zap.LoggerTo(GinkgoWriter, true)

	fakeSignal = make(chan struct{}, 1)
	var ok bool
	exitSignals, ok = signals.SetupExitHandlersWithExistingStopCh(logger, fakeSignal).(signals.TestExitSignals)
	if !ok {
		err = errors.New("failed to create DebugExitSignals")
	}
	Expect(err).ToNot(HaveOccurred())

	By("bootstrapping test environment")
	crdPaths := []string{
		filepath.Join("..", "..", "helm", "atlas.feature.flag", "crd", "v1beta1"),
	}
	testEnv, restConfig, k8sClient, kubeConfigPath, err = StartKubernetesEnvironment(crdPaths)
	Expect(err).ToNot(HaveOccurred())

	Info("creating feature-flag server instance")

	featureFlagController, err = CreateFeatureFlagController(exitSignals, restConfig, logger)
	Expect(err).ToNot(HaveOccurred())

	featureFlagGRPCServer, err = CreateFeatureFlagService(exitSignals, kubeConfigPath, featureFlagController.GetClient())
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err := featureFlagGRPCServer.Serve()
		Expect(err).ToNot(HaveOccurred())
	}()

	featureFlagClient, err = CreateFeatureFlagClient(featureFlagGRPCServer.BindAddress, true)
	Expect(err).ToNot(HaveOccurred())

	time.Sleep(10 * time.Second)
	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	if !retainKubeConfig {
		os.Remove(kubeConfigPath)
	}
	time.Sleep(2 * time.Second)
	gexec.KillAndWait(5 * time.Second)
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
