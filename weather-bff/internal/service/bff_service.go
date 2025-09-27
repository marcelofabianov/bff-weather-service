package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/marcelofabianov/fault"

	"github.com/marcelofabianov/weather-bff/internal/model"
	"github.com/marcelofabianov/weather-bff/internal/port"
)

var (
	ErrWeatherServiceComms = fault.New("communication with weather-service failed", fault.WithCode(fault.InfraError))
)

type BffService struct {
	weatherServiceClient port.WeatherServiceClient
}

func NewBffService(client port.WeatherServiceClient) *BffService {
	return &BffService{
		weatherServiceClient: client,
	}
}

func (s *BffService) GetWeather(ctx context.Context, cep string) (*model.WeatherResponse, error) {
	resp, err := s.weatherServiceClient.GetWeatherForCep(ctx, cep)
	if err != nil {
		return nil, fault.Wrap(err, ErrWeatherServiceComms.Message, fault.WithCode(ErrWeatherServiceComms.Code))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fault.Wrap(err, ErrWeatherServiceComms.Message, fault.WithCode(ErrWeatherServiceComms.Code))
		}

		var faultErr fault.Error
		if err := json.Unmarshal(body, &faultErr); err != nil {
			return nil, fault.Wrap(err, ErrWeatherServiceComms.Message, fault.WithCode(ErrWeatherServiceComms.Code))
		}

		return nil, &faultErr
	}

	var weatherResponse model.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, fault.Wrap(err, ErrWeatherServiceComms.Message, fault.WithCode(ErrWeatherServiceComms.Code))
	}

	return &weatherResponse, nil
}
