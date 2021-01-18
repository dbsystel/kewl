package httpext

import (
	"net/http"

	"github.com/go-logr/logr"
)

// HealthzResponse is the response to be returned when requesting the health of the server
var HealthzResponse = []byte("he's not dead, jim")

var _ http.Handler = &healthzHandler{}

// healthzHandler is the handler for the handling health requests
type healthzHandler struct {
	logger logr.Logger
}

func (h *healthzHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-type", "text/plain")
	if _, err := writer.Write(HealthzResponse); err != nil {
		h.logger.Error(err, "could not write health response")
	}
	writer.WriteHeader(http.StatusOK)
}
