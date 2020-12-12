package grpcwrap

import (
	"fmt"
	"unsafe"

	"github.com/DataWorkbench/glog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	protojsonMarshal = protojson.MarshalOptions{EmitUnpopulated: true}
)

func pbMsgToString(logger *glog.Logger, i interface{}) string {
	if p, ok := i.(proto.Message); ok {
		b, err := protojsonMarshal.Marshal(p)
		if err == nil {
			return *(*string)(unsafe.Pointer(&b))
		}
		logger.Error().Error("marshal proto.Message error", err).Fire()
	}

	if p, ok := i.(fmt.Stringer); ok {
		return p.String()
	}

	return fmt.Sprintf("%+v", i)
}
