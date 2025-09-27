package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcelofabianov/fault"

	"github.com/marcelofabianov/weather-bff/internal/model"
	"github.com/marcelofabianov/weather-bff/internal/port"
	"github.com/marcelofabianov/weather-bff/pkg/web"
)

type ZipcodeHandler struct {
	service   port.BffService
	validator port.Validator
}

func NewZipcodeHandler(service port.BffService, validator port.Validator) *ZipcodeHandler {
	return &ZipcodeHandler{
		service:   service,
		validator: validator,
	}
}

func (h *ZipcodeHandler) HandleCep(w http.ResponseWriter, r *http.Request) {
	var input model.CepInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		web.Error(w, r, fault.New("bad request: invalid json body", fault.WithCode(fault.Invalid)))
		return
	}

	if err := h.validator.Validate(&input); err != nil {
		web.Error(w, r, err)
		return
	}

	weather, err := h.service.GetWeather(r.Context(), input.Cep)
	if err != nil {
		web.Error(w, r, err)
		return
	}

	web.Success(w, r, http.StatusOK, weather)
}

func (h *ZipcodeHandler) RegisterRoutes(router *chi.Mux) {
	router.Post("/weather", h.HandleCep)
}
