package num

import (
	"reflect"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

// Bundle is a struct that represents a bundle of integer fields.
// It provides methods for adding validation rules to the integer fields.
type Bundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
	h        shared.Helper
}

// NewNumBundle creates a new intBundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewNumBundle(deps shared.BundleDependencies) *Bundle {
	return &Bundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

func ptrDeref[T numbers](value any) (T, bool) {
	v, ok := value.(**T)
	if !ok {
		return 0, false
	}
	if *v == nil {
		return 0, false
	}
	return **v, true
}

func deref[T numbers](value any) (T, bool) {
	val, ok := value.(*T)
	if !ok {
		return 0, false
	}
	return *val, true
}

var dereferenceFuncCache = map[reflect.Type]func(value any) (any, bool){
	reflect.TypeOf(new(int)): func(value any) (any, bool) {
		return deref[int](value)
	},
	reflect.TypeOf(new(int8)): func(value any) (any, bool) {
		return deref[int8](value)
	},
	reflect.TypeOf(new(int16)): func(value any) (any, bool) {
		return deref[int16](value)
	},
	reflect.TypeOf(new(int32)): func(value any) (any, bool) {
		return deref[int32](value)
	},
	reflect.TypeOf(new(int64)): func(value any) (any, bool) {
		return deref[int64](value)
	},
	reflect.TypeOf(new(uint)): func(value any) (any, bool) {
		return deref[uint](value)
	},
	reflect.TypeOf(new(uint8)): func(value any) (any, bool) {
		return deref[uint8](value)
	},
	reflect.TypeOf(new(uint16)): func(value any) (any, bool) {
		return deref[uint16](value)
	},
	reflect.TypeOf(new(uint32)): func(value any) (any, bool) {
		return deref[uint32](value)
	},
	reflect.TypeOf(new(uint64)): func(value any) (any, bool) {
		return deref[uint64](value)
	},
	reflect.TypeOf(new(float32)): func(value any) (any, bool) {
		return deref[float32](value)
	},
	reflect.TypeOf(new(float64)): func(value any) (any, bool) {
		return deref[float64](value)
	},
	reflect.TypeOf(new(*int)): func(value any) (any, bool) {
		return ptrDeref[int](value)
	},
	reflect.TypeOf(new(*int8)): func(value any) (any, bool) {
		return ptrDeref[int8](value)
	},
	reflect.TypeOf(new(*int16)): func(value any) (any, bool) {
		return ptrDeref[int16](value)
	},
	reflect.TypeOf(new(*int32)): func(value any) (any, bool) {
		return ptrDeref[int32](value)
	},
	reflect.TypeOf(new(*int64)): func(value any) (any, bool) {
		return ptrDeref[int64](value)
	},
	reflect.TypeOf(new(*uint)): func(value any) (any, bool) {
		return ptrDeref[uint](value)
	},
	reflect.TypeOf(new(*uint8)): func(value any) (any, bool) {
		return ptrDeref[uint8](value)
	},
	reflect.TypeOf(new(*uint16)): func(value any) (any, bool) {
		return ptrDeref[uint16](value)
	},
	reflect.TypeOf(new(*uint32)): func(value any) (any, bool) {
		return ptrDeref[uint32](value)
	},
	reflect.TypeOf(new(*uint64)): func(value any) (any, bool) {
		return ptrDeref[uint64](value)
	},
	reflect.TypeOf(new(*float32)): func(value any) (any, bool) {
		return ptrDeref[float32](value)
	},
	reflect.TypeOf(new(*float64)): func(value any) (any, bool) {
		return ptrDeref[float64](value)
	},
}

type baseConfiguratorParams[T numbers] struct {
	Field    fmap.Field
	Helper   shared.Helper
	AppendFn func(fn shared.FieldValidationFn)
}

func newBaseConfigurator[T numbers](p baseConfiguratorParams[T], derefFn func(value any) (any, bool)) *baseConfigurator[T] {
	mk := shared.NewSimpleFieldFnMaker(shared.SimpleFieldFnMakerParams[T]{
		GetValue: func(value any) (T, bool) {
			val, ok := derefFn(value)
			return val.(T), ok
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

// Number returns a FieldConfigurator instance for an int field.
// It takes a pointer to an integer field as an argument.
func (i *Bundle) Number(fieldPtr any) BaseConfigurator {
	field, err := i.storage.GetFieldByPtr(i.obj, fieldPtr)
	if err != nil {
		panic(err)
	}
	derefFn, ok := dereferenceFuncCache[reflect.PointerTo(field.GetType())]
	if !ok {
		panic("unsupported number field type")
	}
	switch field.GetDereferencedType().Kind() {
	case reflect.Int:
		return newBaseConfigurator(baseConfiguratorParams[int]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Int8:
		return newBaseConfigurator(baseConfiguratorParams[int8]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Int16:
		return newBaseConfigurator(baseConfiguratorParams[int16]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Int32:
		return newBaseConfigurator(baseConfiguratorParams[int32]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Int64:
		return newBaseConfigurator(baseConfiguratorParams[int64]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Uint:
		return newBaseConfigurator(baseConfiguratorParams[uint]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Uint8:
		return newBaseConfigurator(baseConfiguratorParams[uint8]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Uint16:
		return newBaseConfigurator(baseConfiguratorParams[uint16]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Uint32:
		return newBaseConfigurator(baseConfiguratorParams[uint32]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Uint64:
		return newBaseConfigurator(baseConfiguratorParams[uint64]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Float32:
		return newBaseConfigurator(baseConfiguratorParams[float32]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	case reflect.Float64:
		return newBaseConfigurator(baseConfiguratorParams[float64]{
			Field:  field,
			Helper: i.h,
			AppendFn: func(fn shared.FieldValidationFn) {
				i.appendFn(field, fn)
			},
		}, derefFn)
	default:
		panic("unsupported number field type")
	}
	return nil
}

//// Number returns a FieldConfigurator instance for an int field.
//// It takes a pointer to an integer field as an argument.
//func (i *Bundle) NumberSlice(fieldPtr any) *shared.SliceFieldConfigurator {
//	field, err := i.storage.GetFieldByPtr(i.obj, fieldPtr)
//	if err != nil {
//		panic(err)
//	}
//
//	typeOf := field.GetType()
//
//	shared.AddSliceMutation(fieldPtr)
//}
