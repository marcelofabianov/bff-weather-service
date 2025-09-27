package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/marcelofabianov/fault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherServiceClient struct {
	mock.Mock
}

func (m *MockWeatherServiceClient) GetWeatherForCep(ctx context.Context, cep string) (*http.Response, error) {
	args := m.Called(ctx, cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestBffService_GetWeather(t *testing.T) {
	mockClient := new(MockWeatherServiceClient)
	bffService := NewBffService(mockClient)
	ctx := context.Background()

	t.Run("should return weather on success", func(t *testing.T) {
		jsonResponse := `{"city":"Goiânia","temp_C":31.0,"temp_F":87.8,"temp_K":304.0}`
		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(jsonResponse))),
		}
		mockClient.On("GetWeatherForCep", ctx, "74305460").Return(mockResponse, nil).Once()

		weather, err := bffService.GetWeather(ctx, "74305460")

		assert.NoError(t, err)
		assert.NotNil(t, weather)
		assert.Equal(t, "Goiânia", weather.City)
		assert.Equal(t, 31.0, weather.TempC)
		mockClient.AssertExpectations(t)
	})

	t.Run("should return proxied error when downstream service returns a fault error", func(t *testing.T) {
		jsonError := `{"message":"can not find zipcode","code":"not_found"}`
		mockResponse := &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(bytes.NewReader([]byte(jsonError))),
		}
		mockClient.On("GetWeatherForCep", ctx, "00000000").Return(mockResponse, nil).Once()

		weather, err := bffService.GetWeather(ctx, "00000000")

		assert.Nil(t, weather)
		assert.Error(t, err)
		assert.True(t, fault.IsNotFound(err))
		assert.Equal(t, "can not find zipcode", err.Error())
		mockClient.AssertExpectations(t)
	})

	t.Run("should return infra error on communication failure", func(t *testing.T) {
		mockClient.On("GetWeatherForCep", ctx, "12345678").Return(nil, errors.New("connection refused")).Once()

		weather, err := bffService.GetWeather(ctx, "12345678")

		assert.Nil(t, weather)
		assert.Error(t, err)
		assert.True(t, fault.IsInfraError(err))
		mockClient.AssertExpectations(t)
	})
}
