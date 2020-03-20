package query

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

// Define a single validator to do all of the validations for us.
var v = validator.New()

// ValidatedQueryParam extracts a query parameter and validates it.
func ValidatedQueryParam(ctx echo.Context, name, validationTag string) (string, error) {
	value := ctx.QueryParam(name)

	// Validate the value.
	if err := v.Var(value, validationTag); err != nil {
		return "", err
	}

	return value, nil
}
