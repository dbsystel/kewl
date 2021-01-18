package mutation

import (
	"reflect"

	"github.com/dbsystel/kewl/pkg/panicutils"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func invokeMutator(mutator Mutator, newObject, oldObject runtime.Object) (result error) {
	if mutator == nil {
		return nil
	}
	// Make sure we handle panics, since we really do not want to break the whole mutation
	defer panicutils.RecoverToErrorAndHandle(func(err error) {
		result = errors.Wrapf(err, "panic in mutator: %v", mutator.Name())
	})
	// Delegate the mutation
	if err := mutator.Mutate(newObject, oldObject); err != nil {
		return errors.Wrapf(err, "mutator '%v' failed to handle: %v", mutator.Name(), reflect.TypeOf(newObject))
	}

	return nil
}
