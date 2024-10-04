package configloader

import (
	"reflect"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

func parseFlags(value reflect.Value) error {
	err := handleFlags(value)
	if err != nil {
		return err
	}

	flag.Parse()

	return nil
}

//nolint:funlen // too long
func handleFlags(value reflect.Value) error {
	typ := value.Type()

	for i := range typ.NumField() {
		field := typ.Field(i)
		val := value.Field(i)
		tag := field.Tag

		if !value.CanSet() {
			continue
		}

		var tagsData fieldTagsData

		// Parsing tags
		tagsData.Name = tag.Get("name")
		if tagsData.Name == "" {
			tagsData.Name = strings.ToLower(field.Name)
		}
		tagsData.Short = tag.Get("short")
		tagsData.Default = tag.Get("default")
		tagsData.Desc = tag.Get("desc")
		if tagsData.Desc == "" {
			tagsData.Desc = field.Name
		}

		var err error

		//nolint:exhaustive // all supported types
		switch field.Type.Kind() {
		case reflect.String:
			handleStringFlag(val, tagsData)
		case reflect.Int:
			err = handleIntFlag(val, tagsData)
		case reflect.Int8:
			err = handleInt8Flag(val, tagsData)
		case reflect.Int16:
			err = handleInt16Flag(val, tagsData)
		case reflect.Int32:
			err = handleInt32Flag(val, tagsData)
		case reflect.Int64:
			err = handleInt64Flag(val, tagsData)
		case reflect.Uint:
			err = handleUintFlag(val, tagsData)
		case reflect.Uint8:
			err = handleUint8Flag(val, tagsData)
		case reflect.Uint16:
			err = handleUint16Flag(val, tagsData)
		case reflect.Uint32:
			err = handleUint32Flag(val, tagsData)
		case reflect.Uint64:
			err = handleUint64Flag(val, tagsData)
		case reflect.Float32:
			err = handleFloat32Flag(val, tagsData)
		case reflect.Float64:
			err = handleFloat64Flag(val, tagsData)
		case reflect.Bool:
			err = handleBoolFlag(val, tagsData)
		// case reflect.Slice:
		// 	err = handleSliceFlag(val, tagsData)
		case reflect.Struct:
			err = handleFlags(val)
		default:
			err = UnsupportedTypeError{field}
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// func handleSliceFlag(_ reflect.Value, _ fieldTagsData) error {
// 	return errors.New("Not implemented yet")
// }

func handleStringFlag(val reflect.Value, tags fieldTagsData) {
	if tags.Short == "" {
		flag.StringVar(
			val.Addr().Interface().(*string),
			tags.Name,
			tags.Default,
			tags.Desc,
		)
	} else {
		flag.StringVarP(
			val.Addr().Interface().(*string),
			tags.Name,
			tags.Short,
			tags.Default,
			tags.Desc,
		)
	}
}

func handleIntFlag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseInt(tags.Default, 10, 0)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.IntVar(
			val.Addr().Interface().(*int),
			tags.Name,
			int(def),
			tags.Desc,
		)
	} else {
		flag.IntVarP(
			val.Addr().Interface().(*int),
			tags.Name,
			tags.Short,
			int(def),
			tags.Desc,
		)
	}
	return nil
}

func handleInt8Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseInt(tags.Default, 10, 8)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Int8Var(
			val.Addr().Interface().(*int8),
			tags.Name,
			int8(def),
			tags.Desc,
		)
	} else {
		flag.Int8VarP(
			val.Addr().Interface().(*int8),
			tags.Name,
			tags.Short,
			int8(def),
			tags.Desc,
		)
	}
	return nil
}

func handleInt16Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseInt(tags.Default, 10, 16)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Int16Var(
			val.Addr().Interface().(*int16),
			tags.Name,
			int16(def),
			tags.Desc,
		)
	} else {
		flag.Int16VarP(
			val.Addr().Interface().(*int16),
			tags.Name,
			tags.Short,
			int16(def),
			tags.Desc,
		)
	}
	return nil
}

