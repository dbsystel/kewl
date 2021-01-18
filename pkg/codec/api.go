package codec

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// SchemeExtension extends an existing scheme
type SchemeExtension interface {
	AddToScheme(scheme *runtime.Scheme) error
}

// SchemeRegistry manages a runtime.Scheme
type SchemeRegistry interface {
	Register(extension SchemeExtension) error
}

// Deserializer deserializes a known schema.GroupVersionKind from a byte slice
type Deserializer interface {
	SchemeRegistry
	Deserialize(groupVersionKind schema.GroupVersionKind, bytes []byte) (runtime.Object, error)
}

// NewDeserializer creates a new Deserializer for the provided scheme
func NewDeserializer(scheme *runtime.Scheme) Deserializer {
	return &deserializerImpl{
		schemeRegistryImpl: schemeRegistryImpl{scheme},
		decoder:            serializer.NewCodecFactory(scheme).UniversalDeserializer(),
	}
}
