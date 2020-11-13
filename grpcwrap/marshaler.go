package grpcwrap

import (
	"fmt"

	"github.com/DataWorkbench/glog"
	"github.com/Yu-33/gohelper/conv"
	"github.com/golang/protobuf/jsonpb"
	deproto "github.com/golang/protobuf/proto" // the package will be Deprecated
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	protojsonMarshal = protojson.MarshalOptions{}
	jsonpbMarshaler  = jsonpb.Marshaler{OrigName: true, EmitDefaults: true}
)

func pbMsgToString(logger *glog.Logger, i interface{}) string {
	// priority of use protojson
	if p, ok := i.(proto.Message); ok {
		b, err := protojsonMarshal.Marshal(p)
		if err != nil {
			logger.Error().Error("marshal proto.Message v2 error", err).Fire()
			return fmt.Sprintf("%+v", i)
		}
		return conv.BytesToString(b)
	}

	// compatible with old version protobuf
	if p, ok := i.(deproto.Message); ok {
		s, err := jsonpbMarshaler.MarshalToString(p)
		if err != nil {
			logger.Error().Error("marshal proto.Message v1 error", err).Fire()
			s = p.String()
		}
		return s
	}
	return fmt.Sprintf("%+v", i)
}
