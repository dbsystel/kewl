package examples_test

import (
	"github.com/dbsystel/kewl/examples"
	"github.com/dbsystel/kewl/pkg/mutation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PodMutator", func() {
	sut := examples.PodMutator
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
