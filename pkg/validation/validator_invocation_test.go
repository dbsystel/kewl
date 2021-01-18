package validation_test

import (
	"github.com/dbsystel/kewl/pkg/validation"
	"github.com/dbsystel/kewl/testing/corev1_test"
	"github.com/dbsystel/kewl/testing/uncurry"
	"github.com/dbsystel/kewl/testing/validation_test"

	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator Invocation", func() {
	var scheme *runtime.Scheme
	BeforeEach(func() {
		scheme = runtime.NewScheme()
	})
	It("should add the descriptors to the scheme", func() {
		Expect(validation_test.PodValidator.AddToScheme(scheme)).To(Not(HaveOccurred()))
		Expect(scheme.AllKnownTypes()).To(Not(BeEmpty()))
	})
	It("should handle panics", func() {
		Expect(uncurry.Error2(validation.Validate(corev1_test.PanicPod.AsCoreV1(), nil, validation_test.PodValidator))).
			To(HaveOccurred())
	})
	It("should pass errors", func() {
		Expect(uncurry.Error2(validation.Validate(corev1_test.ErrorPod.AsCoreV1(), nil, validation_test.PodValidator))).
			To(HaveOccurred())
	})
	It("should return the validation failures with fields", func() {
		Expect(validation.Validate(corev1_test.InvalidPod.AsCoreV1(), corev1_test.ValidPod.AsCoreV1(), validation_test.PodValidator)).
			To(Equal(validation_test.PodValidator.InvalidStatusCause()))
	})
	It("should return no validation failures on valid object", func() {
		Expect(validation.Validate(corev1_test.ValidPod.AsCoreV1(), corev1_test.InvalidPod.AsCoreV1(), validation_test.PodValidator)).
			To(BeEmpty())
	})
	It("should ignore nil validators", func() {
		Expect(validation.Validate(corev1_test.ValidPod.AsCoreV1(), nil, validation_test.PodValidator, nil)).
			To(BeEmpty())
	})
})
