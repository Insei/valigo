package str

import (
	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
)

type StringBundle struct {
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	storage  fmap.Storage
	obj      any
}

func NewStringBundle(obj any, appendFn func(field fmap.Field, fn shared.FieldValidationFn), fields fmap.Storage) *StringBundle {
	return &StringBundle{appendFn: appendFn, storage: fields, obj: obj}
}

func (s *StringBundle) String(field *string) StringBuilder[string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringBuilder[string]{
		field:    fmapField,
		appendFn: s.appendFn,
	}
}

func (s *StringBundle) StringPtr(field **string) StringBuilder[*string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringBuilder[*string]{
		field:    fmapField,
		appendFn: s.appendFn,
	}
}
