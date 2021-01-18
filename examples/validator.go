package examples

import (
	"github.com/dbsystel/kewl/pkg/validation"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// PodValidator is an example validation.Validator which validates a corev1.Pod name
var PodValidator validation.Validator = &validator{}

type validator struct{}

func (t *validator) AddToScheme(scheme *runtime.Scheme) error {
	// NOTE: for each object to be validated, you need to invoke the scheme builder from k8s
	return corev1.AddToScheme(scheme)
}

func (t *validator) Name() string {
	// NOTE: choose whatever fancy name
	return "PodValidator"
}

func (t *validator) Validate(obj, _ runtime.Object, results validation.ResultCollector) error {
	if pod, ok := obj.(*corev1.Pod); ok {
		// NOTE: generally a good pattern is to delegate the cast object in case there's multiple validations (keeps code clean)
		t.validatePod(pod, results)
	}
	return nil
}

func (t *validator) validatePod(pod *corev1.Pod, results validation.ResultCollector) {
	// NOTE: since we validate the name field, we append a suffix name
	t.validatePodName(pod.Name, results.AppendField(".name"))
	// NOTE: same counts for labels (you get the idea)
	t.validateLabels(pod.Labels, results.AppendField(".labels"))
}

func (t *validator) validatePodName(name string, results validation.ResultCollector) {
	if name == "cat" {
		// NOTE: this adds the failure for the field we have set before the call
		results.AddFailure("switch to dog necessary")
	}
}

func (t *validator) validateLabels(labels map[string]string, results validation.ResultCollector) {
	for k, v := range labels {
		if v == "cat" {
			// NOTE: this is how you add a failure for a field the simple way
			results.AppendField("." + k).AddFailure("cats have bad habbits")
		}
	}
}
