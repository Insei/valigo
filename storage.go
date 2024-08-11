package valigo

import (
	"context"
	"reflect"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
)

type structValidationFn func(ctx context.Context, h shared.Helper, obj any) []shared.Error

type storage struct {
	validators map[reflect.Type][]structValidationFn
}

func (s *storage) newOnStruct(temp any, enabler func(context.Context, any) bool, fn shared.FieldValidationFn) {
	t := reflect.TypeOf(temp)
	fnNew := func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
		if enabler != nil && !enabler(ctx, obj) {
			return nil
		}
		return fn(ctx, h, obj)
	}
	s.validators[t] = append(s.validators[t], fnNew)
}

func (s *storage) newOnField(temp any, enabler func(context.Context, any) bool) func(field fmap.Field, fn shared.FieldValidationFn) {
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

func newStorage() *storage {
	return &storage{
		validators: make(map[reflect.Type][]structValidationFn),
	}
}
