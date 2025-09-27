package port

import "context"

type ViaCepClient interface {
	GetLocation(ctx context.Context, zipcode string) (string, error)
}

type WeatherApiClient interface {
	GetTemperature(ctx context.Context, city string) (float64, error)
}
