package httpext

import (
	"net/http"

	"github.com/dbsystel/kewl/pkg/metering"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const HeaderTraceID = "X-Trace-Id"

// Request is a decorator for http.Request
type Request struct {
	*http.Request
	logger logr.Logger
}

// ResponseWriter is a decorator for http.ResponseWriter
type ResponseWriter interface {
	http.ResponseWriter
	// SendResponse sends a response with status code handling internal errors
	SendResponse(statusCode int, body []byte)
	// SendResponseString sends a string response
	SendResponseString(statusCode int, body string)
	// SendJSON marshals the provided object to JSON
	SendJSON(statusCode int, obj interface{})
	// SendJSONBytes sends a marshaled json object
	SendJSONBytes(statusCode int, body []byte)
	// HandleInternalError will handle error occurred sending http.StatusInternalServerError, also logging the error
	HandleInternalError(err error)
	// Status yields the sent status response
	Status() int
}

// Metrics provides access the http metrics
type Metrics interface {
	// Requests provides the metering.Summary for the http requests
	Requests() metering.Summary
}

// MetricsRegistry extends metering.MetricsRegistry for accessing HTTP metrics
type MetricsRegistry interface {
	metering.MetricsRegistry
	// HTTP provides the Metrics for the http subsystem
	HTTP() Metrics
}

// Server is an extended http.Server
type Server struct {
	*http.Server
	*http.ServeMux
	logger   logr.Logger
	Metrics  MetricsRegistry
	listenFn func() error
}

// Handler is a function which handles the decorated ResponseWriter and Request
type Handler func(writer ResponseWriter, request *Request)

// NewServer creates a server from a logr.Logger and a Config
func NewServer(logger logr.Logger, config *Config) (*Server, error) {
	mux := http.NewServeMux()
	// Create the tls config
	tlsConfig, err := CreateTLSConfig(config.TLSConfig)
	if err != nil {
		return nil, errors.Wrap(err, "could not create tls configuration")
	}

	// Create the http server
	decorated := &http.Server{
		Addr:              config.Addr,
		Handler:           mux,
		TLSConfig:         tlsConfig,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		MaxHeaderBytes:    config.MaxHeaderBytes,
	}

	// Add prometheus Metrics and healthz
	promRegistry := prometheus.NewRegistry()
	mux.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{Registry: promRegistry}))
	mux.Handle("/healthz", &healthzHandler{logger})

	var listenFn func() error
	if config.TLSConfig.PrivateKeyFile != "" || config.TLSConfig.PublicKeyFile != "" {
		listenFn = func() error {
			return decorated.ListenAndServeTLS(config.TLSConfig.PublicKeyFile, config.TLSConfig.PrivateKeyFile)
		}
	} else {
		listenFn = func() error {
			return decorated.ListenAndServe()
		}
	}

	return &Server{
		Server:   decorated,
		ServeMux: mux,
		logger:   logger,
		Metrics:  &metricsRegistryImpl{MetricsRegistry: metering.NewRegistry(promRegistry), namespace: config.ServerNamespace},
		listenFn: listenFn,
	}, nil
}
