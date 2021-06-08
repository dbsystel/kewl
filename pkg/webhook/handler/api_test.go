package handler_test

import (
	"encoding/json"
	"testing"

	"github.com/go-logr/logr"

	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/testing/admission_test"

	"github.com/dbsystel/kewl/pkg/webhook/facade"
	"github.com/dbsystel/kewl/pkg/webhook/handler"
	admission "k8s.io/api/admission/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func InvokeHandler(admRevHandler handler.AdmissionReview, review *admission_test.V1AdmissionReview) error {
	facadeV1 := facade.V1((*admission.AdmissionReview)(review))
	if err := admRevHandler.HandleReview(&logr.DiscardLogger{}, facadeV1); err != nil {
		return err
	}
	bytes, err := facadeV1.Marshal()
	panicutils.PanicIfError(err)
	panicutils.PanicIfError(json.Unmarshal(bytes, review))
	return nil
}

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "webhook.handler Suite")
}
