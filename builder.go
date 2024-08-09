package valigo

import (
	"github.com/insei/fmap/v3"
)

type builder[T any] struct {
	*stringBundle
	obj     any
	storage *storage
}

func (b *builder[T]) When(fn func(obj *T) bool) Builder[T] {
	return configure[T](b.obj, func(obj any) bool {
		return fn(obj.(*T))
	}, b.storage)
}

func configure[T any](obj any, enabler func(obj any) bool, storage *storage) Builder[T] {
	fieldsStorage, _ := fmap.GetFrom(obj)
	sb := newStringBundle(obj, helper, storage.new(obj, enabler), fieldsStorage)
	return &builder[T]{stringBundle: sb, obj: obj, storage: storage}
}

var helper Helper
