package uuid

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
	"slices"
	"strings"
)

const (
	requiredLocaleKey = "validation:uuid:Should be fulfilled"
	anyOfLocaleKey    = "validation:uuid:Only %s values is allowed"
)

type uuidBuilder[T uuid.UUID | *uuid.UUID] struct {
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	h        shared.Helper
}

func (s *uuidBuilder[T]) Required() Builder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch uuidVal := value.(type) {
		case *uuid.UUID:
			if *uuidVal == uuid.Nil {
				return []shared.Error{h.ErrorT(ctx, s.field, *uuidVal, requiredLocaleKey)}
			}
		case **uuid.UUID:
			if *uuidVal == nil {
				return []shared.Error{h.ErrorT(ctx, s.field, *uuidVal, requiredLocaleKey)}
			}
			if **uuidVal == uuid.Nil {
				return []shared.Error{h.ErrorT(ctx, s.field, **uuidVal, requiredLocaleKey)}
			}
		}

		return nil
	})

	return s
}

// AnyOf checks if the string value is one of the allowed values.
func (s *uuidBuilder[T]) AnyOf(allowed ...uuid.UUID) Builder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch uuidVal := value.(type) {
		case *uuid.UUID:
			if !slices.Contains(allowed, *uuidVal) {
				allowedStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(allowed)), ","), "[]")
				return []shared.Error{h.ErrorT(ctx, s.field, *uuidVal, anyOfLocaleKey, allowedStr)}
			}
		case **uuid.UUID:
			if *uuidVal == nil || !slices.Contains(allowed, **uuidVal) {
				allowedStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(allowed)), ","), "[]")
				return []shared.Error{h.ErrorT(ctx, s.field, "", anyOfLocaleKey, allowedStr)}
			}
		}
		return nil
	})
	return s
}

// Custom allows for custom validation logic to be applied to the string value.
func (s *uuidBuilder[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) Builder[T] {
	customHelper := shared.NewFieldCustomHelper(s.field, s.h)
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(*T))
	})
	return s
}

// When allows for conditional validation logic to be applied to the string value.
func (s *uuidBuilder[T]) When(whenFn func(ctx context.Context, value *T) bool) Builder[T] {
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
