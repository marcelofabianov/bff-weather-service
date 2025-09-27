package di

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"github.com/marcelofabianov/weather-bff/internal/handler"
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
	handler *handler.ZipcodeHandler,
	router *chi.Mux,
	server *http.Server,
) {
	handler.RegisterRoutes(router)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting BFF server", "address", server.Addr)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("failed to start BFF server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping BFF server")
			return server.Shutdown(ctx)
		},
	})
}
