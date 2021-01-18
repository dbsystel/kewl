package metering_test

import (
	"github.com/dbsystel/kewl/pkg/metering"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelX = "x"
	labelY = "y"
	labelZ = "z"
)

var _ = Describe("Summary", func() {
	startLabels := prometheus.Labels{labelX: "1", labelY: "2"}
	var sut metering.Summary
	BeforeEach(func() {
		sut = metering.NewRegistry(prometheus.NewRegistry()).NewSummary("test", &prometheus.SummaryOpts{
			Namespace: "prom",
			Subsystem: "metering",
			Name:      "test",
			Help:      "test",
		}, labelX, labelY, labelZ)
	})
	It("should meter a call correctly", func() {
		sut.StartMetering(startLabels)(prometheus.Labels{labelZ: "3"})
		Expect(sut.MetricVec().GetMetricWithLabelValues("1", "2", "3")).To(Not(BeNil()))
	})
	It("should handle problems with the labels", func() {
		Context("invalid label when starting the metering", func() {
			Expect(func() {
				sut.StartMetering(prometheus.Labels{labelX: "1", "a": "2"})
			}).To(Panic())
		})
		Context("invalid label when finishing the metering", func() {
			finish := sut.StartMetering(startLabels)
			Expect(func() {
				finish(prometheus.Labels{"a": "3"})
			}).To(Panic())
		})
	})
})
