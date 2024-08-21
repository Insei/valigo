package shared

import (
	"context"

	"github.com/insei/fmap/v3"
)

// FieldValidationFn is a function type that represents a field validation function.
// It takes a context, a Helper implementation, and a value, and returns a slice of errors.
type FieldValidationFn func(ctx context.Context, h Helper, v any) []Error

// BundleDependencies is a struct that represents a bundle of dependencies for field validation.
type BundleDependencies struct {
	// Object is the object being validated.
	Object any
	// Helper is the Helper implementation used for validation.
	Helper Helper
	// AppendFn is a function that appends a field validation function to the bundle.
	// It takes a fmap.Field and a FieldValidationFn as arguments.
	AppendFn func(field fmap.Field, fn FieldValidationFn)
	// Fields is the storage for the fields being validated.
	Fields fmap.Storage
}
