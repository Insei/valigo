package shared

import (
	"context"

	"github.com/insei/fmap/v3"
)

type FieldConfigurator[T any] struct {
	GetValue func(value any) (T, bool)
	Field    fmap.Field
	Helper   Helper
	appendFn func(field fmap.Field, fn FieldValidationFn)
}

func (i *FieldConfigurator[T]) validate(ctx context.Context, h Helper, v T, fn func(v T) bool, format string, args ...any) []Error {
	isValid := fn(v)
	if !isValid {
		return []Error{h.ErrorT(ctx, i.Field, v, format, args...)}
	}
	return nil
}

func (i *FieldConfigurator[T]) makeValidationFn(validationFn func(v T) bool, format string, args ...any) FieldValidationFn {
	return func(ctx context.Context, h Helper, val any) []Error {
		v, isValid := i.GetValue(val)
		if !isValid {
			return []Error{i.Helper.ErrorT(ctx, i.Field, val, format, args...)}
		}
		isValid = validationFn(v)
		if !isValid {
			return []Error{h.ErrorT(ctx, i.Field, v, format, args...)}
		}
		return nil
	}
}
func (i *FieldConfigurator[T]) Append(validationFn func(v T) bool, format string, args ...any) {
	i.appendFn(i.Field, i.makeValidationFn(validationFn, format, args...))
}

func (i *FieldConfigurator[T]) CustomAppend(fn FieldValidationFn) {
	i.appendFn(i.Field, fn)
}

func (i *FieldConfigurator[T]) NewWithWhen(whenFn func(ctx context.Context, value any) bool) *FieldConfigurator[T] {
	return &FieldConfigurator[T]{
		GetValue: i.GetValue,
		Field:    i.Field,
		Helper:   i.Helper,
		appendFn: func(field fmap.Field, fn FieldValidationFn) {
			fnWithEnabler := func(ctx context.Context, h Helper, v any) []Error {
				if !whenFn(ctx, v) {
					return nil
				}
				return fn(ctx, h, v)
			}
			i.appendFn(field, fnWithEnabler)
		},
	}
}

type FieldConfiguratorParams[T any] struct {
	GetValueFn func(value any) (T, bool)
	Field      fmap.Field
	Helper     Helper
	AppendFn   func(field fmap.Field, fn FieldValidationFn)
}

func NewFieldConfigurator[T any](params FieldConfiguratorParams[T]) *FieldConfigurator[T] {
	return &FieldConfigurator[T]{
		GetValue: params.GetValueFn,
		Field:    params.Field,
		Helper:   params.Helper,
		appendFn: params.AppendFn,
	}
}
