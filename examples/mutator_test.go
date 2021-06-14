package examples_test

import (
	"github.com/dbsystel/kewl/examples"
	"github.com/dbsystel/kewl/pkg/mutation"
	"github.com/dbsystel/kewl/pkg/webhook/integtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PodMutator", func() {
	sut := examples.PodMutator
	Context("unit test", func() {
		It("should change a pod with name cat", func() {
			Expect(mutation.Mutate(NewPod("cat"), nil, sut)).
				To(BeEquivalentTo(NewPod("dog")))
		})
		It("should not change a pod with name dog", func() {
			Expect(mutation.Mutate(NewPod("dog"), nil, sut)).To(BeNil())
		})
		It("should change a pod label value cat", func() {
			Expect(mutation.Mutate(NewPod("rabbit", "dog", "cat"), nil, sut)).
				To(BeEquivalentTo(NewPod("rabbit", "dog", "dog")))
		})
		It("should not change a pod labels not cat", func() {
			Expect(mutation.Mutate(NewPod("rabbit", "dog", "dog"), nil, sut)).
				To(BeNil())
		})
	})
	Context("integration test", func() {
		it, err := integtest.NewMutation("test-namespace", sut)
		Expect(err).NotTo(HaveOccurred())
		It("should mutate the pod", func() {
			response, err := it.InvokeFromFile("invalid-pod.yaml", "")
			Expect(response, err).NotTo(BeNil())
			Expect(response.PatchMaps()).To(ConsistOf(map[string]string{
				"op":    "replace",
				"path":  "/metadata/name",
				"value": "dog",
			}))
		})
	})
})
