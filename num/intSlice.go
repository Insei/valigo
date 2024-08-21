package num

import (
	"context"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type intSliceBuilder[T []int | *[]int] struct {
	h        shared.Helper
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
}

// Min checks if each integer in the slice has a minimum number.
func (s *intSliceBuilder[T]) Min(minNum int) IntSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch intSliceVal := value.(type) {
		case *[]int:
			for i, _ := range *intSliceVal {
				if (*intSliceVal)[i] < minNum {
					return []shared.Error{h.ErrorT(ctx, s.field, (*intSliceVal)[i], minLocaleKey, minNum)}
				}
			}
		case **[]int:
			if *intSliceVal == nil {
				return []shared.Error{h.ErrorT(ctx, s.field, *intSliceVal, minLocaleKey, minNum)}
			}
			for i, _ := range **intSliceVal {
				if (**intSliceVal)[i] < minNum {
					return []shared.Error{h.ErrorT(ctx, s.field, (**intSliceVal)[i], minLocaleKey, minNum)}
				}
			}
		}
		return nil
	})
	return s
}

// Max checks if each integer in the slice has a maximum number.
func (s *intSliceBuilder[T]) Max(maxNum int) IntSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch intSliceVal := value.(type) {
		case *[]int:
			for i, _ := range *intSliceVal {
				if (*intSliceVal)[i] > maxNum {
					return []shared.Error{h.ErrorT(ctx, s.field, (*intSliceVal)[i], maxLocaleKey, maxNum)}
				}
			}
		case **[]int:
			if *intSliceVal == nil {
				return []shared.Error{h.ErrorT(ctx, s.field, *intSliceVal, maxLocaleKey, maxNum)}
			}
			for i, _ := range **intSliceVal {
				if (**intSliceVal)[i] > maxNum {
					return []shared.Error{h.ErrorT(ctx, s.field, (**intSliceVal)[i], maxLocaleKey, maxNum)}
				}
			}
		}
		return nil
	})
	return s
}

// Required checks if the integer slice is not empty.
func (s *intSliceBuilder[T]) Required() IntSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch intSliceVal := value.(type) {
		case *[]int:
			if len(*intSliceVal) < 1 {
				return []shared.Error{h.ErrorT(ctx, s.field, intSliceVal, requiredLocaleKey)}
			}
		case **[]int:
			if *intSliceVal == nil || len(**intSliceVal) < 1 {
				return []shared.Error{h.ErrorT(ctx, s.field, *intSliceVal, requiredLocaleKey)}
			}
		}
		return nil
	})
	return s
}

// Custom allows for custom validation logic.
func (s *intSliceBuilder[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) IntSliceBuilder[T] {
	customHelper := shared.NewFieldCustomHelper(s.field, s.h)
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(*T))
	})
	return s
}

// When allows for conditional validation based on a given condition.
func (s *intSliceBuilder[T]) When(whenFn func(ctx context.Context, value *T) bool) IntSliceBuilder[T] {
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
