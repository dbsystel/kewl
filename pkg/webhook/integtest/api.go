package integtest

import (
	"reflect"

	"github.com/dbsystel/kewl/pkg/mutation"
	"github.com/dbsystel/kewl/pkg/validation"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// Interface is the interface for an integration test fixture
type Interface interface {
	Invoke(obj, oldObj runtime.Object) (*AdmissionResponse, error)
	InvokeFromFile(objPath, oldObjPath string) (*AdmissionResponse, error)
}

// NewMutation creates a new Interface for mutation testing
func NewMutation(defaultNamespace string, mutators ...mutation.Mutator) (Interface, error) {
	result, err := newIntegrationWithPath(defaultNamespace, "/mutate")
	if err != nil {
		return nil, err
	}
	for _, mutator := range mutators {
		if err = result.sut.AddMutator(mutator); err != nil {
			return nil, errors.Wrapf(err, "could not add mutator: %v", reflect.TypeOf(mutator))
		}
		// Note: error is handled in server already
		_ = mutator.AddToScheme(result.scheme)
	}
	return result, nil
}

// NewValidation creates a new Interface for validation testing
func NewValidation(defaultNamespace string, validators ...validation.Validator) (Interface, error) {
	result, err := newIntegrationWithPath(defaultNamespace, "/validate")
	if err != nil {
		return nil, err
	}
	for _, validator := range validators {
		if err = result.sut.AddValidator(validator); err != nil {
			return nil, errors.Wrapf(err, "could not add validator: %v", reflect.TypeOf(validators))
		}
		// Note: error is handled in server already
		_ = validator.AddToScheme(result.scheme)
	}
	return result, err
}
