package mutation_test

import (
	"github.com/dbsystel/kewl/pkg/mutation"
	"github.com/dbsystel/kewl/testing/corev1_test"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// PodMutator is an example mutation.Mutator which mutates a corev1.Pod name
var PodMutator mutation.Mutator = &mutator{}

type mutator struct{}

func (t *mutator) Name() string {
	return "PodMutator"
}

func (t *mutator) Mutate(newObject, _ runtime.Object) error {
	if pod, ok := newObject.(*corev1.Pod); ok {
		if pod.Name == corev1_test.PanicPod.Name {
			panic("platsch")
		}
		if pod.Name == corev1_test.ErrorPod.Name {
			return errors.New("blubb")
		}
		if pod.Name == corev1_test.InvalidPod.Name {
			pod.Name = corev1_test.ValidPod.Name
		}
	}
	return nil
}

func (t *mutator) AddToScheme(scheme *runtime.Scheme) error {
	return corev1.AddToScheme(scheme)
}
