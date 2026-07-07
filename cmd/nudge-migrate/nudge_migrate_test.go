package main

import (
	"bytes"
	"context"

	"github.com/go-logr/logr"
	applicationapiv1alpha1 "github.com/konflux-ci/application-api/api/v1alpha1"
	integrationv1beta2 "github.com/konflux-ci/integration-service/api/v1beta2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tonglil/buflogr"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func makeComponent(namespace, name string, nudgesRef []string) *applicationapiv1alpha1.Component {
	return &applicationapiv1alpha1.Component{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: applicationapiv1alpha1.ComponentSpec{
			ComponentName: name,
			Application:   "test-app",
			Source: applicationapiv1alpha1.ComponentSource{
				ComponentSourceUnion: applicationapiv1alpha1.ComponentSourceUnion{
					GitSource: &applicationapiv1alpha1.GitSource{
						URL: "https://github.com/example/" + name,
					},
				},
			},
			BuildNudgesRef: nudgesRef,
		},
	}
}

func makeNudgeConfig(namespace string, nudges []integrationv1beta2.NudgeRelationship) *integrationv1beta2.NudgeConfig {
	return &integrationv1beta2.NudgeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      integrationv1beta2.NudgeConfigSingletonName,
			Namespace: namespace,
		},
		Spec: integrationv1beta2.NudgeConfigSpec{
			Nudges: nudges,
		},
	}
}

func makeNamespace(name string, labels map[string]string) *core.Namespace {
	return &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
}

func getNudgeConfig(ctx context.Context, cl client.Client, namespace string) (*integrationv1beta2.NudgeConfig, error) {
	nc := &integrationv1beta2.NudgeConfig{}
	err := cl.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      integrationv1beta2.NudgeConfigSingletonName,
	}, nc)
	return nc, err
}

