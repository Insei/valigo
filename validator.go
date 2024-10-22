package valigo

import (
	"context"
	"reflect"

	"github.com/insei/valigo/shared"
)

// Validator is a struct that holds a storage and a helper object.
type Validator struct {
	storage *storage
	helper  *helper
}

// ValidateTyped validates an object of any type using validators from the storage.
// It takes a context.Context and an object as input and returns a slice of shared.Error objects.
// If no validators are found for the object's type, it returns nil.
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

// Validate is similar to ValidateTyped, but it returns a slice of error
// objects instead of shared.Error objects.
func (v *Validator) Validate(ctx context.Context, obj any) []error {
	errsTyped := v.ValidateTyped(ctx, obj)
	errs := make([]error, len(errsTyped))
	for i, err := range errsTyped {
		errs[i] = err
	}
	return errs
}

// GetHelper returns the helper object associated with the Validator instance.
func (v *Validator) GetHelper() shared.Helper {
	return v.helper
}

// New creates a new Validator instance with default values for storage and helper.
// It also applies any options passed to the function to the new instance.
func New(opts ...Option) *Validator {
	v := &Validator{
		storage: newStorage(),
		helper:  newHelper(),
	}
	for _, opt := range opts {
		opt.apply(v)
	}
	return v
}

// Configure configures a Validator instance for a specific type T, i.e.: var t *T.
// It creates a new instance of type T, allocates values for all fields recursively,
// creates a Configurator instance for the type T, calls the provided function fn with the Configurator
// instance and the instance of type T, and appends any user-defined validators to the Configurator instance.
func Configure[T any](v *Validator, fn func(builder Configurator[T], m *T)) {
	model := new(T)
	// allocate all pointer fields values recursively
	mustZero(model)
	b := configure[T](v, model, nil)
	// Append users validators
	fn(b, model)
}
