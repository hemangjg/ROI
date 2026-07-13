package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

// Validate runs struct tag validation.
func Validate[T any](value T) error {
	if err := v.Struct(value); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}