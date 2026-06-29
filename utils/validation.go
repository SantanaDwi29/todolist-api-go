package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError formats go-playground validator errors into user-friendly strings.
func FormatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var errs []string
		for _, fe := range ve {
			errs = append(errs, fmt.Sprintf("Field '%s' failed on '%s' validation", fe.Field(), fe.Tag()))
		}
		return strings.Join(errs, ", ")
	}
	return err.Error()
}
