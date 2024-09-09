package num

import (
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

// NewIntBundle creates a new intBundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewIntBundle(deps shared.BundleDependencies) *Bundle {
	return &Bundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

func newValueConfigurator[T numbers](field fmap.Field,
	appendFn func(field fmap.Field, fn shared.FieldValidationFn),
	h shared.Helper) *valueConfigurator[T] {
	return &valueConfigurator[T]{
		baseConfigurator: &baseConfigurator[T]{
			field:    field,
			appendFn: appendFn,
			h:        h,
			dereference: func(value any) (T, bool) {
				val, ok := value.(*T)
				if !ok {
					return 0, false
				}
				return *val, true
			},
		},
	}
}

func newPtrConfigurator[T numbers](field fmap.Field,
	appendFn func(field fmap.Field, fn shared.FieldValidationFn),
	h shared.Helper) *ptrConfigurator[T] {
	return &ptrConfigurator[T]{
		baseConfigurator: &baseConfigurator[T]{
			field:    field,
			appendFn: appendFn,
			h:        h,
			dereference: func(value any) (T, bool) {
				val, ok := value.(**T)
				if !ok {
					return 0, false
				}
				if *val == nil {
					return 0, false
				}
				return **val, true
			},
		},
	}
}

// Int returns a Configurator instance for an int field.
// It takes a pointer to an integer field as an argument.
func (i *Bundle) Int(field *int) ValueConfigurator[int] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return newValueConfigurator[int](fmapField, i.appendFn, i.h)
}

// IntPtr returns a Configurator instance for a pointer to an int field.
// It takes a pointer to a pointer to an integer field as an argument.
func (i *Bundle) IntPtr(field **int) PtrConfigurator[int] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return newPtrConfigurator[int](fmapField, i.appendFn, i.h)
}

// Int8 returns a Configurator instance for an int8 field.
// It takes a pointer to an int8 field as an argument.
func (i *Bundle) Int8(field *int8) ValueConfigurator[int8] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return newValueConfigurator[int8](fmapField, i.appendFn, i.h)
}

// Int8Ptr returns a Configurator instance for a pointer to an int8 field.
// It takes a pointer to a pointer to an int8 field as an argument.
func (i *Bundle) Int8Ptr(field **int8) PtrConfigurator[int8] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return newPtrConfigurator[int8](fmapField, i.appendFn, i.h)
}

// Int16 returns a Configurator instance for an Int16 field.
// It takes a pointer to an Int16 field as an argument.
func (i *Bundle) Int16(field *int16) ValueConfigurator[int16] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return newValueConfigurator[int16](fmapField, i.appendFn, i.h)
}

// Int16Ptr returns a Configurator instance for a pointer to an Int16 field.
// It takes a pointer to a pointer to an Int16 field as an argument.
func (i *Bundle) Int16Ptr(field **int16) PtrConfigurator[int16] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return newPtrConfigurator[int16](fmapField, i.appendFn, i.h)
}

// Int32 returns a Configurator instance for an int32 field.
// It takes a pointer to an int32 field as an argument.
func (i *Bundle) Int32(field *int32) ValueConfigurator[int32] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return newValueConfigurator[int32](fmapField, i.appendFn, i.h)
}

// Int32Ptr returns a Configurator instance for a pointer to an int32 field.
// It takes a pointer to a pointer to an int32 field as an argument.
func (i *Bundle) Int32Ptr(field **int32) PtrConfigurator[int32] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return newPtrConfigurator[int32](fmapField, i.appendFn, i.h)
}

// IntSlice returns an IntSliceBuilder instance for an integer slice field.
// It takes a pointer to an integer slice field as an argument.
func (i *Bundle) IntSlice(field *[]int) IntSliceBuilder[[]int] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return &intSliceBuilder[[]int]{
		field:    fmapField,
		appendFn: i.appendFn,
		h:        i.h,
	}
}

// IntSlicePtr returns an IntSliceBuilder instance for a pointer to an integer slice field.
// It takes a pointer to a pointer to an integer slice field as an argument.
func (i *Bundle) IntSlicePtr(field **[]int) IntSliceBuilder[*[]int] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return &intSliceBuilder[*[]int]{
		field:    fmapField,
		appendFn: i.appendFn,
		h:        i.h,
	}
}
