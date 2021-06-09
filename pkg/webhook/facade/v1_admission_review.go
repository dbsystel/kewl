// nolint:dupl
package facade

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/types"

	"github.com/pkg/errors"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var v1PatchTypeJSONPatch = v1.PatchTypeJSONPatch

func V1AdmissionReviewFromBytes(bytes []byte) (AdmissionReview, error) {
	target := &v1.AdmissionReview{}
	if err := json.Unmarshal(bytes, target); err != nil {
		return nil, errors.Wrap(err, "got an admission review v1, but could not serialize it")
	}
	return V1(target), nil
}

type v1AdmissionReview struct {
	target  *v1.AdmissionReview
	request *v1AdmissionReviewRequest
}

// Review decorator functions
var _ AdmissionReview = &v1AdmissionReview{}

func (v *v1AdmissionReview) Marshal() ([]byte, error) {
	return json.Marshal(v.target)
}

func (v *v1AdmissionReview) Request() AdmissionRequest {
	if v.request == nil && v.target.Request != nil {
		v.request = &v1AdmissionReviewRequest{target: v.target.Request}
	}
	return v.request
}

func (v *v1AdmissionReview) Response() AdmissionResponse {
	return v
}

func (v *v1AdmissionReview) Version() string {
	return v1.SchemeGroupVersion.Version
}

func (v *v1AdmissionReview) ClearRequest() {
	v.target.Request = nil
}

// Request decorator functions
var _ AdmissionRequest = &v1AdmissionReviewRequest{}

type v1AdmissionReviewRequest struct {
	target *v1.AdmissionRequest
}

func (v *v1AdmissionReviewRequest) ResourceID() types.NamespacedName {
	return types.NamespacedName{Namespace: v.target.Namespace, Name: v.target.Name}
}

func (v *v1AdmissionReviewRequest) Namespace() string {
	return v.target.Namespace
}

func (v *v1AdmissionReviewRequest) Kind() metav1.GroupVersionKind {
	return v.target.Kind
}

func (v *v1AdmissionReviewRequest) Object() *runtime.RawExtension {
	return &v.target.Object
}

func (v *v1AdmissionReviewRequest) OldObject() *runtime.RawExtension {
	return &v.target.OldObject
}

func (v *v1AdmissionReviewRequest) ResourceKind() metav1.GroupVersionResource {
	return v.target.Resource
}

// Response decorator functions
var _ AdmissionResponse = &v1AdmissionReview{}

func (v *v1AdmissionReview) withResponse(handler func(response *v1.AdmissionResponse)) {
	if v.target.Response == nil {
		if v.target.Request == nil {
			return
		}
		v.target.Response = &v1.AdmissionResponse{UID: v.target.Request.UID}
	}
	handler(v.target.Response)
}

func (v *v1AdmissionReview) Allow() {
	v.withResponse(func(response *v1.AdmissionResponse) {
		response.Allowed = true
	})
}

func (v *v1AdmissionReview) Deny(status *metav1.Status) {
	v.withResponse(func(response *v1.AdmissionResponse) {
		response.Allowed = false
		response.Result = status
	})
}

func (v *v1AdmissionReview) PatchJSON(bytes []byte) {
	v.withResponse(func(response *v1.AdmissionResponse) {
		response.Allowed = true
		if len(bytes) > 0 {
			response.PatchType = &v1PatchTypeJSONPatch
			response.Patch = bytes
		}
	})
}

func (v *v1AdmissionReview) IsSet() bool {
	return v.target.Response != nil
}

func (v *v1AdmissionReview) ResponseType() ResponseType {
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
