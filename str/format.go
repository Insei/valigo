package str

import (
	"net/netip"
	"strings"
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
