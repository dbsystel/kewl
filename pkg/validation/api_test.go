package validation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttpExt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "validation Suite")
}
