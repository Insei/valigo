package str

import (
	"strings"
	"testing"
	"time"
)

func check(t *testing.T, name string, fn func(string) bool, valid, invalid []string) {
	t.Helper()
	for _, s := range valid {
		if !fn(s) {
			t.Errorf("%s(%q) = false, want true", name, s)
		}
	}
	for _, s := range invalid {
		if fn(s) {
			t.Errorf("%s(%q) = true, want false", name, s)
		}
	}
}

func TestIsURL(t *testing.T) {
	check(t, "isURL", isURL, []string{
		"https://example.com",
		"http://example.com/",
		"https://example.com/path/to?q=1&b=2#frag",
		"https://example.com:8443/path",
		"ftp://user:pass@files.example.com/dir",
		"https://127.0.0.1:80",
		"https://[::1]:8080/x",
		"https://[2001:db8::1]",
		"git+ssh://git.example.com/repo.git",
		"https://example.com/caf%C3%A9",
		"https://example.com/путь",
		"http://localhost",
	}, []string{
		"",
		"example.com",
		"/relative/path",
		"https://",
		"https:///path",
		"1http://example.com",
		"ht~tp://example.com",
		"https://exa mple.com",
		"https://example.com/a b",
		"https://example.com/\x7f",
		"https://example.com/%zz",
		"https://example.com/%a",
		"https://example.com:99999",
		"https://example.com:8a",
		"https://example.com:123456",
		"https://[::1",
		"https://[127.0.0.1]",
		"https://[nope]",
		"https://-bad-.com",
		"https://us er@example.com",
	})
}

func TestIsHostname(t *testing.T) {
	check(t, "isHostname", isHostname, []string{
		"example.com",
		"localhost",
		"1.example",
		"a-b.c-d.example",
		strings.Repeat("a", 63) + ".com",
	}, []string{
		"",
		".example.com",
		"example..com",
		"example.com.",
		"-example.com",
		"example-.com",
		"exa_mple.com",
		strings.Repeat("a", 64) + ".com",
		strings.Repeat("a.", 127) + "aa",
	})
}

func TestIsIP(t *testing.T) {
	check(t, "isIP", isIP, []string{"192.168.0.1", "::1", "2001:db8::1"}, []string{"", "256.0.0.1", "example.com", "1.2.3"})
	check(t, "isIPv4", isIPv4, []string{"192.168.0.1", "0.0.0.0"}, []string{"::1", "2001:db8::1", "1.2.3.4.5"})
	check(t, "isIPv6", isIPv6, []string{"::1", "2001:db8::1", "::ffff:192.168.0.1"}, []string{"192.168.0.1", "2001:db8:::1"})
}

func TestIsUUID(t *testing.T) {
	check(t, "isUUID", isUUID, []string{
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		"6BA7B810-9DAD-11D1-80B4-00C04FD430C8",
		"00000000-0000-0000-0000-000000000000",
	}, []string{
		"",
		"6ba7b8109dad11d180b400c04fd430c8",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c",
		"6ba7b810-9dad-11d1-80b4_00c04fd430c8",
		"6ba7b810-9dad-11d1-80b4-00c04fd430cg",
		"{6ba7b810-9dad-11d1-80b4-00c04fd430c8}",
	})
}

func TestIsAlphanumeric(t *testing.T) {
	check(t, "isAlphanumeric", isAlphanumeric, []string{"abc", "ABC123", "0"}, []string{"", "a-b", "a b", "яя", "a_1"})
}

func TestFormatsDontAllocate(t *testing.T) {
	for name, fn := range map[string]func() bool{
		"url":      func() bool { return isURL("https://user@example.com:8443/path?q=1#f") },
		"url ipv6": func() bool { return isURL("https://[2001:db8::1]:8080/x") },
		"hostname": func() bool { return isHostname("api.example.com") },
		"ip":       func() bool { return isIP("2001:db8::1") },
		"uuid":     func() bool { return isUUID("6ba7b810-9dad-11d1-80b4-00c04fd430c8") },
	} {
		if allocs := testing.AllocsPerRun(100, func() { fn() }); allocs != 0 {
			t.Errorf("%s allocates %v times", name, allocs)
		}
	}
}

