package handler_test

import (
	"net/http"

	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/testing/admission_test"
	"github.com/dbsystel/kewl/testing/validation_test"

	"github.com/dbsystel/kewl/pkg/validation"
	"github.com/dbsystel/kewl/pkg/webhook/handler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ validation.Validator = &brokenValidator{}

type brokenValidator struct{}

func (f *brokenValidator) AddToScheme(_ *runtime.Scheme) error {
	return errors.New("meh")
}
func (f *brokenValidator) Name() string {
	return "brokenMutator"
}
func (f *brokenValidator) Validate(_, _ runtime.Object, _ validation.ResultCollector) error {
	return errors.New("meh")
}

var _ = Describe("Validator test", func() {
	var sut handler.Validator
	BeforeEach(func() {
		sut = handler.NewValidator()
		panicutils.PanicIfError(sut.AddValidator(validation_test.PodValidator))
	})
	It("should skip nil validators", func() {
		Expect(sut.AddValidator(nil)).To(Not(HaveOccurred()))
		review := admission_test.V1ValidPod()
		Expect(InvokeHandler(sut, review)).To(Not(HaveOccurred()))
	})
	It("should propagate error from registration", func() {
		Expect(sut.AddValidator(&brokenValidator{})).To(HaveOccurred())
	})
	It("should propagate validator errors", func() {
		review := admission_test.V1ErrorPod()
		Expect(InvokeHandler(sut, review)).To(HaveOccurred())
	})
	It("should not allow invalid objects", func() {
		review := admission_test.V1InvalidPod()
		Expect(InvokeHandler(sut, review)).To(Not(HaveOccurred()))
		Expect(review.Response.Allowed).To(BeFalse())
		Expect(review.Response.Result.Code).To(Equal(int32(http.StatusUnprocessableEntity)))
		Expect(review.Response.Result.Details.Causes).To(BeEquivalentTo([]v1.StatusCause{
			{Type: "invalid", Message: "not valid", Field: ".name"},
		}))
	})
	It("should set allowed to true on valid object", func() {
		review := admission_test.V1ValidPod()
		Expect(InvokeHandler(sut, review)).To(Not(HaveOccurred()))
		Expect(review.Response.Allowed).To(BeTrue())
	})
})
