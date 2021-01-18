package mutation

import (
	"reflect"

	kewl "github.com/dbsystel/kewl/pkg"
	"github.com/dbsystel/kewl/pkg/codec"
	"k8s.io/apimachinery/pkg/runtime"
)

// Mutator is the interface which is to be used to mutate a k8s runtime.Object
type Mutator interface {
	kewl.NamedObject
	codec.SchemeExtension
	// Mutate mutates the new object provided having access on the old object as well
	Mutate(newObject, oldObject runtime.Object) error
}

// Mutate creates a copy of the provided runtime.Object and runs the provided Mutator on it
// will return the changed copy if changes have been made or nil otherwise
func Mutate(newObject, oldObject runtime.Object, mutators ...Mutator) (runtime.Object, error) {
	newObjectCopy := newObject.DeepCopyObject()
	for _, mutator := range mutators {
		err := invokeMutator(mutator, newObjectCopy, oldObject)
		if err != nil {
			return nil, err
		}
	}
	if reflect.DeepEqual(newObject, newObjectCopy) {
		return nil, nil
	}
	return newObjectCopy, nil
}
