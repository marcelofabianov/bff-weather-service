package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/marcelofabianov/weather-api/internal/port"
	"github.com/marcelofabianov/weather-api/pkg/web"
)

type WeatherHandler struct {
	service port.WeatherService
}

func NewWeatherHandler(service port.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: service}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	zipcode := chi.URLParam(r, "zipcode")

	weather, err := h.service.GetWeatherByZipcode(zipcode)
	if err != nil {
		web.Error(w, r, err)
		return
	}

	web.Success(w, r, http.StatusOK, weather)
}

func (h *WeatherHandler) RegisterRoutes(router *chi.Mux) {
	router.Get("/weather/{zipcode}", h.GetWeather)
}
