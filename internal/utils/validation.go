package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ValidateStruct(v any) error {
	if err := Validate.Struct(v); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
