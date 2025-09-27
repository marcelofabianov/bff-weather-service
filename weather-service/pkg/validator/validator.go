package validator

import (
	"errors"
	"fmt"

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

func (v *Validator) Validate(data any) *fault.Error {
	err := v.validate.Struct(data)
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return fault.NewInternalError(err, nil)
	}

	details := make([]*fault.Error, 0, len(validationErrs))
	for _, fieldErr := range validationErrs {
		details = append(details, fault.New(
			fmt.Sprintf("validation failed on field '%s'", fieldErr.Field()),
			fault.WithCode(fault.Invalid),
			fault.WithContext("field", fieldErr.Field()),
			fault.WithContext("tag", fieldErr.Tag()),
			fault.WithContext("param", fieldErr.Param()),
		))
	}

	return fault.New(
		"Request validation failed",
		fault.WithWrappedErr(fault.ErrValidation),
		fault.WithCode(fault.Invalid),
		fault.WithDetails(details...),
	)
}