func handleInt32Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseInt(tags.Default, 10, 32)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Int32Var(
			val.Addr().Interface().(*int32),
			tags.Name,
			int32(def),
			tags.Desc,
		)
	} else {
		flag.Int32VarP(
			val.Addr().Interface().(*int32),
			tags.Name,
			tags.Short,
			int32(def),
			tags.Desc,
		)
	}
	return nil
}

func handleInt64Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseInt(tags.Default, 10, 64)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Int64Var(
			val.Addr().Interface().(*int64),
			tags.Name,
			def,
			tags.Desc,
		)
	} else {
		flag.Int64VarP(
			val.Addr().Interface().(*int64),
			tags.Name,
			tags.Short,
			def,
			tags.Desc,
		)
	}
	return nil
}

func handleBoolFlag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseBool(tags.Default)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.BoolVar(
			val.Addr().Interface().(*bool),
			tags.Name,
			def,
			tags.Desc,
		)
	} else {
		flag.BoolVarP(
			val.Addr().Interface().(*bool),
			tags.Name,
			tags.Short,
			def,
			tags.Desc,
		)
	}
	return nil
}

func handleFloat32Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseFloat(tags.Default, 32)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Float32Var(
			val.Addr().Interface().(*float32),
			tags.Name,
			float32(def),
			tags.Desc,
		)
	} else {
		flag.Float32VarP(
			val.Addr().Interface().(*float32),
			tags.Name,
			tags.Short,
			float32(def),
			tags.Desc,
		)
	}
	return nil
}

func handleFloat64Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseFloat(tags.Default, 64)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Float64Var(
			val.Addr().Interface().(*float64),
			tags.Name,
			def,
			tags.Desc,
		)
	} else {
		flag.Float64VarP(
			val.Addr().Interface().(*float64),
			tags.Name,
			tags.Short,
			def,
			tags.Desc,
		)
	}
	return nil
}

func handleUintFlag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseUint(tags.Default, 10, 64)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.UintVar(
			val.Addr().Interface().(*uint),
			tags.Name,
			uint(def),
			tags.Desc,
		)
	} else {
		flag.UintVarP(
			val.Addr().Interface().(*uint),
			tags.Name,
			tags.Short,
			uint(def),
			tags.Desc,
		)
	}
	return nil
}

func handleUint8Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseUint(tags.Default, 10, 8)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Uint8Var(
			val.Addr().Interface().(*uint8),
			tags.Name,
			uint8(def),
			tags.Desc,
		)
	} else {
		flag.Uint8VarP(
			val.Addr().Interface().(*uint8),
			tags.Name,
			tags.Short,
			uint8(def),
			tags.Desc,
		)
	}
	return nil
}

func handleUint16Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseUint(tags.Default, 10, 16)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Uint16Var(
			val.Addr().Interface().(*uint16),
			tags.Name,
			uint16(def),
			tags.Desc,
		)
	} else {
		flag.Uint16VarP(
			val.Addr().Interface().(*uint16),
			tags.Name,
			tags.Short,
			uint16(def),
			tags.Desc,
		)
	}
	return nil
}

func handleUint32Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseUint(tags.Default, 10, 32)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Uint32Var(
			val.Addr().Interface().(*uint32),
			tags.Name,
			uint32(def),
			tags.Desc,
		)
	} else {
		flag.Uint32VarP(
			val.Addr().Interface().(*uint32),
			tags.Name,
			tags.Short,
			uint32(def),
			tags.Desc,
		)
	}
	return nil
}

func handleUint64Flag(val reflect.Value, tags fieldTagsData) error {
	def, err := strconv.ParseUint(tags.Default, 10, 64)
	if err != nil {
		return err
	}
	if tags.Short == "" {
		flag.Uint64Var(
			val.Addr().Interface().(*uint64),
			tags.Name,
			def,
			tags.Desc,
		)
	} else {
		flag.Uint64VarP(
			val.Addr().Interface().(*uint64),
			tags.Name,
			tags.Short,
			def,
			tags.Desc,
		)
	}
	return nil
}
