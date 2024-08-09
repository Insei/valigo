package valigo

import (
	"github.com/insei/fmap/v3"
)

type builder[T any] struct {
	*stringBundle
	obj any
	v   *Validator
}

func (b *builder[T]) When(fn func(obj *T) bool) Builder[T] {
	return configure[T](b.v, b.obj, func(obj any) bool {
		return fn(obj.(*T))
	})
}

func configure[T any](v *Validator, obj any, enabler func(obj any) bool) Builder[T] {
	fieldsStorage, _ := fmap.GetFrom(obj)
	sb := newStringBundle(obj, v.storage.new(obj, enabler), fieldsStorage)
	return &builder[T]{stringBundle: sb, obj: obj, v: v}
}
