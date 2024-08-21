package num

import (
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

// UintBundle is a struct that represents a bundle of uint fields.
// It provides methods for adding validation rules to the uint fields.
type UintBundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
	h        shared.Helper
}

// NewUintBundle creates a new uintBundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewUintBundle(deps shared.BundleDependencies) *UintBundle {
	return &UintBundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

// Uint returns an intBuilder instance for an uint field.
// It takes a pointer to an uint (*uint) field as an argument.
func (i *UintBundle) Uint(field any) UintBuilder[uint] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return &uintBuilder[uint]{
		field:    fmapField,
		appendFn: i.appendFn,
		h:        i.h,
	}
}

// UintPtr returns an intBuilder instance for a pointer to an uint field.
// It takes a pointer to a pointer to an uint (**uint) field as an argument.
func (i *UintBundle) UintPtr(field any) UintBuilder[*uint] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return &uintBuilder[*uint]{
		field:    fmapField,
		appendFn: i.appendFn,
		h:        i.h,
	}
}

// UintSlice returns an UintSliceBuilder instance for an uint slice field.
// It takes a pointer to an uint slice (*[]uint) field as an argument.
func (i *UintBundle) UintSlice(field any) UintSliceBuilder[[]uint] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return &uintSliceBuilder[[]uint]{
		field:    fmapField,
		appendFn: i.appendFn,
		h:        i.h,
	}
}

// UintSlicePtr returns an UintSliceBuilder instance for a pointer to an uint slice field.
// It takes a pointer to a pointer to an uint slice (**[]uint) field as an argument.
func (i *UintBundle) UintSlicePtr(field any) UintSliceBuilder[*[]uint] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return &uintSliceBuilder[*[]uint]{
		field:    fmapField,
		appendFn: i.appendFn,
		h:        i.h,
	}
}
