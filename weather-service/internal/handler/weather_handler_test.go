package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/marcelofabianov/weather-api/internal/model"
	"github.com/marcelofabianov/weather-api/internal/service"
)

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetWeatherByZipcode(zipcode string) (*model.Weather, error) {
	args := m.Called(zipcode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Weather), args.Error(1)
}

func TestWeatherHandler_GetWeather(t *testing.T) {
	t.Run("should return 200 OK with weather data on success", func(t *testing.T) {
		mockService := new(MockWeatherService)
		handler := NewWeatherHandler(mockService)
		router := chi.NewMux()
		handler.RegisterRoutes(router)

		expectedWeather := model.NewWeather("Goiânia", 25.0)
		mockService.On("GetWeatherByZipcode", "74305460").Return(expectedWeather, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/weather/74305460", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"city":"Goiânia", "temp_C": 25.0, "temp_F": 77.0, "temp_K": 298.0}`, rr.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("should return 422 Unprocessable Entity for invalid zipcode", func(t *testing.T) {
		mockService := new(MockWeatherService)
		handler := NewWeatherHandler(mockService)
		router := chi.NewMux()
		handler.RegisterRoutes(router)

		mockService.On("GetWeatherByZipcode", "123").Return(nil, service.ErrInvalidZipcode).Once()

		req, _ := http.NewRequest(http.MethodGet, "/weather/123", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.Contains(t, rr.Body.String(), `"code":"domain_violation"`)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 404 Not Found when zipcode is not found", func(t *testing.T) {
		mockService := new(MockWeatherService)
		handler := NewWeatherHandler(mockService)
		router := chi.NewMux()
		handler.RegisterRoutes(router)

		mockService.On("GetWeatherByZipcode", "00000000").Return(nil, service.ErrZipcodeNotFound).Once()

		req, _ := http.NewRequest(http.MethodGet, "/weather/00000000", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), `"code":"not_found"`)
		mockService.AssertExpectations(t)
	})

	t.Run("should return 500 Internal Server Error for other errors", func(t *testing.T) {
		mockService := new(MockWeatherService)
		handler := NewWeatherHandler(mockService)
		router := chi.NewMux()
		handler.RegisterRoutes(router)

		genericError := errors.New("unexpected database error")
		mockService.On("GetWeatherByZipcode", "11111111").Return(nil, genericError).Once()

		req, _ := http.NewRequest(http.MethodGet, "/weather/11111111", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), `"code":"internal_error"`)
		mockService.AssertExpectations(t)
	})
}