func TestMoreFormats(t *testing.T) {
	check(t, "isBase64", isBase64, []string{"Zm9v", "Zm9vYg==", "Zm9vYmE=", "+/+/"}, []string{"", "Zm9vYg", "Zm9vY===", "Zm9v_-==", "Zm9vYg=a", "Z"})
	check(t, "isBase64URL", isBase64URL, []string{"Zm9v", "Zm9vYg", "Zm9vYg==", "-_-_"}, []string{"", "Zm9v+/", "Zm9vY", "Zm9vYg="})
	check(t, "isJSON", isJSON, []string{`{"a":1}`, `[1,2]`, `"x"`, `null`, `3.14`}, []string{"", `{"a":}`, `{a:1}`, `[1,2`})
	check(t, "semver", semverRegexp.MatchString, []string{"1.2.3", "0.0.1", "1.2.3-rc.1", "1.2.3+build.5", "1.2.3-alpha.1+exp.sha.5114f85"},
		[]string{"", "v1.2.3", "1.2", "01.2.3", "1.2.3-", "1.2.3-01"})
	check(t, "isCIDR", isCIDR, []string{"10.0.0.0/8", "192.168.1.0/24", "2001:db8::/32"}, []string{"", "10.0.0.0", "10.0.0.0/33", "2001:db8::/129"})
	check(t, "isMAC", isMAC, []string{"00:1a:2b:3c:4d:5e", "00-1A-2B-3C-4D-5E", "00:1a:2b:3c:4d:5e:6f:70", "001a.2b3c.4d5e", "001a.2b3c.4d5e.6f70"},
		[]string{"", "00:1a:2b:3c:4d", "00:1a-2b:3c:4d:5e", "00;1a;2b;3c;4d;5e", "00:1a:2b:3c:4d:5g", "001a.2b3c-4d5e", "001a.2b3c.4d5g"})
	check(t, "isE164", isE164, []string{"+14155550123", "+79001234567", "+12"}, []string{"", "14155550123", "+0123", "+1", "+1234567890123456", "+1-415"})
	check(t, "isHexColor", isHexColor, []string{"#fff", "#FFFA", "#a1b2c3", "#a1b2c3d4"}, []string{"", "fff", "#ff", "#ggg", "#a1b2c3d"})
	check(t, "isASCII", isASCII, []string{"abc", "a b!~"}, []string{"", "яя", "café"})
	check(t, "isLowercase", isLowercase, []string{"abc", "a-1", "яя"}, []string{"", "aBc", "Яя"})
	check(t, "isUppercase", isUppercase, []string{"ABC", "A-1", "ЯЯ"}, []string{"", "AbC", "яЯ"})
	check(t, "isDatetime", isDatetime(time.DateOnly), []string{"2026-05-01"}, []string{"", "2026-13-01", "01.05.2026"})
}

func TestMoreFormatsDontAllocate(t *testing.T) {
	for name, fn := range map[string]func() bool{
		"base64":   func() bool { return isBase64("Zm9vYg==") },
		"json":     func() bool { return isJSON(`{"a": [1, {"b": null}]}`) },
		"semver":   func() bool { return semverRegexp.MatchString("1.2.3-rc.1") },
		"cidr":     func() bool { return isCIDR("2001:db8::/32") },
		"mac":      func() bool { return isMAC("00:1a:2b:3c:4d:5e") },
		"datetime": func() bool { return isDatetime(time.RFC3339)("2026-05-01T10:00:00+03:00") },
		"lower":    func() bool { return isLowercase("abc") },
	} {
		if allocs := testing.AllocsPerRun(100, func() { fn() }); allocs != 0 {
			t.Errorf("%s allocates %v times", name, allocs)
		}
	}
}
