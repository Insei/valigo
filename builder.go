package valigo

import (
	"context"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
)

type builder[T any] struct {
	*str.StringBundle
	obj       any
	v         *Validator
	enablerFn func(ctx context.Context, obj any) bool
}

func (b *builder[T]) When(fn func(ctx context.Context, obj *T) bool) Builder[T] {
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

type errorTFn func(ctx context.Context, ptrToField, fieldValue any, localeKey string, args ...any) shared.Error

func (e errorTFn) ErrorT(ctx context.Context, ptrToFieldValue, fieldValue any, localeKey string, args ...any) shared.Error {
	return e(ctx, ptrToFieldValue, fieldValue, localeKey, args...)
}

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
	b.v.storage.newOnStruct(b.obj, b.enablerFn, fnConvert)
}

func configure[T any](v *Validator, obj any, enabler func(ctx context.Context, obj any) bool) Builder[T] {
	fieldsStorage, _ := fmap.GetFrom(obj)
	sb := str.NewStringBundle(obj, v.GetHelper(), v.storage.newOnField(obj, enabler), fieldsStorage)
	return &builder[T]{StringBundle: sb, obj: obj, v: v, enablerFn: enabler}
}
