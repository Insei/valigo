package valigo

import (
	"context"
	"reflect"

	"github.com/insei/valigo/helper"
)

type Validator struct {
	storage *storage
	helper  *helper.Helper
}

func (v *Validator) Validate(ctx context.Context, obj any) []error {
	validators, ok := v.storage.validators[reflect.TypeOf(obj)]
	if !ok {
		return nil
	}
	var errs []error
	for _, validator := range validators {
		errs = append(errs, validator(ctx, v.helper, obj)...)
	}
	return errs
}

func New(opts ...Option) *Validator {
	v := &Validator{
		storage: newStorage(),
		helper:  helper.NewHelper(),
	}
	for _, opt := range opts {
		opt.apply(v)
	}
	return v
}

func Configure[T any](v *Validator, fn func(builder Builder[T], obj *T)) {
	obj := new(T)
	// allocate all pointers fields values recursively
	mustZero(obj)
	b := configure[T](v, obj, nil)
	// Append users validators
	fn(b, obj)
}
