package uuid

import (
	"github.com/google/uuid"
	"github.com/insei/fmap/v3"
	"reflect"

	"github.com/insei/valigo/shared"
)

// UUIDBundle is a struct that represents a bundle of uuid fields.
// It provides methods for adding validation rules to the uuid fields.
type UUIDBundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
	h        shared.Helper
}

// NewUUIDBundle creates a new Bundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewUUIDBundle(deps shared.BundleDependencies) *UUIDBundle {
	return &UUIDBundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

type baseConfiguratorParams struct {
	Field    fmap.Field
	Helper   shared.Helper
	AppendFn func(fn shared.FieldValidationFn)
}

func newBaseConfigurator(p baseConfiguratorParams, derefFn func(value any) (uuid.UUID, bool)) *baseConfigurator {
	mk := shared.NewSimpleFieldFnMaker(shared.SimpleFieldFnMakerParams[uuid.UUID]{
		GetValue: func(value any) (uuid.UUID, bool) {
			val, ok := derefFn(value)
			return val, ok
		},
		Field:  p.Field,
		Helper: p.Helper,
	})
	return &baseConfigurator{
		field: p.Field,
		h:     p.Helper,
		c: shared.NewFieldConfigurator[uuid.UUID](shared.FieldConfiguratorParams[uuid.UUID]{
			Maker:    mk,
			AppendFn: p.AppendFn,
		}),
	}
}

func ptrDeref(value any) (uuid.UUID, bool) {
	val, ok := value.(**uuid.UUID)
	if !ok {
		return uuid.Nil, ok
	}
	return **val, ok
}

func deref(value any) (uuid.UUID, bool) {
	val, ok := value.(*uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}
	return *val, true
}

// UUID returns a FieldConfigurator instance for an uuid field.
// It takes a pointer to an uuid field as an argument.
func (i *UUIDBundle) UUID(fieldPtr any) BaseConfigurator {
	field, err := i.storage.GetFieldByPtr(i.obj, fieldPtr)

	var derefFn func(value any) (uuid.UUID, bool)
	switch reflect.PointerTo(field.GetType()) {
	case reflect.TypeOf(new(uuid.UUID)):
		derefFn = deref
	case reflect.TypeOf(new(*uuid.UUID)):
		derefFn = ptrDeref
	}

	if err != nil {
		panic(err)
	}
	return newBaseConfigurator(baseConfiguratorParams{
		Field:  field,
		Helper: i.h,
		AppendFn: func(fn shared.FieldValidationFn) {
			i.appendFn(field, fn)
		},
	}, derefFn)
}
