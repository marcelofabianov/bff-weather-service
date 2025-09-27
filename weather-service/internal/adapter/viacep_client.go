package adapter

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jpillora/backoff"
	"github.com/marcelofabianov/fault"
	"github.com/sony/gobreaker"

	"github.com/marcelofabianov/weather-api/config"
)

type ViaCepResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

type ViaCepClient struct {
	BaseURL    string
	breaker    *gobreaker.CircuitBreaker
	resilience *config.ResilienceConfig
	logger     *slog.Logger
}

func NewViaCepClient(
	breaker *gobreaker.CircuitBreaker,
	resilienceCfg *config.ResilienceConfig,
	logger *slog.Logger,
) *ViaCepClient {
	return &ViaCepClient{
		BaseURL:    "https://viacep.com.br/ws",
		breaker:    breaker,
		resilience: resilienceCfg,
		logger:     logger.With("adapter", "viacep_client"),
	}
}

func (c *ViaCepClient) GetLocation(zipcode string) (string, error) {
	body, err := c.breaker.Execute(func() (interface{}, error) {
		b := &backoff.Backoff{
			Min:    c.resilience.RetryInitialBackoff,
			Max:    c.resilience.RetryMaxBackoff,
			Factor: 2,
			Jitter: true,
		}
		var lastErr error

		for i := 0; i < c.resilience.RetryMaxAttempts; i++ {
			requestURL := fmt.Sprintf("%s/%s/json/", c.BaseURL, zipcode)
			resp, err := http.Get(requestURL)

			if err != nil {
				lastErr = fault.Wrap(err, ErrExternalAPICall.Message, fault.WithCode(ErrExternalAPICall.Code), fault.WithContext("url", requestURL))
			} else {
				if resp.StatusCode == http.StatusOK {
					var data ViaCepResponse
					if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
						resp.Body.Close()
						if data.Erro {
							return nil, ErrLocationNotFound
						}
						return &data, nil
					}
					lastErr = ErrExternalAPIParse
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
		if fault.IsNotFound(err) {
			return "", err
		}
		c.logger.Error("Failed to get location from ViaCEP after all retries", "error", err)
		return "", err
	}

	data := body.(*ViaCepResponse)
	return data.Localidade, nil
}
