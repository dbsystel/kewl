package metering_test

import (
	"github.com/dbsystel/kewl/pkg/metering"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
)

var _ = Describe("MetricsRegistry", func() {
	var sut metering.MetricsRegistry
	BeforeEach(func() {
		sut = metering.NewRegistry(prometheus.NewRegistry())
	})
	It("should return nil for unknown", func() {
		Expect(sut.WithAlias("test")).To(BeNil())
	})
	It("should create a new summary and provide it by alias", func() {
		summary := sut.NewSummary("test", &prometheus.SummaryOpts{Name: "test"})
		Expect(summary).NotTo(BeNil())
		Expect(sut.WithAlias("test")).To(Equal(summary))
	})
	It("should panic in case we register twice", func() {
		sut.NewSummary("test", &prometheus.SummaryOpts{Name: "test"})
		Expect(func() { sut.NewSummary("test", &prometheus.SummaryOpts{Name: "test"}) }).To(Panic())
	})
	It("should create a new summary and provide it by alias or create a new one", func() {
		summary := sut.WithAliasOrCreate("test", &prometheus.SummaryOpts{Name: "test"})
		Expect(summary).NotTo(BeNil())
		Expect(summary).NotTo(Equal(sut.WithAliasOrCreate("test3", &prometheus.SummaryOpts{Name: "test3"})))
		Expect(sut.WithAliasOrCreate("test", &prometheus.SummaryOpts{Name: "test"})).To(Equal(summary))
	})

})
