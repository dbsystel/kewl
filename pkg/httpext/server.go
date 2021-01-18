package httpext

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/dbsystel/kewl/pkg/panicutils"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/util/uuid"
)

// HandleExt adds a Handler for the provided path
func (s *Server) HandleExt(path string, handler Handler) {
	s.ServeMux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		// Start metering
		finishMetering := s.Metrics.HTTP().Requests().StartMetering(
			prometheus.Labels{labelMethod: request.Method, labelPath: path},
		)
		// Decorate writer and request
		decoratedWriter, decoratedRequest := s.decorate(writer, request)
		// And meter the request
		defer func() {
			finishMetering(prometheus.Labels{labelStatus: strconv.Itoa(decoratedWriter.Status())})
		}()
		// Make sure we handle panics
		defer panicutils.RecoverToErrorAndHandle(func(err error) {
			decoratedWriter.HandleInternalError(errors.Wrapf(err, "request handler panicked for path: %v", path))
		})
		// Delegate the handling
		handler(decoratedWriter, decoratedRequest)
	})
}

// RunWaitingForSig runs the server waiting for one of the provided os.Signal
func (s *Server) RunWaitingForSig(signals ...os.Signal) chan error {
	result := make(chan error)

	errChan := make(chan error)

	// Start the server in the background
	go func() {
		s.logger.Info("web-server starting")
		// On error => write the error directly to the channel
		if err := s.listenFn(); err != nil {
			errChan <- errors.Wrap(err, "could not start web-server")
		}
	}()

	// Make sure we stop the server in case of an error or an os.Signal
	signalChan := make(chan os.Signal, 1)
	go func() {
		select {
		// In case an error has occurred, push it to the result channel
		case err := <-errChan:
			result <- err
		// In case an os.Signal was received, stop the server and push the result to the result channel
		case sig := <-signalChan:
			s.logger.Info("web-server shutting down due to signal", "signal", sig)
			if err := s.Shutdown(context.Background()); err != nil {
				result <- errors.Wrap(err, "could not shutdown web-server")
			} else {
				result <- nil
			}
		}
	}()
	signal.Notify(signalChan, signals...)
	return result
}

// decorate creates the decorators for http.ResponseWriter and *http.Request
func (s *Server) decorate(writer http.ResponseWriter, request *http.Request) (ResponseWriter, *Request) {
	// Generate a trace id for the request and add it to the logger and response header for tracing errors
	traceID := string(uuid.NewUUID())
	logger := s.logger.WithValues(HeaderTraceID, traceID)
	writer.Header().Add(HeaderTraceID, traceID)
	return &responseWriterImpl{ResponseWriter: writer, traceID: traceID, logger: logger},
		&Request{Request: request, logger: logger}
}
