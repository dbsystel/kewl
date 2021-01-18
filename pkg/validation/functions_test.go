package validation_test

import (
	"github.com/dbsystel/kewl/pkg/validation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NonNil validation", func() {
	var collector validation.ResultCollector
	BeforeEach(func() {
		collector = validation.NewResultCollector()
	})
	It("should add a failure in case nil is encountered", func() {
		validation.EnsureNonNil(nil, collector)
		Expect(collector.Failures()).To(Not(BeEmpty()))
	})
	It("should complete without failures in case non nil is encountered", func() {
		validation.EnsureNonNil(2, collector)
		Expect(collector.Failures()).To(BeEmpty())
	})
})
