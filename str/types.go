package str

import (
	"context"
	"regexp"

	"github.com/insei/valigo/shared"
)

// StringBuilder is a builder interface for string fields.
// It provides methods for adding validation rules to a string field.
type StringBuilder[T string | *string] interface {
	// Trim removes leading and trailing whitespace from the string.
	Trim() StringBuilder[T]

	// Required checks if the string is not empty.
	Required() StringBuilder[T]

	// AnyOf checks if the string is one of the allowed values.
	AnyOf(allowed ...string) StringBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) StringBuilder[T]

	// Regexp checks if the string matches the given regular expression.
	Regexp(regexp *regexp.Regexp, opts ...RegexpOption) StringBuilder[T]

	// MaxLen checks if the string length is not greater than the given maximum length.
	MaxLen(int) StringBuilder[T]

	// MinLen checks if the string length is not less than the given minimum length.
	MinLen(int) StringBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(f func(ctx context.Context, value *T) bool) StringBuilder[T]
}

// StringSliceBuilder is a builder interface for string slice fields.
// It provides methods for adding validation rules to a string slice field.
type StringSliceBuilder[T []string | *[]string] interface {
	// Trim removes leading and trailing whitespace from each string in the slice.
	Trim() StringSliceBuilder[T]

	// Required checks if the slice is not empty.
	Required() StringSliceBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) StringSliceBuilder[T]

	// Regexp checks if each string in the slice matches the given regular expression.
	Regexp(regexp *regexp.Regexp, opts ...RegexpOption) StringSliceBuilder[T]

	// MaxLen checks if the length of each string in the slice is not greater than the given maximum length.
	MaxLen(int) StringSliceBuilder[T]

	// MinLen checks if the length of each string in the slice is not less than the given minimum length.
	MinLen(int) StringSliceBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(f func(ctx context.Context, value *T) bool) StringSliceBuilder[T]
}

// StringsBundleBuilder is a builder interface for a bundle of string fields.
// It provides methods for adding string fields to the bundle.
type StringsBundleBuilder interface {
	// String adds a string field to the bundle.
	String(field *string) StringBuilder[string]

	// StringPtr adds a pointer to a string field to the bundle.
	StringPtr(field **string) StringBuilder[*string]

	// StringSlice adds a string slice field to the bundle.
	StringSlice(field *[]string) StringSliceBuilder[[]string]

	// StringSlicePtr adds a pointer to a string slice field to the bundle.
	StringSlicePtr(field **[]string) StringSliceBuilder[*[]string]
}
