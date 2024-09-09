package num

import (
	"context"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
)

var _ PtrConfigurator[int] = &ptrConfigurator[int]{}

type ptrConfigurator[T numbers] struct {
	*baseConfigurator[T]
}

func (i *ptrConfigurator[T]) Required() PtrConfigurator[T] {
	i.baseConfigurator.Required()
	return i
}

func (i *ptrConfigurator[T]) AnyOf(allowed ...T) PtrConfigurator[T] {
	i.baseConfigurator.AnyOf(allowed...)
	return i
}

func (i *ptrConfigurator[T]) AnyOfInterval(begin, end T) PtrConfigurator[T] {
	i.baseConfigurator.AnyOfInterval(begin, end)
	return i
}

func (i *ptrConfigurator[T]) Max(val T) PtrConfigurator[T] {
	i.baseConfigurator.Max(val)
	return i
}

func (i *ptrConfigurator[T]) Min(val T) PtrConfigurator[T] {
	i.baseConfigurator.Min(val)
	return i
}

// Custom allows for custom validation logic to be applied to the integer value.
func (i *ptrConfigurator[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value **T) []shared.Error) PtrConfigurator[T] {
	customHelper := shared.NewFieldCustomHelper(i.field, i.h)
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(**T))
	})
	return i
}

// When allows for conditional validation logic to be applied to the integer value.
func (i *ptrConfigurator[T]) When(whenFn func(ctx context.Context, value **T) bool) PtrConfigurator[T] {
	if whenFn == nil {
		return i
	}
	i.appendFn = func(field fmap.Field, fn shared.FieldValidationFn) {
		fnWithEnabler := func(ctx context.Context, h shared.Helper, v any) []shared.Error {
			if !whenFn(ctx, v.(**T)) {
				return nil
			}
			return fn(ctx, h, v)
		}
		i.appendFn(field, fnWithEnabler)
	}
	return i
}
