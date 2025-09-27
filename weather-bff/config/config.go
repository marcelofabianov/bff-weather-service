package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger  LoggerConfig  `mapstructure:"logger"`
	Server  ServerConfig  `mapstructure:"server"`
	Clients ClientsConfig `mapstructure:"clients"`
	OTEL    OTELConfig    `mapstructure:"otel"`
}

type OTELConfig struct {
	ServiceName          string `mapstructure:"service_name"`
	ExporterOTLPEndpoint string `mapstructure:"exporter_otlp_endpoint"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type ServerConfig struct {
	API  APIConfig  `mapstructure:"api"`
	CORS CORSConfig `mapstructure:"cors"`
}

type APIConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	RateLimit    int           `mapstructure:"rate_limit"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	MaxBodySize  int           `mapstructure:"maxbodysize"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowedorigins"`
	AllowedMethods   []string `mapstructure:"allowedmethods"`
	AllowedHeaders   []string `mapstructure:"allowedheaders"`
	ExposedHeaders   []string `mapstructure:"exposedheaders"`
	AllowCredentials bool     `mapstructure:"allowcredentials"`
}

type WeatherServiceConfig struct {
	URL string `mapstructure:"url"`
}

type ClientsConfig struct {
	WeatherService WeatherServiceConfig `mapstructure:"weatherservice"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetDefault("otel.service_name", "weather-bff")
	v.SetDefault("otel.exporter_otlp_endpoint", "otel-collector:4317")
	v.SetDefault("logger.level", "info")
	v.SetDefault("server.api.host", "0.0.0.0")
	v.SetDefault("server.api.port", 8080)
	v.SetDefault("server.api.rate_limit", 100)
	v.SetDefault("server.api.read_timeout", "5s")
	v.SetDefault("server.api.write_timeout", "10s")
	v.SetDefault("server.api.idle_timeout", "120s")
	v.SetDefault("server.api.maxbodysize", 1048576) // 1MB
	v.SetDefault("server.cors.allowedorigins", []string{"*"})
	v.SetDefault("server.cors.allowedmethods", []string{"GET", "POST"})
	v.SetDefault("server.cors.allowedheaders", []string{"Content-Type", "Authorization"})
	v.SetDefault("server.cors.exposedheaders", []string{})
	v.SetDefault("server.cors.allowcredentials", true)
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
