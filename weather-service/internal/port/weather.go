package port

import (
	"context"
	"github.com/marcelofabianov/weather-server/internal/model"
)

type WeatherService interface {
	GetWeatherByZipcode(ctx context.Context, zipcode string) (*model.Weather, error)
}
