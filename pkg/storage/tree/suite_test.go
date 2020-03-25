package tree

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Storage Tree Suite")
}

var _ = BeforeSuite(func(done Done) {
	logrus.SetOutput(GinkgoWriter)
	logrus.SetLevel(logrus.TraceLevel)
	close(done)
})
