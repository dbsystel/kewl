package validation

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/dbsystel/kewl/pkg/panicutils"

	"github.com/pkg/errors"
)

// invokeValidator will safely invoke the validator, handling panics
func invokeValidator(validator Validator, newObject, oldObject runtime.Object, results ResultCollector) (result error) {
	if validator == nil {
		return nil
	}
	// Make sure we handle panics, since we really do not want to break the whole validation
	defer panicutils.RecoverToErrorAndHandle(func(err error) {
		result = errors.Wrapf(err, "panic in validator: %v", validator.Name())
	})
	// Delegate the validation
	if err := validator.Validate(newObject, oldObject, results); err != nil {
		return errors.Wrapf(err, "validator '%v' failed to handle: %v", validator.Name(), reflect.TypeOf(newObject))
	}

	return nil
}
