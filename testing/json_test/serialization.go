package json_test

import (
	"encoding/json"

	"github.com/dbsystel/kewl/pkg/panicutils"
)

func MarshalJSONOrPanic(obj interface{}) []byte {
	result, err := json.Marshal(obj)
	panicutils.PanicIfError(err)
	return result
}

func UnmarshalJSONOrPanic(src []byte, obj interface{}) {
	if len(src) == 0 {
		return
	}
	panicutils.PanicIfError(json.Unmarshal(src, obj))
}
