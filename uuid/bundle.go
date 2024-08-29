package uuid

import (
	"github.com/google/uuid"
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

// UUIDBundle is a struct that represents a bundle of string fields.
// It provides methods for adding validation rules to the string fields.
type UUIDBundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
	h        shared.Helper
}

// NewUUIDBundle creates a new StringBundle instance.
// It takes a BundleDependencies object as an argument, which provides the necessary dependencies.
func NewUUIDBundle(deps shared.BundleDependencies) *UUIDBundle {
	return &UUIDBundle{
		appendFn: deps.AppendFn,
		storage:  deps.Fields,
		obj:      deps.Object,
		h:        deps.Helper,
	}
}

// UUID returns a UUIDBuilder instance for a string field.
// It takes a pointer to a string field as an argument.
func (s *UUIDBundle) UUID(field *uuid.UUID) UUIDBuilder[uuid.UUID] {
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

// UUIDPtr returns a UUIDBuilder instance for a pointer to a string field.
// It takes a pointer to a pointer to a string field as an argument.
func (s *UUIDBundle) UUIDPtr(field **uuid.UUID) UUIDBuilder[*uuid.UUID] {
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

// UUIDSlice returns a UUIDSliceBuilder instance for a string slice field.
// It takes a pointer to a string slice field as an argument.
func (s *UUIDBundle) UUIDSlice(field *[]uuid.UUID) UUIDSliceBuilder[[]uuid.UUID] {
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

// UUIDSlicePtr returns a UUIDSliceBuilder instance for a pointer to a string slice field.
// It takes a pointer to a pointer to a string slice field as an argument.
func (s *UUIDBundle) UUIDSlicePtr(field **[]uuid.UUID) UUIDSliceBuilder[*[]uuid.UUID] {
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
