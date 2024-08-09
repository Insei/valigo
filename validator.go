package valigo

import "reflect"

type Validator struct {
	storage *storage
}

func (v *Validator) Validate(obj any) []error {
	cache, ok := v.storage.getCache(reflect.TypeOf(obj))
	if !ok {
		return nil
	}
	var errs []error = nil
	for _, ch := range cache {
		errs = append(errs, ch.fn(helper, obj, ch.field.GetPtr(obj))...)
	}
	return errs
}

func New(opts ...Option) *Validator {
	v := &Validator{
		storage: newStorage(),
	}
	for _, opt := range opts {
		opt.apply(v)
	}
	return v
}

func Configure[T any](v *Validator, fn func(builder Builder[T], temp *T)) {
	obj := new(T)
	// allocate all pointers fields values recursively
	MustZero(obj)
	b := configure[T](obj, nil, v.storage)
	// Append users validators
	fn(b, obj)
	// save validators functions to cache
	v.storage.toCache(reflect.TypeOf(obj))
}
