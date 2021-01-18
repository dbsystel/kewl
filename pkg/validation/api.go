package validation

import (
	kewl "github.com/dbsystel/kewl/pkg"
	"github.com/dbsystel/kewl/pkg/codec"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Validator denotes a validator for any interface{}
type Validator interface {
	kewl.NamedObject
	codec.SchemeExtension
	// Validate validates the runtime.Object representation collecting the result using the ResultCollector or yielding an
	// error if the validator fails to run
	Validate(oldObject, newObject runtime.Object, results ResultCollector) error
}

// ResultCollector denotes a collector for failures which can be used descending an object hierarchy
type ResultCollector interface {
	// Failures returns the collected v1.StatusCause each representing a validation error
	Failures() []v1.StatusCause
	// AddFailure adds a message denoting the a failed validation the current field
	AddFailure(message string)
	// AppendField appends the provided suffix to the field and returns a new ResultCollector for the field
	AppendField(suffix string) ResultCollector
}

// Validate validates the provided object using the provided Validator
func Validate(newObject, oldObject runtime.Object, validators ...Validator) ([]v1.StatusCause, error) {
	resultCollector := NewResultCollector()
	for _, validator := range validators {
		err := invokeValidator(validator, newObject, oldObject, resultCollector)
		if err != nil {
			return nil, err
		}
	}
	return resultCollector.Failures(), nil
}

func NewResultCollector() ResultCollector {
	return &resultCollectorImpl{origin: "", messagesRef: &statusCauseSliceRef{}}
}
