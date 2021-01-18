package handler

import (
	"github.com/dbsystel/kewl/pkg/codec"
	"github.com/dbsystel/kewl/pkg/webhook/facade"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ UnmarshalReqObj = &unmarshalReqObjImpl{}

type unmarshalReqObjImpl struct {
	codec.Deserializer
}

func (u *unmarshalReqObjImpl) HandlerType() Type {
	return TypeOther
}

func (u *unmarshalReqObjImpl) HandleReview(_ logr.Logger, review facade.AdmissionReview) error {
	request := review.Request()
	kind := request.Kind()
	schemaKind := schema.GroupVersionKind{Group: kind.Group, Version: kind.Version, Kind: kind.Kind}
	if err := u.deserializeRawExtension(schemaKind, request.Object()); err != nil {
		// If we don't know the object, allow the request, it's not their fault :)
		if runtime.IsNotRegisteredError(err) {
			review.Response().Allow()
			return nil
		}
		return errors.Wrapf(err, "could not deserialize request object: %v", kind)
	}
	if err := u.deserializeRawExtension(schemaKind, request.OldObject()); err != nil {
		return errors.Wrapf(err, "could not deserialize request old object: %v", kind)
	}
	return nil
}

func (u *unmarshalReqObjImpl) deserializeRawExtension(gvk schema.GroupVersionKind, ext *runtime.RawExtension) error {
	if len(ext.Raw) == 0 {
		return nil
	}
	deserialized, err := u.Deserialize(gvk, ext.Raw)
	if err != nil {
		return err
	}
	ext.Object = deserialized
	return nil
}
