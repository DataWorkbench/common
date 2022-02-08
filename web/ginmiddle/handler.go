package ginmiddle

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/yu31/protoc-plugin/xgo/pkg/protovalidator"

	"github.com/DataWorkbench/common/gtrace"
	"github.com/DataWorkbench/common/qerror"
)

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := qerror.NewResponse(qerror.NotAcceptable, gtrace.IdFromContext(GetStdContext(c)))
		c.JSON(resp.Status, resp)
	}
}

func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := qerror.NewResponse(qerror.MethodNotAllowed, gtrace.IdFromContext(GetStdContext(c)))
		c.JSON(resp.Status, resp)
	}
}

// ErrorHandler processing the error.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Call next handler
		c.Next()

		// Handle gin Context errors
		ginErr := c.Errors.Last()
		if ginErr == nil {
			return
		}

		if IsWekSocket(c) && c.Writer.Written() {
			return
		}

		var err *qerror.Error

		switch et := ginErr.Err.(type) {
		case *protovalidator.ValidateError:
			err = protovalidatorToQError(et)
		case validator.ValidationErrors:
			err = validatorErrorsToQError(et)
		case qerror.Error:
			err = &et
		case *qerror.Error:
			err = et
		default:
			err = qerror.FromGRPC(et)
			if err == nil {
				err = qerror.Internal
			}
		}

		resp := qerror.NewResponse(err, gtrace.IdFromContext(GetStdContext(c)))
		c.AbortWithStatusJSON(resp.Status, resp)
	}
}

func protovalidatorToQError(v *protovalidator.ValidateError) (err *qerror.Error) {
	err = qerror.ParameterValidationError.Format(v.Error())
	return
}

func validatorErrorsToQError(errs validator.ValidationErrors) (err *qerror.Error) {
	if len(errs) == 0 {
		return qerror.Internal
	}
	field := errs[0]

	name := field.Field()
	if name == "" {
		name = field.StructField()
	}

	switch field.Tag() {
	case "required":
		err = qerror.ParamsIsEmpty.Format(name)
	case "eq", "gt", "gte", "lt", "lte", "ne":
		switch field.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			err = qerror.InvalidParamsLength.Format(name, field.Tag(), field.Param(), reflect.ValueOf(field.Value()).Len())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			err = qerror.InvalidParamsValue.Format(name, field.Tag(), field.Param(), field.Value())
		default:
			err = qerror.InvalidParams.Format(name)
		}
	default:
		err = qerror.InvalidParams.Format(name)
	}

	return
}
