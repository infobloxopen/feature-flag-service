/*
Copyright 2020 Infoblox.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/Infoblox-CTO/atlas-app-definition-controller/tests/helpers"
	featureflagv1 "github.com/Infoblox-CTO/atlas.feature.flag/api/v1"
	. "github.com/Infoblox-CTO/atlas.feature.flag/api/v1/testing"
	featureflagpb "github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/golang/protobuf/ptypes/empty"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Testing feature-flag service", func() {
	var err error
	featureFlagName := "ff-app-infra"
	featureFlagValue := "2.4.1"
	featureFlagOverrideValue := "2.5.0"

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	Context("With k8s client", func() {
		It("it should allow FeatureFlag to be created", func() {
			ctx, acancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer acancel()

			ff := &featureflagv1.FeatureFlag{
				ObjectMeta: metav1.ObjectMeta{
					Name:      featureFlagName,
					Namespace: namespace,
				},
				Spec: featureflagv1.FeatureFlagSpec{
					Value: featureFlagValue,
				},
			}
			key, _ := ctrlclient.ObjectKeyFromObject(ff)
			err := k8sClient.Get(ctx, key, ff)
			Expect(err).To(HaveOccurred())

			err = k8sClient.Create(ctx, ff)
			Expect(err).ToNot(HaveOccurred())

			err = k8sClient.Get(ctx, key, ff)
			Expect(err).ToNot(HaveOccurred())
		})

		It("it should allow FeatureFlagOverride to be created", func() {
			ctx, acancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer acancel()

			ffo := &featureflagv1.FeatureFlagOverride{
				ObjectMeta: metav1.ObjectMeta{
					Name:      featureFlagName,
					Namespace: namespace,
				},
				Spec: featureflagv1.FeatureFlagOverrideSpec{
					FeatureName: featureFlagName,
					Value:       featureFlagOverrideValue,
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"test": "test"},
					},
				},
			}
			err := k8sClient.Create(ctx, ffo)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("With gRPC client", func() {

		BeforeEach(func() {
			err = CreateOrUpdateFeatureFlag(k8sClient, featureFlagName, namespace, featureFlagName, featureFlagValue)
			Expect(err).ToNot(HaveOccurred())

			labels := map[string]string{
				"ophid": "test",
			}
			err = CreateOrUpdateFeatureFlagOverride(k8sClient, featureFlagName, namespace, featureFlagName, featureFlagOverrideValue, 100, labels)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(250 * time.Millisecond)
		})

		AfterEach(func() {
			err = DeleteFeatureFlag(k8sClient, featureFlagName, namespace)
			Expect(err).ToNot(HaveOccurred())

			err = DeleteFeatureFlagOverride(k8sClient, featureFlagName, namespace)
			Expect(err).ToNot(HaveOccurred())
			time.Sleep(250 * time.Millisecond)
		})

		It("it should respond to GetVersion", func() {
			res, err := featureFlagClient.GetVersion(context.Background(), &empty.Empty{})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Version).ShouldNot(Equal(""))
		})

		It("valid FeatureFlags should be listable", func() {
			ctx, acancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer acancel()

			res, err := featureFlagClient.List(ctx, &featureflagpb.ListFeatureFlagsRequest{
				Labels: map[string]string{},
			})
			Expect(err).ToNot(HaveOccurred())
			Info(res.Results)
			Expect(len(res.Results)).Should(BeNumerically(">=", 1))
		})

		It("a valid FeatureFlag should be Read", func() {
			ctx, acancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer acancel()

			res, err := featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: featureFlagName,
				Labels:      map[string]string{},
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.FeatureName).Should(Equal(featureFlagName))
			Expect(res.Result.Value).Should(Equal(featureFlagValue))
		})

		It("a non-existent FeatureFlag should fail to Read", func() {
			ctx, acancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer acancel()

			_, err := featureFlagClient.Read(ctx, &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: "i-do-not-exist",
				Labels:      map[string]string{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("rpc error: code = Unknown desc = FeatureFlag or FeatureFlagOverride not found for FeatureName"))
		})

		It("a FeatureFlag should be overridden by FeatureFlagOverride", func() {
			ctx, acancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer acancel()

			name := "overridetest"
			value := "value"
			overrideValue1 := "value1"
			overrideValue2 := "value2"
			overrideValue3 := "value3"
			testLabels := map[string]string{
				"my-label": "test-label",
			}
			req := &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: name,
				Labels:      testLabels,
			}

			reqNoLabels := &featureflagpb.ReadFeatureFlagRequest{
				FeatureName: name,
			}
			// set initial feature value
			err := CreateOrUpdateFeatureFlag(k8sClient, name, namespace, name, value)
			Expect(err).ToNot(HaveOccurred())

			// set override value
			err = CreateOrUpdateFeatureFlagOverride(k8sClient, overrideValue1, namespace, name, overrideValue1, 100, testLabels)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(250 * time.Microsecond)
			// confirm override is returned
			res, err := featureFlagClient.Read(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.Value).Should(Equal(overrideValue1))

			time.Sleep(250 * time.Microsecond)
			// confirm default is returned if no labels match
			res, err = featureFlagClient.Read(ctx, reqNoLabels)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.Value).Should(Equal(value))

			time.Sleep(250 * time.Microsecond)
			// add a higher priority override
			err = CreateOrUpdateFeatureFlagOverride(k8sClient, overrideValue3, namespace, name, overrideValue3, 300, testLabels)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(250 * time.Microsecond)
			// confirm higher priority override is returned
			res, err = featureFlagClient.Read(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.Value).Should(Equal(overrideValue3))

			time.Sleep(250 * time.Microsecond)
			// add a lower priority override
			err = CreateOrUpdateFeatureFlagOverride(k8sClient, overrideValue2, namespace, name, overrideValue2, 200, testLabels)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(2250 * time.Microsecond)
			// confirm the higher priority override is returned
			res, err = featureFlagClient.Read(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Info("result", res.Result)
			Expect(res.Result.Value).Should(Equal(overrideValue3))

			time.Sleep(250 * time.Microsecond)
			// remove each override from highest to lowest and verify
			err = DeleteFeatureFlagOverride(k8sClient, overrideValue3, namespace)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(250 * time.Microsecond)
			res, err = featureFlagClient.Read(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.Value).Should(Equal(overrideValue2))

			time.Sleep(250 * time.Microsecond)
			err = DeleteFeatureFlagOverride(k8sClient, overrideValue2, namespace)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(250 * time.Microsecond)
			res, err = featureFlagClient.Read(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.Value).Should(Equal(overrideValue1))

			time.Sleep(250 * time.Microsecond)
			err = DeleteFeatureFlagOverride(k8sClient, overrideValue1, namespace)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(250 * time.Microsecond)
			res, err = featureFlagClient.Read(ctx, req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Result.Value).Should(Equal(value))
		})
	})
})
