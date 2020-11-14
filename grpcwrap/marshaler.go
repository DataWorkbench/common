package grpcwrap

import (
	"fmt"

	"github.com/DataWorkbench/glog"
	"github.com/Yu-33/gohelper/conv"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	protojsonMarshal = protojson.MarshalOptions{}
)

func pbMsgToString(logger *glog.Logger, i interface{}) string {
	if p, ok := i.(proto.Message); ok {
		b, err := protojsonMarshal.Marshal(p)
		if err == nil {
			return conv.BytesToString(b)
		}
		logger.Error().Error("marshal proto.Message error", err).Fire()
	}

	if p, ok := i.(fmt.Stringer); ok {
		return p.String()
	}

	return fmt.Sprintf("%+v", i)
}
