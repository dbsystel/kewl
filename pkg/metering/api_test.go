package metering_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMetering(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "metering Suite")
}
