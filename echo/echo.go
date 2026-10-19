// Package echo validates the requests of echo (github.com/labstack/echo/v4) by valigo.
//
// New is the echo.Validator, c.Validate validates by valigo:
//
//	e.Validator = vecho.New(validator)
//
// echo doesn't pass the request to the validator, so its messages are in the default language;
// Validate validates with the languages of the Accept-Language header of the request.
package echo

import (
	"context"

	decho "github.com/labstack/echo/v4"

	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/translator"
)

// Validator is the echo.Validator by valigo.
type Validator struct {
	v *valigo.Validator
}

// New returns the echo.Validator by the validator.
func New(v *valigo.Validator) *Validator {
	return &Validator{v: v}
}

// Validate validates i, it returns shared.Errors or nil.
func (e *Validator) Validate(i any) error {
	return errorOf(e.v.ValidateTyped(context.Background(), i))
}

// Validate validates obj with the languages of the Accept-Language header of the request, it returns shared.Errors or nil.
func Validate(c decho.Context, v *valigo.Validator, obj any) error {
	r := c.Request()
	ctx := translator.ContextWithAcceptLanguage(r.Context(), r.Header.Get("Accept-Language"))
	return errorOf(v.ValidateTyped(ctx, obj))
}

func errorOf(errs []shared.Error) error {
	if len(errs) == 0 {
		return nil
	}
	return shared.Errors(errs)
}
