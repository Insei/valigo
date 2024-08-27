package uuid

import (
	"github.com/google/uuid"
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

// Bundle is a struct that represents a bundle of string fields.
// It provides methods for adding validation rules to the string fields.
type Bundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
	h        shared.Helper
}

// NewBundle creates a new StringBundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewBundle(deps shared.BundleDependencies) *Bundle {
	return &Bundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

// Uuid returns a Builder instance for a string field.
// It takes a pointer to a string field as an argument.
func (s *Bundle) Uuid(field *uuid.UUID) Builder[uuid.UUID] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &uuidBuilder[uuid.UUID]{
		field:    fmapField,
		appendFn: s.appendFn,
		h:        s.h,
	}
}

// UuidPtr returns a Builder instance for a pointer to a string field.
// It takes a pointer to a pointer to a string field as an argument.
func (s *Bundle) UuidPtr(field **uuid.UUID) Builder[*uuid.UUID] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &uuidBuilder[*uuid.UUID]{
		field:    fmapField,
		appendFn: s.appendFn,
		h:        s.h,
	}
}

// UuidSlice returns a SliceBuilder instance for a string slice field.
// It takes a pointer to a string slice field as an argument.
func (s *Bundle) UuidSlice(field *[]uuid.UUID) SliceBuilder[[]uuid.UUID] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &uuidSliceBuilder[[]uuid.UUID]{
		field:    fmapField,
		appendFn: s.appendFn,
		h:        s.h,
	}
}

// UuidSlicePtr returns a SliceBuilder instance for a pointer to a string slice field.
// It takes a pointer to a pointer to a string slice field as an argument.
func (s *Bundle) UuidSlicePtr(field **[]uuid.UUID) SliceBuilder[*[]uuid.UUID] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &uuidSliceBuilder[*[]uuid.UUID]{
		field:    fmapField,
		appendFn: s.appendFn,
		h:        s.h,
	}
}
