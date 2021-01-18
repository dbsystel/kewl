package handler_test

import (
	"github.com/dbsystel/kewl/pkg/mutation"
	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/pkg/webhook/handler"
	"github.com/dbsystel/kewl/testing/admission_test"
	"github.com/dbsystel/kewl/testing/mutation_test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ mutation.Mutator = &brokenMutator{}

type brokenMutator struct{}

func (f brokenMutator) Mutate(_, _ runtime.Object) error {
	return errors.New("blubb")
}
func (f brokenMutator) AddToScheme(*runtime.Scheme) error {
	return errors.New("meh")
}
func (f brokenMutator) Name() string {
	return "fail"
}

var _ = Describe("Mutator test", func() {
	var sut handler.Mutator
	BeforeEach(func() {
		sut = handler.NewMutator()
		panicutils.PanicIfError(sut.AddMutator(mutation_test.PodMutator))
	})
	It("should skip nil mutators", func() {
		Expect(sut.AddMutator(nil)).To(Not(HaveOccurred()))
		Expect(InvokeHandler(sut, admission_test.V1ValidPod())).To(Not(HaveOccurred()))
	})
	It("should propagate error from registration", func() {
		Expect(sut.AddMutator(&brokenMutator{})).To(HaveOccurred())
	})
	It("should provide no failures, if object unknown", func() {
		Expect(InvokeHandler(sut, admission_test.V1BadPod())).NotTo(HaveOccurred())
	})
	It("should not set a patch, if unchanged", func() {
		review := admission_test.V1ValidPod()
		Expect(InvokeHandler(sut, review)).To(Not(HaveOccurred()))
		Expect(review.Response.Allowed).To(BeTrue())
		Expect(review.Response.PatchType).To(BeNil())
		Expect(review.Response.Patch).To(BeNil())
	})
	It("should  set a correct patch, if mutated", func() {
		review := admission_test.V1InvalidPod()
		patchData := []byte("[{\"op\":\"replace\",\"path\":\"/metadata/name\",\"value\":\"valid\"}]")
		Expect(InvokeHandler(sut, review)).To(Not(HaveOccurred()))
		Expect(review.Response.Allowed).To(BeTrue())
		Expect(review.Response.PatchType).To(Not(BeNil()))
		Expect(review.Response.Patch).To(Equal(patchData))
	})
	It("should propagate mutator errors", func() {
		review := admission_test.V1ErrorPod()
		Expect(InvokeHandler(sut, review)).To(HaveOccurred())
	})
})
