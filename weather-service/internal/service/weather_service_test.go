package service

import (
	"context"
	"errors"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockViaCepClient struct {
	mock.Mock
}

func (m *MockViaCepClient) GetLocation(ctx context.Context, zipcode string) (string, error) {
	args := m.Called(ctx, zipcode)
	return args.String(0), args.Error(1)
}

type MockWeatherApiClient struct {
	mock.Mock
}

func (m *MockWeatherApiClient) GetTemperature(ctx context.Context, city string) (float64, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(float64), args.Error(1)
}

func TestWeatherService_GetWeatherByZipcode(t *testing.T) {
	mockViaCep := new(MockViaCepClient)
	mockWeatherApi := new(MockWeatherApiClient)
	weatherService := NewWeatherService(mockViaCep, mockWeatherApi)

	t.Run("should return weather on success", func(t *testing.T) {
		mockViaCep.On("GetLocation", mock.Anything, "74305460").Return("Goiânia", nil).Once()
		mockWeatherApi.On("GetTemperature", mock.Anything, "Goiânia").Return(25.0, nil).Once()

		weather, err := weatherService.GetWeatherByZipcode(context.Background(), "74305460")

		assert.NoError(t, err)
		assert.NotNil(t, weather)
		assert.Equal(t, "Goiânia", weather.City)
		assert.Equal(t, 25.0, weather.TempC)
		assert.Equal(t, 77.0, weather.TempF)
		assert.Equal(t, 298.0, weather.TempK)
		mockViaCep.AssertExpectations(t)
		mockWeatherApi.AssertExpectations(t)
	})

	t.Run("should return error for invalid zipcode", func(t *testing.T) {
		weather, err := weatherService.GetWeatherByZipcode(context.Background(), "123")

		assert.Nil(t, weather)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidZipcode)
	})

	t.Run("should return error when zipcode is not found", func(t *testing.T) {
		notFoundErr := fault.New("not found", fault.WithCode(fault.NotFound))
		mockViaCep.On("GetLocation", mock.Anything, "00000000").Return("", notFoundErr).Once()

		weather, err := weatherService.GetWeatherByZipcode(context.Background(), "00000000")

		assert.Nil(t, weather)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrZipcodeNotFound)
		mockViaCep.AssertExpectations(t)
	})

	t.Run("should return error when weather is not found for a valid city", func(t *testing.T) {
		mockViaCep.On("GetLocation", mock.Anything, "74305460").Return("Goiânia", nil).Once()
		mockWeatherApi.On("GetTemperature", mock.Anything, "Goiânia").Return(0.0, errors.New("weather api failed")).Once()

		weather, err := weatherService.GetWeatherByZipcode(context.Background(), "74305460")

		assert.Nil(t, weather)
		assert.Error(t, err)
		assert.True(t, fault.IsInternal(err), "error should have internal code")
		assert.Contains(t, err.Error(), ErrWeatherNotFound.Message)
		mockViaCep.AssertExpectations(t)
		mockWeatherApi.AssertExpectations(t)
	})
}
