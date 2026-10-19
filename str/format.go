package str

import (
	"encoding/json"
	"net/netip"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// isURL reports whether s is an absolute URL: scheme://[userinfo@]host[:port][path][?query][#fragment].
// The host is a hostname, an IPv4 address or an IPv6 address in brackets. Whitespace and control characters
// are not allowed, percent escapes must be valid. It doesn't allocate on valid URLs.
func isURL(s string) bool {
	scheme, rest, ok := strings.Cut(s, "://")
	if !ok || !isScheme(scheme) {
		return false
	}
	authority, tail := rest, ""
	if end := strings.IndexAny(rest, "/?#"); end >= 0 {
		authority, tail = rest[:end], rest[end:]
	}
	if i := strings.LastIndexByte(authority, '@'); i >= 0 {
		if !isURLText(authority[:i]) {
			return false
		}
		authority = authority[i+1:]
	}
	host, port := splitHostPort(authority)
	if !isPort(port) {
		return false
	}
	if strings.HasPrefix(host, "[") {
		if !strings.HasSuffix(host, "]") {
			return false
		}
		addr, err := netip.ParseAddr(host[1 : len(host)-1])
		return err == nil && addr.Is6() && isURLText(tail)
	}
	return (isHostname(host) || isIPv4(host)) && isURLText(tail)
}

// isScheme reports whether s is a URL scheme: a letter followed by letters, digits, "+", "-" or ".".
func isScheme(s string) bool {
	if s == "" || !isASCIILetter(s[0]) {
		return false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if !isASCIILetter(c) && !isASCIIDigit(c) && c != '+' && c != '-' && c != '.' {
			return false
		}
	}
	return true
}

// splitHostPort splits the authority to the host and the port, the port is empty if there is no port.
func splitHostPort(authority string) (host, port string) {
	i := strings.LastIndexByte(authority, ':')
	// a colon inside the brackets of an IPv6 address isn't a port separator
	if i < 0 || strings.LastIndexByte(authority, ']') > i {
		return authority, ""
	}
	return authority[:i], authority[i+1:]
}

// isPort reports whether s is empty or a decimal port number up to 65535.
func isPort(s string) bool {
	if len(s) > 5 {
		return false
	}
	n := 0
	for i := 0; i < len(s); i++ {
		if !isASCIIDigit(s[i]) {
			return false
		}
		n = n*10 + int(s[i]-'0')
	}
	return n <= 65535
}

// isURLText reports whether s has no whitespace or control characters and its percent escapes are valid.
func isURLText(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c <= ' ' || c == 0x7f:
			return false
		case c == '%':
			if i+2 >= len(s) || !isHex(s[i+1]) || !isHex(s[i+2]) {
				return false
			}
			i += 2
		}
	}
	return true
}

// isHostname reports whether s is a hostname by RFC 1123: dot separated labels of letters, digits and hyphens,
// a label is 1-63 characters and doesn't start or end with a hyphen, the hostname is up to 253 characters.
func isHostname(s string) bool {
	if s == "" || len(s) > 253 {
		return false
	}
	for label := range strings.SplitSeq(s, ".") {
		if label == "" || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for i := 0; i < len(label); i++ {
			c := label[i]
			if !isASCIILetter(c) && !isASCIIDigit(c) && c != '-' {
				return false
			}
		}
	}
	return true
}

// isIP reports whether s is an IPv4 or IPv6 address.
func isIP(s string) bool {
	_, err := netip.ParseAddr(s)
	return err == nil
}

// isIPv4 reports whether s is an IPv4 address.
func isIPv4(s string) bool {
	addr, err := netip.ParseAddr(s)
	return err == nil && addr.Is4()
}

// isIPv6 reports whether s is an IPv6 address.
func isIPv6(s string) bool {
	addr, err := netip.ParseAddr(s)
	return err == nil && addr.Is6()
}

