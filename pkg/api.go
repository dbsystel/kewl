package pkg

const Namespace = "webhook"

// NamedObject is an interface for all objects providing a name
type NamedObject interface {
	// Name denotes the name of the validator
	Name() string
}
