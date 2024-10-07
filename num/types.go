package num

import (
	"context"

	"github.com/insei/valigo/shared"
)

type numbers interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float64 | float32
}

type ptrNumbers interface {
	*int | *int8 | *int16 | *int32 | *int64 | *uint | *uint8 | *uint16 | *uint32 | *uint64 | *float64 | *float32
}

type BaseConfigurator interface {
	// Required checks if the integer is not empty.
	Required() BaseConfigurator

	// AnyOf checks if the integer is one of the allowed values.
	AnyOf(allowed ...any) BaseConfigurator

	// AnyOfInterval checks if the integer is one of the allowed values intervals.
	AnyOfInterval(begin, end any) BaseConfigurator

	// Max checks if the integer is not greater than the given maximum number.
	Max(any) BaseConfigurator

	// Min checks if the integer is not less than the given minimum number.
	Min(any) BaseConfigurator

	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error) BaseConfigurator
	When(whenFn func(ctx context.Context, value any) bool) BaseConfigurator
}

// BundleConfigurator is a builder interface for a bundle of uint fields.
// It provides methods for adding integer fields to the bundle.
type BundleConfigurator interface {
	Number(fieldPtr any) BaseConfigurator
}
