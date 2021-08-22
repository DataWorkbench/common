package ginvalidator

import (
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func New() binding.StructValidator {
	sv := &structValidator{
		validate: validator.New(),
	}
	sv.validate.SetTagName("binding")
	sv.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("params")
		if tag == "" {
			tag = field.Tag.Get("json")
		}
		name := strings.SplitN(tag, ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return sv
}

// structValidator for impls web/binding.StructValidator to set default and validate struct
// see: https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L52
type structValidator struct {
	validate *validator.Validate
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *structValidator) ValidateStruct(obj interface{}) error {
	value := reflect.ValueOf(obj)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	if valueType == reflect.Struct {
		// Set default value
		if err := defaults.Set(obj); err != nil {
			return err
		}
		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}

// Engine returns the underlying validator engine which powers the default
// Validator scheduler. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (v *structValidator) Engine() interface{} {
	return v.validate
}
