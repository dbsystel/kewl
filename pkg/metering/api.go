package metering

import (
	"github.com/prometheus/client_golang/prometheus"
)

// FinishFn is the function which will finish the metring for a request using the provided counter labels
type FinishFn func(counterLabels prometheus.Labels)

// Summary is provides a facade for metering using a counter and a summary
type Summary interface {
	// StartMetering starts the metering for the provided prometheus.Labels returning a FinishFn to end the metering
	StartMetering(labels prometheus.Labels) FinishFn
	// MetricVec returns the prometheus.MetricVec for the summary
	MetricVec() *prometheus.MetricVec
}

// MetricsRegistry is a facade for registering metrics
type MetricsRegistry interface {
	// PromRegistry returns the prometheus registry
	PromRegistry() *prometheus.Registry
	// NewSummary gets or creates and registers a new Summary
	NewSummary(alias string, promOpts *prometheus.SummaryOpts, labelNames ...string) Summary
	// WithAlias returns the Summary with the provided alias
	WithAlias(alias string) Summary
	// WithAliasOrCreate tries to get a metering.Summary by the provided alias, or create it using the parameters
	WithAliasOrCreate(alias string, promOpts *prometheus.SummaryOpts, labelNames ...string) Summary
}

// NewRegistry creates a new MetricsRegistry backed by a new prometheus.Registry
func NewRegistry(registry *prometheus.Registry) MetricsRegistry {
	return &promRegistryFacade{Registry: registry}
}
