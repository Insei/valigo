package str

import (
	"strings"
	"testing"
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
