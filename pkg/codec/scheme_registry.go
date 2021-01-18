package codec

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ SchemeRegistry = &schemeRegistryImpl{}

// schemeRegistryImpl is the implementation of SchemeRegistry
type schemeRegistryImpl struct {
	scheme *runtime.Scheme
}

func (s *schemeRegistryImpl) Register(extension SchemeExtension) error {
	if extension == nil {
		return errors.New("extension cannot be nil")
	}
	return extension.AddToScheme(s.scheme)
}
