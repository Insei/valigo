# Strings

`str.Of` converts a `string` field, a `*string` field or a field of a named string type to the string configurator:

```go
type Status string

type User struct {
	Name     string
	Nickname *string
	Status   Status
}

b.Field(&u.Name).As(str.Of).Trim().Required().MaxLen(64)
b.Field(&u.Nickname).As(str.Of).MinLen(3) // a nil Nickname is skipped
b.Field(&u.Status).As(str.Of).AnyOf("active", "blocked")
```

Slices of strings are converted by `str.Slice`, see [Slices](slices.md).

```go
b.Slice(&u.Emails).As(str.Slice).MaxLen(3).Email()
```

## Rules

| Rule | Fails when | Message |
|---|---|---|
| `Required()` | the value is absent or empty | Should be fulfilled |
| `MinLen(n)` | the length in bytes is less than `n` | Cannot be shorter than %d characters |
| `MaxLen(n)` | the length in bytes is greater than `n` | Cannot be longer than %d characters |
| `MinRunes(n)` | the number of characters (runes) is less than `n` | Cannot be shorter than %d characters |
| `MaxRunes(n)` | the number of characters (runes) is greater than `n` | Cannot be longer than %d characters |
| `AnyOf(values...)` | the value isn't one of `values` | Only %s values is allowed |
| `Prefix(s)`, `Suffix(s)`, `Contains(s)` | the value doesn't start with, end with or contain `s` | Should start with %s, Should end with %s, Should contain %s |
| `Regexp(re, opts...)` | the value doesn't match `re` | Doesn't match required regexp pattern |
| `Email()` | the value isn't an email address | Should be email address |
| `URL()` | the value isn't an absolute URL with a scheme and a host | Should be URL |
| `IP()`, `IPv4()`, `IPv6()` | the value isn't an IP address of the version | Should be IP address, Should be IPv4 address, Should be IPv6 address |
| `Hostname()` | the value isn't a hostname by RFC 1123 | Should be hostname |
| `UUID()` | the value isn't a UUID in the canonical form, any case | Should be UUID |
| `Alphanumeric()` | the value is empty or has characters other than ASCII letters and digits | Should contain only letters and digits |
| `Custom(fn)` | `fn` returns errors | your messages, see [Custom rules](custom.md) |
| `When(fn)` | — | the next rules are applied only if `fn` returns true, see [Conditions](conditions.md) |

`MinLen` and `MaxLen` count bytes, like `len`: `"яя"` has the length 4. `MinRunes` and `MaxRunes` count characters: `"яя"` has 2.

## Formats

The format rules don't allocate on valid values, like the other rules, and fail on empty strings: use [`When`](conditions.md#when-on-a-field) for optional ones.

| Rule | Valid | Invalid |
|---|---|---|
| `URL()` | `https://example.com/path?q=1`, `http://localhost:8080`, `https://[::1]/x`, `ftp://user@files.example.com` | `example.com`, `/relative`, `https://exa mple.com`, `https://example.com:99999` |
| `Hostname()` | `example.com`, `localhost`, `1.example` | `-example.com`, `exa_mple.com`, `example.com.` |
| `IP()` | `192.168.0.1`, `2001:db8::1` | `256.0.0.1`, `1.2.3` |
| `UUID()` | `6ba7b810-9dad-11d1-80b4-00c04fd430c8` | `6ba7b8109dad11d180b400c04fd430c8`, `{6ba7b810-...}` |

```go
package main

import (
	"context"
	"fmt"

	"github.com/insei/valigo"
	"github.com/insei/valigo/str"
)

type Webhook struct {
	URL      string
	Host     string
	ClientIP string
	Secret   string
}

func main() {
	v := valigo.New()
	valigo.Configure[Webhook](v, func(b valigo.Configurator[Webhook], w *Webhook) {
		b.Field(&w.URL).As(str.Of).Required().URL().Prefix("https://")
		b.Field(&w.Host).As(str.Of).Hostname()
		b.Field(&w.ClientIP).As(str.Of).When(func(_ context.Context, ip string) bool { return ip != "" }).IP()
		b.Field(&w.Secret).As(str.Of).Alphanumeric().MinRunes(16)
	})

	for _, err := range v.Validate(context.Background(), &Webhook{URL: "http://example.com", Host: "my_host", Secret: "s3cr3t!"}) {
		fmt.Println(err)
	}
}
```

Output:

```text
Should start with https:// (URL: http://example.com)
Should be hostname (Host: my_host)
Should contain only letters and digits (Secret: s3cr3t!)
Cannot be shorter than 16 characters (Secret: s3cr3t!)
```

## Trim

`Trim()` removes the leading and trailing whitespace of the value **in the validated object**. Rules run in the order of the chain, so the rules after `Trim` check the trimmed value:

```go
b.Field(&u.Name).As(str.Of).Trim().Required() // "   " fails, "  Ann " becomes "Ann"
```

The struct is validated in place, `Validate(ctx, &user)` trims `user.Name`. See the copy semantics of the values passed by value in [Errors](errors.md#validate-and-validatetyped).

## Regexp message

`str.WithRegexpLocaleKey` replaces the message of a `Regexp` rule, the key can have its own [translations](translations.md):

```go
package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/insei/valigo"
	"github.com/insei/valigo/str"
)

type Product struct {
	SKU string
}

func main() {
	v := valigo.New()
	valigo.Configure[Product](v, func(b valigo.Configurator[Product], p *Product) {
		b.Field(&p.SKU).As(str.Of).Regexp(regexp.MustCompile(`^[A-Z]{3}-\d{4}$`), str.WithRegexpLocaleKey("Should look like ABC-1234"))
	})
	fmt.Println(v.Validate(context.Background(), &Product{SKU: "abc-1"}))
}
```

Output:

```text
[Should look like ABC-1234 (SKU: abc-1)]
```
