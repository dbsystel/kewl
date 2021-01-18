package corev1_test

import (
	"github.com/dbsystel/kewl/testing"
	"github.com/dbsystel/kewl/testing/json_test"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Namespace for the pods
const Namespace = "test"

// PodTypeMeta is the metav1.TypeMeta for corev1.Pod
var PodTypeMeta = metav1.TypeMeta{
	Kind:       "Pod",
	APIVersion: corev1.SchemeGroupVersion.Group + "/" + corev1.SchemeGroupVersion.Version,
}

// PodKind is the metav1.GroupVersionKind for corev1.Pod
var PodKind = metav1.GroupVersionKind{Group: corev1.SchemeGroupVersion.Group, Version: corev1.SchemeGroupVersion.Version, Kind: "Pod"}

var _ testing.Reviewable = &Pod{}
var _ runtime.Object = &Pod{}

// Pod extends corev1.Pod for testing
type Pod corev1.Pod

func (t *Pod) MustMarshal() []byte {
	return json_test.MarshalJSONOrPanic((*corev1.Pod)(t))
}

// RawExtension marshals the pod and creates a runtime.RawExtension for which Raw is set to the marshaled JSON
func (t *Pod) RawExtension() runtime.RawExtension {
	return runtime.RawExtension{Raw: t.MustMarshal()}
}

func (t *Pod) DeepCopyObject() runtime.Object {
	return (*Pod)((*corev1.Pod)(t).DeepCopy())
}

func (t *Pod) AsCoreV1() *corev1.Pod {
	if t == nil {
		return nil
	}
	return (*corev1.Pod)(t)
}

// NewPod creates a new Pod with the provided name
func NewPod(name string) *Pod {
	return &Pod{TypeMeta: PodTypeMeta, ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: Namespace}}
}

// NewBrokenPod creates a new broken Pod with the provided name
func NewBrokenPod(name string) *Pod {
	return &Pod{TypeMeta: metav1.TypeMeta{
		Kind:       "BrokenPod",
		APIVersion: PodTypeMeta.APIVersion,
	}, ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: Namespace}}
}

// ErrorPod is a pod which should create an error on handling
var ErrorPod = NewPod("error")

// PanicPod is a pod which should create a panic on handling
var PanicPod = NewPod("panic")

// ValidPod which is considered valid on handling
var ValidPod = NewPod("valid")

// InvalidPod which is considered invalid on handling
var InvalidPod = NewPod("invalid")

// BadPod is a pod which does not serialize correctly
var BadPod = NewBrokenPod("broken")
