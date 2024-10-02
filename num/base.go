package num

import (
	"context"
	"fmt"
	"reflect"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
)

const (
	minLocaleKey          = "validation:num:Cannot be less than %v"
	maxLocaleKey          = "validation:num:Cannot be greater than %v"
	requiredLocaleKey     = "validation:num:Should be fulfilled"
	anyOfLocaleKey        = "validation:num:Only %v values is allowed"
	anyOfIntervalLocalKey = "validation:num:Only interval[%v - %v] is allowed"
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

var _ BaseConfigurator = &baseConfigurator[int]{}

type baseConfigurator[T numbers] struct {
	c     *shared.FieldConfigurator[T]
	field fmap.Field
	h     shared.Helper
}

// Max checks if the integer exceeds the maximum allowed number.
func (i *baseConfigurator[T]) Max(maxNum any) BaseConfigurator {
	if i.field.GetDereferencedType() != reflect.TypeOf(maxNum) {
		panic(fmt.Sprintf("field dereferenced type is %s, but maxNum type is %s", i.field.GetDereferencedType().String(), reflect.TypeOf(maxNum).String()))
	}
	i.c.Append(func(v T) bool {
		return maxT[T](v, maxNum.(T))
	}, maxLocaleKey, maxNum)
	return i
}

// Min checks if the integer is less than the minimum allowed number.
func (i *baseConfigurator[T]) Min(minNum any) BaseConfigurator {
	if i.field.GetDereferencedType() != reflect.TypeOf(minNum) {
		panic(fmt.Sprintf("field dereferenced type is %s, but minNum type is %s", i.field.GetDereferencedType().String(), reflect.TypeOf(minNum).String()))
	}
	i.c.Append(func(v T) bool {
		return minT[T](v, minNum.(T))
	}, minLocaleKey, minNum)
	return i
}

// Required checks if the integer value is not empty.
func (i *baseConfigurator[T]) Required() BaseConfigurator {
	i.c.Append(func(v T) bool {
		return true
	}, requiredLocaleKey)
	return i
}

// AnyOf checks if the integer value is one of the allowed values.
func (i *baseConfigurator[T]) AnyOf(allowed ...any) BaseConfigurator {
	if i.field.GetDereferencedType().Elem() != reflect.TypeOf(allowed).Elem() {
		panic(fmt.Sprintf("field dereferenced type is %s, but minNum type is %s", i.field.GetDereferencedType().String(), reflect.TypeOf(allowed).Elem().String()))
	}
	i.c.Append(func(v T) bool {
		return true
		//return anyOfT[T](v, allowed.([]T))
	}, anyOfLocaleKey, allowed)
	return i
}

// AnyOfInterval checks if the integer value is one of the allowed values intervals.
func (i *baseConfigurator[T]) AnyOfInterval(begin, end any) BaseConfigurator {
	if i.field.GetDereferencedType() != reflect.TypeOf(begin) ||
		i.field.GetDereferencedType() != reflect.TypeOf(end) {
		panic(fmt.Sprintf("field dereferenced type is %s, but begin and end type is %s", i.field.GetDereferencedType().String(), reflect.TypeOf(reflect.TypeOf(end)).String()))
	}
	i.c.Append(func(v T) bool {
		return anyOfIntervalT[T](v, begin.(T), end.(T))
	}, anyOfIntervalLocalKey, begin, end)
	return i
}

// Custom allows for custom validation logic to be applied to the integer value.
func (i *baseConfigurator[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error) BaseConfigurator {
	customHelper := shared.NewFieldCustomHelper(i.field, i.h)
	i.c.CustomAppend(func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value)
	})
	return i
}

// When allows for conditional validation logic to be applied to the integer value.
func (i *baseConfigurator[T]) When(whenFn func(ctx context.Context, value any) bool) BaseConfigurator {
	if whenFn == nil {
		return i
	}
	base := i.c.NewWithWhen(func(ctx context.Context, value any) bool {
		v, ok := value.(**T)
		if !ok {
			return false
		}
		return whenFn(ctx, v)
	})
	return &baseConfigurator[T]{
		c:     base,
		field: i.field,
		h:     i.h,
	}
}
