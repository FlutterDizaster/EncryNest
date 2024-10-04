package config

import (
	flag "github.com/spf13/pflag"
)

func Load(dest any) error {
	var configFile string
	flag.StringVarP(&configFile, "config", "c", "", "config file path")

	// TODO: define all other flags

	flag.Parse()

	return nil
}
