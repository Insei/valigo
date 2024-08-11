package shared

import (
	"context"

	"github.com/insei/fmap/v3"
)

type FieldValidationFn func(ctx context.Context, h Helper, v any) []Error

type BundleDependencies struct {
	Object   any
	Helper   Helper
	AppendFn func(field fmap.Field, fn FieldValidationFn)
	Fields   fmap.Storage
}
