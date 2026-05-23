package config

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
}

type Config struct {
	Host            string
	Port            int
	ExternalBaseUrl string
	UserPath        string
	TodosPath       string
}

func LoadConfig(ctx context.Context) (*Config, error) {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	cfg := &Config{
		Host:            viper.GetString("host"),
		Port:            viper.GetInt("port"),
		ExternalBaseUrl: viper.GetString("externalBaseUrl"),
		UserPath:        viper.GetString("userPath"),
		TodosPath:       viper.GetString("todosPath"),
	}

	return cfg, nil
}
