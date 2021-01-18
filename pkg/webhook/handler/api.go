package handler

import (
	"github.com/dbsystel/kewl/pkg/codec"
	"github.com/dbsystel/kewl/pkg/mutation"
	"github.com/dbsystel/kewl/pkg/validation"
	"github.com/dbsystel/kewl/pkg/webhook/facade"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

// Type denotes a type name for the handler
type Type string

const (
	// TypeValidation marks a handler for validation
	TypeValidation Type = "validation"
	// TypeMutation marks a handler mutation
	TypeMutation Type = "mutation"
	// TypeOther marks a handler which does neither (like unmarshalling)
	TypeOther Type = "other"
)

// AdmissionReview is a handler for v1beta1.AdmissionReview
type AdmissionReview interface {
	// HandleReview handles the facade.AdmissionReview using the logr.Logger
	HandleReview(logger logr.Logger, review facade.AdmissionReview) error
	// Type returns the Type to identify the web hook
	HandlerType() Type
}

// UnmarshalReqObj is an AdmissionReview which handles the deserialization of the object in v1beta1.AdmissionRequest
type UnmarshalReqObj interface {
	codec.SchemeRegistry
	AdmissionReview
}

// Validator is an AdmissionReview for validating an object in v1beta1.AdmissionRequest
type Validator interface {
	AdmissionReview
	// AddValidator adds a validation.Validator to be used to validate the v1beta1.AdmissionRequest
	AddValidator(validator validation.Validator) error
}

// Mutator is a AdmissionReview for mutating an object in v1beta1.AdmissionRequest
type Mutator interface {
	AdmissionReview
	// AddMutator adds a mutation.Mutator to be used to mutate the v1beta1.AdmissionRequest
	AddMutator(mutator mutation.Mutator) error
}

// NewUnmarshalReqObj creates a new UnmarshalReqObj using a new runtime.Scheme
func NewUnmarshalReqObj() UnmarshalReqObj {
	return &unmarshalReqObjImpl{codec.NewDeserializer(runtime.NewScheme())}
}

// NewValidator creates a new Validator
func NewValidator() Validator {
	return &validatorImpl{unmarshaller: NewUnmarshalReqObj()}
}

// NewMutator creates a new Mutator and the provided metering.Summary for metering
func NewMutator() Mutator {
	return &mutatorImpl{unmarshaller: NewUnmarshalReqObj()}
}
