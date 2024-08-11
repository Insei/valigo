package valigo

import (
	"context"
	"reflect"

	"github.com/insei/valigo/shared"
)

type Validator struct {
	storage *storage
	helper  *shared.Helper
}

func (v *Validator) ValidateTyped(ctx context.Context, obj any) []shared.Error {
	validators, ok := v.storage.validators[reflect.TypeOf(obj)]
	if !ok {
		return nil
	}
	var errs []shared.Error
	for _, validator := range validators {
		errs = append(errs, validator(ctx, v.helper, obj)...)
	}
	return errs
}

func (v *Validator) Validate(ctx context.Context, obj any) []error {
	errsTyped := v.ValidateTyped(ctx, obj)
	errs := make([]error, len(errsTyped))
	for i, err := range errsTyped {
		errs[i] = err
	}
	return errs
}

func New(opts ...Option) *Validator {
	v := &Validator{
		storage: newStorage(),
		helper:  shared.NewHelper(),
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
