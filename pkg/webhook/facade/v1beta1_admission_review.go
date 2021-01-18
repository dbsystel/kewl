package facade

import (
	"encoding/json"

	"github.com/pkg/errors"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func v1beta1AdmissionReviewFromBytes(bytes []byte) (AdmissionReview, error) {
	target := &v1beta1.AdmissionReview{}
	if err := json.Unmarshal(bytes, target); err != nil {
		return nil, errors.Wrap(err, "got an admission review v1beta1, but could not serialize it")
	}
	return &v1beta1AdmissionReview{target}, nil
}

var v1beta1PatchTypeJSONPatch = v1beta1.PatchTypeJSONPatch

type v1beta1AdmissionReview struct {
	target *v1beta1.AdmissionReview
}

// Review decorator functions
var _ AdmissionReview = &v1beta1AdmissionReview{}

func (v *v1beta1AdmissionReview) Marshal() ([]byte, error) {
	return json.Marshal(v.target)
}

func (v *v1beta1AdmissionReview) Request() AdmissionRequest {
	return v
}

func (v *v1beta1AdmissionReview) Response() AdmissionResponse {
	return v
}

func (v *v1beta1AdmissionReview) Kind() metav1.GroupVersionKind {
	return v.target.Request.Kind
}

// Request decorator functions
var _ AdmissionRequest = &v1beta1AdmissionReview{}

func (v *v1beta1AdmissionReview) Namespace() string {
	return v.target.Request.Namespace
}

func (v *v1beta1AdmissionReview) Object() *runtime.RawExtension {
	return &v.target.Request.Object
}

func (v *v1beta1AdmissionReview) OldObject() *runtime.RawExtension {
	return &v.target.Request.OldObject
}

func (v *v1beta1AdmissionReview) Resource() metav1.GroupVersionResource {
	return v.target.Request.Resource
}

func (v v1beta1AdmissionReview) Version() string {
	return v1beta1.SchemeGroupVersion.Version
}

// Response decorator functions
var _ AdmissionResponse = &v1beta1AdmissionReview{}

func (v *v1beta1AdmissionReview) withResponse(handler func(response *v1beta1.AdmissionResponse)) {
	if v.target.Response == nil {
		v.target.Response = &v1beta1.AdmissionResponse{}
	}
	handler(v.target.Response)
}

func (v *v1beta1AdmissionReview) Allow() {
	v.withResponse(func(response *v1beta1.AdmissionResponse) {
		response.Allowed = true
	})
}

func (v *v1beta1AdmissionReview) Deny(status *metav1.Status) {
	v.withResponse(func(response *v1beta1.AdmissionResponse) {
		response.Allowed = false
		response.Result = status
	})
}

func (v *v1beta1AdmissionReview) PatchJSON(bytes []byte) {
	v.withResponse(func(response *v1beta1.AdmissionResponse) {
		response.Allowed = true
		if len(bytes) > 0 {
			response.PatchType = &v1beta1PatchTypeJSONPatch
			response.Patch = bytes
		}
	})
}

func (v *v1beta1AdmissionReview) ResponseType() ResponseType {
	if v.target.Response == nil {
		return AdmissionError
	}
	if len(v.target.Response.Patch) > 0 {
		return AdmissionMutated
	}
	if v.target.Response.Allowed {
		return AdmissionAllowed
	}
	return AdmissionDenied
}