// isUUID reports whether s is a UUID in the canonical form xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx, any case.
func isUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if s[i] != '-' {
				return false
			}
			continue
		}
		if !isHex(s[i]) {
			return false
		}
	}
	return true
}

// isAlphanumeric reports whether s is not empty and has only ASCII letters and digits.
func isAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !isASCIILetter(s[i]) && !isASCIIDigit(s[i]) {
			return false
		}
	}
	return true
}

func isASCIILetter(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isHex(c byte) bool {
	return isASCIIDigit(c) || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// isBase64 reports whether s is base64 of RFC 4648 with padding, the standard alphabet.
func isBase64(s string) bool {
	return isBase64Of(s, '+', '/', true)
}

// isBase64URL reports whether s is base64 of RFC 4648 with the URL and filename safe alphabet, padding is optional.
func isBase64URL(s string) bool {
	return isBase64Of(s, '-', '_', false)
}

func isBase64Of(s string, c62, c63 byte, padded bool) bool {
	if s == "" {
		return false
	}
	data := strings.TrimRight(s, "=")
	padding := len(s) - len(data)
	if padding > 2 || (padded || padding > 0) && len(s)%4 != 0 || len(data)%4 == 1 {
		return false
	}
	for i := 0; i < len(data); i++ {
		c := data[i]
		if !isASCIILetter(c) && !isASCIIDigit(c) && c != c62 && c != c63 {
			return false
		}
	}
	return true
}

// isJSON reports whether s is a valid JSON value.
func isJSON(s string) bool {
	return s != "" && json.Valid(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// semverRegexp is the regular expression of semver.org for the versions without the "v" prefix.
var semverRegexp = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)` +
	`(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?` +
	`(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

// isCIDR reports whether s is an IPv4 or IPv6 prefix in the CIDR notation, 10.0.0.0/8.
func isCIDR(s string) bool {
	_, err := netip.ParsePrefix(s)
	return err == nil
}

// isMAC reports whether s is a MAC address of 6 or 8 bytes: 00:1a:2b:3c:4d:5e, 00-1a-2b-3c-4d-5e or 001a.2b3c.4d5e.
func isMAC(s string) bool {
	switch len(s) {
	case 17, 23:
		sep := s[2]
		if sep != ':' && sep != '-' {
			return false
		}
		for i := 0; i < len(s); i++ {
			if i%3 == 2 {
				if s[i] != sep {
					return false
				}
			} else if !isHex(s[i]) {
				return false
			}
		}
		return true
	case 14, 19:
		for i := 0; i < len(s); i++ {
			if i%5 == 4 {
				if s[i] != '.' {
					return false
				}
			} else if !isHex(s[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// isE164 reports whether s is a phone number of E.164: + and up to 15 digits, the first one isn't 0.
func isE164(s string) bool {
	if len(s) < 3 || len(s) > 16 || s[0] != '+' || s[1] == '0' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !isASCIIDigit(s[i]) {
			return false
		}
	}
	return true
}

// isHexColor reports whether s is a hex color: #rgb, #rgba, #rrggbb or #rrggbbaa.
func isHexColor(s string) bool {
	if len(s) != 4 && len(s) != 5 && len(s) != 7 && len(s) != 9 || s[0] != '#' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !isHex(s[i]) {
			return false
		}
	}
	return true
}

// isASCII reports whether s is not empty and has only ASCII characters.
func isASCII(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}

// isLowercase reports whether s is not empty and has no upper case letters.
func isLowercase(s string) bool {
	return s != "" && !strings.ContainsFunc(s, unicode.IsUpper)
}

// isUppercase reports whether s is not empty and has no lower case letters.
func isUppercase(s string) bool {
	return s != "" && !strings.ContainsFunc(s, unicode.IsLower)
}

// isDatetime returns a function that reports whether s is a date or a time of the layout of time.Parse.
func isDatetime(layout string) func(s string) bool {
	return func(s string) bool {
		_, err := time.Parse(layout, s)
		return err == nil
	}
}
