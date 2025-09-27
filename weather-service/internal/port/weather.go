package port

import "github.com/marcelofabianov/weather-api/internal/model"

type WeatherService interface {
	GetWeatherByZipcode(zipcode string) (*model.Weather, error)
}
