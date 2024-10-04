package configloader

import (
	"os"
	"reflect"
	"strconv"
)

//nolint:funlen,gocognit,gocyclo,cyclop // too long
func parseENVs(val reflect.Value) error {
	typ := val.Type()

	for i := range typ.NumField() {
		field := typ.Field(i)
		envName := field.Tag.Get("env")
		value := val.Field(i)

		envVal, exists := os.LookupEnv(envName)

		if !exists || !value.CanSet() {
			continue
		}

		//nolint:exhaustive // all supported types
		switch field.Type.Kind() {
		case reflect.String:
			value.SetString(envVal)

		case reflect.Int:
			intVar, err := strconv.ParseInt(envVal, 10, 0)
			if err != nil {
				return err
			}
			value.SetInt(intVar)

		case reflect.Int8:
			intVar, err := strconv.ParseInt(envVal, 10, 8)
			if err != nil {
				return err
			}
			value.SetInt(intVar)

		case reflect.Int16:
			intVar, err := strconv.ParseInt(envVal, 10, 16)
			if err != nil {
				return err
			}
			value.SetInt(intVar)

		case reflect.Int32:
			intVar, err := strconv.ParseInt(envVal, 10, 32)
			if err != nil {
				return err
			}
			value.SetInt(intVar)

		case reflect.Int64:
			intVar, err := strconv.ParseInt(envVal, 10, 64)
			if err != nil {
				return err
			}
			value.SetInt(intVar)

		case reflect.Uint:
			uintVar, err := strconv.ParseUint(envVal, 10, 0)
			if err != nil {
				return err
			}
			value.SetUint(uintVar)

		case reflect.Uint8:
			uintVar, err := strconv.ParseUint(envVal, 10, 8)
			if err != nil {
				return err
			}
			value.SetUint(uintVar)

		case reflect.Uint16:
			uintVar, err := strconv.ParseUint(envVal, 10, 16)
			if err != nil {
				return err
			}
			value.SetUint(uintVar)

		case reflect.Uint32:
			uintVar, err := strconv.ParseUint(envVal, 10, 32)
			if err != nil {
				return err
			}
			value.SetUint(uintVar)

		case reflect.Uint64:
			uintVar, err := strconv.ParseUint(envVal, 10, 64)
			if err != nil {
				return err
			}
			value.SetUint(uintVar)

		case reflect.Float32:
			floatVal, err := strconv.ParseFloat(envVal, 32)
			if err != nil {
				return err
			}
			value.SetFloat(floatVal)

		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(envVal, 64)
			if err != nil {
				return err
			}
			value.SetFloat(floatVal)

		case reflect.Bool:
			boolVar, err := strconv.ParseBool(envVal)
			if err != nil {
				return err
			}
			value.SetBool(boolVar)

		case reflect.Slice:
			return UnsupportedTypeError{field}
		case reflect.Struct:
			err := parseENVs(value)
			if err != nil {
				return err
			}
		default:
			return UnsupportedTypeError{field}
		}
	}

	return nil
}
