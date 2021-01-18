package httpext_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/dbsystel/kewl/testing/json_test"
)

type TestFacade struct {
	*http.ServeMux
}

type ResponseRecoverFacade struct {
	*httptest.ResponseRecorder
}

func (f *TestFacade) RequestBytes(method, path string, data []byte) *ResponseRecoverFacade {
	response, request := httptest.NewRecorder(), httptest.NewRequest(method, path, bytes.NewBuffer(data))
	f.ServeHTTP(response, request)
	return &ResponseRecoverFacade{response}
}

func (f *TestFacade) RequestJSON(method, path string, obj interface{}) *ResponseRecoverFacade {
	return f.RequestBytes(method, path, json_test.MarshalJSONOrPanic(obj))
}

func (f *TestFacade) Healthz() *ResponseRecoverFacade {
	return f.RequestBytes("GET", "/healthz", nil)
}

func (f *TestFacade) Metrics() string {
	return f.RequestBytes("GET", "/metrics", nil).String()
}

func (f *ResponseRecoverFacade) Bytes() []byte {
	return f.Body.Bytes()
}

func (f *ResponseRecoverFacade) String() string {
	return f.Body.String()
}

func (f *ResponseRecoverFacade) JSON(obj interface{}) interface{} {
	json_test.UnmarshalJSONOrPanic(f.Bytes(), obj)
	return obj
}
