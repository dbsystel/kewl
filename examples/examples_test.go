package examples_test

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExamples(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "examples Suite")
}

func NewPod(name string, labelKeyValuePairs ...string) *corev1.Pod {
	labels := make(map[string]string, len(labelKeyValuePairs)/2)
	for i := 0; i < len(labelKeyValuePairs); i += 2 {
		labels[labelKeyValuePairs[i]] = labelKeyValuePairs[i+1]
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Labels: labels, Name: name},
	}
}
