package valigo

import (
	"context"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
	"github.com/insei/valigo/uuid"
)

// builder represents a builder for validators.
// It is generic and can be used to build validators for any type.
type builder[T any] struct {
	*str.StringBundle
	*num.Bundle
	*uuid.UUIDBundle
	obj       any
	v         *Validator
	enablerFn func(ctx context.Context, obj any) bool
}

// When adds a condition to the builder.
// It takes a function that returns a boolean and returns a builder.
// The function is called with the context and the object being validated.
// If the function returns true, the validation is enabled.
func (b *builder[T]) When(fn func(ctx context.Context, obj *T) bool) Configurator[T] {
	return configure[T](b.v, b.obj, func(ctx context.Context, obj any) bool {
		enablerFn := fn
		if b.enablerFn != nil {
			enablerFn = func(ctx context.Context, obj *T) bool {
				if b.enablerFn(ctx, obj) && fn(ctx, obj) {
					return true
				}
				return false
			}
		}
		return enablerFn(ctx, obj.(*T))
	})
}

// errorTFn represents a function that returns an error.
// It takes a context, a pointer to a field, a field value, a locale key, and optional arguments as input,
// and returns a shared.Error.
type errorTFn func(ctx context.Context, ptrToField, fieldValue any, localeKey string, args ...any) shared.Error

// ErrorT calls the errorTFn function and returns a shared.Error.
// It takes a context, a pointer to a field, a field value, a locale key, and optional arguments as input,
// and returns a shared.Error.
func (e errorTFn) ErrorT(ctx context.Context, ptrToFieldValue, fieldValue any, localeKey string, args ...any) shared.Error {
	return e(ctx, ptrToFieldValue, fieldValue, localeKey, args...)
}

// Custom adds a custom validation function to the builder.
// It takes a function that takes a context, a helper, and an object as input,
// and returns a slice of shared.Error.
// The helper is used to create errors.
func (b *builder[T]) Custom(fn func(ctx context.Context, h shared.StructCustomHelper, obj *T) []shared.Error) {
	newHFn := func(obj any, h shared.Helper) errorTFn {
		return func(ctx context.Context, ptrToField, fieldValue any, localeKey string, args ...any) shared.Error {
			fields, err := fmap.GetFrom(obj)
			if err != nil {
				panic(err)
			}
			field, err := fields.GetFieldByPtr(obj, ptrToField)
			if err != nil {
				panic(err)
			}
			return h.ErrorT(ctx, field, fieldValue, localeKey, args...)
		}
	}
	fnConvert := func(ctx context.Context, h shared.Helper, objAny any) []shared.Error {
		return fn(ctx, newHFn(objAny, h), objAny.(*T))
	}
	b.v.storage.newOnStructAppend(b.obj, b.enablerFn, fnConvert)
}

func (b *builder[T]) Slice(sliceFieldPtr any) *num.StringSliceFieldConfigurator {
	fields, err := fmap.GetFrom(b.obj)
	if err != nil {
		panic(err)
	}
	field, err := fields.GetFieldByPtr(b.obj, sliceFieldPtr)
	if err != nil {
		panic(err)
	}
	ok := shared.AddSliceMutation[string](sliceFieldPtr)
	if !ok {
		panic(err)
	}
	return &num.StringSliceFieldConfigurator{
		SliceFieldConfigurator: shared.NewSliceFieldConfigurator[string](shared.FieldConfiguratorParams[[]*any]{
			GetValueFn: func(value any) ([]*any, bool) {
				val, ok := shared.ConvertSliceValue[string](value)
				return val, ok
			},
			Field:    field,
			Helper:   b.v.GetHelper(),
			AppendFn: b.v.storage.newOnFieldAppend(b.obj, nil),
		})}
}

// configure creates a new builder with the given validator, object, and enabler function.
// It takes a validator, an object, and an enabler function as input,
// and returns a builder.
func configure[T any](v *Validator, obj any, enabler func(ctx context.Context, obj any) bool) *builder[T] {
	fields, err := fmap.GetFrom(obj)
	if err != nil {
		panic(err)
	}
	bundleDeps := shared.BundleDependencies{
		Object:   obj,
		Helper:   v.GetHelper(),
		AppendFn: v.storage.newOnFieldAppend(obj, enabler),
		Fields:   fields,
	}
	sb := str.NewStringBundle(bundleDeps)
	nb := num.NewNumBundle(bundleDeps)
	ub := uuid.NewUUIDBundle(bundleDeps)
	return &builder[T]{
		StringBundle: sb,
		Bundle:       nb,
		UUIDBundle:   ub,
		obj:          obj,
		v:            v,
		enablerFn:    enabler,
	}
}
