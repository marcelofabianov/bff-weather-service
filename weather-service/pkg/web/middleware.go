package web

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type contextKey string

const LoggerCtxKey = contextKey("logger")

func SlogLoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())

			ctxLogger := logger.With("request_id", requestID)

			ctx := context.WithValue(r.Context(), LoggerCtxKey, ctxLogger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerCtxKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
