package handler

import (
	"encoding/json"

	"github.com/dbsystel/kewl/pkg/mutation"
	"github.com/dbsystel/kewl/pkg/webhook/facade"
	"github.com/go-logr/logr"
	"github.com/mattbaird/jsonpatch"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ Mutator = &mutatorImpl{}

type mutatorImpl struct {
	unmarshaller UnmarshalReqObj
	mutators     []mutation.Mutator
}

func (m *mutatorImpl) HandlerType() Type {
	return TypeMutation
}

func (m *mutatorImpl) HandleReview(logger logr.Logger, review facade.AdmissionReview) error {
	request := review.Request()
	logger.Info("mutation hook start", "resource", request.Resource())
	if err := m.unmarshaller.HandleReview(logger, review); err != nil {
		return errors.Wrap(err, "unable to handle request object")
	}

	object, oldObject := request.Object().Object, request.OldObject().Object

	if object == nil && oldObject == nil {
		logger.Info("mutation hook skipped, object type not registered", "resource", request.Resource())
		return nil
	}

	// Mutate the object and create the patch bytes
	patchBytes, err := m.mutateAndCreatePatch(logger, object, oldObject)
	if err != nil {
		return err
	}
	// Make sure we allow the admission
	review.Response().Allow()

	if len(patchBytes) > 0 {
		review.Response().PatchJSON(patchBytes)
	}

	logger.Info("mutation hook complete", "resource", request.Resource())
	return nil
}

func (m *mutatorImpl) AddMutator(mutator mutation.Mutator) error {
	if mutator == nil {
		return nil
	}
	if err := m.unmarshaller.Register(mutator); err != nil {
		return errors.Wrapf(err, "could not register mutator: %v", mutator.Name())
	}
	m.mutators = append(m.mutators, mutator)
	return nil
}

func (m *mutatorImpl) mutateAndCreatePatch(logger logr.Logger, originalObj, oldObject runtime.Object) ([]byte, error) {
	mutated, err := mutation.Mutate(originalObj, oldObject, m.mutators...)
	if mutated == nil || err != nil {
		return nil, err
	}
	logger.Info("resource has been mutated, creating patch", "kind", originalObj.GetObjectKind())
	return m.createPatch(originalObj, mutated)
}

func (m *mutatorImpl) createPatch(original, mutated runtime.Object) ([]byte, error) {
	originalJSON, err := json.Marshal(original)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal original object")
	}
	mutatedJSON, err := json.Marshal(mutated)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal mutated object")
	}
	patch, err := jsonpatch.CreatePatch(originalJSON, mutatedJSON)
	if err != nil {
		return nil, errors.Wrap(err, "could not create merge patch for mutation")
	}
	if len(patch) == 0 {
		return nil, nil
	}
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshall patch data")
	}
	return patchBytes, nil
}
