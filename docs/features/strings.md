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
| `MinLen(n)` | the length is less than `n` | Cannot be shorter than %d characters |
| `MaxLen(n)` | the length is greater than `n` | Cannot be longer than %d characters |
| `AnyOf(values...)` | the value isn't one of `values` | Only %s values is allowed |
| `Regexp(re, opts...)` | the value doesn't match `re` | Doesn't match required regexp pattern |
| `Email()` | the value isn't an email address | Should be email address |
| `Custom(fn)` | `fn` returns errors | your messages, see [Custom rules](custom.md) |
| `When(fn)` | — | the next rules are applied only if `fn` returns true, see [Conditions](conditions.md) |

The lengths are in bytes, like `len`: `"яя"` has the length 4.

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
