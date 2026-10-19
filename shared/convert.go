package shared

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"
)

// As converts the field to a typed configurator by conv, e.g. f.As(str.Of) or f.As(num.Ptr).
// conv is a function with the type constraint of the values it supports,
// so a conversion that doesn't fit the field type doesn't compile.
func (f *Field[T]) As[R any](conv func(f *Field[T]) R) R {
	return conv(f)
}

// As converts the slice configurator to a typed one by conv, e.g. s.As(str.Slice).
// conv is a function with the type constraint of the elements it supports,
// so a conversion that doesn't fit the element type doesn't compile.
func (s *SliceFieldConfigurator[E]) As[R any](conv func(s *SliceFieldConfigurator[E]) R) R {
	return conv(s)
}

// Deref converts the rules of a pointer value to the rules of the pointed value, a nil pointer is absent.
// Use it in conversions of pointer fields, e.g. num.Ptr is num.Of of the dereferenced field.
func Deref[T any](f *Field[*T]) *Field[T] {
	add := f.add
	return &Field[T]{field: f.field, h: f.h, fields: f.fields, schema: f.schema, add: func(fn checkFn[T]) {
		add(func(ctx context.Context, h Helper, obj any, v **T, errs []Error) []Error {
			var p *T
			if v != nil {
				p = *v
			}
			return fn(ctx, h, obj, p, errs)
		})
	}}
}

// Retype converts the rules of T to the rules of U with the same memory layout,
// e.g. of a named string type to string. It panics if the layouts differ.
func Retype[T, U any](f *Field[T]) *Field[U] {
	if same, ok := any(f).(*Field[U]); ok {
		return same
	}
	mustLayout[U](f.field, reflect.TypeFor[T]())
	add := f.add
	return &Field[U]{field: f.field, h: f.h, fields: f.fields, schema: f.schema, add: func(fn checkFn[U]) {
		add(func(ctx context.Context, h Helper, obj any, v *T, errs []Error) []Error {
			return fn(ctx, h, obj, (*U)(unsafe.Pointer(v)), errs)
		})
	}}
}

// RetypeSlice converts the slice configurator of E to the one of U, where E is U, *U
// or a type with the same memory layout. It panics if the layouts differ.
func RetypeSlice[E, U any](s *SliceFieldConfigurator[E]) *SliceFieldConfigurator[U] {
	if same, ok := any(s).(*SliceFieldConfigurator[U]); ok {
		return same
	}
	return &SliceFieldConfigurator[U]{field: s.field, h: s.h, fields: s.fields, get: elemsGetter[U](s.field), add: s.add, schema: s.schema}
}

// SlicePtr converts the field of a pointer to a slice, *[]E, to the slice configurator, a nil pointer is absent:
// b.Field(&m.Tags).As(shared.SlicePtr). It must be called on the field returned by the configurator.
func SlicePtr[E any](f *Field[*[]E]) *SliceFieldConfigurator[E] {
	if f.root == nil {
		panic(fmt.Sprintf("valigo: SlicePtr of field %s must be called before When and other conversions", f.field.GetStructPath()))
	}
	return NewSliceFieldConfigurator[E](f.field, f.h, f.fields, f.root, f.rootDescribe)
}

// SliceRules are the rules of the slice itself for the typed slice configurators, S is the type of the configurator.
// Embed it to get Required, MinLen, MaxLen, Custom and When that return the typed configurator.
type SliceRules[S, E any] struct {
	s    *SliceFieldConfigurator[E]
	self S
	wrap func(s *SliceFieldConfigurator[E]) S
}

// NewSliceRules creates SliceRules of the typed configurator self built over s,
// wrap builds the typed configurator over another slice configurator, e.g. the conditional one of When.
func NewSliceRules[S, E any](s *SliceFieldConfigurator[E], self S, wrap func(s *SliceFieldConfigurator[E]) S) SliceRules[S, E] {
	return SliceRules[S, E]{s: s, self: self, wrap: wrap}
}

// Slice returns the slice configurator the rules are built over, e.g. for its EachField.
func (r SliceRules[S, E]) Slice() *SliceFieldConfigurator[E] {
	return r.s
}

// Required checks that the slice is present and not empty.
func (r SliceRules[S, E]) Required() S {
	r.s.Required()
	return r.self
}

// MinLen checks that the slice contains at least minLen elements. Absent slice is skipped.
func (r SliceRules[S, E]) MinLen(minLen int) S {
	r.s.MinLen(minLen)
	return r.self
}

// MaxLen checks that the slice contains no more than maxLen elements. Absent slice is skipped.
func (r SliceRules[S, E]) MaxLen(maxLen int) S {
	r.s.MaxLen(maxLen)
	return r.self
}

// Custom allows for custom validation logic on the whole slice. Absent slice is skipped.
func (r SliceRules[S, E]) Custom(f func(ctx context.Context, h *FieldCustomHelper, value Elems[E]) []Error) S {
	r.s.Custom(f)
	return r.self
}

// When returns the typed configurator whose rules are applied only if whenFn returns true.
func (r SliceRules[S, E]) When(whenFn func(ctx context.Context, value Elems[E]) bool) S {
	if whenFn == nil {
		return r.self
	}
	return r.wrap(r.s.When(whenFn))
}
