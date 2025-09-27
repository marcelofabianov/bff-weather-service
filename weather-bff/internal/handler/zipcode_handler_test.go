package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/marcelofabianov/weather-bff/internal/model"
)

type MockBffService struct {
	mock.Mock
}

func (m *MockBffService) GetWeather(ctx context.Context, cep string) (*model.WeatherResponse, error) {
	args := m.Called(ctx, cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WeatherResponse), args.Error(1)
}

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(data any) error {
	args := m.Called(data)
	return args.Error(0)
}

func TestZipcodeHandler_HandleCep(t *testing.T) {
	mockService := new(MockBffService)
	mockValidator := new(MockValidator)
	handler := NewZipcodeHandler(mockService, mockValidator)
	router := chi.NewMux()
	handler.RegisterRoutes(router)

	t.Run("should return 200 OK on success", func(t *testing.T) {
		cepInput := &model.CepInput{Cep: "74305460"}
		jsonInput := `{"cep": "74305460"}`
		mockValidator.On("Validate", cepInput).Return(nil).Once()
		mockService.On("GetWeather", mock.Anything, "74305460").Return(&model.WeatherResponse{City: "Goiânia", TempC: 31.0}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/weather", strings.NewReader(jsonInput))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"city":"Goiânia"`)
		mockService.AssertExpectations(t)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return 422 when validation fails", func(t *testing.T) {
		cepInput := &model.CepInput{Cep: "123"}
		jsonInput := `{"cep": "123"}`
		validationErr := fault.New("validation failed", fault.WithCode(fault.DomainViolation))
		mockValidator.On("Validate", cepInput).Return(validationErr).Once()

		req, _ := http.NewRequest(http.MethodPost, "/weather", strings.NewReader(jsonInput))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.Contains(t, rr.Body.String(), `"code":"domain_violation"`)
		mockValidator.AssertExpectations(t)
	})

	t.Run("should return 400 when json is malformed", func(t *testing.T) {
		jsonInput := `{"cep": "12345678"`

		req, _ := http.NewRequest(http.MethodPost, "/weather", strings.NewReader(jsonInput))
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
