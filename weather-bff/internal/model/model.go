package model

type CepInput struct {
	Cep string `json:"cep" validate:"required,len=8,numeric"`
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}
