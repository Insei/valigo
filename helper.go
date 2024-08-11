package valigo

import (
	"context"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/translator"
)

type helper struct {
	t                translator.Translator
	getFieldLocation func(field fmap.Field) string
}

func (h *helper) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	location := h.getFieldLocation(field)
	msg := h.t.T(ctx, localeKey, args...)
	return shared.Error{
		Location: location,
		Message:  msg,
		Value:    value,
	}
}

func newHelper() *helper {
	return &helper{
		t: translator.New(),
		getFieldLocation: func(field fmap.Field) string {
			return field.GetStructPath()
		},
	}
}
