package shared

import (
	"context"

	"github.com/insei/fmap/v3"
)

// Helper interface defines a contract for error handling.
type Helper interface {
	// ErrorT method returns an error based on the provided context, field, value, locale key, and arguments.
	ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) Error
}

// StructCustomHelper interface defines a contract for custom error handling in structs.
type StructCustomHelper interface {
	// ErrorT method returns an error based on the provided context,
	// pointer to field value, field value, locale key, and arguments.
	ErrorT(ctx context.Context, ptrToFieldValue, fieldValue any, localeKey string, args ...any) Error
}

// FieldCustomHelper struct implements the Helper interface and provides custom error handling for fields.
type FieldCustomHelper struct {
	// h is the underlying Helper implementation.
	h Helper
	// field is the fmap.Field instance associated with this helper.
	field fmap.Field
}

// ErrorT method implements the Helper interface and delegates
// error handling to the underlying Helper implementation.
func (h *FieldCustomHelper) ErrorT(ctx context.Context, value any, localeKey string, args ...any) Error {
	return h.h.ErrorT(ctx, h.field, value, localeKey, args...)
}

// NewFieldCustomHelper function creates a new FieldCustomHelper
// instance with the provided field and Helper implementation.
func NewFieldCustomHelper(field fmap.Field, h Helper) *FieldCustomHelper {
	return &FieldCustomHelper{
		field: field,
		h:     h,
	}
}
