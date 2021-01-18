package webhook

import (
	"github.com/dbsystel/kewl/pkg/httpext"
	"github.com/dbsystel/kewl/pkg/metering"
	"github.com/dbsystel/kewl/pkg/webhook/handler"
	"github.com/prometheus/client_golang/prometheus"
)

var _ MetricsRegistry = &metricsRegistryImpl{}

const (
	labelAdmissionReviewVersion = "review_version"
	labelObjGroup               = "obj_group"
	labelObjVersion             = "obj_version"
	labelObjKind                = "obj_kind"
	labelObjNamespace           = "obj_namespace"
	labelResult                 = "result"
	metricsAliasWebHook         = "webhook."
	metricsAliasValidation      = metricsAliasWebHook + handler.TypeValidation
	metricsAliasMutation        = metricsAliasWebHook + handler.TypeMutation
)

// metricsRegistryImpl facades a metering.MetricsRegistry adding simple access to the server Metrics
type metricsRegistryImpl struct {
	httpext.MetricsRegistry
	namespace string
}

func (m *metricsRegistryImpl) Validation() metering.Summary {
	return m.WithAliasOrCreate(string(metricsAliasValidation), &prometheus.SummaryOpts{
		Namespace: m.namespace,
		Subsystem: "handler",
		Name:      string(handler.TypeValidation),
		Help:      "statics about the validation webhook",
	}, labelAdmissionReviewVersion, labelObjGroup, labelObjVersion, labelObjKind, labelObjNamespace, labelResult)
}

func (m *metricsRegistryImpl) Mutation() metering.Summary {
	return m.WithAliasOrCreate(string(metricsAliasMutation), &prometheus.SummaryOpts{
		Namespace: m.namespace,
		Subsystem: "handler",
		Name:      string(handler.TypeMutation),
		Help:      "statics about the mutation webhook",
	}, labelAdmissionReviewVersion, labelObjGroup, labelObjVersion, labelObjKind, labelObjNamespace, labelResult)
}

func (m *metricsRegistryImpl) WebHook() Metrics {
	return m
}
