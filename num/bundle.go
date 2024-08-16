package num

import (
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

// IntBundle is a struct that represents a bundle of integer fields.
// It provides methods for adding validation rules to the integer fields.
type IntBundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
	h        shared.Helper
}

// NewIntBundle creates a new intBundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewIntBundle(deps shared.BundleDependencies) *IntBundle {
	return &IntBundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

// Int returns an intBuilder instance for an integer field.
// It takes a pointer to an integer field as an argument.
func (i *IntBundle) Int(field *int) IntBuilder[int] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return &intBuilder[int]{
		field:    fmapField,
		appendFn: i.appendFn,
		h:        i.h,
	}
}

// IntPtr returns an intBuilder instance for a pointer to an integer field.
// It takes a pointer to a pointer to an integer field as an argument.
func (i *IntBundle) IntPtr(field **int) IntBuilder[*int] {
	fmapField, err := i.storage.GetFieldByPtr(i.obj, field)
	if err != nil {
		panic(err)
	}
	return &intBuilder[*int]{
		field:    fmapField,
		appendFn: i.appendFn,
		h:        i.h,
	}
}

// IntSlice returns an IntSliceBuilder instance for an integer slice field.
// It takes a pointer to an integer slice field as an argument.
func (i *IntBundle) IntSlice(field *[]int) IntSliceBuilder[[]int] {
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
func (i *IntBundle) IntSlicePtr(field **[]int) IntSliceBuilder[*[]int] {
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
