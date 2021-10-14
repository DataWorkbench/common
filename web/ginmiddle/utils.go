package ginmiddle

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/DataWorkbench/gproto/pkg/model"
)

// ParseOpName to parse the operation name from func name.
func ParseOpName(i interface{}) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	fields := strings.Split(funcName, "/")
	opName := strings.Split(fields[len(fields)-1], ".")[1]
	return opName
}

// ParseOpType to parse the operation type from http Method.
func ParseOpType(method string) model.Operation_Type {
	var opType model.Operation_Type
	switch method {
	case http.MethodGet, http.MethodHead:
		opType = model.Operation_Read
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		opType = model.Operation_Write
	default:
		panic("unsupported operation type")
	}
	return opType
}
