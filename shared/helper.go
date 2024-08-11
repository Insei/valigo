package shared

import (
	"context"

	"github.com/insei/fmap/v3"
)

type Helper interface {
	ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) Error
}

type StructCustomHelper interface {
	ErrorT(ctx context.Context, ptrToFieldValue, fieldValue any, localeKey string, args ...any) Error
}

type FieldCustomHelper struct {
	h     Helper
	field fmap.Field
}

func (h *FieldCustomHelper) ErrorT(ctx context.Context, value any, localeKey string, args ...any) Error {
	return h.h.ErrorT(ctx, h.field, value, localeKey, args...)
}

func NewFieldCustomHelper(field fmap.Field, h Helper) *FieldCustomHelper {
	return &FieldCustomHelper{
		field: field,
		h:     h,
	}
}
