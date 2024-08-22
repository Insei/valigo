package num

import (
	"context"

	"golang.org/x/exp/constraints"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type intSliceBuilder[T []int | *[]int | []int8 | *[]int8 | []int16 | *[]int16 | []int32 | *[]int32 | []int64 | *[]int64 |
	[]uint | *[]uint | []uint8 | *[]uint8 | []uint16 | *[]uint16 | []uint32 | *[]uint32 | []uint64 | *[]uint64] struct {
	h        shared.Helper
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
}

func compare[T constraints.Ordered](a T, b T) bool {
	return a > b
}

func maxSliceT[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](valuePtr *[]T, valuePtr2 **[]T, max int) (any, bool) {
	if (valuePtr == nil || len(*valuePtr) == 0) && (*valuePtr2 == nil || len(**valuePtr2) == 0) {
		return nil, false
	}
	if valuePtr != nil {
		for i := 0; i < len(*valuePtr); i++ {
			if compare((*valuePtr)[i], T(max)) {
				return (*valuePtr)[i], false
			}
		}
		return nil, true
	}
	for i := 0; i < len(**valuePtr2); i++ {
		if compare((**valuePtr2)[i], T(max)) {
			return (**valuePtr2)[i], false
		}
	}
	return nil, true
}

func minSliceT[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](valuePtr *[]T, valuePtr2 **[]T, min int) (any, bool) {
	if (valuePtr == nil || len(*valuePtr) == 0) && (*valuePtr2 == nil || len(**valuePtr2) == 0) {
		return nil, false
	}
	if valuePtr != nil {
		for i := 0; i < len(*valuePtr); i++ {
			if compare(T(min), (*valuePtr)[i]) {
				return (*valuePtr)[i], false
			}
		}
		return nil, true
	}
	for i := 0; i < len(**valuePtr2); i++ {
		if compare(T(min), (**valuePtr2)[i]) {
			return (**valuePtr2)[i], false
		}
	}
	return nil, true
}

func requiredSliceT[T ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](valuePtr *[]T, valuePtr2 **[]T) bool {
	if valuePtr == nil && *valuePtr2 == nil {
		return false
	}
	if valuePtr != nil {
		return len(*valuePtr) > 0
	}
	if *valuePtr2 == nil || len(**valuePtr2) == 0 {
		return false
	}
	return true
}

// Max checks if each integer in the slice has a maximum number.
func (s *intSliceBuilder[T]) Max(maxNum int) IntSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		var v any
		switch typed := value.(type) {
		case *[]int:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]int:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]int8:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]int8:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]int16:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]int16:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]int32:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]int32:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]int64:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]int64:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]uint:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]uint:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]uint8:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]uint8:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]uint16:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]uint16:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]uint32:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]uint32:
			v, isValid = maxSliceT(nil, typed, maxNum)
		case *[]uint64:
			v, isValid = maxSliceT(typed, nil, maxNum)
		case **[]uint64:
			v, isValid = maxSliceT(nil, typed, maxNum)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, s.field, v, maxLocaleKey, maxNum)}
		}
		return nil
	})
	return s
}

// Min checks if each integer in the slice has a minimum number.
func (s *intSliceBuilder[T]) Min(minNum int) IntSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		var v any
		switch typed := value.(type) {
		case *[]int:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]int:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]int8:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]int8:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]int16:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]int16:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]int32:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]int32:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]int64:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]int64:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]uint:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]uint:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]uint8:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]uint8:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]uint16:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]uint16:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]uint32:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]uint32:
			v, isValid = minSliceT(nil, typed, minNum)
		case *[]uint64:
			v, isValid = minSliceT(typed, nil, minNum)
		case **[]uint64:
			v, isValid = minSliceT(nil, typed, minNum)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, s.field, v, minLocaleKey, minNum)}
		}
		return nil
	})
	return s
}

// Required checks if the integer slice is not empty.
func (s *intSliceBuilder[T]) Required() IntSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		var v any
		switch typed := value.(type) {
		case *[]int:
			isValid = requiredSliceT(typed, nil)
		case **[]int:
			isValid = requiredSliceT(nil, typed)
		case *[]int8:
			isValid = requiredSliceT(typed, nil)
		case **[]int8:
			isValid = requiredSliceT(nil, typed)
		case *[]int16:
			isValid = requiredSliceT(typed, nil)
		case **[]int16:
			isValid = requiredSliceT(nil, typed)
		case *[]int32:
			isValid = requiredSliceT(typed, nil)
		case **[]int32:
			isValid = requiredSliceT(nil, typed)
		case *[]int64:
			isValid = requiredSliceT(typed, nil)
		case **[]int64:
			isValid = requiredSliceT(nil, typed)
		case *[]uint:
			isValid = requiredSliceT(typed, nil)
		case **[]uint:
			isValid = requiredSliceT(nil, typed)
		case *[]uint8:
			isValid = requiredSliceT(typed, nil)
		case **[]uint8:
			isValid = requiredSliceT(nil, typed)
		case *[]uint16:
			isValid = requiredSliceT(typed, nil)
		case **[]uint16:
			isValid = requiredSliceT(nil, typed)
		case *[]uint32:
			isValid = requiredSliceT(typed, nil)
		case **[]uint32:
			isValid = requiredSliceT(nil, typed)
		case *[]uint64:
			isValid = requiredSliceT(typed, nil)
		case **[]uint64:
			isValid = requiredSliceT(nil, typed)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, s.field, v, invalidLocaleKey)}
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
