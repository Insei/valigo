package num

import (
	"context"

	"github.com/insei/valigo/shared"
)

// UintBuilder is a builder interface for uint fields.
// It provides methods for adding validation rules to an uint field.
type UintBuilder[T uint8 | *uint8 | uint16 | *uint16 | uint32 | *uint32 | uint64 | *uint64 | uint | *uint] interface {
	// Required checks if the uint is not empty.
	Required() UintBuilder[T]

	// AnyOf checks if the uint is one of the allowed values.
	AnyOf(allowed ...uint) UintBuilder[T]

	// AnyOfInterval checks if the uint is one of the allowed values intervals.
	AnyOfInterval(begin, end uint) UintBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) UintBuilder[T]

	// Max checks if the uint is not greater than the given maximum number.
	Max(uint) UintBuilder[T]

	// Min checks if the uint is not less than the given minimum number.
	Min(uint) UintBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(f func(ctx context.Context, value *T) bool) UintBuilder[T]
}

// UintSliceBuilder is a builder interface for uint slice fields.
// It provides methods for adding validation rules to an integer slice field.
type UintSliceBuilder[T []uint8 | *[]uint8 | []uint16 | *[]uint16 | []uint32 | *[]uint32 | []uint64 | *[]uint64 | []uint | *[]uint] interface {
	// Required checks if the slice is not empty.
	Required() UintSliceBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) UintSliceBuilder[T]

	// Max checks if the of each uint in the slice is not greater than the given maximum number.
	Max(uint) UintSliceBuilder[T]

	// Min checks if the of each uint in the slice is not less than the given minimum number.
	Min(uint) UintSliceBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(f func(ctx context.Context, value *T) bool) UintSliceBuilder[T]
}

// UintBundleBuilder is a builder interface for a bundle of uint fields.
// It provides methods for adding uint fields to the bundle.
type UintBundleBuilder interface {
	// Uint adds an uint field to the bundle.
	Uint(field any) UintBuilder[uint]

	// UintPtr adds a pointer to an uint field to the bundle.
	UintPtr(field any) UintBuilder[*uint]

	// UintSlice adds an uint slice field to the bundle.
	UintSlice(field any) UintSliceBuilder[[]uint]

	// UintSlicePtr adds a pointer to an uint slice field to the bundle.
	UintSlicePtr(field any) UintSliceBuilder[*[]uint]
}
