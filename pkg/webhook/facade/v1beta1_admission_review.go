// nolint:dupl
package facade

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/types"

	"github.com/pkg/errors"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var v1beta1PatchTypeJSONPatch = v1beta1.PatchTypeJSONPatch

func V1beta1AdmissionReviewFromBytes(bytes []byte) (AdmissionReview, error) {
	target := &v1beta1.AdmissionReview{}
	if err := json.Unmarshal(bytes, target); err != nil {
		return nil, errors.Wrap(err, "got an admission review v1beta1, but could not serialize it")
	}
	return V1Beta1(target), nil
}

type v1beta1AdmissionReview struct {
	target  *v1beta1.AdmissionReview
	request *v1beta1AdmissionReviewRequest
}

// Review decorator functions
var _ AdmissionReview = &v1beta1AdmissionReview{}

func (v *v1beta1AdmissionReview) Marshal() ([]byte, error) {
	return json.Marshal(v.target)
}

func (v *v1beta1AdmissionReview) Request() AdmissionRequest {
	if v.request == nil && v.target.Request != nil {
		v.request = &v1beta1AdmissionReviewRequest{target: v.target.Request}
	}
	return v.request
}

func (v *v1beta1AdmissionReview) Response() AdmissionResponse {
	return v
}

func (v *v1beta1AdmissionReview) Version() string {
	return v1beta1.SchemeGroupVersion.Version
}

func (v *v1beta1AdmissionReview) ClearRequest() {
	v.target.Request = nil
}

// Request decorator functions
var _ AdmissionRequest = &v1beta1AdmissionReviewRequest{}

type v1beta1AdmissionReviewRequest struct {
	target *v1beta1.AdmissionRequest
}

func (v *v1beta1AdmissionReviewRequest) ResourceID() types.NamespacedName {
	return types.NamespacedName{Namespace: v.target.Namespace, Name: v.target.Name}
}

func (v *v1beta1AdmissionReviewRequest) Namespace() string {
	return v.target.Namespace
}

func (v *v1beta1AdmissionReviewRequest) Kind() metav1.GroupVersionKind {
	return v.target.Kind
}

func (v *v1beta1AdmissionReviewRequest) Object() *runtime.RawExtension {
	return &v.target.Object
}

func (v *v1beta1AdmissionReviewRequest) OldObject() *runtime.RawExtension {
	return &v.target.OldObject
}

func (v *v1beta1AdmissionReviewRequest) ResourceKind() metav1.GroupVersionResource {
	return v.target.Resource
}

// Response decorator functions
var _ AdmissionResponse = &v1beta1AdmissionReview{}

func (v *v1beta1AdmissionReview) withResponse(handler func(response *v1beta1.AdmissionResponse)) {
	if v.target.Response == nil {
		if v.target.Request == nil {
			return
		}
		v.target.Response = &v1beta1.AdmissionResponse{UID: v.target.Request.UID}
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

func (v *v1beta1AdmissionReview) IsSet() bool {
	return v.target.Response != nil
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
