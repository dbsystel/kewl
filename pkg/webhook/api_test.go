package webhook_test

import (
	"testing"

	"github.com/dbsystel/kewl/pkg/webhook/facade"
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebhook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "webhook Suite")
}

type ObjectWithKind interface {
	metav1.Object
	schema.ObjectKind
}

func prometheusLabels(version string, obj ObjectWithKind, responseType facade.ResponseType) prometheus.Labels {
	gvk := obj.GroupVersionKind()
	return prometheus.Labels{
		"obj_group": gvk.Group, "obj_version": gvk.Version, "obj_kind": gvk.Kind,
		"obj_namespace": obj.GetNamespace(), "result": string(responseType),
		"review_version": version,
	}
}
