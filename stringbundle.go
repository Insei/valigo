package valigo

import (
	"context"

	"github.com/insei/fmap/v3"
)

type stringBundle struct {
	appendFn func(field fmap.Field, fn func(ctx context.Context, h *Helper, v any) []error)
	storage  fmap.Storage
	obj      any
}

func newStringBundle(obj any, appendFn func(field fmap.Field, fn func(ctx context.Context, h *Helper, v any) []error), fields fmap.Storage) *stringBundle {
	return &stringBundle{appendFn: appendFn, storage: fields, obj: obj}
}

func (s *stringBundle) String(field *string) StringBuilder[string] {
	fmapField, err := s.storage.GetFieldByPtr(s.obj, field)
	if err != nil {
		panic(err)
	}
	return &stringBuilder[string]{
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
		field:    fmapField,
		appendFn: s.appendFn,
	}
}
