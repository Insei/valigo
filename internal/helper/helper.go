// Package helper builds the errors of the rules with the translated messages and the locations of the fields.
package helper

import (
	"context"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/translator"
)

// Helper is the shared.Helper of a validator.
type Helper struct {
	// Translator translates the messages.
	Translator translator.Translator
	// Location returns the location of a field in the errors.
	Location func(field fmap.Field) string
}

// ErrorT returns a shared.Error with the given location, message, and value.
// The message is translated using the translator, the code of the error is the locale key.
func (h *Helper) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	location := h.Location(field)
	msg := h.Translator.T(ctx, localeKey, args...)
	return shared.Error{
		Location: location,
		Message:  msg,
		Value:    value,
		Code:     localeKey,
	}
}

// FieldLocation returns the location of the field in errors.
func (h *Helper) FieldLocation(field fmap.Field) string {
	return h.Location(field)
}

// New returns a helper with the default translator and the struct paths of the fields as their locations.
func New() *Helper {
	return &Helper{
		Translator: translator.New(),
		Location: func(field fmap.Field) string {
			return field.GetStructPath()
		},
	}
}
