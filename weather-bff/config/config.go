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

type WeatherServiceConfig struct {
	URL string `mapstructure:"url"`
}

type ClientsConfig struct {
	WeatherService WeatherServiceConfig `mapstructure:"weatherservice"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetDefault("logger.level", "info")
	v.SetDefault("server.port", 8080)

	v.SetDefault("clients.weatherservice.url", "http://weather-service:8080")

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
