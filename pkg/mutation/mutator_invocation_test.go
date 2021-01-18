package mutation_test

import (
	"github.com/dbsystel/kewl/pkg/mutation"
	"github.com/dbsystel/kewl/testing/corev1_test"
	"github.com/dbsystel/kewl/testing/mutation_test"
	"github.com/dbsystel/kewl/testing/uncurry"
	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mutator Invocation", func() {
	var scheme *runtime.Scheme
	BeforeEach(func() {
		scheme = runtime.NewScheme()
	})
	It("should add the descriptors to the scheme", func() {
		Expect(mutation_test.PodMutator.AddToScheme(scheme)).To(Not(HaveOccurred()))
		Expect(scheme.AllKnownTypes()).To(Not(BeEmpty()))
	})
	It("should handle panics", func() {
		Expect(uncurry.Error2(mutation.Mutate(corev1_test.PanicPod.AsCoreV1(), nil, mutation_test.PodMutator))).
			To(HaveOccurred())
	})
	It("should pass errors", func() {
		Expect(uncurry.Error2(mutation.Mutate(corev1_test.ErrorPod.AsCoreV1(), nil, mutation_test.PodMutator))).
			To(HaveOccurred())
	})
	It("should return mutated object on change", func() {
		Expect(mutation.Mutate(corev1_test.InvalidPod.AsCoreV1(), corev1_test.ValidPod.AsCoreV1(), mutation_test.PodMutator)).
			To(BeEquivalentTo(corev1_test.ValidPod))
	})

	It("should return nil on unchanged", func() {
		Expect(mutation.Mutate(corev1_test.ValidPod.AsCoreV1(), corev1_test.InvalidPod.AsCoreV1(), mutation_test.PodMutator)).
			To(BeNil())
	})
	It("should ignore nil mutators", func() {
		Expect(mutation.Mutate(corev1_test.ValidPod.AsCoreV1(), nil, mutation_test.PodMutator, nil)).
			To(BeNil())
	})
})
