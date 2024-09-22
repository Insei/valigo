package valigo

import (
	"context"
	"reflect"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

// structValidationFn represents a function that validates a struct.
// It takes a context, a helper, and an object as input, and returns a slice of errors.
type structValidationFn func(ctx context.Context, h shared.Helper, obj any) []shared.Error

// storage represents the storage for struct validators.
type storage struct {
	// Validators is a map that stores validators for each struct type.
	// The key is the reflect.Type of the struct, and the value is a slice of structValidationFn.
	validators map[reflect.Type][]structValidationFn
}

// newOnStructAppend adds a new struct validator to the storage.
// It takes a temporary object, an enabler function, and a field validation function as input.
// The enabler function is optional and can be used to conditionally enable the validation.
// The field validation function is called for each field of the struct.
func (s *storage) newOnStructAppend(temp any, enabler func(context.Context, any) bool, fn shared.FieldValidationFn) {
	t := reflect.TypeOf(temp)
	fnNew := func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
		if enabler != nil && !enabler(ctx, obj) {
			return nil
		}
		return fn(ctx, h, obj)
	}
	s.validators[t] = append(s.validators[t], fnNew)
}

// newOnFieldAppend adds a new field validator to the storage.
// It takes a temporary object, an enabler function, and a field validation function as input.
// The enabler function is optional and can be used to conditionally enable the validation.
// The field validation function is called for each field of the struct.
func (s *storage) newOnFieldAppend(temp any, enabler func(context.Context, any) bool) func(field fmap.Field, fn shared.FieldValidationFn) {
	t := reflect.TypeOf(temp)
	return func(field fmap.Field, fn shared.FieldValidationFn) {
		fnNew := func(ctx context.Context, h shared.Helper, obj, v any) []shared.Error {
			if enabler != nil && !enabler(ctx, obj) {
				return nil
			}
			return fn(ctx, h, v)
		}
		_, ok := s.validators[t]
		if !ok {
			s.validators[t] = make([]structValidationFn, 0)
		}
		structValidatorFn := func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
			return fnNew(ctx, h, obj, field.GetPtr(obj))
		}
		s.validators[t] = append(s.validators[t], structValidatorFn)
	}
}

// newStorage creates a new storage object.
func newStorage() *storage {
	return &storage{
		validators: make(map[reflect.Type][]structValidationFn),
	}
}
