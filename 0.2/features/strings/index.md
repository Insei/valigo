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

Slices of strings are converted by `str.Slice`, see [Slices](https://insei.github.io/valigo/0.2/features/slices/index.md).

```go
b.Slice(&u.Emails).As(str.Slice).MaxLen(3).Email()
```

## Rules

| Rule                                          | Fails when                                                                                                                                           | Message                                                                                                                                 |
| --------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| `Required()`                                  | the value is absent or empty                                                                                                                         | Should be fulfilled                                                                                                                     |
| `MinLen(n)`                                   | the length in bytes is less than `n`                                                                                                                 | Cannot be shorter than %d characters                                                                                                    |
| `MaxLen(n)`                                   | the length in bytes is greater than `n`                                                                                                              | Cannot be longer than %d characters                                                                                                     |
| `MinRunes(n)`                                 | the number of characters (runes) is less than `n`                                                                                                    | Cannot be shorter than %d characters                                                                                                    |
| `MaxRunes(n)`                                 | the number of characters (runes) is greater than `n`                                                                                                 | Cannot be longer than %d characters                                                                                                     |
| `AnyOf(values...)`                            | the value isn't one of `values`                                                                                                                      | Only %s values is allowed                                                                                                               |
| `EqualField(&m.X)`                            | the value isn't equal to the value of the string field `X`, see [Rules on other fields](https://insei.github.io/valigo/0.2/features/fields/index.md) | Should be equal to %s                                                                                                                   |
| `NotEqualField(&m.X)`                         | the value is equal to the value of the string field `X`                                                                                              | Cannot be equal to %s                                                                                                                   |
| `EqualFold(s)`                                | the value isn't equal to `s` ignoring case, by Unicode case folding                                                                                  | Should be equal to %s ignoring case                                                                                                     |
| `Prefix(s)`, `Suffix(s)`, `Contains(s)`       | the value doesn't start with, end with or contain `s`                                                                                                | Should start with %s, Should end with %s, Should contain %s                                                                             |
| `NotPrefix(s)`, `NotSuffix(s)`, `Excludes(s)` | the value starts with, ends with or contains `s`                                                                                                     | Cannot start with %s, Cannot end with %s, Cannot contain %s                                                                             |
| `ContainsAny(chars)`, `ExcludesAny(chars)`    | the value has none / one of the characters of `chars`                                                                                                | Should contain any of the characters %s, Cannot contain the characters %s                                                               |
| `Regexp(re, opts...)`                         | the value doesn't match `re`                                                                                                                         | Doesn't match required regexp pattern                                                                                                   |
| `Email()`                                     | the value isn't an email address                                                                                                                     | Should be email address                                                                                                                 |
| `URL()`                                       | the value isn't an absolute URL with a scheme and a host                                                                                             | Should be URL                                                                                                                           |
| `IP()`, `IPv4()`, `IPv6()`                    | the value isn't an IP address of the version                                                                                                         | Should be IP address, Should be IPv4 address, Should be IPv6 address                                                                    |
| `Hostname()`                                  | the value isn't a hostname by RFC 1123                                                                                                               | Should be hostname                                                                                                                      |
| `UUID()`                                      | the value isn't a UUID in the canonical form, any case                                                                                               | Should be UUID                                                                                                                          |
| `Alphanumeric()`                              | the value is empty or has characters other than ASCII letters and digits                                                                             | Should contain only letters and digits                                                                                                  |
| `Alpha()`                                     | the value is empty or has characters other than ASCII letters                                                                                        | Should contain only letters                                                                                                             |
| `Digits()`                                    | the value is empty or has characters other than ASCII digits, `0042`                                                                                 | Should contain only digits                                                                                                              |
| `Numeric()`                                   | the value isn't a decimal number with an optional sign and fraction, `-12`, `+3.5`                                                                   | Should be a number                                                                                                                      |
| `Hexadecimal()`                               | the value isn't a hexadecimal number with an optional `0x`, `ff`, `0x1A`                                                                             | Should be hexadecimal                                                                                                                   |
| `ASCII()`                                     | the value is empty or has non-ASCII characters                                                                                                       | Should contain only ASCII characters                                                                                                    |
| `Lowercase()`, `Uppercase()`                  | the value is empty or has upper / lower case letters                                                                                                 | Should be lowercase, Should be uppercase                                                                                                |
| `Base64()`, `Base64URL()`                     | the value isn't base64 of RFC 4648 with the standard alphabet and padding / with the URL alphabet, padding optional                                  | Should be base64 encoded, Should be base64url encoded                                                                                   |
| `JSON()`                                      | the value isn't a valid JSON value                                                                                                                   | Should be JSON                                                                                                                          |
| `Semver()`                                    | the value isn't a semantic version of semver.org, `1.2.3-rc.1+build.5`, without `v`                                                                  | Should be semantic version                                                                                                              |
| `CIDR()`                                      | the value isn't an IPv4 or IPv6 prefix, `10.0.0.0/8`                                                                                                 | Should be CIDR                                                                                                                          |
| `MAC()`                                       | the value isn't a MAC address of 6 or 8 bytes, `00:1a:2b:3c:4d:5e`, `00-1a-...`, `001a.2b3c.4d5e`                                                    | Should be MAC address                                                                                                                   |
| `E164()`                                      | the value isn't a phone number of E.164, `+14155550123`                                                                                              | Should be E.164 phone number                                                                                                            |
| `HexColor()`                                  | the value isn't `#rgb`, `#rgba`, `#rrggbb` or `#rrggbbaa`                                                                                            | Should be hex color                                                                                                                     |
| `Datetime(layout)`                            | the value isn't parsed by `time.Parse(layout, value)`                                                                                                | Should be date and time in %s format                                                                                                    |
| `CountryCode()`                               | the value isn't an upper case code of ISO 3166-1 alpha-2, `US`                                                                                       | Should be ISO 3166-1 alpha-2 country code                                                                                               |
| `CurrencyCode()`                              | the value isn't an upper case active code of ISO 4217, `EUR`                                                                                         | Should be ISO 4217 currency code                                                                                                        |
| `Luhn()`                                      | the value isn't 2 or more digits with a valid Luhn checksum                                                                                          | Should have valid Luhn checksum                                                                                                         |
| `CreditCard()`                                | the value isn't 12-19 digits, grouped by spaces or hyphens, with a valid Luhn checksum                                                               | Should be credit card number                                                                                                            |
| `Custom(fn)`                                  | `fn` returns errors                                                                                                                                  | your messages, see [Custom rules](https://insei.github.io/valigo/0.2/features/custom/index.md)                                          |
| `When(fn)`                                    | —                                                                                                                                                    | the next rules are applied only if `fn` returns true, see [Conditions](https://insei.github.io/valigo/0.2/features/conditions/index.md) |

`MinLen` and `MaxLen` count bytes, like `len`: `"яя"` has the length 4. `MinRunes` and `MaxRunes` count characters: `"яя"` has 2.

## Formats

The format rules don't allocate on valid values, like the other rules, and fail on empty strings: use [`When`](https://insei.github.io/valigo/0.2/features/conditions/#when-on-a-field) for optional ones.

| Rule         | Valid                                                                                                      | Invalid                                                                         |
| ------------ | ---------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| `URL()`      | `https://example.com/path?q=1`, `http://localhost:8080`, `https://[::1]/x`, `ftp://user@files.example.com` | `example.com`, `/relative`, `https://exa mple.com`, `https://example.com:99999` |
| `Hostname()` | `example.com`, `localhost`, `1.example`                                                                    | `-example.com`, `exa_mple.com`, `example.com.`                                  |
| `IP()`       | `192.168.0.1`, `2001:db8::1`                                                                               | `256.0.0.1`, `1.2.3`                                                            |
| `UUID()`     | `6ba7b810-9dad-11d1-80b4-00c04fd430c8`                                                                     | `6ba7b8109dad11d180b400c04fd430c8`, `{6ba7b810-...}`                            |

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/str"
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

The struct is validated in place, `Validate(ctx, &user)` trims `user.Name`. See the copy semantics of the values passed by value in [Errors](https://insei.github.io/valigo/0.2/features/errors/#validate-and-validatetyped).

## Regexp message

`str.WithRegexpLocaleKey` replaces the message of a `Regexp` rule, the key can have its own [translations](https://insei.github.io/valigo/0.2/features/translations/index.md):

```go
package main

import (
    "context"
    "fmt"
    "regexp"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/str"
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
