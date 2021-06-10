package handler_test

import (
	"github.com/dbsystel/kewl/pkg/codec"
	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/pkg/webhook/handler"
	"github.com/dbsystel/kewl/testing/admission_test"
	"github.com/dbsystel/kewl/testing/corev1_test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ codec.SchemeExtension = &corev1Extension{}

type corev1Extension struct{}

func (c *corev1Extension) AddToScheme(scheme *runtime.Scheme) error {
	return corev1.AddToScheme(scheme)
}

var _ = Describe("UnmarshalReqObj test", func() {
	var sut handler.UnmarshalReqObj
	BeforeEach(func() {
		sut = handler.NewUnmarshalReqObj()
		panicutils.PanicIfError(sut.Register(&corev1Extension{}))
	})
	It("should not error, if no old object is present", func() {
		Expect(InvokeHandler(sut, admission_test.V1ValidPod())).To(Not(HaveOccurred()))
	})
	It("should not add a failure in case object is not known", func() {
		Expect(InvokeHandler(sut, admission_test.V1BadPod())).NotTo(HaveOccurred())
	})
	It("should deserialize both objects and set them on review", func() {
		review := admission_test.V1InvalidPod()
		Expect(InvokeHandler(sut, review)).To(Not(HaveOccurred()))
		Expect(review.Request.Object.Object).To(BeEquivalentTo(corev1_test.InvalidPod))
		Expect(review.Request.OldObject.Object).To(BeEquivalentTo(corev1_test.ValidPod))
	})
	It("should set namespace from review", func() {
		review := admission_test.V1DetachedPod()
		review.Request.Namespace = "blubb"
		Expect(InvokeHandler(sut, review)).To(Not(HaveOccurred()))
		Expect(review.Request.Object.Object).To(BeAssignableToTypeOf(&corev1.Pod{}))
		Expect(review.Request.Object.Object.(*corev1.Pod).Namespace).To(Equal(review.Request.Namespace))
		Expect(review.Request.OldObject.Object).To(BeAssignableToTypeOf(&corev1.Pod{}))
		Expect(review.Request.OldObject.Object.(*corev1.Pod).Namespace).To(Equal(review.Request.Namespace))
	})
	It("should provide handler type", func() {
		Expect(sut.HandlerType()).To(Equal(handler.TypeOther))
	})
})
