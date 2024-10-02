package shared

import (
	"context"

	"github.com/insei/fmap/v3"
)

type ValidationFnMaker[T any] interface {
	Make(validationFn func(v T) bool, format string, args ...any) FieldValidationFn
}

type simpleFieldFnMaker[T any] struct {
	getValue func(value any) (T, bool)
	field    fmap.Field
	helper   Helper
}

type SimpleFieldFnMakerParams[T any] struct {
	GetValue func(value any) (T, bool)
	Field    fmap.Field
	Helper   Helper
}

func NewSimpleFieldFnMaker[T any](p SimpleFieldFnMakerParams[T]) ValidationFnMaker[T] {
	return &simpleFieldFnMaker[T]{
		getValue: p.GetValue,
		field:    p.Field,
		helper:   p.Helper,
	}
}

func (i *simpleFieldFnMaker[T]) Make(validationFn func(v T) bool, format string, args ...any) FieldValidationFn {
	return func(ctx context.Context, h Helper, val any) []Error {
		v, isValid := i.getValue(val)
		if !isValid {
			return []Error{i.helper.ErrorT(ctx, i.field, val, format, args...)}
		}
		isValid = validationFn(v)
		if !isValid {
			return []Error{h.ErrorT(ctx, i.field, v, format, args...)}
		}
		return nil
	}
}

type FieldConfigurator[T any] struct {
	appendFn func(fn FieldValidationFn)
	mk       ValidationFnMaker[T]
}

func (i *FieldConfigurator[T]) Append(validationFn func(v T) bool, format string, args ...any) {
	i.appendFn(i.mk.Make(validationFn, format, args...))
}

func (i *FieldConfigurator[T]) CustomAppend(fn FieldValidationFn) {
	i.appendFn(fn)
}

func (i *FieldConfigurator[T]) NewWithWhen(whenFn func(ctx context.Context, value any) bool) *FieldConfigurator[T] {
	return &FieldConfigurator[T]{
		appendFn: func(fn FieldValidationFn) {
			fnWithEnabler := func(ctx context.Context, h Helper, v any) []Error {
				if !whenFn(ctx, v) {
					return nil
				}
				return fn(ctx, h, v)
			}
			i.appendFn(fnWithEnabler)
		},
	}
}

type FieldConfiguratorParams[T any] struct {
	Maker    ValidationFnMaker[T]
	AppendFn func(fn FieldValidationFn)
}

func NewFieldConfigurator[T any](p FieldConfiguratorParams[T]) *FieldConfigurator[T] {
	return &FieldConfigurator[T]{
		appendFn: p.AppendFn,
		mk:       p.Maker,
	}
}
