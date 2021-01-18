package httpext_test

import (
	"net/http"
	"strconv"
	"syscall"
	"time"

	"github.com/dbsystel/kewl/testing/httpext_test"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dbsystel/kewl/pkg/panicutils"

	"github.com/dbsystel/kewl/pkg/httpext"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	dto "github.com/prometheus/client_model/go"
)

type testJSON struct {
	Int    int    `json:"int"`
	String string `json:"string"`
}

var _ = Describe("Server", func() {

	var fixture *httpext_test.Fixture
	BeforeEach(func() {
		fixture = httpext_test.NewFixture()
	})
	requestMetric := func(method, path string, status int) *dto.Summary {
		result, err := fixture.Sut.Metrics.HTTP().Requests().MetricVec().GetMetricWith(prometheus.Labels{
			"method": method, "path": path, "status": strconv.Itoa(status),
		})
		Expect(err).NotTo(HaveOccurred())
		dto := &dto.Metric{}
		Expect(result.Write(dto)).NotTo(HaveOccurred())
		return dto.Summary
	}
	It("serves /healthz", func() {
		Expect(fixture.Test.Healthz().String()).To(BeEquivalentTo(httpext.HealthzResponse))
	})
	It("serves /metrics", func() {
		Expect(fixture.Test.Metrics()).NotTo(BeEmpty())
	})
	It("adds a trace header for the extended handlers", func() {
		fixture.Sut.HandleExt("/json", func(writer httpext.ResponseWriter, request *httpext.Request) {
			writer.SendResponse(http.StatusNoContent, nil)
		})
		response := fixture.Test.RequestBytes("PUT", "/json", nil)
		Expect(response.Code).To(Equal(http.StatusNoContent))
		Expect(response.Body.String()).To(BeEmpty())
		Expect(response.Header().Get(httpext.HeaderTraceID)).To(Not(BeEmpty()))
	})
	It("reads and serves JSON via decorated handler", func() {
		requestObj := &testJSON{Int: 1, String: "uff"}
		fixture.Sut.HandleExt("/json", func(writer httpext.ResponseWriter, request *httpext.Request) {
			reqObj := &testJSON{}
			panicutils.PanicIfError(request.ReadBodyAsJSON(reqObj))
			writer.SendJSON(http.StatusCreated, reqObj)
		})
		response := fixture.Test.RequestJSON("PUT", "/json", requestObj)
		Expect(response.Code).To(Equal(http.StatusCreated))
		Expect(response.JSON(&testJSON{})).To(BeEquivalentTo(requestObj))
		Expect(response.Header().Get(httpext.HeaderContentType)).To(Equal("application/json"))
		metric := requestMetric("PUT", "/json", http.StatusCreated)
		Expect(*metric.SampleCount).To(BeNumerically(">", 0))
		Expect(*metric.SampleSum).To(BeNumerically(">", 0))
	})
	It("errors when reading empty json request", func() {
		fixture.Sut.HandleExt("/emptyJson", func(writer httpext.ResponseWriter, request *httpext.Request) {
			Expect(request.ReadBodyAsJSON(&testJSON{})).To(HaveOccurred())
		})
		response := fixture.Test.RequestBytes("PUT", "/emptyJson", nil)
		Expect(response).NotTo(BeNil())
	})
	It("errors when reading invalid json from request", func() {
		fixture.Sut.HandleExt("/invalidJson", func(writer httpext.ResponseWriter, request *httpext.Request) {
			Expect(request.ReadBodyAsJSON(&testJSON{})).To(HaveOccurred())
		})
		response := fixture.Test.RequestBytes("PUT", "/invalidJson", nil)
		Expect(response).NotTo(BeNil())
	})
	It("reads bytes and serves text via decorated handler", func() {
		reqBody := []byte("uff")
		resBody := "meh"
		fixture.Sut.HandleExt("/text", func(writer httpext.ResponseWriter, request *httpext.Request) {
			Expect(request.ReadBody()).To(Equal(reqBody))
			writer.SendResponseString(http.StatusOK, resBody)
		})
		response := fixture.Test.RequestBytes("PUT", "/text", reqBody)
		Expect(response.Code).To(Equal(http.StatusOK))
		Expect(response.Body.String()).To(Equal(resBody))
		Expect(response.Header().Get(httpext.HeaderContentType)).To(Equal("text/plain"))
		metric := requestMetric("PUT", "/text", http.StatusOK)
		Expect(*metric.SampleCount).To(BeNumerically(">", 0))
		Expect(*metric.SampleSum).To(BeNumerically(">", 0))
	})
	It("handles internal server errors", func() {
		fixture.Sut.HandleExt("/text", func(writer httpext.ResponseWriter, request *httpext.Request) {
			writer.HandleInternalError(errors.New("yo"))
		})
		response := fixture.Test.RequestBytes("PUT", "/text", nil)
		Expect(response.Code).To(Equal(http.StatusInternalServerError))
		Expect(response.Body.String()).To(Not(BeEmpty()))
		metric := requestMetric("PUT", "/text", http.StatusInternalServerError)
		Expect(*metric.SampleCount).To(BeNumerically(">", 0))
		Expect(*metric.SampleSum).To(BeNumerically(">", 0))
	})
	It("provides request logger", func() {
		fixture.Sut.HandleExt("/log", func(writer httpext.ResponseWriter, request *httpext.Request) {
			Expect(request.Logger()).To(Not(BeNil()))
		})
		response := fixture.Test.RequestBytes("PUT", "/log", nil)
		Expect(response).NotTo(BeNil())
	})
	It("provides set status", func() {
		fixture.Sut.HandleExt("/log", func(writer httpext.ResponseWriter, request *httpext.Request) {
			writer.SendResponse(http.StatusTeapot, nil)
			Expect(writer.Status()).To(Equal(http.StatusTeapot))
		})
		response := fixture.Test.RequestBytes("PUT", "/log", nil)
		Expect(response).NotTo(BeNil())
		metric := requestMetric("PUT", "/log", http.StatusTeapot)
		Expect(*metric.SampleCount).To(BeNumerically(">", 0))
		Expect(*metric.SampleSum).To(BeNumerically(">", 0))
	})
	It("handles panics", func() {
		fixture.Sut.HandleExt("/panic", func(writer httpext.ResponseWriter, request *httpext.Request) {
			panic(errors.New("check diapers"))
		})
		response := fixture.Test.RequestBytes("PUT", "/panic", nil)
		Expect(response.Code).To(Equal(http.StatusInternalServerError))
		metric := requestMetric("PUT", "/panic", http.StatusInternalServerError)
		Expect(*metric.SampleCount).To(BeNumerically(">", 0))
		Expect(*metric.SampleSum).To(BeNumerically(">", 0))
	})
	It("records the requests using prometheus", func() {
		fixture.Sut.HandleExt("/nyan", func(writer httpext.ResponseWriter, request *httpext.Request) {
			writer.SendResponseString(http.StatusOK, "cat")
		})
		_ = fixture.Test.RequestBytes("GET", "/nyan", nil)
		metric := requestMetric("GET", "/nyan", http.StatusOK)
		Expect(*metric.SampleCount).To(BeNumerically(">", 0))
		Expect(*metric.SampleSum).To(BeNumerically(">", 0))

	})
	It("should abort running on signal http", func() {
		// Run the server
		errChan := fixture.Sut.RunWaitingForSig(syscall.SIGUSR1)
		// After 1 second, try to shut it down
		<-time.After(1 * time.Second)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

		// Handle the result
		select {
		case err := <-errChan:
			// Expect no errors have occurred
			Expect(err).To(Not(HaveOccurred()))
		case <-time.After(2 * time.Second):
			// Make sure we time out
			Fail("server was not stopped after signal")
		}
	})
	It("should abort running on signal https", func() {
		fixture = httpext_test.NewFixtureHTTP()
		// Run the server
		errChan := fixture.Sut.RunWaitingForSig(syscall.SIGUSR1)
		// After 1 second, try to shut it down
		<-time.After(1 * time.Second)
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

		// Handle the result
		select {
		case err := <-errChan:
			// Expect no errors have occurred
			Expect(err).To(Not(HaveOccurred()))
		case <-time.After(2 * time.Second):
			// Make sure we time out
			Fail("server was not stopped after signal")
		}
	})
})
