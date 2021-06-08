package admission_test

import (
	"github.com/dbsystel/kewl/testing"
	"github.com/dbsystel/kewl/testing/corev1_test"
	"github.com/dbsystel/kewl/testing/json_test"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// V1Beta1TypeMeta is the metav1.TypeMeta for v1beta1.AdmissionReview
var V1Beta1TypeMeta = metav1.TypeMeta{
	Kind:       "AdmissionReview",
	APIVersion: v1beta1.SchemeGroupVersion.Group + "/" + v1beta1.SchemeGroupVersion.Version,
}
var _ testing.Marshalable = &V1Beta1AdmissionReview{}

type V1Beta1AdmissionReview v1beta1.AdmissionReview

func (v *V1Beta1AdmissionReview) MustMarshal() []byte {
	return json_test.MarshalJSONOrPanic((*v1beta1.AdmissionReview)(v))
}

// nolint:dupl
func NewV1Beta1Review(obj, oldObj testing.Reviewable) func() *V1Beta1AdmissionReview {
	rawExt, oldRawExt := rawExts(obj, oldObj)
	return func() *V1Beta1AdmissionReview {
		kind := obj.GetObjectKind().GroupVersionKind()
		return &V1Beta1AdmissionReview{
			TypeMeta: V1Beta1TypeMeta,
			Request: &v1beta1.AdmissionRequest{
				Kind: metav1.GroupVersionKind{Group: kind.Group, Version: kind.Version, Kind: kind.Kind},
				Resource: metav1.GroupVersionResource{
					Group:    kind.Group,
					Version:  kind.Version,
					Resource: obj.GetNamespace() + "/" + obj.GetName(),
				},
				Name:      obj.GetName(),
				Namespace: obj.GetNamespace(),
				Object:    rawExt,
				OldObject: oldRawExt,
			},
			Response: nil,
		}
	}
}

// V1Beta1ErrorPod is the v1beta1.AdmissionReview for corev1_test.ErrorPod
var V1Beta1ErrorPod = NewV1Beta1Review(corev1_test.ErrorPod, nil)

// V1Beta1ValidPod is the v1beta1.AdmissionReview for corev1_test.ValidPod
var V1Beta1ValidPod = NewV1Beta1Review(corev1_test.ValidPod, nil)

// V1Beta1InvalidPod is the v1beta1.AdmissionReview for corev1_test.InvalidPod
var V1Beta1InvalidPod = NewV1Beta1Review(corev1_test.InvalidPod, corev1_test.ValidPod)

// V1Beta1BadPod is the v1.AdmissionReview for corev1_test.BadPod
var V1Beta1BadPod = NewV1Beta1Review(corev1_test.BadPod, nil)
