package str

import (
	"context"

	"github.com/insei/valigo/helper"
)

type StringBuilder[T string | *string] interface {
	Trim() StringBuilder[T]
	Required() StringBuilder[T]
	AnyOf(vals ...string) StringBuilder[T]
	Custom(f func(ctx context.Context, h *helper.Helper, value *T) []error) StringBuilder[T]
	//MaxLen(uint) StringBuilder[T]
	//MinLen(uint) StringBuilder[T]
	When(f func(ctx context.Context, value *T) bool) StringBuilder[T]
}

type StringsBundleBuilder interface {
	String(field *string) StringBuilder[string]
	StringPtr(field **string) StringBuilder[*string]
}
