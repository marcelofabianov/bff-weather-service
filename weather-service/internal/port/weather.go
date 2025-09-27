package port

import "github.com/marcelofabianov/weather-server/internal/model"

type WeatherService interface {
	GetWeatherByZipcode(zipcode string) (*model.Weather, error)
}
