package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/marcelofabianov/weather-server/internal/port"
	"github.com/marcelofabianov/weather-server/pkg/web"
)

type WeatherHandler struct {
	service port.WeatherService
	tracer  trace.Tracer
}

func NewWeatherHandler(service port.WeatherService, tracer trace.Tracer) *WeatherHandler {
	return &WeatherHandler{service: service, tracer: tracer}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tracer.Start(r.Context(), "WeatherHandler.GetWeather")
	defer span.End()

	zipcode := chi.URLParam(r, "zipcode")
	span.SetAttributes(attribute.String("zipcode", zipcode))

	weather, err := h.service.GetWeatherByZipcode(ctx, zipcode)
	if err != nil {
		web.Error(w, r, err)
		return
	}

	web.Success(w, r, http.StatusOK, weather)
}

func (h *WeatherHandler) RegisterRoutes(router *chi.Mux) {
	router.Get("/weather/{zipcode}", h.GetWeather)
}
