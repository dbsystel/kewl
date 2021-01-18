package webhook

import (
	"github.com/dbsystel/kewl/pkg/httpext"
	"github.com/dbsystel/kewl/pkg/metering"
	"github.com/dbsystel/kewl/pkg/webhook/handler"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

// MetricsRegistry is an extension to httpext.MetricsRegistry
type MetricsRegistry interface {
	httpext.MetricsRegistry
	WebHook() Metrics
}

// Metrics provided metrics about the webhooks
type Metrics interface {
	// Validation provides the metering.Summary for validation web hooks
	Validation() metering.Summary
	// Mutation provides the metering.Summary for mutation web hooks
	Mutation() metering.Summary
}

// Server is the webhook server
type Server struct {
	*httpext.Server // Inherit from our extended httpext.Server
	// Metrics provides the MetricsRegistry to access metrics
	Metrics MetricsRegistry
	// logger is the logger to be used for this server
	logger logr.Logger
	// validator is the handler.Validator to delegate the validation to (if any validation.Validator was set)
	validator handler.Validator
	// mutator is the handler.Mutator to delegate the mutation to (if any mutation.Mutator was set)
	mutator handler.Mutator
}

// NewServer creates a new webhook Server based on the provided logr.Logger, Config
func NewServer(logger logr.Logger, config *httpext.Config) (*Server, error) {
	config.ServerNamespace = "webhook"
	decorated, err := httpext.NewServer(logger, config)
	if err != nil {
		return nil, errors.Wrap(err, "could not create http server")
	}
	return &Server{
		Server: decorated, logger: logger,
		Metrics: &metricsRegistryImpl{MetricsRegistry: decorated.Metrics, namespace: config.ServerNamespace},
	}, nil
}
