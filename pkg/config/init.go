package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConfig(configPath string) error {
	viper.SetConfigType("yaml")
	f, err := os.Open(configPath)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot initialize config file: [%s]\n", configPath))
	}
	defer f.Close()

	if err := viper.ReadConfig(f); err != nil {
		return errors.New(fmt.Sprintf("cannot read config file: [%s]\n", configPath))
	}

	return nil
}
