package di

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/marcelofabianov/weather-server/config"
	"github.com/marcelofabianov/weather-server/internal/adapter"
	"github.com/marcelofabianov/weather-server/internal/handler"
	"github.com/marcelofabianov/weather-server/internal/port"
	"github.com/marcelofabianov/weather-server/internal/service"
	"github.com/marcelofabianov/weather-server/pkg/logger"
	"github.com/marcelofabianov/weather-server/pkg/web"
)

var AppModule = fx.Options(
	ConfigModule,
	LoggerModule,
	OtelModule,
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
		func(cfg *config.Config) *config.OTELConfig { return &cfg.OTEL },
		func(clientsCfg *config.ClientsConfig) *config.WeatherAPIConfig { return &clientsCfg.WeatherAPI },
	),
)

var OtelModule = fx.Module("otel",
	fx.Provide(func(lc fx.Lifecycle, cfg *config.OTELConfig) (*tracesdk.TracerProvider, error) {
		ctx := context.Background()

		res, err := resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceNameKey.String(cfg.ServiceName),
				semconv.ServiceVersionKey.String("1.0.0"),
			),
		)
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		conn, err := grpc.NewClient(cfg.ExporterOTLPEndpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}

		traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, err
		}

		bsp := tracesdk.NewBatchSpanProcessor(traceExporter)
		tracerProvider := tracesdk.NewTracerProvider(
			tracesdk.WithSampler(tracesdk.AlwaysSample()),
			tracesdk.WithResource(res),
			tracesdk.WithSpanProcessor(bsp),
		)

		otel.SetTracerProvider(tracerProvider)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return tracerProvider.Shutdown(ctx)
			},
		})

		return tracerProvider, nil
	}),
	fx.Provide(func(tp *tracesdk.TracerProvider) trace.Tracer {
		return tp.Tracer("github.com/marcelofabianov/weather-server")
	}),
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
			func(breaker *gobreaker.CircuitBreaker, resilienceCfg *config.ResilienceConfig, logger *slog.Logger, tracer trace.Tracer) *adapter.ViaCepClient {
				return adapter.NewViaCepClient(breaker, resilienceCfg, logger, tracer)
			},
			fx.ParamTags(`name:"viaCepBreaker"`),
			fx.As(new(port.ViaCepClient)),
		),
		fx.Annotate(
			func(breaker *gobreaker.CircuitBreaker, clientCfg *config.WeatherAPIConfig, resilienceCfg *config.ResilienceConfig, logger *slog.Logger, tracer trace.Tracer) *adapter.WeatherApiClient {
				return adapter.NewWeatherApiClient(clientCfg, breaker, resilienceCfg, logger, tracer)
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