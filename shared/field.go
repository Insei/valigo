package shared

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/insei/fmap/v3"
)

// ValidationFn validates an object, obj is a pointer to the configured struct.
type ValidationFn func(ctx context.Context, h Helper, obj any) []Error

// checkFn checks a single value: a struct field or a slice element.
// v is nil when the value is absent, i.e. a nil pointer field or a field of a nil nested struct.
type checkFn[T any] func(ctx context.Context, h Helper, v *T) []Error

// Field is a typed rules set for a single value: a struct field or every element of a slice.
type Field[T any] struct {
	field fmap.Field
	h     Helper
	add   func(fn checkFn[T])
}

// NewValueField creates a Field for the struct field of type T or *T,
// named types with the same underlying layout are allowed.
// add registers a validation function of the configured struct.
func NewValueField[T any](field fmap.Field, h Helper, add func(fn ValidationFn)) *Field[T] {
	return &Field[T]{field: field, h: h, add: fieldSource(add, valueGetter[T](field))}
}

// Check adds a rule that fails when ok returns false. Absent values are skipped.
func (f *Field[T]) Check(ok func(v T) bool, localeKey string, args ...any) {
	field := f.field
	f.add(func(ctx context.Context, h Helper, v *T) []Error {
		if v == nil || ok(*v) {
			return nil
		}
		return []Error{h.ErrorT(ctx, field, *v, localeKey, args...)}
	})
}

// Required adds a rule that fails when the value is absent or ok returns false.
func (f *Field[T]) Required(ok func(v T) bool, localeKey string) {
	field := f.field
	f.add(func(ctx context.Context, h Helper, v *T) []Error {
		if v != nil && ok(*v) {
			return nil
		}
		var val any
		if v != nil {
			val = *v
		}
		return []Error{h.ErrorT(ctx, field, val, localeKey)}
	})
}

// Modify adds a step that changes the value in place, e.g. trims a string. Absent values are skipped.
func (f *Field[T]) Modify(fn func(v *T)) {
	f.add(func(_ context.Context, _ Helper, v *T) []Error {
		if v != nil {
			fn(v)
		}
		return nil
	})
}

// Custom adds a user defined rule. Absent values are skipped.
func (f *Field[T]) Custom(fn func(ctx context.Context, h *FieldCustomHelper, value T) []Error) {
	customHelper := NewFieldCustomHelper(f.field, f.h)
	f.add(func(ctx context.Context, _ Helper, v *T) []Error {
		if v == nil {
			return nil
		}
		return fn(ctx, customHelper, *v)
	})
}

// When returns a Field whose rules are applied only if cond returns true.
// When the value is absent, cond is not called and the rules, including Required, are skipped.
func (f *Field[T]) When(cond func(ctx context.Context, value T) bool) *Field[T] {
	if cond == nil {
		return f
	}
	add := f.add
	return &Field[T]{field: f.field, h: f.h, add: func(fn checkFn[T]) {
		add(func(ctx context.Context, h Helper, v *T) []Error {
			if v == nil || !cond(ctx, *v) {
				return nil
			}
			return fn(ctx, h, v)
		})
	}}
}

// fieldSource applies checks to the struct field value read by get.
// All checks of the field are registered as a single validation function, so the value is read once.
func fieldSource[T any](add func(fn ValidationFn), get func(obj any) *T) func(fn checkFn[T]) {
	var checks []checkFn[T]
	return func(fn checkFn[T]) {
		checks = append(checks, fn)
		if len(checks) > 1 {
			return
		}
		add(func(ctx context.Context, h Helper, obj any) []Error {
			v := get(obj)
			var errs []Error
			for _, check := range checks {
				errs = append(errs, check(ctx, h, v)...)
			}
			return errs
		})
	}
}

// valueGetter returns a function that reads a pointer to the T value of the field from obj.
// The field type must be T or *T, named types with the same underlying layout are allowed.
// The returned function returns nil when the value is absent.
func valueGetter[T any](field fmap.Field) func(obj any) *T {
	isPtr := field.GetType().Kind() == reflect.Pointer && reflect.TypeFor[T]().Kind() != reflect.Pointer
	t := field.GetType()
	if isPtr {
		t = t.Elem()
	}
	mustLayout[T](field, t)
	if isPtr {
		return func(obj any) *T {
			p := (**T)(dataPtr(field.GetPtr(obj)))
			if p == nil {
				return nil
			}
			return *p
		}
	}
	return func(obj any) *T {
		return (*T)(dataPtr(field.GetPtr(obj)))
	}
}

// dataPtr returns the data pointer of the interface value, i.e. the pointer itself for pointer values.
func dataPtr(v any) unsafe.Pointer {
	return (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[1]
}

// mustLayout panics if values of type t can't be read as T.
func mustLayout[T any](field fmap.Field, t reflect.Type) {
	want := reflect.TypeFor[T]()
	if t.Kind() != want.Kind() || t.Size() != want.Size() {
		panic(fmt.Sprintf("valigo: field %s has type %s, %s expected", field.GetStructPath(), field.GetType(), want))
	}
}
