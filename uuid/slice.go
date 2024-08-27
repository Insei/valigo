package uuid

import (
	"context"
	"github.com/google/uuid"
	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
)

type uuidSliceBuilder[T []uuid.UUID | *[]uuid.UUID] struct {
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	h        shared.Helper
}

// Required checks if the string slice is not empty.
func (s *uuidSliceBuilder[T]) Required() SliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch uuidSliceVal := value.(type) {
		case *[]uuid.UUID:
			if *uuidSliceVal == nil {
				return []shared.Error{h.ErrorT(ctx, s.field, uuidSliceVal, requiredLocaleKey)}
			}
			for i, _ := range *uuidSliceVal {
				if (*uuidSliceVal)[i] == uuid.Nil {
					return []shared.Error{h.ErrorT(ctx, s.field, (*uuidSliceVal)[i], requiredLocaleKey)}
				}
			}
		case **[]uuid.UUID:
			if *uuidSliceVal == nil {
				return []shared.Error{h.ErrorT(ctx, s.field, uuidSliceVal, requiredLocaleKey)}
			}
			for i, _ := range **uuidSliceVal {
				if (**uuidSliceVal)[i] == uuid.Nil {
					return []shared.Error{h.ErrorT(ctx, s.field, (**uuidSliceVal)[i], requiredLocaleKey)}
				}
			}
		}
		return nil
	})
	return s
}

// Custom allows for custom validation logic.
func (s *uuidSliceBuilder[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) SliceBuilder[T] {
	customHelper := shared.NewFieldCustomHelper(s.field, s.h)
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(*T))
	})
	return s
}

// When allows for conditional validation based on a given condition.
func (s *uuidSliceBuilder[T]) When(whenFn func(ctx context.Context, value *T) bool) SliceBuilder[T] {
	if whenFn == nil {
		return s
	}
	s.appendFn = func(field fmap.Field, fn shared.FieldValidationFn) {
		fnWithEnabler := func(ctx context.Context, h shared.Helper, v any) []shared.Error {
			if !whenFn(ctx, v.(*T)) {
				return nil
			}
			return fn(ctx, h, v)
		}
		s.appendFn(field, fnWithEnabler)
	}
	return s
}
