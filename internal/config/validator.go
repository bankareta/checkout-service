package config

import (
	"checkout-service/internal/environtment"

	"github.com/go-playground/validator/v10"
)

func NewValidator(viper *environtment.Config) *validator.Validate {
	return validator.New()
}
