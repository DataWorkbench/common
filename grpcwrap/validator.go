package grpcwrap

import (
	"github.com/DataWorkbench/glog"
	"github.com/yu31/proto-go-plugin/pkg/protodefaults"
	"github.com/yu31/proto-go-plugin/pkg/protovalidator"
)

//// validator defines interface grpc_validator.validator
//// See https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/validator/validator.go#L14
//type validator interface {
//	Validate() error
//}

// validateRequestParameters helper for validate the request arguments
func validateRequestParameters(i interface{}, logger *glog.Logger) error {
	// Set defaults values.
	protodefaults.CallDefaultsIfExists(i)

	if v, ok := i.(protovalidator.Validator); ok {
		if err := v.Validate(); err != nil {
			logger.Error().Error("grpc invalid request parameters", err).Fire()
			return err
		}
		return nil
	}
	logger.Warn().Msg("grpc request message not implement validator").Fire()
	return nil
}

// validateReplyParameters helper for validate the reply arguments
func validateReplyParameters(i interface{}, logger *glog.Logger) error {
	// Set defaults values.
	protodefaults.CallDefaultsIfExists(i)

	if v, ok := i.(protovalidator.Validator); ok {
		if err := v.Validate(); err != nil {
			logger.Error().Error("grpc invalid reply parameters", err).Fire()
			return err
		}
		return nil
	}
	logger.Warn().Msg("grpc reply message not implement validator").Fire()
	return nil
}
