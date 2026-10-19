package str

import (
	"context"
	"regexp"
	"slices"
	"strings"
	"time"
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
	equalFieldKey      = "validation:string:Should be equal to %s"
	base64LocaleKey    = "validation:string:Should be base64 encoded"
	base64URLLocaleKey = "validation:string:Should be base64url encoded"
	jsonLocaleKey      = "validation:string:Should be JSON"
	semverLocaleKey    = "validation:string:Should be semantic version"
	cidrLocaleKey      = "validation:string:Should be CIDR"
	macLocaleKey       = "validation:string:Should be MAC address"
	e164LocaleKey      = "validation:string:Should be E.164 phone number"
	hexColorLocaleKey  = "validation:string:Should be hex color"
	asciiLocaleKey     = "validation:string:Should contain only ASCII characters"
	lowercaseLocaleKey = "validation:string:Should be lowercase"
	uppercaseLocaleKey = "validation:string:Should be uppercase"
	datetimeLocaleKey  = "validation:string:Should be date and time in %s format"
	countryLocaleKey   = "validation:string:Should be ISO 3166-1 alpha-2 country code"
	currencyLocaleKey  = "validation:string:Should be ISO 4217 currency code"
	luhnLocaleKey      = "validation:string:Should have valid Luhn checksum"
	cardLocaleKey      = "validation:string:Should be credit card number"
)

var emailRegexp = regexp.MustCompile(`^[^\s,@#$%^&*!()]+@([a-zA-Z0-9]+[.])+[a-zA-Z]{2,8}$`)

// datetimeFormats are the formats of JSON Schema of the layouts of Datetime.
var datetimeFormats = map[string]string{time.RFC3339: "date-time", time.RFC3339Nano: "date-time", time.DateOnly: "date"}

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
	i.f.Schema("maxLength", maxLen) // the bytes are not less than the characters, the schema is not stricter
	i.f.CheckMsg(func(v string) bool { return len(v) <= maxLen }, shared.M("string.max_len", maxLengthLocaleKey, shared.P("max", maxLen)))
	return i
}

// MinLen checks if the string length is not less than the given minimum length.
func (i *BaseConfigurator) MinLen(minLen int) *BaseConfigurator {
	i.f.CheckMsg(func(v string) bool { return len(v) >= minLen }, shared.M("string.min_len", minLengthLocaleKey, shared.P("min", minLen)))
	return i
}

// MinRunes checks if the string has at least minLen characters, unlike MinLen it counts runes, not bytes.
func (i *BaseConfigurator) MinRunes(minLen int) *BaseConfigurator {
	i.f.Schema("minLength", minLen)
	i.f.CheckMsg(func(v string) bool { return utf8.RuneCountInString(v) >= minLen }, shared.M("string.min_runes", minLengthLocaleKey, shared.P("min", minLen)))
	return i
}

// MaxRunes checks if the string has no more than maxLen characters, unlike MaxLen it counts runes, not bytes.
func (i *BaseConfigurator) MaxRunes(maxLen int) *BaseConfigurator {
	i.f.Schema("maxLength", maxLen)
	i.f.CheckMsg(func(v string) bool { return utf8.RuneCountInString(v) <= maxLen }, shared.M("string.max_runes", maxLengthLocaleKey, shared.P("max", maxLen)))
	return i
}

// Required checks if the string is present and not empty.
func (i *BaseConfigurator) Required() *BaseConfigurator {
	i.f.Schema(shared.KeywordRequired, true).Schema("minLength", 1)
	i.f.RequiredMsg(func(v string) bool { return len(v) > 0 }, shared.M("string.required", requiredLocaleKey))
	return i
}

// Custom allows for custom validation logic to be applied to the string value.
func (i *BaseConfigurator) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value string) []shared.Error) *BaseConfigurator {
	i.f.Custom(f)
	return i
}

// Regexp checks if the string value matches the given regular expression.
func (i *BaseConfigurator) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) *BaseConfigurator {
	i.f.Schema("pattern", regexp.String())
	options := regexpOptions{
		localeKey: regexpLocaleKey,
	}
	for _, opt := range opts {
		opt.apply(&options)
	}
	i.f.CheckMsg(regexp.MatchString, shared.M("string.regexp", options.localeKey))
	return i
}

// AnyOf checks if the string value is one of the allowed values.
func (i *BaseConfigurator) AnyOf(allowed ...string) *BaseConfigurator {
	i.f.Schema("enum", shared.Values(allowed))
	i.f.CheckMsg(func(v string) bool { return slices.Contains(allowed, v) }, shared.M("string.any_of", anyOfLocaleKey, shared.P("values", allowed)))
	return i
}

