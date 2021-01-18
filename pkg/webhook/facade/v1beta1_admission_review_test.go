//nolint:dupl
package facade_test

import (
	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/pkg/webhook/facade"
	"github.com/dbsystel/kewl/testing/admission_test"
	"github.com/dbsystel/kewl/testing/validation_test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("v1beta1AdmissionReview test", func() {
	var sut facade.AdmissionReview
	BeforeEach(func() {
		review, err := facade.AdmissionReviewFrom(admission_test.V1Beta1ValidPod().MustMarshal())
		panicutils.PanicIfError(err)
		sut = review
	})
	v1beta1AdmissionReview := func() *v1beta1.AdmissionReview {
		return K8sAdmissionReview(sut, &v1beta1.AdmissionReview{}).(*v1beta1.AdmissionReview)
	}
	It("should facade the request correctly", func() {
		expected := admission_test.V1Beta1ValidPod().Request
		Expect(sut.Request().Kind()).To(BeEquivalentTo(expected.Kind))
		Expect(sut.Request().Object()).To(BeEquivalentTo(&expected.Object))
		Expect(sut.Request().OldObject()).To(BeEquivalentTo(&expected.OldObject))
		Expect(sut.Request().Resource()).To(BeEquivalentTo(expected.Resource))
		Expect(sut.Request().Namespace()).To(BeEquivalentTo(expected.Namespace))
		Expect(sut.Response().ResponseType()).To(Equal(facade.AdmissionError))
	})
	Context("facade response", func() {
		It("should apply allow correctly", func() {
			sut.Response().Allow()
			result := v1beta1AdmissionReview()
			Expect(result.Response).NotTo(BeNil())
			Expect(result.Response.Allowed).To(BeTrue())
			Expect(sut.Response().ResponseType()).To(Equal(facade.AdmissionAllowed))
		})
		It("should apply deny correctly", func() {
			status := &metav1.Status{
				Status: "xy", Message: "z", Reason: "nope",
				Details: &metav1.StatusDetails{
					Name: "foo", Group: "bar", Kind: "baz", UID: "123",
					Causes:            validation_test.PodValidator.InvalidStatusCause(),
					RetryAfterSeconds: 123,
				},
				Code: 666,
			}
			sut.Response().Deny(status)
			result := v1beta1AdmissionReview()
			Expect(result.Response).NotTo(BeNil())
			Expect(result.Response.Allowed).To(BeFalse())
			Expect(result.Response.Result).To(BeEquivalentTo(status))
			Expect(sut.Response().ResponseType()).To(Equal(facade.AdmissionDenied))
		})
		It("should apply patch correctly", func() {
			sut.Response().PatchJSON([]byte("{}"))
			result := v1beta1AdmissionReview()
			Expect(result.Response).NotTo(BeNil())
			Expect(result.Response.Allowed).To(BeTrue())
			Expect(result.Response.PatchType).NotTo(BeNil())
			Expect(result.Response.Patch).NotTo(BeEmpty())
			Expect(sut.Response().ResponseType()).To(Equal(facade.AdmissionMutated))
		})
	})
})
