package adapter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/marcelofabianov/weather-bff/config"
)

type WeatherServiceClient struct {
	BaseURL string
	client  *http.Client
}

func NewWeatherServiceClient(cfg *config.WeatherServiceConfig) *WeatherServiceClient {
	return &WeatherServiceClient{
		BaseURL: cfg.URL,
		client:  &http.Client{},
	}
}

func (c *WeatherServiceClient) GetWeatherForCep(ctx context.Context, cep string) (*http.Response, error) {
	requestURL := fmt.Sprintf("%s/weather/%s", c.BaseURL, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
