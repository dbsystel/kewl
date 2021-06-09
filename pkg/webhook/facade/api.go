package facade

import (
	"encoding/json"
	"fmt"
	"reflect"

	v1 "k8s.io/api/admission/v1"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ResponseType string

const (
	// AdmissionAllowed denotes an allowed admission
	AdmissionAllowed ResponseType = "allowed"
	// AdmissionDenied denotes a denied admission
	AdmissionDenied ResponseType = "denied"
	// AdmissionMutated denotes a mutated admission
	AdmissionMutated ResponseType = "mutated"
	// AdmissionError denotes an erroneous
	AdmissionError ResponseType = "error"
	// AdmissionClientError denotes an error in the request
	AdmissionClientError ResponseType = "client_error"
)

const InvalidAdmissionReviewMsg = "invalid k8s admission review"

// AdmissionResponse facades the AdmissionReview response
type AdmissionResponse interface {
	// Allow allows the review
	Allow()
	// Deny denies the review using the provided status
	Deny(status *metav1.Status)
	// PatchJSON will apply a json patch to the response
	PatchJSON(bytes []byte)
	// ResponseType returns the ResponseType for statistics
	ResponseType() ResponseType
	// IsSet returns true if the response is set
	IsSet() bool
}

// AdmissionRequest facades the AdmissionReview request
type AdmissionRequest interface {
	// Kind returns the metav1.GroupVersionKind of the request object
	Kind() metav1.GroupVersionKind
	// Object returns the runtime.RawExtension representing the request object
	Object() *runtime.RawExtension
	// OldObject returns the runtime.RawExtension representing the request old object
	OldObject() *runtime.RawExtension
	// Resource returns the metav1.GroupVersionResource for the request object
	Resource() metav1.GroupVersionResource
	// Namespace returns the name of the namespace which is source to this request
	Namespace() string
}

// AdmissionReview is a facade for the admission review (to deal with the type-safety of different versions)
type AdmissionReview interface {
	// ClearRequest clears the request
	ClearRequest()
	// Request returns the request
	Request() AdmissionRequest
	// Response returns the response
	Response() AdmissionResponse
	// Version returns the version string of the admission review itself
	Version() string
	// Marshal marshals the object again
	Marshal() ([]byte, error)
}

// AdmissionReviewFrom tries to unmarshal an admission review form the provided bytes and create a decorator for it
func AdmissionReviewFrom(bytes []byte) (AdmissionReview, error) {
	// Unmarshall meta data first
	typeMeta := &metav1.TypeMeta{}
	if err := json.Unmarshal(bytes, typeMeta); err != nil {
		return nil, fmt.Errorf("%v - %v", InvalidAdmissionReviewMsg, "missing type metadata")
	}
	kind := typeMeta.GroupVersionKind()
	if kind.Kind != reflect.TypeOf(v1.AdmissionReview{}).Name() {
		return nil, fmt.Errorf("%v - %v: %v", InvalidAdmissionReviewMsg, "invalid object GroupVersionKind", kind)
	}
	if kind.Group == v1.SchemeGroupVersion.Group && kind.Version == v1.SchemeGroupVersion.Version {
		return V1AdmissionReviewFromBytes(bytes)
	}
	if kind.Group == v1beta1.SchemeGroupVersion.Group && kind.Version == v1beta1.SchemeGroupVersion.Version {
		return V1beta1AdmissionReviewFromBytes(bytes)
	}
	return nil, fmt.Errorf("could not create facade for: %v", kind)
}

// V1 creates an admission review for v1
func V1(target *v1.AdmissionReview) AdmissionReview {
	return &v1AdmissionReview{target: target}
}

// V1Beta1 creates an admission review for v1beta1
func V1Beta1(target *v1beta1.AdmissionReview) AdmissionReview {
	return &v1beta1AdmissionReview{target: target}
}
