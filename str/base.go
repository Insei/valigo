package str

import (
	"context"
	"regexp"
	"slices"
	"strings"

	"github.com/insei/valigo/shared"
)

const (
	minLengthLocaleKey = "validation:string:Cannot be shorter than %d characters"
	maxLengthLocaleKey = "validation:string:Cannot be longer than %d characters"
	requiredLocaleKey  = "validation:string:Should be fulfilled"
	regexpLocaleKey    = "validation:string:Doesn't match required regexp pattern"
	anyOfLocaleKey     = "validation:string:Only %s values is allowed"
	emailLocaleKey     = "validation:string:Should be email address"
)

var emailRegexp = regexp.MustCompile(`^[^\s,@#$%^&*!()]+@([a-zA-Z0-9]+[.])+[a-zA-Z]{2,8}$`)

// BaseConfigurator configures rules for a string value.
type BaseConfigurator struct {
	f *shared.Field[string]
}

// NewBaseConfigurator creates a new BaseConfigurator.
func NewBaseConfigurator(f *shared.Field[string]) *BaseConfigurator {
	return &BaseConfigurator{f: f}
}

// Of converts a field of a string, *string or named string type to the string configurator:
// b.Field(&m.Name).As(str.Of). A nil *string is absent.
func Of[S ~string | *string](f *shared.Field[S]) *BaseConfigurator {
	if p, ok := any(f).(*shared.Field[*string]); ok {
		return NewBaseConfigurator(shared.Deref(p))
	}
	return NewBaseConfigurator(shared.Retype[S, string](f))
}

// Trim removes leading and trailing whitespace from the string value.
func (i *BaseConfigurator) Trim() *BaseConfigurator {
	i.f.Modify(func(v *string) {
		*v = strings.TrimSpace(*v)
	})
	return i
}

// MaxLen checks if the string length is not greater than the given maximum length.
func (i *BaseConfigurator) MaxLen(maxLen int) *BaseConfigurator {
	i.f.Check(func(v string) bool { return len(v) <= maxLen }, maxLengthLocaleKey, maxLen)
	return i
}

// MinLen checks if the string length is not less than the given minimum length.
func (i *BaseConfigurator) MinLen(minLen int) *BaseConfigurator {
	i.f.Check(func(v string) bool { return len(v) >= minLen }, minLengthLocaleKey, minLen)
	return i
}

// Required checks if the string is present and not empty.
func (i *BaseConfigurator) Required() *BaseConfigurator {
	i.f.Required(func(v string) bool { return len(v) > 0 }, requiredLocaleKey)
	return i
}

// Custom allows for custom validation logic to be applied to the string value.
func (i *BaseConfigurator) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value string) []shared.Error) *BaseConfigurator {
	i.f.Custom(f)
	return i
}

// Regexp checks if the string value matches the given regular expression.
func (i *BaseConfigurator) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) *BaseConfigurator {
	options := regexpOptions{
		localeKey: regexpLocaleKey,
	}
	for _, opt := range opts {
		opt.apply(&options)
	}
	i.f.Check(regexp.MatchString, options.localeKey)
	return i
}

// AnyOf checks if the string value is one of the allowed values.
func (i *BaseConfigurator) AnyOf(allowed ...string) *BaseConfigurator {
	i.f.Check(func(v string) bool { return slices.Contains(allowed, v) }, anyOfLocaleKey, allowed)
	return i
}

// Email checks if the string value is an email address.
func (i *BaseConfigurator) Email() *BaseConfigurator {
	i.f.Check(emailRegexp.MatchString, emailLocaleKey)
	return i
}

// When allows for conditional validation logic to be applied to the string value.
func (i *BaseConfigurator) When(whenFn func(ctx context.Context, value string) bool) *BaseConfigurator {
	return NewBaseConfigurator(i.f.When(whenFn))
}
