package service

import (
	"regexp"

	"github.com/marcelofabianov/fault"

	"github.com/marcelofabianov/weather-api/internal/model"
	"github.com/marcelofabianov/weather-api/internal/port"
)

var (
	ErrInvalidZipcode  = fault.New("invalid zipcode", fault.WithCode(fault.DomainViolation))
	ErrZipcodeNotFound = fault.New("can not find zipcode", fault.WithCode(fault.NotFound))
	ErrWeatherNotFound = fault.New("can not find weather for location", fault.WithCode(fault.Internal))
)

type WeatherService struct {
	viaCepClient     port.ViaCepClient
	weatherApiClient port.WeatherApiClient
}

func NewWeatherService(viaCepClient port.ViaCepClient, weatherApiClient port.WeatherApiClient) *WeatherService {
	return &WeatherService{
		viaCepClient:     viaCepClient,
		weatherApiClient: weatherApiClient,
	}
}

func (s *WeatherService) GetWeatherByZipcode(zipcode string) (*model.Weather, error) {
	if !s.isValidZipcode(zipcode) {
		return nil, ErrInvalidZipcode
	}

	city, err := s.viaCepClient.GetLocation(zipcode)
	if err != nil {
		if fault.IsNotFound(err) {
			return nil, ErrZipcodeNotFound
		}
		return nil, err
	}

	tempC, err := s.weatherApiClient.GetTemperature(city)
	if err != nil {
		return nil, fault.Wrap(err, ErrWeatherNotFound.Message, fault.WithCode(ErrWeatherNotFound.Code))
	}

	return model.NewWeather(city, tempC), nil
}

func (s *WeatherService) isValidZipcode(zipcode string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(zipcode)
}
