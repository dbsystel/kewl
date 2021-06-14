package integtest

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
)

// AdmissionResponse is a decorator for admissionv1.AdmissionResponse
type AdmissionResponse struct {
	admissionv1.AdmissionResponse
}

// PatchMaps returns the patches as maps
func (a *AdmissionResponse) PatchMaps() (result []map[string]string, err error) {
	if a.PatchType == nil || len(a.Patch) == 0 {
		return nil, nil
	}
	err = json.Unmarshal(a.Patch, &result)
	return result, err
}
