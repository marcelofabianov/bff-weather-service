package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	General    GeneralConfig    `mapstructure:"general"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Server     ServerConfig     `mapstructure:"server"`
	Resilience ResilienceConfig `mapstructure:"resilience"`
	Clients    ClientsConfig    `mapstructure:"clients"`
}

type GeneralConfig struct {
	Env string `mapstructure:"env"`
	TZ  string `mapstructure:"tz"`
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

type ResilienceConfig struct {
	RetryMaxAttempts    int           `mapstructure:"retry_max_attempts"`
	RetryInitialBackoff time.Duration `mapstructure:"retry_initial_backoff"`
	RetryMaxBackoff     time.Duration `mapstructure:"retry_max_backoff"`
	BreakerMaxFailures  uint32        `mapstructure:"breaker_max_failures"`
	BreakerTimeout      time.Duration `mapstructure:"breaker_timeout"`
}

type WeatherAPIConfig struct {
	URL string `mapstructure:"url"`
	Key string `mapstructure:"key"`
}

type ClientsConfig struct {
	WeatherAPI WeatherAPIConfig `mapstructure:"weatherapi"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// General
	v.SetDefault("general.env", "development")
	v.SetDefault("general.tz", "UTC")

	// Logger
	v.SetDefault("logger.level", "info")

	// Server
	v.SetDefault("server.api.host", "0.0.0.0")
	v.SetDefault("server.api.port", 8080)
	v.SetDefault("server.api.rate_limit", 100)
	v.SetDefault("server.api.read_timeout", "5s")
	v.SetDefault("server.api.write_timeout", "10s")
	v.SetDefault("server.api.idle_timeout", "120s")
	v.SetDefault("server.api.maxbodysize", 1048576) // 1MB

	// CORS
	v.SetDefault("server.cors.allowedorigins", []string{"*"})
	v.SetDefault("server.cors.allowedmethods", []string{"GET", "POST"})
	v.SetDefault("server.cors.allowedheaders", []string{"Content-Type", "Authorization"})
	v.SetDefault("server.cors.exposedheaders", []string{})
	v.SetDefault("server.cors.allowcredentials", true)

	// Resilience
	v.SetDefault("resilience.retry_max_attempts", 3)
	v.SetDefault("resilience.retry_initial_backoff", "100ms")
	v.SetDefault("resilience.retry_max_backoff", "2s")
	v.SetDefault("resilience.breaker_max_failures", 5)
	v.SetDefault("resilience.breaker_timeout", "30s")

	// Clients
	v.SetDefault("clients.weatherapi.url", "http://api.weatherapi.com/v1")
	v.SetDefault("clients.weatherapi.key", "")

	// Viper setup
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(path)

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
