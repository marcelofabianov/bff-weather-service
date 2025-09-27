package adapter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/jpillora/backoff"
	"github.com/marcelofabianov/fault"
	"github.com/sony/gobreaker"

	"github.com/marcelofabianov/weather-api/config"
)

type WeatherApiResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherApiClient struct {
	BaseURL    string
	ApiKey     string
	breaker    *gobreaker.CircuitBreaker
	resilience *config.ResilienceConfig
	logger     *slog.Logger
}

func NewWeatherApiClient(
	clientCfg *config.WeatherAPIConfig,
	breaker *gobreaker.CircuitBreaker,
	resilienceCfg *config.ResilienceConfig,
	logger *slog.Logger,
) *WeatherApiClient {
	return &WeatherApiClient{
		BaseURL:    clientCfg.URL,
		ApiKey:     clientCfg.Key,
		breaker:    breaker,
		resilience: resilienceCfg,
		logger:     logger.With("adapter", "weatherapi_client"),
	}
}

func (c *WeatherApiClient) GetTemperature(city string) (float64, error) {
	body, err := c.breaker.Execute(func() (interface{}, error) {
		b := &backoff.Backoff{
			Min:    c.resilience.RetryInitialBackoff,
			Max:    c.resilience.RetryMaxBackoff,
			Factor: 2,
			Jitter: true,
		}
		var lastErr error

		for i := 0; i < c.resilience.RetryMaxAttempts; i++ {
			encodedCity := url.QueryEscape(city)
			requestURL := fmt.Sprintf("%s/current.json?key=%s&q=%s", c.BaseURL, c.ApiKey, encodedCity)
			resp, err := http.Get(requestURL)

			if err != nil {
				lastErr = fault.Wrap(err, ErrExternalAPICall.Message, fault.WithCode(ErrExternalAPICall.Code), fault.WithContext("url", requestURL))
			} else {
				if resp.StatusCode == http.StatusOK {
					var data WeatherApiResponse
					if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
						resp.Body.Close()
						return &data, nil
					}
					lastErr = fault.Wrap(err, ErrExternalAPIParse.Message, fault.WithCode(ErrExternalAPIParse.Code))
				} else {
					lastErr = fault.New(
						ErrExternalAPIUnacceptableStatusCode.Message,
						fault.WithCode(ErrExternalAPIUnacceptableStatusCode.Code),
						fault.WithContext("url", requestURL),
						fault.WithContext("status_code", resp.StatusCode),
					)
				}
				resp.Body.Close()
			}

			if i < c.resilience.RetryMaxAttempts-1 {
				duration := b.Duration()
				c.logger.Warn(
					"API call attempt failed, retrying...",
					"attempt", i+1,
					"error", lastErr.Error(),
					"backoff", duration.String(),
				)
				time.Sleep(duration)
			}
		}
		return nil, lastErr
	})

	if err != nil {
		c.logger.Error("Failed to get temperature from WeatherAPI after all retries", "error", err)
		return 0, err
	}

	data := body.(*WeatherApiResponse)
	return data.Current.TempC, nil
}
