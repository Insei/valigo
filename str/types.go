package str

import (
	"context"
	"github.com/insei/valigo/shared"
	"regexp"
)

type BaseConfigurator interface {
	// Trim removes leading and trailing whitespace from the string.
	Trim() BaseConfigurator

	// Required checks if the string is not empty.
	Required() BaseConfigurator

	// AnyOf checks if the string is one of the allowed values.
	AnyOf(allowed ...string) BaseConfigurator

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error) BaseConfigurator

	// Regexp checks if the string matches the given regular expression.
	Regexp(regexp *regexp.Regexp, opts ...RegexpOption) BaseConfigurator

	// MaxLen checks if the string length is not greater than the given maximum length.
	MaxLen(int) BaseConfigurator

	// MinLen checks if the string length is not less than the given minimum length.
	MinLen(int) BaseConfigurator

	// Email checks is the string is email address
	Email() BaseConfigurator

	// When allows for conditional validation based on a given condition.
	When(whenFn func(ctx context.Context, value any) bool) BaseConfigurator
}

// StringBundleConfigurator is a builder interface for a bundle of string fields.
// It provides methods for adding string fields to the bundle.
type StringBundleConfigurator interface {
	String(fieldPtr any) BaseConfigurator
}
