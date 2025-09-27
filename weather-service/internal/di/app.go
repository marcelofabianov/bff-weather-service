package di

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"github.com/marcelofabianov/weather-api/internal/handler"
)

func NewApp() *fx.App {
	return fx.New(
		AppModule,
		fx.Invoke(registerHooks),
	)
}

func registerHooks(
	lifecycle fx.Lifecycle,
	logger *slog.Logger,
	handler *handler.WeatherHandler,
	router *chi.Mux,
	server *http.Server,
) {
	handler.RegisterRoutes(router)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting server", "address", server.Addr)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("failed to start server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping server")
			return server.Shutdown(ctx)
		},
	})
}
