package e2e

import (
	"context"
	"time"

	. "github.com/Infoblox-CTO/atlas.feature.flag/api/v1/testing"
	featureflagpb "github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
)

var _ = Describe("Testing storage tree", func() {
	var err error
	_ = err
	var (
		featureID          = "noa"
		accountOverride    = "acc-40-noa"
		accountBMOverride  = "acc-40-depl-baremetal-noa"
		accountK8SOverride = "acc-40-k8s-enabled-noa"

		featureFlagValue        = "v1.0.0"
		accountOverrideValue    = "v2.0.0"
		accountBMOverrideValue  = "v3.0.0"
		accountK8SOverrideValue = "v4.0.0"

		accountOverridePriority    = 100
		accountBMOverridePriority  = 200
		accountK8SOverridePriority = 300

		emptyLabelSet      = map[string]string{}
		accountLabelSet    = map[string]string{"account_id": "40"}
		accountBMLabelSet  = map[string]string{"account_id": "40", "deployment_type": "BAREMETAL"}
		accountK8SLabelSet = map[string]string{"account_id": "40", "k8s": "true"}
	)

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		err = CreateOrUpdateFeatureFlag(k8sClient, featureID, namespace, featureID, featureFlagValue)
		Expect(err).ToNot(HaveOccurred())

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	})

	JustBeforeEach(func() {
		time.Sleep(250 * time.Millisecond)
	})

	AfterEach(func() {
		err = DeleteFeatureFlag(k8sClient, featureID, namespace)
		Expect(err).ToNot(HaveOccurred())
		cancel()
	})

	Context("With only FeatureID defined", func() {

		It("should be readable via featureName and no labels", func() {
			res, err := featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      emptyLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(featureFlagValue))
		})

		It("should be readable via featureName and labels defined", func() {
			res, err := featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      accountBMLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(featureFlagValue))
		})
	})

	Context("With FeatureID defined and an override", func() {

		BeforeEach(func() {
			err = CreateOrUpdateFeatureFlagOverride(k8sClient, accountOverride, namespace, featureID, accountOverrideValue, accountOverridePriority, accountLabelSet)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err = DeleteFeatureFlagOverride(k8sClient, accountOverride, namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		It("it should match override with exact label", func() {
			res, err := featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      accountLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(accountOverrideValue))
		})

		It("it should match override with extra labels", func() {
			res, err := featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      accountBMLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(accountOverrideValue))

			res, err = featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      accountK8SLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(accountOverrideValue))
		})
	})

	Context("With FeatureID defined and multiple overrides (account and account+deployment_type)", func() {

		BeforeEach(func() {
			err = CreateOrUpdateFeatureFlagOverride(k8sClient, accountOverride, namespace, featureID, accountOverrideValue, accountOverridePriority, accountLabelSet)
			Expect(err).ToNot(HaveOccurred())
			err = CreateOrUpdateFeatureFlagOverride(k8sClient, accountBMOverride, namespace, featureID, accountBMOverrideValue, accountBMOverridePriority, accountBMLabelSet)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			err = DeleteFeatureFlagOverride(k8sClient, accountOverride, namespace)
			Expect(err).ToNot(HaveOccurred())
			err = DeleteFeatureFlagOverride(k8sClient, accountBMOverride, namespace)
			Expect(err).ToNot(HaveOccurred())
		})

		It("it should match correct override with multiple labels", func() {
			res, err := featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      accountBMLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(accountBMOverrideValue))

			res, err = featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      accountK8SLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(accountOverrideValue))
		})

		It("it should match correct override with multiple labels", func() {
			err = CreateOrUpdateFeatureFlagOverride(k8sClient, accountK8SOverride, namespace, featureID, accountK8SOverrideValue, accountK8SOverridePriority, accountK8SLabelSet)
			Expect(err).ToNot(HaveOccurred())

			res, err := featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      accountBMLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(accountBMOverrideValue))

			res, err = featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureID,
				Labels:      accountK8SLabelSet,
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureID))
			Expect(res.Result.Value).Should(Equal(accountK8SOverrideValue))
		})

	})
})
