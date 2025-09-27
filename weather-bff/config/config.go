package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConfig  `mapstructure:"logger"`
	Server  ServerConfig  `mapstructure:"server"`
	Clients ClientsConfig `mapstructure:"clients"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type ServiceBConfig struct {
	URL string `mapstructure:"url"`
}

type ClientsConfig struct {
	ServiceB ServiceBConfig `mapstructure:"serviceb"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetDefault("logger.level", "info")
	v.SetDefault("server.port", 8080)
	v.SetDefault("clients.serviceb.url", "http://localhost:8080")

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
