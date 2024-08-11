package shared

import (
	"context"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/translator"
)

type Helper struct {
	translator.Translator
	GetFieldLocation func(field fmap.Field) string
}

func (h *Helper) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) Error {
	location := h.GetFieldLocation(field)
	msg := h.Translator.T(ctx, localeKey, args...)
	return Error{
		Location: location,
		Message:  msg,
		Value:    value,
	}
}

func NewHelper() *Helper {
	return &Helper{
		Translator: translator.New(),
		GetFieldLocation: func(field fmap.Field) string {
			return field.GetStructPath()
		},
	}
}
