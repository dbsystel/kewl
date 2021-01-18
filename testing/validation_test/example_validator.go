package validation_test

import (
	"github.com/dbsystel/kewl/pkg/validation"
	"github.com/dbsystel/kewl/testing/corev1_test"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ExamplePodValidator interface {
	validation.Validator
	InvalidStatusCause() []v1.StatusCause
}

// PodValidator is an example validation.Validator which validates a corev1.Pod name
var PodValidator ExamplePodValidator = &validator{}

type validator struct{}

func (t *validator) InvalidStatusCause() []v1.StatusCause {
	return []v1.StatusCause{{Type: "invalid", Message: "not valid", Field: ".name"}}
}

func (t *validator) AddToScheme(scheme *runtime.Scheme) error {
	return corev1.AddToScheme(scheme)
}

func (t *validator) Name() string {
	return "PodValidator"
}

func (t *validator) Validate(obj, _ runtime.Object, results validation.ResultCollector) error {
	if pod, ok := obj.(*corev1.Pod); ok {
		if pod.Name == corev1_test.PanicPod.Name {
			panic("platsch")
		}
		if pod.Name == corev1_test.ErrorPod.Name {
			return errors.New("blubb")
		}
		if pod.Name != corev1_test.ValidPod.Name {
			cause := t.InvalidStatusCause()[0]
			results.AppendField(cause.Field).AddFailure(cause.Message)
		}
	}
	return nil
}
