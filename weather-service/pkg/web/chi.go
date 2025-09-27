package web

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/marcelofabianov/weather-api/config"
)

func NewServer(cfg *config.Config, logger *slog.Logger, router *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.API.Host, cfg.Server.API.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.API.ReadTimeout,
		WriteTimeout: cfg.Server.API.WriteTimeout,
		IdleTimeout:  cfg.Server.API.IdleTimeout,
	}
}

func NewRouter(cfg *config.ServerConfig, logger *slog.Logger) *chi.Mux {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(SlogLoggerMiddleware(logger))
	r.Use(httprate.Limit(
		cfg.API.RateLimit,
		1*time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
		httprate.WithResponseHeaders(headersRateLimit()),
	))
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RequestSize(int64(cfg.API.MaxBodySize)))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(setCorsOptions(cfg.CORS)))

	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "deny"))
	r.Use(middleware.SetHeader("X-DNS-Prefetch-Control", "off"))
	r.Use(middleware.SetHeader("X-Download-Options", "noopen"))
	r.Use(middleware.SetHeader("Content-Security-Policy", "default-src 'none'"))
	r.Use(middleware.SetHeader("Referrer-Policy", "no-referrer"))
	r.Use(middleware.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains"))
	r.Use(middleware.SetHeader("Cache-Control", "no-store, no-cache"))

	r.Use(middleware.SetHeader("Cross-Origin-Resource-Policy", "same-origin"))
	r.Use(middleware.SetHeader("Cross-Origin-Opener-Policy", "same-origin"))
	r.Use(middleware.SetHeader("Permissions-Policy", "camera=(), microphone=(), geolocation=()"))

	return r
}

func headersRateLimit() httprate.ResponseHeaders {
	return httprate.ResponseHeaders{
		Limit:     "X-RateLimit-Limit",
		Remaining: "X-RateLimit-Remaining",
		Reset:     "X-RateLimit-Reset",
	}
}

func setCorsOptions(cfg config.CORSConfig) cors.Options {
	return cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		ExposedHeaders:   cfg.ExposedHeaders,
		AllowCredentials: cfg.AllowCredentials,
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "pong",
		"status":  "ok",
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
