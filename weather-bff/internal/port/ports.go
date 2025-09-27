package port

import (
	"context"
	"net/http"

	"github.com/marcelofabianov/weather-bff/internal/model"
)

type BffService interface {
	GetWeather(ctx context.Context, cep string) (*model.WeatherResponse, error)
}

type WeatherServiceClient interface {
	GetWeatherForCep(ctx context.Context, cep string) (*http.Response, error)
}

type Validator interface {
	Validate(data any) error
}
