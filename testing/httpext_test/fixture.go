package httpext_test

import (
	"github.com/dbsystel/kewl/pkg/httpext"
	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/go-logr/logr/testing"
)

// Fixture is a fixture for testing httpext.Server
type Fixture struct {
	// Sut is the system under test
	Sut *httpext.Server
	// Test is the TestFacade for testing requests
	Test *TestFacade
}

// NewFixture creates a new Fixture for testing httpext.Server
func NewFixture() *Fixture {
	return newFixture(false)
}

// NewFixtureHTTPS creates a new Fixture for testing httpext.Server using HTTPS
func NewFixtureHTTP() *Fixture {
	return newFixture(true)
}

// newFixture creates a new Fixture with either HTTPS activated or deactivated
func newFixture(enableHTTPS bool) *Fixture {
	config := TestConfig()
	if enableHTTPS {
		pemFilePair := MustNewTmpPemFilePair()
		panicutils.PanicIfError(pemFilePair.Generate(CACertificate))
		config.TLSConfig.MinVersion = "1.2"
		config.TLSConfig.PublicKeyFile = pemFilePair.PublicKeyPath()
		config.TLSConfig.PrivateKeyFile = pemFilePair.PrivateKeyPath()
	}
	server, err := httpext.NewServer(testing.NullLogger{}, config)
	panicutils.PanicIfError(err)
	return &Fixture{Sut: server, Test: &TestFacade{ServeMux: server.ServeMux}}
}

// TestConfig creates a new httpext.Config for testing
func TestConfig() *httpext.Config {
	config := httpext.NewDefaultConfig()
	config.TLSConfig = httpext.TLSConfig{}
	return &config
}
