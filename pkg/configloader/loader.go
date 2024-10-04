package configloader

import (
	"reflect"

	flag "github.com/spf13/pflag"
)

type fieldTagsData struct {
	Desc    string
	Name    string
	Short   string
	Default string
}

func LoadConfig(cfg any) error {
	val := reflect.ValueOf(cfg).Elem()

	var filePath string
	flag.StringVar(&filePath, "config", "", "path to config file")

	err := parseFlags(val)
	if err != nil {
		return err
	}

	err = parseENVs(val)
	if err != nil {
		return err
	}

	err = parseFile(val, filePath)
	if err != nil {
		return err
	}

	return nil
}
