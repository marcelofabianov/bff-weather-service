package port

type ViaCepClient interface {
	GetLocation(zipcode string) (string, error)
}

type WeatherApiClient interface {
	GetTemperature(city string) (float64, error)
}
