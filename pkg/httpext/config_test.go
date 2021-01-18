package httpext_test

import (
	"flag"
	"os"

	"github.com/dbsystel/kewl/testing/httpext_test"

	"github.com/dbsystel/kewl/pkg/panicutils"

	"github.com/dbsystel/kewl/pkg/httpext"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	It("should parse the default values", func() {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config := httpext.NewDefaultConfig()
		flagSet := flag.NewFlagSet("test", flag.PanicOnError)
		Expect(flagSet.Parse(nil)).NotTo(HaveOccurred())
		httpext.AddFlagsToFlagSet(&config, flagSet)
		Expect(config).To(BeEquivalentTo(httpext.NewDefaultConfig()))
	})
})

var _ = Describe("TLSConfig", func() {
	var pemFilePair httpext_test.PemFilePair
	BeforeEach(func() {
		pemFilePair = httpext_test.MustNewTmpPemFilePair()
		panicutils.PanicIfError(pemFilePair.Generate(httpext_test.CACertificate))
	})
	AfterEach(func() {
		if pemFilePair != nil {
			panicutils.PanicIfError(pemFilePair.Delete())
		}
	})
	It("should error in case the min version is unsupported", func() {
		_, err := httpext.CreateTLSConfig(httpext.TLSConfig{
			MinVersion:  "1.1",
			CaCertsFile: pemFilePair.PublicKeyPath(),
		})
		Expect(err).To(HaveOccurred())
	})
	It("should error in case the pem file does not exist", func() {
		_, err := httpext.CreateTLSConfig(httpext.TLSConfig{
			MinVersion:  "1.2",
			CaCertsFile: "bla",
		})
		Expect(err).To(HaveOccurred())
	})
	It("should error in case the pem file is invalid", func() {
		_, err := httpext.CreateTLSConfig(httpext.TLSConfig{
			MinVersion:  "1.2",
			CaCertsFile: "config_test.go",
		})
		Expect(err).To(HaveOccurred())
	})
	It("should create a valid TLS config from a CA file", func() {
		Expect(httpext.CreateTLSConfig(httpext.TLSConfig{
			MinVersion:  "1.2",
			CaCertsFile: pemFilePair.PublicKeyPath(),
		})).To(Not(BeNil()))
	})
})
