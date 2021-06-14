package integtest

import (
	"bytes"
	"encoding/json"
	"github.com/dbsystel/kewl/pkg/httpext"
	"github.com/dbsystel/kewl/pkg/webhook"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"io/ioutil"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/uuid"
	"net/http"
	"net/http/httptest"
	"reflect"
)

var _ Interface = &fixture{}

type fixture struct {
	sut              *webhook.Server
	scheme           *runtime.Scheme
	path             string
	defaultNamespace string
}

func newIntegrationWithPath(defaultNamespace, path string) (*fixture, error) {
	cfg := &httpext.Config{}
	server, err := webhook.NewServer(logr.DiscardLogger{}, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not create servers")
	}
	return &fixture{sut: server, path: path, scheme: runtime.NewScheme(), defaultNamespace: defaultNamespace}, nil
}

func (i *fixture) Invoke(obj, oldObj runtime.Object) (*AdmissionResponse, error) {
	if obj == nil && oldObj == nil {
		return nil, errors.New("at least one object must be provided")
	}
	operation := i.operation(obj, oldObj)
	var (
		err              error
		marshalledObj    []byte
		marshalledOldObj []byte
	)
	if obj != nil {
		if marshalledObj, err = json.Marshal(obj); err != nil {
			return nil, errors.Wrapf(err, "could not marshal obj of type: %v", reflect.TypeOf(obj))
		}
	}
	if oldObj != nil {
		if marshalledOldObj, err = json.Marshal(oldObj); err != nil {
			return nil, errors.Wrapf(err, "could not marshal obj of type: %v", reflect.TypeOf(oldObj))
		}
	}
	review := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{Kind: "AdmissionReview", APIVersion: admissionv1.SchemeGroupVersion.Identifier()},
		Request: &admissionv1.AdmissionRequest{
			UID:       uuid.NewUUID(),
			Kind:      i.kind(obj, oldObj),
			Name:      "test-review",
			Namespace: i.namespace(obj, oldObj),
			Operation: operation,
			Object:    runtime.RawExtension{Raw: marshalledObj},
			OldObject: runtime.RawExtension{Raw: marshalledOldObj},
		},
	}
	recorder := httptest.NewRecorder()
	body, err := json.Marshal(review)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal the review")
	}
	request := httptest.NewRequest("POST", i.path, bytes.NewBuffer(body))
	i.sut.Server.Server.Handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		return nil, errors.Errorf("review invocation response code invalid: %v", recorder.Code)
	}
	body, err = ioutil.ReadAll(recorder.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response body")
	}
	review = &admissionv1.AdmissionReview{}
	if err = json.Unmarshal(body, &review); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal response body")
	}
	if review.Request != nil {
		return nil, errors.Wrap(err, "review should not have the request set")
	}
	var result *AdmissionResponse
	if review.Response != nil {
		result = &AdmissionResponse{*review.Response}
	}
	return result, nil
}

func (i *fixture) InvokeFromFile(objPath, oldObjPath string) (*AdmissionResponse, error) {
	obj, err := i.load(objPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not load object from path: %v", objPath)
	}
	oldObj, err := i.load(oldObjPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not load old object from path: %v", oldObjPath)
	}
	return i.Invoke(obj, oldObj)
}

func (i *fixture) load(path string) (runtime.Object, error) {
	if len(path) == 0 {
		return nil, nil
	}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read from path: %v", path)
	}
	result, _, err := serializer.NewCodecFactory(i.scheme).UniversalDeserializer().Decode(contents, nil, nil)
	return result, errors.Wrapf(err, "could not decode object from path: %v", path)
}

func (i *fixture) operation(obj, oldObj runtime.Object) admissionv1.Operation {
	if obj == nil {
		if oldObj == nil {
			return ""
		}
		return admissionv1.Delete
	}
	if oldObj != nil {
		return admissionv1.Update
	}
	return admissionv1.Create
}

func (i *fixture) kind(objs ...runtime.Object) (result metav1.GroupVersionKind) {
	for _, candidate := range objs {
		if candidate != nil {
			result.Group, result.Version, result.Kind = i.gvk(candidate)
			break
		}
	}
	return result
}

func (i *fixture) gvk(obj runtime.Object) (string, string, string) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	return gvk.Group, gvk.Version, gvk.Kind
}

func (i *fixture) namespace(objs ...runtime.Object) string {
	for _, obj := range objs {
		metaAcc, ok := obj.(metav1.ObjectMetaAccessor)
		if !ok {
			continue
		}
		objMeta := metaAcc.GetObjectMeta()
		if objMeta == nil {
			continue
		}
		return objMeta.GetNamespace()
	}
	return i.defaultNamespace
}
