package codec

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ Deserializer = &deserializerImpl{}

type deserializerImpl struct {
	schemeRegistryImpl
	decoder runtime.Decoder
}

func (d *deserializerImpl) Deserialize(groupVersionKind schema.GroupVersionKind, bytes []byte) (runtime.Object, error) {
	// Create a new object of the provided kind using the scheme
	obj, err := d.scheme.New(groupVersionKind)
	if err != nil {
		if runtime.IsNotRegisteredError(err) {
			return nil, err
		}
		return nil, errors.Wrapf(err, "could not create an object for: %v", groupVersionKind)
	}
	// Decode the object using the provided decoder
	obj, _, err = d.decoder.Decode(bytes, &groupVersionKind, obj)
	if err != nil {
		return nil, errors.Wrapf(err, "unable decoded object as: %v", groupVersionKind)
	}
	return obj, nil
}