// Email checks if the string value is an email address.
func (i *BaseConfigurator) Email() *BaseConfigurator {
	i.f.Schema("format", "email")
	i.f.CheckMsg(emailRegexp.MatchString, shared.M("string.email", emailLocaleKey))
	return i
}

// URL checks if the string is an absolute URL with a scheme and a host: https://example.com/path?q=1.
// The host is a hostname, an IPv4 address or an IPv6 address in brackets.
func (i *BaseConfigurator) URL() *BaseConfigurator {
	i.f.Schema("format", "uri")
	i.f.CheckMsg(isURL, shared.M("string.url", urlLocaleKey))
	return i
}

// IP checks if the string is an IPv4 or IPv6 address.
func (i *BaseConfigurator) IP() *BaseConfigurator {
	i.f.CheckMsg(isIP, shared.M("string.ip", ipLocaleKey))
	return i
}

// IPv4 checks if the string is an IPv4 address.
func (i *BaseConfigurator) IPv4() *BaseConfigurator {
	i.f.Schema("format", "ipv4")
	i.f.CheckMsg(isIPv4, shared.M("string.ipv4", ipv4LocaleKey))
	return i
}

// IPv6 checks if the string is an IPv6 address.
func (i *BaseConfigurator) IPv6() *BaseConfigurator {
	i.f.Schema("format", "ipv6")
	i.f.CheckMsg(isIPv6, shared.M("string.ipv6", ipv6LocaleKey))
	return i
}

// Hostname checks if the string is a hostname by RFC 1123: example.com, localhost, 1.example.
func (i *BaseConfigurator) Hostname() *BaseConfigurator {
	i.f.Schema("format", "hostname")
	i.f.CheckMsg(isHostname, shared.M("string.hostname", hostnameLocaleKey))
	return i
}

// UUID checks if the string is a UUID in the canonical form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, any case.
func (i *BaseConfigurator) UUID() *BaseConfigurator {
	i.f.Schema("format", "uuid")
	i.f.CheckMsg(isUUID, shared.M("string.uuid", uuidLocaleKey))
	return i
}

// Alphanumeric checks if the string is not empty and has only ASCII letters and digits.
func (i *BaseConfigurator) Alphanumeric() *BaseConfigurator {
	i.f.Schema("pattern", "^[a-zA-Z0-9]+$")
	i.f.CheckMsg(isAlphanumeric, shared.M("string.alphanumeric", alnumLocaleKey))
	return i
}

// EqualField checks if the string is equal to the value of another string field of the object:
// b.Field(&m.PasswordConfirmation).As(str.Of).EqualField(&m.Password).
func (i *BaseConfigurator) EqualField(other *string) *BaseConfigurator {
	i.f.CheckWithMsg(other, func(v, o string) bool { return v == o }, shared.M("string.equal_field", equalFieldKey))
	return i
}

// Prefix checks if the string starts with prefix.
func (i *BaseConfigurator) Prefix(prefix string) *BaseConfigurator {
	i.f.Schema("pattern", "^"+regexp.QuoteMeta(prefix))
	i.f.CheckMsg(func(v string) bool { return strings.HasPrefix(v, prefix) }, shared.M("string.prefix", prefixLocaleKey, shared.P("prefix", prefix)))
	return i
}

// Suffix checks if the string ends with suffix.
func (i *BaseConfigurator) Suffix(suffix string) *BaseConfigurator {
	i.f.Schema("pattern", regexp.QuoteMeta(suffix)+"$")
	i.f.CheckMsg(func(v string) bool { return strings.HasSuffix(v, suffix) }, shared.M("string.suffix", suffixLocaleKey, shared.P("suffix", suffix)))
	return i
}

// Contains checks if the string contains substr.
func (i *BaseConfigurator) Contains(substr string) *BaseConfigurator {
	i.f.Schema("pattern", regexp.QuoteMeta(substr))
	i.f.CheckMsg(func(v string) bool { return strings.Contains(v, substr) }, shared.M("string.contains", containsLocaleKey, shared.P("substring", substr)))
	return i
}

// Base64 checks if the string is base64 of RFC 4648 with padding, the standard alphabet.
func (i *BaseConfigurator) Base64() *BaseConfigurator {
	i.f.Schema("contentEncoding", "base64")
	i.f.CheckMsg(isBase64, shared.M("string.base64", base64LocaleKey))
	return i
}

// Base64URL checks if the string is base64 of RFC 4648 with the URL and filename safe alphabet, padding is optional.
func (i *BaseConfigurator) Base64URL() *BaseConfigurator {
	i.f.Schema("contentEncoding", "base64url")
	i.f.CheckMsg(isBase64URL, shared.M("string.base64url", base64URLLocaleKey))
	return i
}

