package webhook_test

import (
	"net/http"

	"github.com/dbsystel/kewl/testing/admission_test"
	"github.com/dbsystel/kewl/testing/httpext_test"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ReviewFacade struct {
	*httpext_test.TestFacade
}

func (f *ReviewFacade) ValidateV1(pod *admission_test.V1AdmissionReview) *v1.AdmissionResponse {
	review := f.validate(pod, &v1.AdmissionReview{})
	if review == nil {
		return nil
	}
	return review.(*v1.AdmissionReview).Response
}

func (f *ReviewFacade) ValidateV1Beta1(pod *admission_test.V1Beta1AdmissionReview) *v1beta1.AdmissionResponse {
	review := f.validate(pod, &v1beta1.AdmissionReview{})
	if review == nil {
		return nil
	}
	return review.(*v1beta1.AdmissionReview).Response
}

func (f *ReviewFacade) MutateV1(pod *admission_test.V1AdmissionReview) *v1.AdmissionResponse {
	review := f.mutate(pod, &v1.AdmissionReview{})
	if review == nil {
		return nil
	}
	return review.(*v1.AdmissionReview).Response
}

func (f *ReviewFacade) MutateV1Beta1(pod *admission_test.V1Beta1AdmissionReview) *v1beta1.AdmissionResponse {
	review := f.mutate(pod, &v1beta1.AdmissionReview{})
	if review == nil {
		return nil
	}
	return review.(*v1beta1.AdmissionReview).Response
}

func (f *ReviewFacade) validate(review interface{}, result runtime.Object) runtime.Object {
	return f.invokeWebHook("/validate", review, result)
}

func (f *ReviewFacade) mutate(review interface{}, result runtime.Object) runtime.Object {
	return f.invokeWebHook("/mutate", review, result)
}

func (f *ReviewFacade) invokeWebHook(path string, review interface{}, result runtime.Object) runtime.Object {
	response := f.RequestJSON("POST", path, review)
	if response.Code >= http.StatusMultipleChoices {
		return nil
	}
	return response.JSON(result).(runtime.Object)
}
