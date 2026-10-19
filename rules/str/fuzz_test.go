package str

import (
	"encoding/base64"
	"net"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

// The formats are checked against the parsers of the standard library and the regexps of their definitions,
// the patterns are the patterns of JSON Schema of the rules, so the schema describes the same strings.
var (
	uuidPattern     = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	hostnamePattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	e164Pattern     = regexp.MustCompile(`^\+[1-9][0-9]{1,14}$`)
	hexColorPattern = regexp.MustCompile(`^#([0-9a-fA-F]{3,4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	alnumPattern    = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	alphaPattern    = regexp.MustCompile(`^[a-zA-Z]+$`)
	digitsPattern   = regexp.MustCompile(`^[0-9]+$`)
	numericPattern  = regexp.MustCompile(`^[-+]?[0-9]+(\.[0-9]+)?$`)
	hexPattern      = regexp.MustCompile(`^(0[xX])?[0-9a-fA-F]+$`)
)

var fuzzSeeds = []string{
	"", "a", "A1", "0x1f", "-1.5", "+3", "1.", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "example.com", "-a.com", "a..b",
	"+14155550123", "+0123", "#fff", "#ffff1", "QUJD", "QUJ=", "QQ==", "QQ", "a-_b", "00:1a:2b:3c:4d:5e", "001a.2b3c.4d5e",
	"00-1a-2b-3c-4d-5e-6f-70", "https://example.com/x?y=1", "http://[::1]:80", "79927398713", "4111 1111 1111 1111", "é",
}

func seed(f *testing.F) {
	for _, s := range fuzzSeeds {
		f.Add(s)
	}
}

func FuzzPatterns(f *testing.F) {
	seed(f)
	f.Fuzz(func(t *testing.T, s string) {
		for _, c := range []struct {
			name    string
			fn      func(string) bool
			pattern *regexp.Regexp
		}{
			{"uuid", isUUID, uuidPattern}, {"e164", isE164, e164Pattern}, {"hex color", isHexColor, hexColorPattern},
			{"alphanumeric", isAlphanumeric, alnumPattern}, {"alpha", isAlpha, alphaPattern}, {"digits", isDigits, digitsPattern},
			{"numeric", isNumeric, numericPattern}, {"hexadecimal", isHexadecimal, hexPattern},
		} {
			if got, want := c.fn(s), c.pattern.MatchString(s); got != want {
				t.Errorf("%s(%q) = %v, the pattern %v", c.name, s, got, want)
			}
		}
		if got, want := isHostname(s), len(s) <= 253 && hostnamePattern.MatchString(s); got != want {
			t.Errorf("hostname(%q) = %v, want %v", s, got, want)
		}
	})
}

func FuzzEncodings(f *testing.F) {
	seed(f)
	f.Fuzz(func(t *testing.T, s string) {
		// the decoders of the standard library skip new lines
		lines := strings.ContainsAny(s, "\r\n")
		_, err := base64.StdEncoding.DecodeString(s)
		if got, want := isBase64(s), s != "" && !lines && err == nil; got != want {
			t.Errorf("base64(%q) = %v, want %v", s, got, want)
		}
		_, padded := base64.URLEncoding.DecodeString(s)
		_, raw := base64.RawURLEncoding.DecodeString(s)
		if got, want := isBase64URL(s), s != "" && !lines && (padded == nil || raw == nil); got != want {
			t.Errorf("base64url(%q) = %v, want %v", s, got, want)
		}
		// MAC accepts the forms with separators, net.ParseMAC also 12 hex digits without them
		mac, err := net.ParseMAC(s)
		separated := strings.ContainsAny(s, ":-.")
		if got, want := isMAC(s), err == nil && separated && (len(mac) == 6 || len(mac) == 8); got != want {
			t.Errorf("mac(%q) = %v, want %v", s, got, want)
		}
	})
}

// luhnOf is the reference Luhn algorithm of the digits of s.
func luhnOf(digits string) bool {
	sum := 0
	for i := range len(digits) {
		d := int(digits[len(digits)-1-i] - '0')
		if i%2 == 1 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return sum%10 == 0
}

func FuzzChecksums(f *testing.F) {
	seed(f)
	f.Fuzz(func(t *testing.T, s string) {
		if got, want := isLuhn(s), len(s) >= 2 && digitsPattern.MatchString(s) && luhnOf(s); got != want {
			t.Errorf("luhn(%q) = %v, want %v", s, got, want)
		}
		if isCreditCard(s) {
			digits := strings.NewReplacer(" ", "", "-", "").Replace(s)
			if !digitsPattern.MatchString(digits) || len(digits) < 12 || len(digits) > 19 || !luhnOf(digits) {
				t.Errorf("credit card(%q) is valid with the digits %q", s, digits)
			}
		}
	})
}

func FuzzURL(f *testing.F) {
	seed(f)
	f.Fuzz(func(t *testing.T, s string) {
		if !isURL(s) {
			return
		}
		// a valid URL is parsed by the standard library with a scheme and a host
		u, err := url.Parse(s)
		if err != nil || u.Scheme == "" || u.Host == "" {
			t.Errorf("url(%q) is valid, url.Parse: %v, %#v", s, err, u)
		}
	})
}
