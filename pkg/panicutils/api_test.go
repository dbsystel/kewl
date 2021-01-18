package panicutils_test

import (
	"testing"

	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/pkg/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPanicUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "panicutils Suite")
}

const panicMsg = "change diapers"

var _ = Describe("PanicIfError", func() {
	It("should panic on non-nil", func() {
		Expect(func() { panicutils.PanicIfError(errors.New("uff")) }).To(Panic())
	})
	It("should not panic on nil", func() {
		Expect(func() { panicutils.PanicIfError(nil) }).To(Not(Panic()))
	})
})

var _ = Describe("RecoverToErrorAndHandle", func() {
	var panicErr error
	handler := func(err error) {
		panicErr = err
	}
	AfterEach(func() {
		panicErr = nil
	})
	It("should recover panicked error", func() {
		err := errors.New(panicMsg)
		Expect(func() { defer panicutils.RecoverToErrorAndHandle(handler); panic(err) }).To(Not(Panic()))
		Expect(panicErr).To(Equal(err))
	})
	It("should recover simple panic", func() {
		Expect(func() { defer panicutils.RecoverToErrorAndHandle(handler); panic(panicMsg) }).To(Not(Panic()))
		Expect(panicErr).To(Not(BeNil()))
		Expect(panicErr.Error()).To(Equal("panic: " + panicMsg))
	})
	It("should continue on nil recovery", func() {
		Expect(func() { defer panicutils.RecoverToErrorAndHandle(handler) }).To(Not(Panic()))
		Expect(panicErr).To(BeNil())
	})
})

var _ = Describe("PanicToError", func() {
	It("should convert a panic", func() {
		Expect(panicutils.PanicToError("a")).NotTo(BeNil())
	})
	It("should return nils", func() {
		Expect(panicutils.PanicToError(nil)).To(BeNil())
	})
})