var _ = Describe("NudgeMigrate", func() {
	var (
		buf    bytes.Buffer
		logger logr.Logger
		ctx    context.Context
	)

	BeforeEach(func() {
		buf.Reset()
		logger = buflogr.NewWithBuffer(&buf)
		ctx = context.Background()
	})

	Describe("mergeNudges", func() {
		It("appends new entries not present in existing", func() {
			existing := []integrationv1beta2.NudgeRelationship{
				{From: "a", To: "b", Mode: integrationv1beta2.NudgeModeImmediate},
			}
			incoming := []integrationv1beta2.NudgeRelationship{
				{From: "a", To: "c", Mode: integrationv1beta2.NudgeModeImmediate},
			}
			merged, added := mergeNudges(existing, incoming)
			Expect(merged).To(HaveLen(2))
			Expect(added).To(Equal(1))
		})

		It("does not duplicate entries already present", func() {
			existing := []integrationv1beta2.NudgeRelationship{
				{From: "a", To: "b", Mode: integrationv1beta2.NudgeModeImmediate},
			}
			incoming := []integrationv1beta2.NudgeRelationship{
				{From: "a", To: "b", Mode: integrationv1beta2.NudgeModeImmediate},
			}
			merged, added := mergeNudges(existing, incoming)
			Expect(merged).To(HaveLen(1))
			Expect(added).To(Equal(0))
		})

		It("handles empty existing list", func() {
			incoming := []integrationv1beta2.NudgeRelationship{
				{From: "a", To: "b", Mode: integrationv1beta2.NudgeModeImmediate},
			}
			merged, added := mergeNudges(nil, incoming)
			Expect(merged).To(HaveLen(1))
			Expect(added).To(Equal(1))
		})

		It("handles empty incoming list", func() {
			existing := []integrationv1beta2.NudgeRelationship{
				{From: "a", To: "b", Mode: integrationv1beta2.NudgeModeImmediate},
			}
			merged, added := mergeNudges(existing, nil)
			Expect(merged).To(HaveLen(1))
			Expect(added).To(Equal(0))
		})
	})

	Describe("buildNudgeConfig", func() {
		It("creates NudgeConfig with correct name, labels, and annotations", func() {
			nudges := []integrationv1beta2.NudgeRelationship{
				{From: "a", To: "b", Mode: integrationv1beta2.NudgeModeImmediate},
			}
			nc := buildNudgeConfig("test-ns", nudges)
			Expect(nc.Name).To(Equal(integrationv1beta2.NudgeConfigSingletonName))
			Expect(nc.Namespace).To(Equal("test-ns"))
			Expect(nc.Labels).To(HaveKeyWithValue(MigrationLabel, MigrationLabelValue))
			Expect(nc.Annotations).To(HaveKeyWithValue(MigrationAnnotation, MigrationAnnotationValue))
			Expect(nc.Spec.Nudges).To(HaveLen(1))
		})
	})

	Describe("migrateNamespace", func() {
		It("skips namespace with no components", func() {
			cl := fake.NewClientBuilder().WithScheme(scheme).Build()
			s := migrateNamespace(ctx, cl, logger, "empty-ns", false)
			Expect(s.Action).To(Equal("skipped"))
			Expect(s.ComponentsRead).To(Equal(0))
			Expect(s.RelationshipsFound).To(Equal(0))
		})

		It("skips namespace with components but no build-nudges-ref", func() {
			compA := makeComponent("test-ns", "comp-a", nil)
			compB := makeComponent("test-ns", "comp-b", nil)
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB).Build()
			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("skipped"))
			Expect(s.ComponentsRead).To(Equal(2))
			Expect(s.RelationshipsFound).To(Equal(0))
		})

		It("creates NudgeConfig with correct entries for simple migration", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-b", "comp-c"})
			compB := makeComponent("test-ns", "comp-b", nil)
			compC := makeComponent("test-ns", "comp-c", nil)
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB, compC).Build()

			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("created"))
			Expect(s.RelationshipsFound).To(Equal(2))
			Expect(s.RelationshipsMigrated).To(Equal(2))

			nc, err := getNudgeConfig(ctx, cl, "test-ns")
			Expect(err).NotTo(HaveOccurred())
			Expect(nc.Spec.Nudges).To(HaveLen(2))
			Expect(nc.Labels).To(HaveKeyWithValue(MigrationLabel, MigrationLabelValue))
			Expect(nc.Annotations).To(HaveKeyWithValue(MigrationAnnotation, MigrationAnnotationValue))
		})

		It("filters out dangling references to non-existent components", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-b", "comp-missing"})
			compB := makeComponent("test-ns", "comp-b", nil)
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB).Build()

			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("created"))
			Expect(s.RelationshipsFound).To(Equal(1))
			Expect(s.RelationshipsSkippedDangling).To(Equal(1))

			nc, err := getNudgeConfig(ctx, cl, "test-ns")
			Expect(err).NotTo(HaveOccurred())
			Expect(nc.Spec.Nudges).To(HaveLen(1))
			Expect(nc.Spec.Nudges[0].From).To(Equal("comp-a"))
			Expect(nc.Spec.Nudges[0].To).To(Equal("comp-b"))
		})

		It("filters out self-nudges", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-a"})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA).Build()

			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("skipped"))
			Expect(s.RelationshipsFound).To(Equal(0))
		})

		It("deduplicates identical (from, to) pairs from different components", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-b", "comp-b"})
			compB := makeComponent("test-ns", "comp-b", nil)
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB).Build()

			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("created"))
			Expect(s.RelationshipsFound).To(Equal(1))

			nc, err := getNudgeConfig(ctx, cl, "test-ns")
			Expect(err).NotTo(HaveOccurred())
			Expect(nc.Spec.Nudges).To(HaveLen(1))
		})

		It("detects cycles and skips namespace", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-b"})
			compB := makeComponent("test-ns", "comp-b", []string{"comp-a"})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB).Build()

			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("error"))
			Expect(s.Errors).To(HaveLen(1))
			Expect(s.Errors[0]).To(ContainSubstring("cycle"))
		})

		It("merges with existing NudgeConfig, adding only new entries", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-b", "comp-c"})
			compB := makeComponent("test-ns", "comp-b", nil)
			compC := makeComponent("test-ns", "comp-c", nil)
			existingNC := makeNudgeConfig("test-ns", []integrationv1beta2.NudgeRelationship{
				{From: "comp-a", To: "comp-b", Mode: integrationv1beta2.NudgeModeImmediate},
			})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB, compC, existingNC).Build()

			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("updated"))
			Expect(s.RelationshipsMigrated).To(Equal(1))
			Expect(s.RelationshipsAlreadyPresent).To(Equal(1))

			nc, err := getNudgeConfig(ctx, cl, "test-ns")
			Expect(err).NotTo(HaveOccurred())
			Expect(nc.Spec.Nudges).To(HaveLen(2))
			Expect(nc.Labels).To(HaveKeyWithValue(MigrationLabel, MigrationLabelValue))
		})

		It("skips update when all relationships are already present", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-b"})
			compB := makeComponent("test-ns", "comp-b", nil)
			existingNC := makeNudgeConfig("test-ns", []integrationv1beta2.NudgeRelationship{
				{From: "comp-a", To: "comp-b", Mode: integrationv1beta2.NudgeModeImmediate},
			})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB, existingNC).Build()

			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("skipped"))
			Expect(s.RelationshipsAlreadyPresent).To(Equal(1))
			Expect(s.RelationshipsMigrated).To(Equal(0))
		})

		It("preserves user-added entries when merging", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-b"})
			compB := makeComponent("test-ns", "comp-b", nil)
			compC := makeComponent("test-ns", "comp-c", nil)
			existingNC := makeNudgeConfig("test-ns", []integrationv1beta2.NudgeRelationship{
				{From: "comp-c", To: "comp-b", Mode: integrationv1beta2.NudgeModeValidated},
			})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB, compC, existingNC).Build()

			s := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s.Action).To(Equal("updated"))

			nc, err := getNudgeConfig(ctx, cl, "test-ns")
			Expect(err).NotTo(HaveOccurred())
			Expect(nc.Spec.Nudges).To(HaveLen(2))

			modes := map[string]integrationv1beta2.NudgeModeType{}
			for _, n := range nc.Spec.Nudges {
				modes[n.From+"->"+n.To] = n.Mode
			}
			Expect(modes).To(HaveKeyWithValue("comp-c->comp-b", integrationv1beta2.NudgeModeValidated))
			Expect(modes).To(HaveKeyWithValue("comp-a->comp-b", integrationv1beta2.NudgeModeImmediate))
		})

		Context("dry-run mode", func() {
			It("does not create NudgeConfig in dry-run mode", func() {
				compA := makeComponent("test-ns", "comp-a", []string{"comp-b"})
				compB := makeComponent("test-ns", "comp-b", nil)
				cl := fake.NewClientBuilder().WithScheme(scheme).
					WithObjects(compA, compB).Build()

				s := migrateNamespace(ctx, cl, logger, "test-ns", true)
				Expect(s.Action).To(Equal("dry-run:create"))
				Expect(s.RelationshipsMigrated).To(Equal(1))

				_, err := getNudgeConfig(ctx, cl, "test-ns")
				Expect(err).To(HaveOccurred())
			})

			It("does not update NudgeConfig in dry-run mode", func() {
				compA := makeComponent("test-ns", "comp-a", []string{"comp-b", "comp-c"})
				compB := makeComponent("test-ns", "comp-b", nil)
				compC := makeComponent("test-ns", "comp-c", nil)
				existingNC := makeNudgeConfig("test-ns", []integrationv1beta2.NudgeRelationship{
					{From: "comp-a", To: "comp-b", Mode: integrationv1beta2.NudgeModeImmediate},
				})
				cl := fake.NewClientBuilder().WithScheme(scheme).
					WithObjects(compA, compB, compC, existingNC).Build()

				s := migrateNamespace(ctx, cl, logger, "test-ns", true)
				Expect(s.Action).To(Equal("dry-run:update"))
				Expect(s.RelationshipsMigrated).To(Equal(1))

				nc, err := getNudgeConfig(ctx, cl, "test-ns")
				Expect(err).NotTo(HaveOccurred())
				Expect(nc.Spec.Nudges).To(HaveLen(1))
			})
		})

		It("is idempotent — running twice produces the same result", func() {
			compA := makeComponent("test-ns", "comp-a", []string{"comp-b"})
			compB := makeComponent("test-ns", "comp-b", nil)
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB).Build()

			s1 := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s1.Action).To(Equal("created"))

			s2 := migrateNamespace(ctx, cl, logger, "test-ns", false)
			Expect(s2.Action).To(Equal("skipped"))
			Expect(s2.RelationshipsAlreadyPresent).To(Equal(1))
			Expect(s2.RelationshipsMigrated).To(Equal(0))

			nc, err := getNudgeConfig(ctx, cl, "test-ns")
			Expect(err).NotTo(HaveOccurred())
			Expect(nc.Spec.Nudges).To(HaveLen(1))
		})
	})

	Describe("migrateNamespaces", func() {
		It("processes multiple namespaces and returns correct summaries", func() {
			compA := makeComponent("ns1", "comp-a", []string{"comp-b"})
			compB := makeComponent("ns1", "comp-b", nil)
			compC := makeComponent("ns2", "comp-c", nil)
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(compA, compB, compC).Build()

			summaries := migrateNamespaces(ctx, cl, logger, []string{"ns1", "ns2"}, false)
			Expect(summaries).To(HaveLen(2))
			Expect(summaries[0].Action).To(Equal("created"))
			Expect(summaries[1].Action).To(Equal("skipped"))
		})
	})

	Describe("getTenantNamespaces", func() {
		It("discovers namespaces with toolchain tenant label", func() {
			ns := makeNamespace("tenant-1", map[string]string{
				"toolchain.dev.openshift.com/type": "tenant",
			})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(ns).Build()

			result, err := getTenantNamespaces(ctx, cl, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainElement("tenant-1"))
		})

		It("discovers namespaces with konflux user label", func() {
			ns := makeNamespace("user-1", map[string]string{
				"konflux.ci/type": "user",
			})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(ns).Build()

			result, err := getTenantNamespaces(ctx, cl, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainElement("user-1"))
		})

		It("discovers namespaces with konflux-ci tenant label", func() {
			ns := makeNamespace("tenant-2", map[string]string{
				"konflux-ci.dev/type": "tenant",
			})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(ns).Build()

			result, err := getTenantNamespaces(ctx, cl, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainElement("tenant-2"))
		})

		It("deduplicates namespaces matching multiple label selectors", func() {
			ns := makeNamespace("multi-label", map[string]string{
				"toolchain.dev.openshift.com/type": "tenant",
				"konflux-ci.dev/type":              "tenant",
			})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(ns).Build()

			result, err := getTenantNamespaces(ctx, cl, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(1))
			Expect(result).To(ContainElement("multi-label"))
		})

		It("returns sorted namespace names", func() {
			ns1 := makeNamespace("z-ns", map[string]string{"konflux-ci.dev/type": "tenant"})
			ns2 := makeNamespace("a-ns", map[string]string{"konflux-ci.dev/type": "tenant"})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(ns1, ns2).Build()

			result, err := getTenantNamespaces(ctx, cl, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal([]string{"a-ns", "z-ns"}))
		})

		It("returns empty list when no tenant namespaces exist", func() {
			ns := makeNamespace("regular-ns", map[string]string{
				"some-label": "some-value",
			})
			cl := fake.NewClientBuilder().WithScheme(scheme).
				WithObjects(ns).Build()

			result, err := getTenantNamespaces(ctx, cl, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})
	})
})
