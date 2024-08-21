package num

import (
	"context"
	"reflect"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type uintSliceBuilder[T []uint8 | *[]uint8 | []uint16 | *[]uint16 | []uint32 | *[]uint32 | []uint64 | *[]uint64 | []uint | *[]uint] struct {
	h        shared.Helper
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
}

// Min checks if each uint in the slice has a minimum number.
func (s *uintSliceBuilder[T]) Min(minNum uint) UintSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				if uint(v.Index(i).Uint()) < minNum {
					return []shared.Error{h.ErrorT(ctx, s.field, v.Index(i).Uint(), minLocaleKey, minNum)}
				}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, s.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return s
}

// Max checks if each uint in the slice has a maximum number.
func (s *uintSliceBuilder[T]) Max(maxNum uint) UintSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				if uint(v.Index(i).Uint()) > maxNum {
					return []shared.Error{h.ErrorT(ctx, s.field, v.Index(i).Uint(), maxLocaleKey, maxNum)}
				}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, s.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return s
}

// Required checks if the uint slice is not empty.
func (s *uintSliceBuilder[T]) Required() UintSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			if v.Len() < 1 {
				return []shared.Error{h.ErrorT(ctx, s.field, v, requiredLocaleKey)}
			}
		default:
			return []shared.Error{h.ErrorT(ctx, s.field, "", invalidLocaleKey)}
		}
		return nil
	})
	return s
}

// Custom allows for custom validation logic.
func (s *uintSliceBuilder[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) UintSliceBuilder[T] {
	customHelper := shared.NewFieldCustomHelper(s.field, s.h)
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(*T))
	})
	return s
}

// When allows for conditional validation based on a given condition.
func (s *uintSliceBuilder[T]) When(whenFn func(ctx context.Context, value *T) bool) UintSliceBuilder[T] {
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
