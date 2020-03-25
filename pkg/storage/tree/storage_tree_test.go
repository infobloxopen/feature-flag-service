package tree

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	ffv1 "github.com/Infoblox-CTO/atlas.feature.flag/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// "github.com/Infoblox-CTO/atlas.feature.flag/pkg/pb"
	"github.com/Infoblox-CTO/atlas.feature.flag/pkg/storage"
)

var _ = Describe("Testing storage tree", func() {
	var err error
	_ = err
	var tree storage.Storage

	BeforeEach(func() {
		tree = NewInMemoryStorage()
	})

	AfterEach(func() {

	})

	Context("With tree", func() {
		It("it should allow FeatureFlag to be added", func() {
			ff := &ffv1.FeatureFlag{
				Spec: ffv1.FeatureFlagSpec{
					FeatureID: "test",
					Value:     "ffValue",
				},
			}
			tree.Define(ff)
			tree.Find("test", nil)
		})

		It("it should allow FeatureFlagOverride to be added if FeatureFlag exists", func() {
			ff := &ffv1.FeatureFlag{
				Spec: ffv1.FeatureFlagSpec{
					FeatureID: "test",
					Value:     "ffValue",
				},
			}
			tree.Define(ff)

			ffo := &ffv1.FeatureFlagOverride{
				Spec: ffv1.FeatureFlagOverrideSpec{
					OverrideName: "ffo1",
					FeatureID:    "test",
					Value:        "ffoValue",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"account_id": "40",
							"k8s":        "true",
						},
					},
				},
			}
			tree.Override(ffo)
		})

		It("it should match override with exact label", func() {
			ff := &ffv1.FeatureFlag{
				Spec: ffv1.FeatureFlagSpec{
					FeatureID: "test",
					Value:     "ffValue",
				},
			}
			tree.Define(ff)

			ffo := &ffv1.FeatureFlagOverride{
				Spec: ffv1.FeatureFlagOverrideSpec{
					OverrideName: "ffo1",
					FeatureID:    "test",
					Value:        "ffoValue",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"account_id": "40",
						},
					},
				},
			}
			tree.Override(ffo)

			res := tree.Find("test", map[string]string{
				"account_id": "40",
			})
			Expect(res).ToNot(BeNil())
			Expect(res.Value).Should(Equal("ffoValue"))

		})

		It("it should match override with multiple exact labels", func() {
			ff := &ffv1.FeatureFlag{
				Spec: ffv1.FeatureFlagSpec{
					FeatureID: "test",
					Value:     "ffValue",
				},
			}
			tree.Define(ff)

			ffo := &ffv1.FeatureFlagOverride{
				Spec: ffv1.FeatureFlagOverrideSpec{
					OverrideName: "ffo1",
					FeatureID:    "test",
					Value:        "ffoValue",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"account_id": "40",
							"k8s":        "true",
						},
					},
				},
			}
			tree.Override(ffo)

			res := tree.Find("test", map[string]string{
				"account_id": "40",
				"k8s":        "true",
			})
			Expect(res).ToNot(BeNil())
			Expect(res.Value).Should(Equal("ffoValue"))

		})

		It("it should match override with extra labels", func() {
			ff := &ffv1.FeatureFlag{
				Spec: ffv1.FeatureFlagSpec{
					FeatureID: "test",
					Value:     "ffValue",
				},
			}
			tree.Define(ff)

			ffo := &ffv1.FeatureFlagOverride{
				Spec: ffv1.FeatureFlagOverrideSpec{
					OverrideName: "ffo1",
					FeatureID:    "test",
					Value:        "ffoValue",
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"k8s": "true",
						},
					},
				},
			}
			tree.Override(ffo)

			res := tree.Find("test", map[string]string{
				"k8s": "true",
			})
			Expect(res).ToNot(BeNil())
			Expect(res.Value).Should(Equal("ffoValue"))

			res = tree.Find("test", map[string]string{
				"account_id": "40",
				"k8s":        "true",
			})
			Expect(res).ToNot(BeNil())
			Expect(res.Value).Should(Equal("ffoValue"))

			res = tree.Find("test", map[string]string{
				"account_id": "40",
				"k8s":        "true",
				"extra":      "extra",
			})
			Expect(res).ToNot(BeNil())
			Expect(res.Value).Should(Equal("ffoValue"))
		})
	})
})

var _ = Describe("Tests by Siddharth", func() {
	var err error
	_ = err
	var tree storage.Storage

	BeforeEach(func() {
		tree = NewInMemoryStorage()
		ff := &ffv1.FeatureFlag{
			Spec: ffv1.FeatureFlagSpec{
				FeatureID: "noa",
				Value:     "3.3.6",
			},
		}
		tree.Define(ff)

		ffo1 := &ffv1.FeatureFlagOverride{
			Spec: ffv1.FeatureFlagOverrideSpec{
				OverrideName: "ffo1",
				FeatureID:    "noa",
				Value:        "4.0.0",
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"account_id": "40",
					},
				},
				Priority: 100,
			},
		}
		tree.Override(ffo1)

		ffo2 := &ffv1.FeatureFlagOverride{
			Spec: ffv1.FeatureFlagOverrideSpec{
				OverrideName: "ffo2",
				FeatureID:    "noa",
				Value:        "4.0.1",
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"account_id": "40",
						"k8s":        "true",
					},
				},
				Priority: 200,
			},
		}
		tree.Override(ffo2)
	})

	AfterEach(func() {

	})

	Context("With tree containing two overrides", func() {
		It("it should match many labels as possible (sending 2, 1 exist)", func() {
			res := tree.Find("noa", map[string]string{
				"account_id":      "40",
				"deployment_type": "BAREMETAL",
			})
			Expect(res).ToNot(BeNil())
			Expect(res.Value).Should(Equal("4.0.0"))
		})

		It("it should match with as many labels as possible (sending 3, 2 exist)", func() {
			res := tree.Find("noa", map[string]string{
				"account_id":      "40",
				"deployment_type": "BAREMETAL",
				"k8s":             "true",
			})
			Expect(res).ToNot(BeNil())
			Expect(res.Value).Should(Equal("4.0.1"))
		})

		It("it should match with as many labels as possible (sending 2, 2 exist)", func() {
			res := tree.Find("noa", map[string]string{
				"account_id": "40",
				"k8s":        "true",
			})
			Expect(res).ToNot(BeNil())
			Expect(res.Value).Should(Equal("4.0.1"))
		})
	})
})
