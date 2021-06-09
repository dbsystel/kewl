package webhook

import (
	"net/http"
	"strings"

	"github.com/dbsystel/kewl/pkg/webhook/facade"

	"github.com/dbsystel/kewl/pkg/mutation"

	"github.com/dbsystel/kewl/pkg/metering"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dbsystel/kewl/pkg/httpext"
	"github.com/dbsystel/kewl/pkg/validation"
	"github.com/dbsystel/kewl/pkg/webhook/handler"
	"github.com/pkg/errors"
)

// AddValidator adds a validation.Validator to the Server
func (s *Server) AddValidator(validators ...validation.Validator) error {
	// Create a validation handler, if none is set
	if s.validator == nil {
		s.validator = handler.NewValidator()
		s.HandleAdmissionReview("/validate", s.validator, s.Metrics.WebHook().Validation())
	}
	// Register the validators at the handler
	for _, validator := range validators {
		if err := s.validator.AddValidator(validator); err != nil {
			return err
		}
	}
	return nil
}

// AddMutator adds a mutation.Mutator to the Server
func (s *Server) AddMutator(mutators ...mutation.Mutator) error {
	// Create a mutation handler, if non is set
	if s.mutator == nil {
		s.mutator = handler.NewMutator()
		s.HandleAdmissionReview("/mutate", s.mutator, s.Metrics.WebHook().Mutation())
	}
	// Register the mutators at the handler
	for _, mutator := range mutators {
		if err := s.mutator.AddMutator(mutator); err != nil {
			return err
		}
	}
	return nil
}

// HandleAdmissionReview adds a handler.AdmissionReview
func (s *Server) HandleAdmissionReview(path string, admRevHandler handler.AdmissionReview, summary metering.Summary) {
	// Add a handler function
	s.HandleExt(path, func(writer httpext.ResponseWriter, request *httpext.Request) {
		logger := request.Logger().V(1).WithValues("path", path, "handler", admRevHandler.HandlerType())
		logger.Info("reading admission review")
		review := s.readAdmissionReview(writer, request)
		if review == nil {
			logger.Info("no review provided")
			return
		}

		// Ensure we meter a handled request
		finishMetering := s.startMetering(summary, review, writer)
		defer finishMetering()

		logger.Info("handling review")
		// Delegate the review to the admRevHandler
		if err := admRevHandler.HandleReview(request.Logger(), review); err != nil {
			writer.HandleInternalError(errors.Wrap(err, "review handler failed"))
			return
		}

		logger.Info("removing request part of review and sending response")

		// Clear out the request
		review.ClearRequest()
		// Send the response
		bytes, err := review.Marshal()
		if err != nil {
			logger.Error(err, "response could not be marshalled")
			writer.HandleInternalError(errors.Wrap(err, "response review cannot be marshaled"))
			return
		}
		logger.Info("sending response", "size", len(bytes))
		writer.SendJSONBytes(http.StatusOK, bytes)
	})
}

// readAdmissionReview will read and facade the admission review from the request, handling errors
// will return nil, if an error occurred
func (s *Server) readAdmissionReview(writer httpext.ResponseWriter, request *httpext.Request) facade.AdmissionReview {
	// Read the admission review and create a facade
	review, err := s.facadeAdmissionReviewFrom(request)
	// Handle errors when reading
	if err != nil {
		// This is a user error (invalid object representation)
		if strings.HasPrefix(err.Error(), facade.InvalidAdmissionReviewMsg) {
			writer.SendResponseString(http.StatusUnprocessableEntity, err.Error())
			return nil
		}
		// This is a server error
		writer.HandleInternalError(err)
		return nil
	}
	// Handle empty bodies
	if review == nil {
		writer.SendResponseString(http.StatusUnprocessableEntity, "no content provided")
		return nil
	}
	// Handle empty reviews
	if review.Request() == nil {
		writer.SendResponseString(http.StatusUnprocessableEntity, "no request in review provided")
		return nil
	}
	return review
}

// facadeAdmissionReviewFrom reads a facade.AdmissionReview from the request, will return nil, if body is empty
func (s *Server) facadeAdmissionReviewFrom(request *httpext.Request) (facade.AdmissionReview, error) {
	// Try to get a facade.AdmissionReview from the request
	body, err := request.ReadBody()
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	return facade.AdmissionReviewFrom(body)
}

// startMetering starts the metering for the facade.AdmissionReview and httpext.ResponseWriter using metering.Summary
func (s *Server) startMetering(tgt metering.Summary, review facade.AdmissionReview, wr httpext.ResponseWriter) func() {
	kind := review.Request().Kind()
	decoratedFn := tgt.StartMetering(prometheus.Labels{
		labelAdmissionReviewVersion: review.Version(),
		labelObjGroup:               kind.Group, labelObjVersion: kind.Version, labelObjKind: kind.Kind,
		labelObjNamespace: review.Request().Namespace(),
	})
	return func() {
		status := wr.Status()
		// Handle client errors
		if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			decoratedFn(prometheus.Labels{labelResult: string(facade.AdmissionClientError)})
			return
		}
		// Handle server errors
		if status >= http.StatusInternalServerError {
			decoratedFn(prometheus.Labels{labelResult: string(facade.AdmissionError)})
			return
		}
		// No object in the request means we do not handle it, so we do not track it
		if review.Request().Object().Object == nil && review.Request().OldObject().Object == nil {
			return
		}
		// Use the response type from the facade
		decoratedFn(prometheus.Labels{labelResult: string(review.Response().ResponseType())})
	}
}
