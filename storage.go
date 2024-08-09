package valigo

import (
	"context"
	"reflect"

	"github.com/insei/fmap/v3"
)

type storage struct {
	validators map[reflect.Type][]func(ctx context.Context, h *Helper, obj any) []error
}

func (s *storage) newOnStruct(temp any, enabler func(context.Context, any) bool, fn func(ctx context.Context, h *Helper, obj any) []*Error) {
	t := reflect.TypeOf(temp)
	fnNew := func(ctx context.Context, h *Helper, obj any) []error {
		if enabler != nil && !enabler(ctx, obj) {
			return nil
		}
		errsTyped := fn(ctx, h, obj)
		errs := make([]error, len(errsTyped))
		for i, err := range errsTyped {
			errs[i] = err
		}
		return errs
	}
	s.validators[t] = append(s.validators[t], fnNew)
}

func (s *storage) newOnField(temp any, enabler func(context.Context, any) bool) func(field fmap.Field, fn func(ctx context.Context, h *Helper, v any) []error) {
	t := reflect.TypeOf(temp)
	return func(field fmap.Field, fn func(ctx context.Context, h *Helper, v any) []error) {
		fnNew := func(ctx context.Context, h *Helper, obj, v any) []error {
			if enabler != nil && !enabler(ctx, obj) {
				return nil
			}
			return fn(ctx, h, v)
		}
		_, ok := s.validators[t]
		if !ok {
			s.validators[t] = make([]func(ctx context.Context, h *Helper, obj any) []error, 0)
		}
		structValidatorFn := func(ctx context.Context, h *Helper, obj any) []error {
			return fnNew(ctx, h, obj, field.GetPtr(obj))
		}
		s.validators[t] = append(s.validators[t], structValidatorFn)
	}
}

func newStorage() *storage {
	return &storage{
		validators: make(map[reflect.Type][]func(ctx context.Context, h *Helper, obj any) []error),
	}
}
