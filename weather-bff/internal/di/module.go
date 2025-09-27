package di

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
	OtelModule,
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
		func(cfg *config.Config) *config.OTELConfig { return &cfg.OTEL },
		func(clientsCfg *config.ClientsConfig) *config.WeatherServiceConfig { return &clientsCfg.WeatherService },
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
)

var LoggerModule = fx.Module("logger",
	fx.Provide(logger.NewSlogLogger),
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
		web.NewRouter,
		web.NewServer,
	),
)

var UtilsModule = fx.Module("utils",
	fx.Provide(
		fx.Annotate(
			validator.NewValidator,
			fx.As(new(port.Validator)),
		),
	),
)
