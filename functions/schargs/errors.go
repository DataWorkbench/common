package schargs

import (
	"errors"
)

var (
	ErrUnsupportedVariableNested = errors.New("defining nested variables is not supported")
	ErrInvalidVariableFormat     = errors.New("the variable must start with '${' and end with '}'")
	ErrVariableIsEmpty           = errors.New("the name of variable can not be empty")
	//ErrVariableNotDefined        = errors.New("variable not defined")
)
