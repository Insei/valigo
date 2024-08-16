package num

import (
	"context"
	"slices"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

const (
	minLocaleKey      = "validation:int:Cannot be less than %d characters"
	maxLocaleKey      = "validation:int:Cannot be greater than %d characters"
	requiredLocaleKey = "validation:int:Should be fulfilled"
	anyOfLocaleKey    = "validation:int:Only %d values is allowed"
	anyOfInterval     = "validation:int:Only interval[%d - %d] is allowed"
)

type intBuilder[T int | *int] struct {
	h        shared.Helper
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
}

// Max checks if the integer exceeds the maximum allowed number.
func (i *intBuilder[T]) Max(maxNum int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch intVal := value.(type) {
		case *int:
			if *intVal > maxNum {
				return []shared.Error{h.ErrorT(ctx, i.field, *intVal, maxLocaleKey, maxNum)}
			}
		case **int:
			if *intVal == nil || **intVal > maxNum {
				return []shared.Error{h.ErrorT(ctx, i.field, "", maxLocaleKey, maxNum)}
			}
		}
		return nil
	})
	return i
}

// Min checks if the integer is less than the minimum allowed number.
func (i *intBuilder[T]) Min(minNum int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch intVal := value.(type) {
		case *int:
			if *intVal < minNum {
				return []shared.Error{h.ErrorT(ctx, i.field, *intVal, minLocaleKey, minNum)}
			}
		case **int:
			if *intVal == nil || **intVal < minNum {
				return []shared.Error{h.ErrorT(ctx, i.field, "", minLocaleKey, minNum)}
			}
		}
		return nil
	})
	return i
}

// Required checks if the integer value is not empty.
func (i *intBuilder[T]) Required() IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch intVal := value.(type) {
		case *int:
			if *intVal == 0 {
				return []shared.Error{h.ErrorT(ctx, i.field, *intVal, requiredLocaleKey)}
			}
		case **int:
			if intVal == nil || *intVal == nil {
				return []shared.Error{h.ErrorT(ctx, i.field, "", requiredLocaleKey)}
			}
		}
		return nil
	})
	return i
}

// AnyOf checks if the integer value is one of the allowed values.
func (i *intBuilder[T]) AnyOf(allowed ...int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch intVal := value.(type) {
		case *int:
			if !slices.Contains(allowed, *intVal) {
				return []shared.Error{h.ErrorT(ctx, i.field, *intVal, anyOfLocaleKey, allowed)}
			}
		case **int:
			if *intVal == nil || !slices.Contains(allowed, **intVal) {
				return []shared.Error{h.ErrorT(ctx, i.field, "", anyOfLocaleKey, allowed)}
			}
		}
		return nil
	})
	return i
}

// AnyOfInterval checks if the integer value is one of the allowed values intervals.
func (i *intBuilder[T]) AnyOfInterval(begin, end int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch intVal := value.(type) {
		case *int:
			if !(*intVal > begin && *intVal < end) {
				return []shared.Error{h.ErrorT(ctx, i.field, *intVal, anyOfInterval, begin, end)}
			}
		case **int:
			if *intVal == nil || !(**intVal > begin && **intVal < end) {
				return []shared.Error{h.ErrorT(ctx, i.field, "", anyOfInterval, begin, end)}
			}
		}
		return nil
	})
	return i
}

// Custom allows for custom validation logic to be applied to the integer value.
func (i *intBuilder[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) IntBuilder[T] {
	customHelper := shared.NewFieldCustomHelper(i.field, i.h)
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(*T))
	})
	return i
}

// When allows for conditional validation logic to be applied to the integer value.
func (i *intBuilder[T]) When(whenFn func(ctx context.Context, value *T) bool) IntBuilder[T] {
	if whenFn == nil {
		return i
	}
	i.appendFn = func(field fmap.Field, fn shared.FieldValidationFn) {
		fnWithEnabler := func(ctx context.Context, h shared.Helper, v any) []shared.Error {
			if !whenFn(ctx, v.(*T)) {
				return nil
			}
			return fn(ctx, h, v)
		}
		i.appendFn(field, fnWithEnabler)
	}
	return i
}
