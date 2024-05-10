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
	Hostname   string `mapstructure:"hostname"`
	AuthKey    string `mapstructure:"auth_key"`
	DBPath     string `mapstructure:"db_path"`
	DataPath   string `mapstructure:"data_path"`
	ConfigPath string `mapstructure:"-"`
}

func LoadConfig() (Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./data")

	v.BindEnv("hostname", "TO_HOSTNAME")
	v.BindEnv("auth_key", "TO_AUTH_KEY")
	v.BindEnv("db_path", "TO_DB_PATH")
	v.BindEnv("data_path", "TO_DATA_PATH")

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
	config.ConfigPath = v.ConfigFileUsed()

	return config, nil
}

func isNotFound(err error) bool {
	var notFoundError viper.ConfigFileNotFoundError
	if errors.As(err, &notFoundError) {
		return true
	}
	return false
}
