package facade_test

import (
	"testing"

	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/testing/json_test"

	"github.com/dbsystel/kewl/pkg/webhook/facade"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func K8sAdmissionReview(review facade.AdmissionReview, k8sObj interface{}) interface{} {
	bytes, err := review.Marshal()
	panicutils.PanicIfError(err)
	json_test.UnmarshalJSONOrPanic(bytes, k8sObj)
	return k8sObj
}

func TestFacade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "webhook.facade Suite")
}

var _ = Describe("AdmissionReviewFrom test", func() {
	var typeMeta *metav1.TypeMeta
	BeforeEach(func() {
		typeMeta = &metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: v1.SchemeGroupVersion.Group + "/" + v1.SchemeGroupVersion.Version,
		}
	})
	Context("failures", func() {
		It("should fail if the testKind is incorrect", func() {
			typeMeta.Kind = "Foo"
			_, err := facade.AdmissionReviewFrom(json_test.MarshalJSONOrPanic(typeMeta))
			Expect(err).To(HaveOccurred())
		})
		It("should fail if the group is incorrect", func() {
			typeMeta.APIVersion = "Foo/" + v1.SchemeGroupVersion.Version
			_, err := facade.AdmissionReviewFrom(json_test.MarshalJSONOrPanic(typeMeta))
			Expect(err).To(HaveOccurred())
		})
		It("should fail if the version is incorrect", func() {
			typeMeta.APIVersion = v1.SchemeGroupVersion.Group + "/v2"
			_, err := facade.AdmissionReviewFrom(json_test.MarshalJSONOrPanic(typeMeta))
			Expect(err).To(HaveOccurred())
		})
		It("should fail if not a k8s object is present", func() {
			_, err := facade.AdmissionReviewFrom([]byte("a"))
			Expect(err).To(HaveOccurred())
		})
	})
})
