package ginmiddle

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
)

// ParseHandlerFuncName to parse the func name from handler func .
func ParseHandlerFuncName(i interface{}) string {
	funcPath := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	fields := strings.Split(funcPath, "/")
	funcName := strings.Split(fields[len(fields)-1], ".")[1]
	return funcName
}

// ParseAPIPermType to parse the operation permission type from http Method.
func ParseAPIPermType(method string) pbmodel.ProjectAPI_PermType {
	var opType pbmodel.ProjectAPI_PermType
	switch method {
	case http.MethodGet, http.MethodHead:
		opType = pbmodel.ProjectAPI_Read
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		opType = pbmodel.ProjectAPI_Write
	default:
		panic("unsupported api type")
	}
	return opType
}
