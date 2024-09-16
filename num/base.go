package num

import (
	"context"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
)

const (
	minLocaleKey          = "validation:num:Cannot be less than %d characters"
	maxLocaleKey          = "validation:num:Cannot be greater than %d characters"
	requiredLocaleKey     = "validation:num:Should be fulfilled"
	anyOfLocaleKey        = "validation:num:Only %d values is allowed"
	anyOfIntervalLocalKey = "validation:num:Only interval[%d - %d] is allowed"
	invalidLocaleKey      = "validation:num:Invalid value"
)

func minT[T numbers](val T, min T) bool {
	return val > min
}

func maxT[T numbers](val T, max T) bool {
	return val < max
}

func anyOfT[T numbers](val T, allowed []T) bool {
	for _, v := range allowed {
		if v == val {
			return true
		}
	}
	return false
}

func anyOfIntervalT[T numbers](val T, begin, end T) bool {
	return end > val && val > begin
}

var _ BaseConfigurator[int, *baseConfigurator[int]] = &baseConfigurator[int]{}

type baseConfigurator[T numbers] struct {
	h           shared.Helper
	field       fmap.Field
	appendFn    func(field fmap.Field, fn shared.FieldValidationFn)
	dereference func(value any) (T, bool)
}

// Max checks if the integer exceeds the maximum allowed number.
func (i *baseConfigurator[T]) Max(maxNum T) *baseConfigurator[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v, isValid := i.dereference(value)
		if isValid {
			isValid = maxT[T](v, maxNum)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, maxLocaleKey, maxNum)}
		}
		return nil
	})
	return i
}

// Min checks if the integer is less than the minimum allowed number.
func (i *baseConfigurator[T]) Min(minNum T) *baseConfigurator[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v, isValid := i.dereference(value)
		if isValid {
			isValid = minT[T](v, minNum)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, minLocaleKey, minNum)}
		}
		return nil
	})
	return i
}

// Required checks if the integer value is not empty.
func (i *baseConfigurator[T]) Required() *baseConfigurator[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v, isValid := i.dereference(value)
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, requiredLocaleKey)}
		}
		return nil
	})
	return i
}

// AnyOf checks if the integer value is one of the allowed values.
func (i *baseConfigurator[T]) AnyOf(allowed ...T) *baseConfigurator[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v, isValid := i.dereference(value)
		if isValid {
			isValid = anyOfT[T](v, allowed)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, anyOfLocaleKey, allowed)}
		}
		return nil
	})
	return i
}

// AnyOfInterval checks if the integer value is one of the allowed values intervals.
func (i *baseConfigurator[T]) AnyOfInterval(begin, end T) *baseConfigurator[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v, isValid := i.dereference(value)
		if isValid {
			isValid = anyOfIntervalT[T](v, begin, end)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, anyOfIntervalLocalKey, begin, end)}
		}
		return nil
	})
	return i
}
