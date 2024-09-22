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

func newBaseConfigurator[T numbers](params shared.FieldConfiguratorParams[T], derefFn func(value any) (any, bool)) *baseConfigurator[T] {
	params.GetValueFn = func(value any) (T, bool) {
		val, ok := derefFn(value)
		return val.(T), ok
	}
	return &baseConfigurator[T]{
		c: shared.NewFieldConfigurator[T](params),
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
		return newBaseConfigurator(shared.FieldConfiguratorParams[int]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Int8:
		return newBaseConfigurator(shared.FieldConfiguratorParams[int8]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Int16:
		return newBaseConfigurator(shared.FieldConfiguratorParams[int16]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Int32:
		return newBaseConfigurator(shared.FieldConfiguratorParams[int32]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Int64:
		return newBaseConfigurator(shared.FieldConfiguratorParams[int64]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Uint:
		return newBaseConfigurator(shared.FieldConfiguratorParams[uint]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Uint8:
		return newBaseConfigurator(shared.FieldConfiguratorParams[uint8]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Uint16:
		return newBaseConfigurator(shared.FieldConfiguratorParams[uint16]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Uint32:
		return newBaseConfigurator(shared.FieldConfiguratorParams[uint32]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Uint64:
		return newBaseConfigurator(shared.FieldConfiguratorParams[uint64]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Float32:
		return newBaseConfigurator(shared.FieldConfiguratorParams[float32]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
		}, derefFn)
	case reflect.Float64:
		return newBaseConfigurator(shared.FieldConfiguratorParams[float64]{
			Field:    field,
			Helper:   i.h,
			AppendFn: i.appendFn,
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
