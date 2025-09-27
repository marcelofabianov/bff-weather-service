package di

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sony/gobreaker"
	"go.uber.org/fx"

	"github.com/marcelofabianov/weather-api/config"
	"github.com/marcelofabianov/weather-api/internal/adapter"
	"github.com/marcelofabianov/weather-api/internal/handler"
	"github.com/marcelofabianov/weather-api/internal/port"
	"github.com/marcelofabianov/weather-api/internal/service"
	"github.com/marcelofabianov/weather-api/pkg/logger"
	"github.com/marcelofabianov/weather-api/pkg/web"
)

var AppModule = fx.Options(
	ConfigModule,
	LoggerModule,
	AdaptersModule,
	CoreModule,
	HandlersModule,
	WebServerModule,
)

var ConfigModule = fx.Module("config",
	fx.Provide(
		func() (*config.Config, error) {
			return config.LoadConfig(".")
		},
		func(cfg *config.Config) *config.LoggerConfig { return &cfg.Logger },
		func(cfg *config.Config) *config.ServerConfig { return &cfg.Server },
		func(cfg *config.Config) *config.ResilienceConfig { return &cfg.Resilience },
		func(cfg *config.Config) *config.ClientsConfig { return &cfg.Clients },
		func(clientsCfg *config.ClientsConfig) *config.WeatherAPIConfig { return &clientsCfg.WeatherAPI },
	),
)

var LoggerModule = fx.Module("logger",
	fx.Provide(func(cfg *config.LoggerConfig) *slog.Logger {
		return logger.NewSlogLogger(cfg)
	}),
)

var AdaptersModule = fx.Module("adapters",
	fx.Provide(
		fx.Annotate(
			func(cfg *config.ResilienceConfig) *gobreaker.CircuitBreaker {
				st := gobreaker.Settings{
					Name:        "ViaCEP",
					MaxRequests: 3,
					Interval:    0,
					Timeout:     cfg.BreakerTimeout,
					ReadyToTrip: func(counts gobreaker.Counts) bool {
						return counts.ConsecutiveFailures > cfg.BreakerMaxFailures
					},
				}
				return gobreaker.NewCircuitBreaker(st)
			},
			fx.ResultTags(`name:"viaCepBreaker"`),
		),
		fx.Annotate(
			func(cfg *config.ResilienceConfig) *gobreaker.CircuitBreaker {
				st := gobreaker.Settings{
					Name:        "WeatherAPI",
					MaxRequests: 3,
					Interval:    0,
					Timeout:     cfg.BreakerTimeout,
					ReadyToTrip: func(counts gobreaker.Counts) bool {
						return counts.ConsecutiveFailures > cfg.BreakerMaxFailures
					},
				}
				return gobreaker.NewCircuitBreaker(st)
			},
			fx.ResultTags(`name:"weatherApiBreaker"`),
		),
	),
	fx.Provide(
		fx.Annotate(
			func(breaker *gobreaker.CircuitBreaker, resilienceCfg *config.ResilienceConfig, logger *slog.Logger) *adapter.ViaCepClient {
				return adapter.NewViaCepClient(breaker, resilienceCfg, logger)
			},
			fx.ParamTags(`name:"viaCepBreaker"`),
			fx.As(new(port.ViaCepClient)),
		),
		fx.Annotate(
			func(breaker *gobreaker.CircuitBreaker, clientCfg *config.WeatherAPIConfig, resilienceCfg *config.ResilienceConfig, logger *slog.Logger) *adapter.WeatherApiClient {
				return adapter.NewWeatherApiClient(clientCfg, breaker, resilienceCfg, logger)
			},
			fx.ParamTags(`name:"weatherApiBreaker"`),
			fx.As(new(port.WeatherApiClient)),
		),
	),
)

var CoreModule = fx.Module("core",
	fx.Provide(
		fx.Annotate(
			service.NewWeatherService,
			fx.As(new(port.WeatherService)),
		),
	),
)

var HandlersModule = fx.Module("handlers",
	fx.Provide(handler.NewWeatherHandler),
)

var WebServerModule = fx.Module("webserver",
	fx.Provide(
		func(cfg *config.ServerConfig, logger *slog.Logger) *chi.Mux {
			return web.NewRouter(cfg, logger)
		},
		func(cfg *config.Config, logger *slog.Logger, router *chi.Mux) *http.Server {
			return web.NewServer(cfg, logger, router)
		},
	),
)
