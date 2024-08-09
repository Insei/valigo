package valigo

import "github.com/insei/fmap/v3"

type stringBundle struct {
	h        Helper
	appendFn func(field fmap.Field, fn func(h Helper, v any) []error)
	storage  fmap.Storage
	obj      any
}

func newStringBundle(obj any, h Helper, appendFn func(field fmap.Field, fn func(h Helper, v any) []error), fields fmap.Storage) *stringBundle {
	return &stringBundle{appendFn: appendFn, storage: fields, h: h, obj: obj}
}

func (s *stringBundle) String(field *string) StringBuilder[string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringBuilder[string]{
		h:        s.h,
		field:    fmapField,
		appendFn: s.appendFn,
	}
}

func (s *stringBundle) StringPtr(field **string) StringBuilder[*string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringBuilder[*string]{
		h:        s.h,
		field:    fmapField,
		appendFn: s.appendFn,
	}
}
