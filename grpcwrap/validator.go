package grpcwrap

import (
	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// validator defines interface grpc_validator.validator
// See https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/validator/validator.go#L14
type validator interface {
	Validate() error
}

// validateRequestArgument helper for validate the request arguments
func validateRequestArgument(i interface{}, logger *glog.Logger) error {
	if v, ok := i.(validator); ok {
		if err := v.Validate(); err != nil {
			logger.Error().Error("failed validate request", err).Fire()
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil
	}
	logger.Debug().Msg("request argument not implement validator").Fire()
	return nil
}

// validateReplyArgument helper for validate the reply arguments
func validateReplyArgument(i interface{}, logger *glog.Logger) error {
	if v, ok := i.(validator); ok {
		if err := v.Validate(); err != nil {
			logger.Error().Error("failed validate reply", err).Fire()
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil
	}
	logger.Debug().Msg("reply argument not implement validator").Fire()
	return nil
}
