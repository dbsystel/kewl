package testing

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Marshalable is an object which can be easily marshaled
type Marshalable interface {
	MustMarshal() []byte
}

// Reviewable is the interface which contains the minimum methods on an object which can be tested in an admission review
type Reviewable interface {
	Marshalable
	RawExtension() runtime.RawExtension
	GetName() string
	GetNamespace() string
	GetObjectKind() schema.ObjectKind
}
