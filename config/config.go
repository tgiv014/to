package config

import (
	"bytes"
	_ "embed"
	"errors"

	"github.com/spf13/viper"
)

//go:embed default.yaml
var defaultConfigBytes []byte

type Config struct {
	Hostname string
	AuthKey  string
	DBPath   string `mapstructure:"db_path"`
}

func LoadConfig() (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	reader := bytes.NewReader(defaultConfigBytes)
	err := v.ReadConfig(reader)
	if err != nil {
		return Config{}, err
	}

	err = v.MergeInConfig()
	if err != nil && !isNotFound(err) {
		return Config{}, err
	}

	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func isNotFound(err error) bool {
	var notFoundError viper.ConfigFileNotFoundError
	if errors.As(err, &notFoundError) {
		return true
	}
	return false
}
