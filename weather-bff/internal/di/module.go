package di

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"github.com/marcelofabianov/weather-bff/config"
	"github.com/marcelofabianov/weather-bff/internal/adapter"
	"github.com/marcelofabianov/weather-bff/internal/handler"
	"github.com/marcelofabianov/weather-bff/internal/port"
	"github.com/marcelofabianov/weather-bff/internal/service"
	"github.com/marcelofabianov/weather-bff/pkg/logger"
	"github.com/marcelofabianov/weather-bff/pkg/validator"
	"github.com/marcelofabianov/weather-bff/pkg/web"
)

var AppModule = fx.Options(
	ConfigModule,
	LoggerModule,
	UtilsModule,
	AdaptersModule,
	CoreModule,
	HandlersModule,
	WebServerModule,
)

var ConfigModule = fx.Module("config",
	fx.Provide(
		config.LoadConfig,
		func(cfg *config.Config) *config.LoggerConfig { return &cfg.Logger },
		func(cfg *config.Config) *config.ServerConfig { return &cfg.Server },
		func(cfg *config.Config) *config.ClientsConfig { return &cfg.Clients },
		func(clientsCfg *config.ClientsConfig) *config.WeatherServiceConfig { return &clientsCfg.WeatherService },
	),
)

var LoggerModule = fx.Module("logger",
	fx.Provide(logger.NewSlogLogger),
)

var UtilsModule = fx.Module("utils",
	fx.Provide(validator.NewValidator),
)

var AdaptersModule = fx.Module("adapters",
	fx.Provide(
		fx.Annotate(
			adapter.NewWeatherServiceClient,
			fx.As(new(port.WeatherServiceClient)),
		),
	),
)

var CoreModule = fx.Module("core",
	fx.Provide(
		fx.Annotate(
			service.NewBffService,
			fx.As(new(port.BffService)),
		),
	),
)

var HandlersModule = fx.Module("handlers",
	fx.Provide(handler.NewZipcodeHandler),
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
