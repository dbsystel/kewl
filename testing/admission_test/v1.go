package admission_test

import (
	"github.com/dbsystel/kewl/testing"
	"github.com/dbsystel/kewl/testing/corev1_test"
	"github.com/dbsystel/kewl/testing/json_test"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// V1TypeMeta is the metav1.TypeMeta for v1.AdmissionReview
var V1TypeMeta = metav1.TypeMeta{
	Kind:       "AdmissionReview",
	APIVersion: v1.SchemeGroupVersion.Group + "/" + v1.SchemeGroupVersion.Version,
}

var _ testing.Marshalable = &V1AdmissionReview{}

type V1AdmissionReview v1.AdmissionReview

func (v *V1AdmissionReview) MustMarshal() []byte {
	return json_test.MarshalJSONOrPanic((*v1.AdmissionReview)(v))
}

//nolint:dupl
func NewV1Review(obj, oldObj testing.Reviewable) func() *V1AdmissionReview {
	rawExt, oldRawExt := rawExts(obj, oldObj)
	return func() *V1AdmissionReview {
		kind := obj.GetObjectKind().GroupVersionKind()
		return &V1AdmissionReview{
			TypeMeta: V1TypeMeta,
			Request: &v1.AdmissionRequest{
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

func rawExts(obj, oldObj testing.Reviewable) (oldExt, newExt runtime.RawExtension) {
	var rawExt, oldRawExt runtime.RawExtension
	if obj != nil {
		rawExt = obj.RawExtension()
	}
	if oldObj != nil {
		oldRawExt = oldObj.RawExtension()
	}
	return rawExt, oldRawExt
}

// V1ErrorPod is the v1.AdmissionReview for corev1_test.ErrorPod
var V1ErrorPod = NewV1Review(corev1_test.ErrorPod, nil)

// V1ValidPod is the v1.AdmissionReview for corev1_test.ValidPod
var V1ValidPod = NewV1Review(corev1_test.ValidPod, corev1_test.InvalidPod)

// V1InvalidPod is the v1.AdmissionReview for corev1_test.InvalidPod
var V1InvalidPod = NewV1Review(corev1_test.InvalidPod, corev1_test.ValidPod)

// V1BadPod is the v1.AdmissionReview for corev1_test.BadPod containing an unknown object
var V1BadPod = NewV1Review(corev1_test.BadPod, nil)
