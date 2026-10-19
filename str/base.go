package str

import (
	"context"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/insei/valigo/shared"
)

const (
	minLengthLocaleKey = "validation:string:Cannot be shorter than %d characters"
	maxLengthLocaleKey = "validation:string:Cannot be longer than %d characters"
	requiredLocaleKey  = "validation:string:Should be fulfilled"
	regexpLocaleKey    = "validation:string:Doesn't match required regexp pattern"
	anyOfLocaleKey     = "validation:string:Only %s values is allowed"
	emailLocaleKey     = "validation:string:Should be email address"
	urlLocaleKey       = "validation:string:Should be URL"
	ipLocaleKey        = "validation:string:Should be IP address"
	ipv4LocaleKey      = "validation:string:Should be IPv4 address"
	ipv6LocaleKey      = "validation:string:Should be IPv6 address"
	hostnameLocaleKey  = "validation:string:Should be hostname"
	uuidLocaleKey      = "validation:string:Should be UUID"
	alnumLocaleKey     = "validation:string:Should contain only letters and digits"
	prefixLocaleKey    = "validation:string:Should start with %s"
	suffixLocaleKey    = "validation:string:Should end with %s"
	containsLocaleKey  = "validation:string:Should contain %s"
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

// MinRunes checks if the string has at least minLen characters, unlike MinLen it counts runes, not bytes.
func (i *BaseConfigurator) MinRunes(minLen int) *BaseConfigurator {
	i.f.Check(func(v string) bool { return utf8.RuneCountInString(v) >= minLen }, minLengthLocaleKey, minLen)
	return i
}

// MaxRunes checks if the string has no more than maxLen characters, unlike MaxLen it counts runes, not bytes.
func (i *BaseConfigurator) MaxRunes(maxLen int) *BaseConfigurator {
	i.f.Check(func(v string) bool { return utf8.RuneCountInString(v) <= maxLen }, maxLengthLocaleKey, maxLen)
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

// URL checks if the string is an absolute URL with a scheme and a host: https://example.com/path?q=1.
// The host is a hostname, an IPv4 address or an IPv6 address in brackets.
func (i *BaseConfigurator) URL() *BaseConfigurator {
	i.f.Check(isURL, urlLocaleKey)
	return i
}

// IP checks if the string is an IPv4 or IPv6 address.
func (i *BaseConfigurator) IP() *BaseConfigurator {
	i.f.Check(isIP, ipLocaleKey)
	return i
}

// IPv4 checks if the string is an IPv4 address.
func (i *BaseConfigurator) IPv4() *BaseConfigurator {
	i.f.Check(isIPv4, ipv4LocaleKey)
	return i
}

// IPv6 checks if the string is an IPv6 address.
func (i *BaseConfigurator) IPv6() *BaseConfigurator {
	i.f.Check(isIPv6, ipv6LocaleKey)
	return i
}

// Hostname checks if the string is a hostname by RFC 1123: example.com, localhost, 1.example.
func (i *BaseConfigurator) Hostname() *BaseConfigurator {
	i.f.Check(isHostname, hostnameLocaleKey)
	return i
}

// UUID checks if the string is a UUID in the canonical form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, any case.
func (i *BaseConfigurator) UUID() *BaseConfigurator {
	i.f.Check(isUUID, uuidLocaleKey)
	return i
}

// Alphanumeric checks if the string is not empty and has only ASCII letters and digits.
func (i *BaseConfigurator) Alphanumeric() *BaseConfigurator {
	i.f.Check(isAlphanumeric, alnumLocaleKey)
	return i
}

// Prefix checks if the string starts with prefix.
func (i *BaseConfigurator) Prefix(prefix string) *BaseConfigurator {
	i.f.Check(func(v string) bool { return strings.HasPrefix(v, prefix) }, prefixLocaleKey, prefix)
	return i
}

// Suffix checks if the string ends with suffix.
func (i *BaseConfigurator) Suffix(suffix string) *BaseConfigurator {
	i.f.Check(func(v string) bool { return strings.HasSuffix(v, suffix) }, suffixLocaleKey, suffix)
	return i
}

// Contains checks if the string contains substr.
func (i *BaseConfigurator) Contains(substr string) *BaseConfigurator {
	i.f.Check(func(v string) bool { return strings.Contains(v, substr) }, containsLocaleKey, substr)
	return i
}

// When allows for conditional validation logic to be applied to the string value.
func (i *BaseConfigurator) When(whenFn func(ctx context.Context, value string) bool) *BaseConfigurator {
	return NewBaseConfigurator(i.f.When(whenFn))
}
