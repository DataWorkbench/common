package ginvalidator

import (
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/yu31/protoc-plugin/xgo/pkg/protodefaults"
	"github.com/yu31/protoc-plugin/xgo/pkg/protovalidator"
)

func New() binding.StructValidator {
	sv := &structValidator{
		validate: validator.New(),
	}

	sv.validate.SetTagName("binding")
	sv.validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("json")
		if tag == "" {
			return ""
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
	if t, ok := obj.(protodefaults.Defaults); ok {
		t.SetDefaults()
	} else {
		if err := defaults.Set(obj); err != nil {
			if err.Error() == "not a struct pointer" {
				err = nil
			}
			return err
		}
	}

	if t, ok := obj.(protovalidator.Validator); ok {
		if err := t.Validate(); err != nil {
			return err
		}
	} else {
		if err := v.validate.Struct(obj); err != nil {
			if _, x := err.(*validator.InvalidValidationError); x {
				err = nil
			}
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
