package num

import (
	"context"
	"reflect"
	"slices"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type uintBuilder[T uint8 | *uint8 | uint16 | *uint16 | uint32 | *uint32 | uint64 | *uint64 | uint | *uint] struct {
	h        shared.Helper
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
}

// Max checks if the uint exceeds the maximum allowed number.
func (i *uintBuilder[T]) Max(maxNum uint) UintBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			if uint(v.Uint()) > maxNum {
				return []shared.Error{h.ErrorT(ctx, i.field, uint(v.Uint()), maxLocaleKey, maxNum)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return i
}

// Min checks if the uint is less than the minimum allowed number.
func (i *uintBuilder[T]) Min(minNum uint) UintBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			if uint(v.Uint()) < minNum {
				return []shared.Error{h.ErrorT(ctx, i.field, uint(v.Uint()), minLocaleKey, minNum)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return i
}

// Required checks if the uint value is not empty.
func (i *uintBuilder[T]) Required() UintBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
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

// AnyOf checks if the uint value is one of the allowed values.
func (i *uintBuilder[T]) AnyOf(allowed ...uint) UintBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			if !slices.Contains(allowed, uint(v.Uint())) {
				return []shared.Error{h.ErrorT(ctx, i.field, uint(v.Uint()), anyOfLocaleKey, allowed)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return i
}

// AnyOfInterval checks if the uint value is one of the allowed values intervals.
func (i *uintBuilder[T]) AnyOfInterval(begin, end uint) UintBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			if !(uint(v.Uint()) > begin && uint(v.Uint()) < end) {
				return []shared.Error{h.ErrorT(ctx, i.field, uint(v.Uint()), anyOfInterval, begin, end)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, i.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return i
}

// Custom allows for custom validation logic to be applied to the uint value.
func (i *uintBuilder[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) UintBuilder[T] {
	customHelper := shared.NewFieldCustomHelper(i.field, i.h)
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(*T))
	})
	return i
}

// When allows for conditional validation logic to be applied to the uint value.
func (i *uintBuilder[T]) When(whenFn func(ctx context.Context, value *T) bool) UintBuilder[T] {
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
