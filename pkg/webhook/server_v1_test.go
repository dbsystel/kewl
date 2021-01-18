//nolint:dupl
package webhook_test

import (
	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/pkg/webhook/facade"
	"github.com/dbsystel/kewl/testing/admission_test"
	"github.com/dbsystel/kewl/testing/corev1_test"
	"github.com/dbsystel/kewl/testing/mutation_test"
	"github.com/dbsystel/kewl/testing/validation_test"
	"github.com/dbsystel/kewl/testing/webhook_test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	dto "github.com/prometheus/client_model/go"
)

var _ = Describe("Server v1 test", func() {
	var fixture *webhook_test.Fixture
	BeforeEach(func() {
		fixture = webhook_test.NewFixture()
		panicutils.PanicIfError(fixture.Sut.AddValidator(validation_test.PodValidator))
		panicutils.PanicIfError(fixture.Sut.AddMutator(mutation_test.PodMutator))
	})
	validationMetrics := func(version string, obj ObjectWithKind, responseType facade.ResponseType) *dto.Summary {
		result, err := fixture.Sut.Metrics.WebHook().Validation().MetricVec().
			GetMetricWith(prometheusLabels(version, obj, responseType))
		Expect(err).NotTo(HaveOccurred())
		dto := &dto.Metric{}
		Expect(result.Write(dto)).NotTo(HaveOccurred())
		return dto.Summary
	}
	mutationMetrics := func(version string, obj ObjectWithKind, responseType facade.ResponseType) *dto.Summary {
		result, err := fixture.Sut.Metrics.WebHook().Mutation().MetricVec().
			GetMetricWith(prometheusLabels(version, obj, responseType))
		Expect(err).NotTo(HaveOccurred())
		dto := &dto.Metric{}
		Expect(result.Write(dto)).NotTo(HaveOccurred())
		return dto.Summary
	}
	admRevVer := "v1"
	Context("validation", func() {
		It("should handle valid objects correctly", func() {
			response := fixture.Test.ValidateV1(admission_test.V1ValidPod())
			Expect(response).NotTo(BeNil())
			Expect(response.Allowed).To(BeTrue())
			Expect(response.Patch).To(BeNil())
			Expect(response.PatchType).To(BeNil())
			metrics := validationMetrics(admRevVer, corev1_test.ValidPod, facade.AdmissionAllowed)
			Expect(*metrics.SampleCount).To(BeNumerically("==", 1))
			Expect(*metrics.SampleSum).To(BeNumerically(">", 0))
		})
		It("should handle invalid objects correctly", func() {
			response := fixture.Test.ValidateV1(admission_test.V1InvalidPod())
			Expect(response).NotTo(BeNil())
			Expect(response.Allowed).To(BeFalse())
			Expect(response.Patch).To(BeNil())
			Expect(response.PatchType).To(BeNil())
			metrics := validationMetrics(admRevVer, corev1_test.InvalidPod, facade.AdmissionDenied)
			Expect(*metrics.SampleCount).To(BeNumerically("==", 1))
			Expect(*metrics.SampleSum).To(BeNumerically(">", 0))
		})
		It("should handle errors correctly", func() {
			response := fixture.Test.ValidateV1(admission_test.V1ErrorPod())
			Expect(response).To(BeNil())
			metrics := validationMetrics(admRevVer, corev1_test.ErrorPod, facade.AdmissionError)
			Expect(*metrics.SampleCount).To(BeNumerically("==", 1))
			Expect(*metrics.SampleSum).To(BeNumerically(">", 0))
		})
		It("should skip unknown objects", func() {
			response := fixture.Test.ValidateV1(admission_test.V1BadPod())
			Expect(response).NotTo(BeNil())
			Expect(response.Allowed).To(BeTrue())
			metrics := validationMetrics(admRevVer, corev1_test.BadPod, facade.AdmissionAllowed)
			Expect(*metrics.SampleSum).To(BeNumerically("==", 0))
			Expect(*metrics.SampleCount).To(BeNumerically("==", 0))
		})
	})
	Context("mutation", func() {
		It("should not respond with patches if unchanged", func() {
			response := fixture.Test.MutateV1(admission_test.V1ValidPod())
			Expect(response).NotTo(BeNil())
			Expect(response.Allowed).To(BeTrue())
			Expect(response.Patch).To(BeNil())
			Expect(response.PatchType).To(BeNil())
			metrics := mutationMetrics(admRevVer, corev1_test.ValidPod, facade.AdmissionAllowed)
			Expect(*metrics.SampleCount).To(BeNumerically("==", 1))
			Expect(*metrics.SampleSum).To(BeNumerically(">", 0))
		})
		It("should respond with patches if mutated", func() {
			response := fixture.Test.MutateV1(admission_test.V1InvalidPod())
			Expect(response).NotTo(BeNil())
			Expect(response.Allowed).To(BeTrue())
			Expect(response.Patch).NotTo(BeNil())
			Expect(response.PatchType).NotTo(BeNil())
			metrics := mutationMetrics(admRevVer, corev1_test.InvalidPod, facade.AdmissionMutated)
			Expect(*metrics.SampleCount).To(BeNumerically("==", 1))
			Expect(*metrics.SampleSum).To(BeNumerically(">", 0))
		})
		It("should handle errors correctly", func() {
			response := fixture.Test.MutateV1(admission_test.V1ErrorPod())
			Expect(response).To(BeNil())
			metrics := mutationMetrics(admRevVer, corev1_test.ErrorPod, facade.AdmissionError)
			Expect(*metrics.SampleCount).To(BeNumerically("==", 1))
			Expect(*metrics.SampleSum).To(BeNumerically(">", 0))
		})
		It("should skip unknown objects", func() {
			response := fixture.Test.MutateV1(admission_test.V1BadPod())
			Expect(response).NotTo(BeNil())
			Expect(response.Allowed).To(BeTrue())
			metrics := mutationMetrics(admRevVer, corev1_test.BadPod, facade.AdmissionAllowed)
			Expect(*metrics.SampleSum).To(BeNumerically("==", 0))
			Expect(*metrics.SampleCount).To(BeNumerically("==", 0))
		})
	})
})