// JSON checks if the string is a valid JSON value.
func (i *BaseConfigurator) JSON() *BaseConfigurator {
	i.f.CheckMsg(isJSON, shared.M("string.json", jsonLocaleKey))
	return i
}

// Semver checks if the string is a semantic version of semver.org without the "v" prefix: 1.2.3-rc.1+build.5.
func (i *BaseConfigurator) Semver() *BaseConfigurator {
	i.f.Schema("pattern", semverRegexp.String())
	i.f.CheckMsg(semverRegexp.MatchString, shared.M("string.semver", semverLocaleKey))
	return i
}

// CIDR checks if the string is an IPv4 or IPv6 prefix in the CIDR notation: 10.0.0.0/8, 2001:db8::/32.
func (i *BaseConfigurator) CIDR() *BaseConfigurator {
	i.f.CheckMsg(isCIDR, shared.M("string.cidr", cidrLocaleKey))
	return i
}

// MAC checks if the string is a MAC address of 6 or 8 bytes: 00:1a:2b:3c:4d:5e, 00-1a-2b-3c-4d-5e, 001a.2b3c.4d5e.
func (i *BaseConfigurator) MAC() *BaseConfigurator {
	i.f.CheckMsg(isMAC, shared.M("string.mac", macLocaleKey))
	return i
}

// E164 checks if the string is a phone number of E.164: +14155550123.
func (i *BaseConfigurator) E164() *BaseConfigurator {
	i.f.Schema("pattern", `^\+[1-9][0-9]{1,14}$`)
	i.f.CheckMsg(isE164, shared.M("string.e164", e164LocaleKey))
	return i
}

// HexColor checks if the string is a hex color: #rgb, #rgba, #rrggbb or #rrggbbaa.
func (i *BaseConfigurator) HexColor() *BaseConfigurator {
	i.f.Schema("pattern", `^#([0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	i.f.CheckMsg(isHexColor, shared.M("string.hex_color", hexColorLocaleKey))
	return i
}

// ASCII checks if the string is not empty and has only ASCII characters.
func (i *BaseConfigurator) ASCII() *BaseConfigurator {
	i.f.CheckMsg(isASCII, shared.M("string.ascii", asciiLocaleKey))
	return i
}

// Lowercase checks if the string is not empty and has no upper case letters.
func (i *BaseConfigurator) Lowercase() *BaseConfigurator {
	i.f.CheckMsg(isLowercase, shared.M("string.lowercase", lowercaseLocaleKey))
	return i
}

// Uppercase checks if the string is not empty and has no lower case letters.
func (i *BaseConfigurator) Uppercase() *BaseConfigurator {
	i.f.CheckMsg(isUppercase, shared.M("string.uppercase", uppercaseLocaleKey))
	return i
}

// Datetime checks if the string is a date or a time of the layout of time.Parse: Datetime(time.DateOnly).
func (i *BaseConfigurator) Datetime(layout string) *BaseConfigurator {
	if format, ok := datetimeFormats[layout]; ok {
		i.f.Schema("format", format)
	}
	i.f.CheckMsg(isDatetime(layout), shared.M("string.datetime", datetimeLocaleKey, shared.P("layout", layout)))
	return i
}

// CountryCode checks if the string is an upper case code of ISO 3166-1 alpha-2: US, DE.
func (i *BaseConfigurator) CountryCode() *BaseConfigurator {
	i.f.CheckMsg(isCountryCode, shared.M("string.country_code", countryLocaleKey))
	return i
}

// CurrencyCode checks if the string is an upper case active code of ISO 4217: USD, EUR.
func (i *BaseConfigurator) CurrencyCode() *BaseConfigurator {
	i.f.CheckMsg(isCurrencyCode, shared.M("string.currency_code", currencyLocaleKey))
	return i
}

// Luhn checks if the string is at least 2 digits with a valid checksum of the Luhn algorithm.
func (i *BaseConfigurator) Luhn() *BaseConfigurator {
	i.f.CheckMsg(isLuhn, shared.M("string.luhn", luhnLocaleKey))
	return i
}

// CreditCard checks if the string is a card number of 12-19 digits with a valid Luhn checksum,
// the digits can be grouped by spaces or hyphens.
func (i *BaseConfigurator) CreditCard() *BaseConfigurator {
	i.f.CheckMsg(isCreditCard, shared.M("string.credit_card", cardLocaleKey))
	return i
}

// When allows for conditional validation logic to be applied to the string value.
func (i *BaseConfigurator) When(whenFn func(ctx context.Context, value string) bool) *BaseConfigurator {
	return NewBaseConfigurator(i.f.When(whenFn))
}
