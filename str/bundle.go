package str

import (
	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
)

type strPtr interface {
	*string
}

// StringBundle is a struct that represents a bundle of integer fields.
// It provides methods for adding validation rules to the integer fields.
type StringBundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
	h        shared.Helper
}

// NewStringBundle creates a new intBundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewStringBundle(deps shared.BundleDependencies) *StringBundle {
	return &StringBundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

type baseConfiguratorParams[T strPtr] struct {
	Field    fmap.Field
	Helper   shared.Helper
	AppendFn func(fn shared.FieldValidationFn)
}

func newBaseConfigurator[T strPtr](p baseConfiguratorParams[T], derefFn func(value any) (*string, bool)) *baseConfigurator[T] {
	mk := shared.NewSimpleFieldFnMaker(shared.SimpleFieldFnMakerParams[T]{
		GetValue: func(value any) (T, bool) {
			val, ok := derefFn(value)
			return val, ok
		},
		Field:  p.Field,
		Helper: p.Helper,
	})
	return &baseConfigurator[T]{
		field: p.Field,
		h:     p.Helper,
		c: shared.NewFieldConfigurator[T](shared.FieldConfiguratorParams[T]{
			Maker:    mk,
			AppendFn: p.AppendFn,
		}),
	}
}

func deref[T strPtr](value any) (T, bool) {
	val, ok := value.(T)
	if !ok {
		return nil, false
	}
	return val, true
}

// String returns a FieldConfigurator instance for an string field.
// It takes a pointer to a string field as an argument.
func (i *StringBundle) String(fieldPtr any) BaseConfigurator {
	field, err := i.storage.GetFieldByPtr(i.obj, fieldPtr)
	if err != nil {
		panic(err)
	}
	return newBaseConfigurator(baseConfiguratorParams[*string]{
		Field:  field,
		Helper: i.h,
		AppendFn: func(fn shared.FieldValidationFn) {
			i.appendFn(field, fn)
		},
	}, deref)
}
