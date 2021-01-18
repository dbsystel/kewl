package httpext

import (
	"github.com/dbsystel/kewl/pkg/metering"
	"github.com/prometheus/client_golang/prometheus"
)

const labelMethod = "method"
const labelPath = "path"
const labelStatus = "status"

var _ MetricsRegistry = &metricsRegistryImpl{}
var _ Metrics = &metricsRegistryImpl{}

const metricsAliasHTTPRequests = "http.requests"

// metricsRegistryImpl facades a metering.MetricsRegistry adding simple access to the server Metrics
type metricsRegistryImpl struct {
	metering.MetricsRegistry
	namespace string
}

func (m *metricsRegistryImpl) Requests() metering.Summary {
	return m.WithAliasOrCreate(metricsAliasHTTPRequests, &prometheus.SummaryOpts{
		Namespace: m.namespace,
		Subsystem: "http",
		Name:      "request_seconds",
		Help:      "Handled http requests and their response codes",
	}, labelMethod, labelPath, labelStatus)
}

func (m *metricsRegistryImpl) HTTP() Metrics {
	return m
}
