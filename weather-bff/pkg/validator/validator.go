package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/marcelofabianov/fault"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(data any) error {
	err := v.validate.Struct(data)
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return fault.NewValidationErrorFromValidator(validationErrs)
	}

	return fault.NewInternalError(err, nil)
}
