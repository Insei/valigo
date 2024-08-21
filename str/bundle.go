package str

import (
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

// StringBundle is a struct that represents a bundle of string fields.
// It provides methods for adding validation rules to the string fields.
type StringBundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
	h        shared.Helper
}

// NewStringBundle creates a new StringBundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewStringBundle(deps shared.BundleDependencies) *StringBundle {
	return &StringBundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

// String returns a StringBuilder instance for a string field.
// It takes a pointer to a string field as an argument.
func (s *StringBundle) String(field *string) StringBuilder[string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringBuilder[string]{
		field:    fmapField,
		appendFn: s.appendFn,
		h:        s.h,
	}
}

// StringPtr returns a StringBuilder instance for a pointer to a string field.
// It takes a pointer to a pointer to a string field as an argument.
func (s *StringBundle) StringPtr(field **string) StringBuilder[*string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringBuilder[*string]{
		field:    fmapField,
		appendFn: s.appendFn,
		h:        s.h,
	}
}

// StringSlice returns a StringSliceBuilder instance for a string slice field.
// It takes a pointer to a string slice field as an argument.
func (s *StringBundle) StringSlice(field *[]string) StringSliceBuilder[[]string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringSliceBuilder[[]string]{
		field:    fmapField,
		appendFn: s.appendFn,
		h:        s.h,
	}
}

// StringSlicePtr returns a StringSliceBuilder instance for a pointer to a string slice field.
// It takes a pointer to a pointer to a string slice field as an argument.
func (s *StringBundle) StringSlicePtr(field **[]string) StringSliceBuilder[*[]string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringSliceBuilder[*[]string]{
		field:    fmapField,
		appendFn: s.appendFn,
		h:        s.h,
	}
}
