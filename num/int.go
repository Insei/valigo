package num

import (
	"context"
	"reflect"
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
	invalidLocaleKey  = "validation:int:Invalid value"
)

type intBuilder[T int | *int | int8 | *int8 | int16 | *int16 | int32 | *int32 | int64 | *int64] struct {
	h        shared.Helper
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
}

func maxT[T int | int8 | uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64](valuePtr *T, valuePtr2 **T, max int) bool {
	if valuePtr == nil && *valuePtr2 == nil {
		return true
	}
	if valuePtr != nil {
		return *valuePtr > T(0)
	}
	return **valuePtr2 > T(0)
}

// Max checks if the integer exceeds the maximum allowed number.
func (i *intBuilder[T]) Max(maxNum int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		switch typed := value.(type) {
		case *int:
			isValid = maxT(typed, nil, maxNum)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, value, maxLocaleKey, maxNum)}
		}
		//v := reflect.ValueOf(value)
		//if v.Kind() == reflect.Ptr {
		//	v = v.Elem()
		//}
		//if v.Kind() == reflect.Ptr {
		//	v = v.Elem()
		//}
		//switch v.Kind() {
		//case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		//	if int(v.Int()) > maxNum {
		//		return []shared.Error{h.ErrorT(ctx, i.field, int(v.Int()), maxLocaleKey, maxNum)}
		//	}
		//default:
		//	return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
		//}
		//return nil
	})
	return i
}

// Min checks if the integer is less than the minimum allowed number.
func (i *intBuilder[T]) Min(minNum int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if int(v.Int()) < minNum {
				return []shared.Error{h.ErrorT(ctx, i.field, int(v.Int()), minLocaleKey, minNum)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return i
}

// Required checks if the integer value is not empty.
func (i *intBuilder[T]) Required() IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if v.IsZero() {
				return []shared.Error{h.ErrorT(ctx, i.field, "", requiredLocaleKey)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return i
}

// AnyOf checks if the integer value is one of the allowed values.
func (i *intBuilder[T]) AnyOf(allowed ...int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if !slices.Contains(allowed, int(v.Int())) {
				return []shared.Error{h.ErrorT(ctx, i.field, int(v.Int()), anyOfLocaleKey, allowed)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return i
}

// AnyOfInterval checks if the integer value is one of the allowed values intervals.
func (i *intBuilder[T]) AnyOfInterval(begin, end int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if !(int(v.Int()) > begin && int(v.Int()) < end) {
				return []shared.Error{h.ErrorT(ctx, i.field, int(v.Int()), anyOfInterval, begin, end)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
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
