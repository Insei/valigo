package num

import (
	"context"

	"github.com/insei/valigo/shared"
)

// IntBuilder is a builder interface for integer fields.
// It provides methods for adding validation rules to an integer field.
type IntBuilder[T int | *int] interface {
	// Required checks if the integer is not empty.
	Required() IntBuilder[T]

	// AnyOf checks if the integer is one of the allowed values.
	AnyOf(allowed ...int) IntBuilder[T]

	// AnyOfInterval checks if the integer is one of the allowed values intervals.
	AnyOfInterval(begin, end int) IntBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) IntBuilder[T]

	// Max checks if the integer is not greater than the given maximum number.
	Max(int) IntBuilder[T]

	// Min checks if the integer is not less than the given minimum number.
	Min(int) IntBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(f func(ctx context.Context, value *T) bool) IntBuilder[T]
}

// IntSliceBuilder is a builder interface for integer slice fields.
// It provides methods for adding validation rules to an integer slice field.
type IntSliceBuilder[T []int | *[]int] interface {
	// Required checks if the slice is not empty.
	Required() IntSliceBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) IntSliceBuilder[T]

	// Max checks if the of each integer in the slice is not greater than the given maximum number.
	Max(int) IntSliceBuilder[T]

	// Min checks if the of each integer in the slice is not less than the given minimum number.
	Min(int) IntSliceBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(f func(ctx context.Context, value *T) bool) IntSliceBuilder[T]
}

// IntBundleBuilder is a builder interface for a bundle of integer fields.
// It provides methods for adding integer fields to the bundle.
type IntBundleBuilder interface {
	// Int adds an integer field to the bundle.
	Int(field *int) IntBuilder[int]

	// IntPtr adds a pointer to an integer field to the bundle.
	IntPtr(field **int) IntBuilder[*int]

	// IntSlice adds an integer slice field to the bundle.
	IntSlice(field *[]int) IntSliceBuilder[[]int]

	// IntSlicePtr adds a pointer to an integer slice field to the bundle.
	IntSlicePtr(field **[]int) IntSliceBuilder[*[]int]
}
