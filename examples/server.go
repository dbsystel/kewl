package examples

import (
	"github.com/dbsystel/kewl/pkg/httpext"
	"github.com/dbsystel/kewl/pkg/panicutils"
	"github.com/dbsystel/kewl/pkg/webhook"
	"flag"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"syscall"
)

func main() {
	// Create new webhook server
	srv, err := createWebHookServer()
	if err != nil {
		panic(err)
	}

	// Register our validator
	panicutils.PanicIfError(errors.Wrap(srv.AddValidator(PodValidator), "could not add mutator"))

	// Register our mutator
	panicutils.PanicIfError(errors.Wrap(srv.AddMutator(PodMutator), "could not add mutator"))

	// Let this run until we receive SIG_INT or SIG_TERM and handle errors
	if err := <-srv.RunWaitingForSig(syscall.SIGINT, syscall.SIGTERM); err != nil {
		panic(err)
	}
}

func createWebHookServer() (*webhook.Server, error) {
	// Use whatever logr.Logger implementation you fancy
	logger := createMainLogger()
	// Create a new default configuration and use flag to parse it from the command lines
	config := httpext.NewDefaultConfig()
	httpext.AddCommandLineFlags(&config)
	flag.Parse()

	// Create the server itself
	srv, err := webhook.NewServer(logger.WithName("webhook"), &config)
	if err != nil {
		return nil, errors.Wrap(err, "could not create web-server")
	}

	return srv, nil
}

func createMainLogger() logr.Logger {
	zapIt, err := zap.NewProduction()
	if err != nil {
		panic(errors.Wrap(err, "could not create main logger"))
	}
	return zapr.NewLogger(zapIt)
}
