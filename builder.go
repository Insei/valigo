package valigo

import (
	"context"

	"github.com/insei/fmap/v3"
)

type builder[T any] struct {
	*stringBundle
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

func (b *builder[T]) Custom(fn func(ctx context.Context, h *Helper, obj *T) []*Error) {
	fnConvert := func(ctx context.Context, h *Helper, objAny any) []*Error {
		return fn(ctx, h, objAny.(*T))
	}
	b.v.storage.newOnStruct(b.obj, b.enablerFn, fnConvert)
}

func configure[T any](v *Validator, obj any, enabler func(ctx context.Context, obj any) bool) Builder[T] {
	fieldsStorage, _ := fmap.GetFrom(obj)
	sb := newStringBundle(obj, v.storage.newOnField(obj, enabler), fieldsStorage)
	return &builder[T]{stringBundle: sb, obj: obj, v: v, enablerFn: enabler}
}
