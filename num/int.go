package num

import (
	"context"
	"slices"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

const (
	minLocaleKey          = "validation:int:Cannot be less than %d characters"
	maxLocaleKey          = "validation:int:Cannot be greater than %d characters"
	requiredLocaleKey     = "validation:int:Should be fulfilled"
	anyOfLocaleKey        = "validation:int:Only %d values is allowed"
	anyOfIntervalLocalKey = "validation:int:Only interval[%d - %d] is allowed"
	invalidLocaleKey      = "validation:int:Invalid value"
)

type intBuilder[T int | *int | int8 | *int8 | int16 | *int16 | int32 | *int32 | int64 | *int64 |
	uint | *uint | uint8 | *uint8 | uint16 | *uint16 | uint32 | *uint32 | uint64 | *uint64] struct {
	h        shared.Helper
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
}

func maxT[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](valuePtr *T, valuePtr2 **T, max int) (any, bool) {
	if valuePtr == nil && *valuePtr2 == nil {
		return 0, false
	}
	if valuePtr != nil {
		return *valuePtr, T(max) > *valuePtr
	}
	return **valuePtr2, T(max) > **valuePtr2
}

func minT[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](valuePtr *T, valuePtr2 **T, min int) (any, bool) {
	if valuePtr == nil && *valuePtr2 == nil {
		return 0, false
	}
	if valuePtr != nil {
		return *valuePtr, T(min) < *valuePtr
	}
	return **valuePtr2, T(min) < **valuePtr2
}

func requiredT[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](valuePtr *T, valuePtr2 **T) (any, bool) {
	if valuePtr == nil && *valuePtr2 == nil {
		return 0, false
	}
	if valuePtr != nil {
		return *valuePtr, *valuePtr != T(0)
	}
	return **valuePtr2, **valuePtr2 != T(0)
}

func anyOfT[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](valuePtr *T, valuePtr2 **T, allowed []int) (any, bool) {
	if valuePtr == nil && *valuePtr2 == nil {
		return 0, false
	}

	allowedT := make([]T, len(allowed))
	for i, v := range allowed {
		allowedT[i] = T(v)
	}

	if valuePtr != nil {
		return *valuePtr, slices.Contains(allowedT, *valuePtr)
	}
	return **valuePtr2, slices.Contains(allowedT, **valuePtr2)
}

func anyOfIntervalT[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](valuePtr *T, valuePtr2 **T, begin, end int) (any, bool) {
	if valuePtr == nil && *valuePtr2 == nil {
		return 0, false
	}
	if valuePtr != nil {
		return *valuePtr, *valuePtr > T(begin) && *valuePtr < T(end)
	}
	return **valuePtr2, **valuePtr2 > T(begin) && **valuePtr2 < T(end)
}

// Max checks if the integer exceeds the maximum allowed number.
func (i *intBuilder[T]) Max(maxNum int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		var v any
		switch typed := value.(type) {
		case *int:
			v, isValid = maxT(typed, nil, maxNum)
		case **int:
			v, isValid = maxT(nil, typed, maxNum)
		case *int8:
			v, isValid = maxT(typed, nil, maxNum)
		case **int8:
			v, isValid = maxT(nil, typed, maxNum)
		case *int16:
			v, isValid = maxT(typed, nil, maxNum)
		case **int16:
			v, isValid = maxT(nil, typed, maxNum)
		case *int32:
			v, isValid = maxT(typed, nil, maxNum)
		case **int32:
			v, isValid = maxT(nil, typed, maxNum)
		case *int64:
			v, isValid = maxT(typed, nil, maxNum)
		case **int64:
			v, isValid = maxT(nil, typed, maxNum)
		case *uint:
			v, isValid = maxT(typed, nil, maxNum)
		case **uint:
			v, isValid = maxT(nil, typed, maxNum)
		case *uint8:
			v, isValid = maxT(typed, nil, maxNum)
		case **uint8:
			v, isValid = maxT(nil, typed, maxNum)
		case *uint16:
			v, isValid = maxT(typed, nil, maxNum)
		case **uint16:
			v, isValid = maxT(nil, typed, maxNum)
		case *uint32:
			v, isValid = maxT(typed, nil, maxNum)
		case **uint32:
			v, isValid = maxT(nil, typed, maxNum)
		case *uint64:
			v, isValid = maxT(typed, nil, maxNum)
		case **uint64:
			v, isValid = maxT(nil, typed, maxNum)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, maxLocaleKey, maxNum)}
		}
		return nil
	})
	return i
}

