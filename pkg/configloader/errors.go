package configloader

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotImplementedYet = errors.New("not implemented yet")
)

type UnsupportedTypeError struct {
	field reflect.StructField
}

func (e UnsupportedTypeError) Error() string {
	return fmt.Sprintf(
		"unsupported type: %s of variable %s",
		e.field.Type.Kind().String(),
		e.field.Name,
	)
}
