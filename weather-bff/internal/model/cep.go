package model

type CepInput struct {
	Cep string `json:"cep" validate:"required,len=8,numeric"`
}
