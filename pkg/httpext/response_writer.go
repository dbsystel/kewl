package httpext

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/pkg/errors"

	"github.com/go-logr/logr"
)

var _ ResponseWriter = &responseWriterImpl{}

// HeaderContentType denotes the header key to be used to send the content type
const HeaderContentType = "Content-type"

// responseWriterImpl decorates the http.ResponseWriter
type responseWriterImpl struct {
	http.ResponseWriter
	logger  logr.Logger
	status  int
	traceID string
}

func (d *responseWriterImpl) Status() int {
	return d.status
}

func (d *responseWriterImpl) WriteHeader(status int) {
	d.status = status
	d.ResponseWriter.WriteHeader(status)
}

func (d *responseWriterImpl) SendResponse(status int, body []byte) {
	d.WriteHeader(status)
	if len(body) > 0 {
		_, writeErr := d.Write(body)
		d.HandleInternalError(errors.Wrap(writeErr, "unable to write internal error"))
	}
}

func (d *responseWriterImpl) SendResponseString(status int, message string) {
	d.Header().Add(HeaderContentType, "text/plain")
	d.SendResponse(status, []byte(message))
}

func (d *responseWriterImpl) SendJSON(statusCode int, obj interface{}) {
	bytes, err := json.Marshal(obj)
	d.HandleInternalError(errors.Wrapf(err, "could not marshal response to json: %v", reflect.TypeOf(obj)))
	d.SendJSONBytes(statusCode, bytes)
}

func (d *responseWriterImpl) HandleInternalError(err error) {
	if err == nil {
		return
	}
	d.logger.Error(err, "Could not handle request")
	d.SendResponseString(
		http.StatusInternalServerError, fmt.Sprintf("Internal server error, %v: %v", HeaderTraceID, d.traceID),
	)
}

func (d *responseWriterImpl) SendJSONBytes(status int, bytes []byte) {
	d.Header().Add(HeaderContentType, "application/json")
	d.SendResponse(status, bytes)
}
