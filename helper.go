package valigo

import (
	"context"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/translator"
)

// helper is a struct that provides functionality for handling errors and translations.
type helper struct {
	t                translator.Translator
	getFieldLocation func(field fmap.Field) string
}

// ErrorT returns a shared.Error with the given location, message, and value.
// The message is translated using the translator.
func (h *helper) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	location := h.getFieldLocation(field)
	msg := h.t.T(ctx, localeKey, args...)
	return shared.Error{
		Location: location,
		Message:  msg,
		Value:    value,
	}
}

// newHelper returns a new helper with a default translator and getFieldLocation function.
func newHelper() *helper {
	return &helper{
		t: translator.New(),
		getFieldLocation: func(field fmap.Field) string {
			return field.GetStructPath()
		},
	}
}
