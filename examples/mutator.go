package examples

import (
	"github.com/dbsystel/kewl/pkg/mutation"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// PodMutator is an example mutation.Mutator which mutates a corev1.Pod name
var PodMutator mutation.Mutator = &mutator{}

type mutator struct{}

func (t *mutator) AddToScheme(scheme *runtime.Scheme) error {
	// NOTE: for each object to be validated, you need to invoke the scheme builder from k8s
	return corev1.AddToScheme(scheme)
}

func (t *mutator) Name() string {
	// NOTE: choose whatever fancy name
	return "PodMutator"
}

func (t *mutator) Mutate(obj, _ runtime.Object) error {
	if pod, ok := obj.(*corev1.Pod); ok {
		// NOTE: generally a good pattern is to delegate the cast object in case there's multiple validations (keeps code clean)
		t.mutatePod(pod)
	}
	return nil
}

func (t *mutator) mutatePod(pod *corev1.Pod) {
	// NOTE: whatever you'll do here, will be reflected in the patch operation in the admission response
	if pod.Name == "cat" {
		pod.Name = "dog"
	}
	for k, v := range pod.Labels {
		if v == "cat" {
			pod.Labels[k] = "dog"
		}
	}
}
