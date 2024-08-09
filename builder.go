package valigo

import (
	"reflect"

	"github.com/insei/fmap/v3"
)

type builder[T any] struct {
	*stringBundle
	obj any
}

func (b *builder[T]) When(fn func(obj *T) bool) Builder[T] {
	return newBuilder[T](b.obj, func(obj any) bool {
		return fn(obj.(*T))
	})
}

func newBuilder[T any](obj any, enabler func(obj any) bool) Builder[T] {
	fieldsStorage, _ := fmap.GetFrom(obj)
	sb := newStringBundle(obj, helper, storage.new(obj, enabler), fieldsStorage)
	return &builder[T]{stringBundle: sb, obj: obj}
}

var helper Helper

func AddValidation[T any](fn func(builder Builder[T], temp *T)) {
	obj := new(T)
	// allocate all pointers fields values recursively
	MustZero(obj)
	// initialize new builder instance
	b := newBuilder[T](obj, nil)
	// Append users validators
	fn(b, obj)
	// save validators functions to cache
	typeOf := reflect.TypeOf(obj)
	for field, validators := range storage.data[typeOf] {
		for _, validator := range validators {
			cache[typeOf] = append(cache[typeOf], cached{field: field, fn: validator})
		}
	}
}

type cached struct {
	field fmap.Field
	fn    validatorFn
}

var cache = map[reflect.Type][]cached{}

func Validate(obj any) []error {
	caches, ok := cache[reflect.TypeOf(obj)]
	if !ok {
		return nil
	}
	var errs []error = nil
	for _, ch := range caches {
		errs = append(errs, ch.fn(helper, obj, ch.field.GetPtr(obj))...)
	}
	return errs
}
