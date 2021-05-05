package handler

import (
	"net/http"

	"github.com/dbsystel/kewl/pkg/webhook/facade"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/dbsystel/kewl/pkg/validation"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ Validator = &validatorImpl{}

type validatorImpl struct {
	unmarshaller UnmarshalReqObj
	validators   []validation.Validator
}

func (v *validatorImpl) HandlerType() Type {
	return TypeValidation
}

func (v *validatorImpl) HandleReview(logger logr.Logger, review facade.AdmissionReview) error {
	request := review.Request()
	loggerV1 := logger.V(1)
	loggerV1.Info("validation hook start", "resource", request.Resource())
	if err := v.unmarshaller.HandleReview(logger, review); err != nil {
		return errors.Wrap(err, "unable to handle request object")
	}

	object, oldObject := request.Object().Object, request.OldObject().Object

	if object == nil {
		logger.Info("validation hook skipped, object type not registered", "resource", request.Resource())
		return nil
	}

	// Issue the validation for the request object
	failures, err := validation.Validate(object, oldObject, v.validators...)
	if err != nil {
		return err
	}
	// Handle valid status
	if len(failures) == 0 {
		review.Response().Allow()
		return nil
	}
	// Convert the fail into a status
	status := v.convertFailuresToStatus(object, failures)
	review.Response().Deny(status)
	loggerV1.Info("validation hook complete", "resource", request.Resource())
	return nil
}

func (v *validatorImpl) AddValidator(validator validation.Validator) error {
	if validator == nil {
		return nil
	}
	err := v.unmarshaller.Register(validator)
	if err != nil {
		return errors.Wrapf(err, "could not add validator: %v", validator.Name())
	}
	v.validators = append(v.validators, validator)
	return nil
}

func (v *validatorImpl) convertFailuresToStatus(obj runtime.Object, failures []metav1.StatusCause) *metav1.Status {
	result := &metav1.Status{}
	// No failures => that's perfectly alright
	if len(failures) == 0 {
		return result
	}
	result.Code = http.StatusUnprocessableEntity
	result.Reason = "invalid"
	// Add details
	result.Details = &metav1.StatusDetails{Causes: failures}
	if objKind := obj.GetObjectKind(); objKind != nil {
		kind := objKind.GroupVersionKind()
		result.Details.Group = kind.Group
		result.Details.Kind = kind.Kind
	}
	if meta, ok := obj.(metav1.Object); ok {
		result.Details.Name = meta.GetName()
	}
	return result
}
