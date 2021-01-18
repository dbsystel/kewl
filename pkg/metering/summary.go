package metering

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var _ Summary = &summaryImpl{}

type summaryImpl struct {
	summary *prometheus.SummaryVec
}

func (c *summaryImpl) MetricVec() *prometheus.MetricVec {
	return c.summary.MetricVec
}

func (c *summaryImpl) StartMetering(labels prometheus.Labels) FinishFn {
	// Create a simple timer
	start := time.Now()
	summary := c.summary.MustCurryWith(labels)
	// Return a function which will meter the time
	return func(finishLabels prometheus.Labels) {
		duration := time.Since(start)
		summary.With(finishLabels).Observe(duration.Seconds())
	}
}
