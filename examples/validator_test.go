package examples_test

import (
	"github.com/dbsystel/kewl/examples"
	"github.com/dbsystel/kewl/pkg/validation"
	"github.com/dbsystel/kewl/pkg/webhook/integtest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PodValidator", func() {
	sut := examples.PodValidator
	Context("unit test", func() {
		It("should add an error for name cat", func() {
			Expect(validation.Validate(NewPod("cat"), nil, sut)).
				To(HaveLen(1))
		})
		It("should not find errors for name dog", func() {
			Expect(validation.Validate(NewPod("dog"), nil, sut)).To(BeEmpty())
		})
		It("should add an error for label value cat", func() {
			Expect(validation.Validate(NewPod("rabbit", "dog", "cat"), nil, sut)).
				To(HaveLen(1))
		})
		It("should add an error for label value dog", func() {
			Expect(validation.Validate(NewPod("rabbit", "dog", "dog"), nil, sut)).
				To(BeEmpty())
		})
	})
	Context("integration test", func() {
		it, err := integtest.NewValidation("test-namespace", sut)
		Expect(err).NotTo(HaveOccurred())
		It("should mutate the pod", func() {
			response, err := it.InvokeFromFile("invalid-pod.yaml", "")
			Expect(response, err).NotTo(BeNil())
			Expect(response.Allowed).To(BeFalse())
		})
	})
})
