package webhook_test

import (
	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/pkg/webhook"
	"github.com/dbsystel/kewl/testing/httpext_test"
	"github.com/go-logr/logr"
)

type Fixture struct {
	Sut  *webhook.Server
	Test *ReviewFacade
}

func NewFixture() *Fixture {
	server, err := webhook.NewServer(&logr.DiscardLogger{}, httpext_test.TestConfig())
	panicutils.PanicIfError(err)
	return &Fixture{Sut: server, Test: &ReviewFacade{&httpext_test.TestFacade{ServeMux: server.ServeMux}}}
}