// Min checks if the integer is less than the minimum allowed number.
func (i *intBuilder[T]) Min(minNum int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		var v any
		switch typed := value.(type) {
		case *int:
			v, isValid = minT(typed, nil, minNum)
		case **int:
			v, isValid = minT(nil, typed, minNum)
		case *int8:
			v, isValid = minT(typed, nil, minNum)
		case **int8:
			v, isValid = minT(nil, typed, minNum)
		case *int16:
			v, isValid = minT(typed, nil, minNum)
		case **int16:
			v, isValid = minT(nil, typed, minNum)
		case *int32:
			v, isValid = minT(typed, nil, minNum)
		case **int32:
			v, isValid = minT(nil, typed, minNum)
		case *int64:
			v, isValid = minT(typed, nil, minNum)
		case **int64:
			v, isValid = minT(nil, typed, minNum)
		case *uint:
			v, isValid = minT(typed, nil, minNum)
		case **uint:
			v, isValid = minT(nil, typed, minNum)
		case *uint8:
			v, isValid = minT(typed, nil, minNum)
		case **uint8:
			v, isValid = minT(nil, typed, minNum)
		case *uint16:
			v, isValid = minT(typed, nil, minNum)
		case **uint16:
			v, isValid = minT(nil, typed, minNum)
		case *uint32:
			v, isValid = minT(typed, nil, minNum)
		case **uint32:
			v, isValid = minT(nil, typed, minNum)
		case *uint64:
			v, isValid = minT(typed, nil, minNum)
		case **uint64:
			v, isValid = minT(nil, typed, minNum)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, minLocaleKey, minNum)}
		}
		return nil
	})
	return i
}

// Required checks if the integer value is not empty.
func (i *intBuilder[T]) Required() IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		var v any
		switch typed := value.(type) {
		case *int:
			v, isValid = requiredT(typed, nil)
		case **int:
			v, isValid = requiredT(nil, typed)
		case *int8:
			v, isValid = requiredT(typed, nil)
		case **int8:
			v, isValid = requiredT(nil, typed)
		case *int16:
			v, isValid = requiredT(typed, nil)
		case **int16:
			v, isValid = requiredT(nil, typed)
		case *int32:
			v, isValid = requiredT(typed, nil)
		case **int32:
			v, isValid = requiredT(nil, typed)
		case *int64:
			v, isValid = requiredT(typed, nil)
		case **int64:
			v, isValid = requiredT(nil, typed)
		case *uint:
			v, isValid = requiredT(typed, nil)
		case **uint:
			v, isValid = requiredT(nil, typed)
		case *uint8:
			v, isValid = requiredT(typed, nil)
		case **uint8:
			v, isValid = requiredT(nil, typed)
		case *uint16:
			v, isValid = requiredT(typed, nil)
		case **uint16:
			v, isValid = requiredT(nil, typed)
		case *uint32:
			v, isValid = requiredT(typed, nil)
		case **uint32:
			v, isValid = requiredT(nil, typed)
		case *uint64:
			v, isValid = requiredT(typed, nil)
		case **uint64:
			v, isValid = requiredT(nil, typed)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, requiredLocaleKey)}
		}
		return nil
	})
	return i
}

// AnyOf checks if the integer value is one of the allowed values.
func (i *intBuilder[T]) AnyOf(allowed ...int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		var v any
		switch typed := value.(type) {
		case *int:
			v, isValid = anyOfT(typed, nil, allowed)
		case **int:
			v, isValid = anyOfT(nil, typed, allowed)
		case *int8:
			v, isValid = anyOfT(typed, nil, allowed)
		case **int8:
			v, isValid = anyOfT(nil, typed, allowed)
		case *int16:
			v, isValid = anyOfT(typed, nil, allowed)
		case **int16:
			v, isValid = anyOfT(nil, typed, allowed)
		case *int32:
			v, isValid = anyOfT(typed, nil, allowed)
		case **int32:
			v, isValid = anyOfT(nil, typed, allowed)
		case *int64:
			v, isValid = anyOfT(typed, nil, allowed)
		case **int64:
			v, isValid = anyOfT(nil, typed, allowed)
		case *uint:
			v, isValid = anyOfT(typed, nil, allowed)
		case **uint:
			v, isValid = anyOfT(nil, typed, allowed)
		case *uint8:
			v, isValid = anyOfT(typed, nil, allowed)
		case **uint8:
			v, isValid = anyOfT(nil, typed, allowed)
		case *uint16:
			v, isValid = anyOfT(typed, nil, allowed)
		case **uint16:
			v, isValid = anyOfT(nil, typed, allowed)
		case *uint32:
			v, isValid = anyOfT(typed, nil, allowed)
		case **uint32:
			v, isValid = anyOfT(nil, typed, allowed)
		case *uint64:
			v, isValid = anyOfT(typed, nil, allowed)
		case **uint64:
			v, isValid = anyOfT(nil, typed, allowed)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, anyOfLocaleKey, allowed)}
		}
		return nil
	})
	return i
}

// AnyOfInterval checks if the integer value is one of the allowed values intervals.
func (i *intBuilder[T]) AnyOfInterval(begin, end int) IntBuilder[T] {
	i.appendFn(i.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		isValid := true
		var v any
		switch typed := value.(type) {
		case *int:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **int:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *int8:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **int8:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *int16:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **int16:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *int32:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **int32:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *int64:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **int64:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *uint:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **uint:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *uint8:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **uint8:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *uint16:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **uint16:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *uint32:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **uint32:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		case *uint64:
			v, isValid = anyOfIntervalT(typed, nil, begin, end)
		case **uint64:
			v, isValid = anyOfIntervalT(nil, typed, begin, end)
		}
		if !isValid {
			return []shared.Error{h.ErrorT(ctx, i.field, v, anyOfIntervalLocalKey, begin, end)}
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
