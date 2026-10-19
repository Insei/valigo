# Custom rules

## Field rules

`Custom` of a field configurator adds a rule with your logic. The function gets the typed value of the field and a helper that builds errors with the location of the field:

```go
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/insei/valigo"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
)

type Account struct {
	Login string
	Age   int
}

func main() {
	v := valigo.New()
	valigo.Configure[Account](v, func(b valigo.Configurator[Account], a *Account) {
		b.Field(&a.Login).As(str.Of).Required().Custom(func(ctx context.Context, h *shared.FieldCustomHelper, login string) []shared.Error {
			if strings.EqualFold(login, "admin") {
				return []shared.Error{h.ErrorT(ctx, login, "Login %q is reserved", login)}
			}
			return nil
		})
		b.Field(&a.Age).As(num.Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, age int) []shared.Error {
			if age%2 != 0 {
				return []shared.Error{h.ErrorT(ctx, age, "Should be even")}
			}
			return nil
		})
	})

	fmt.Println(v.Validate(context.Background(), &Account{Login: "Admin", Age: 3}))
}
```

Output:

```text
[Login "Admin" is reserved (Login: Admin) Should be even (Age: 3)]
```

`h.ErrorT(ctx, value, localeKey, args...)` translates the message: `localeKey` is looked up in the [translations](translations.md), a key without a translation is used as is, as a format for `args`. The [code](errors.md#codes-for-clients) of the error is the key. For a stable code and named parameters build the message by `shared.M` and return `h.ErrorMsg(ctx, value, msg)`:

```go
return []shared.Error{h.ErrorMsg(ctx, login, shared.M("login.reserved", "Login %q is reserved", shared.P("login", login)))}
```

The values of `shared.P` are the format arguments of the message in their order. A message without the checked value can be created once on the configuration, its parameters are shared by the errors.

The function isn't called for an [absent value](absent.md).

## Rules on slice elements

The `Each` of `Slice` gives a `*shared.Field[E]` of the elements. Besides `Custom` and `When` it has `Check`, the base of all built-in rules: a predicate and a message for the elements that don't match it.

```go
b.Slice(&o.Quantities).Each(func(q *shared.Field[int]) {
	q.Check(func(v int) bool { return v > 0 }, "Should be positive")
})
```

`Required(ok, localeKey)` fails for absent elements and the elements that don't match `ok`, `Modify(fn)` changes the elements in place. See [Slices](slices.md).

## Struct rules

`Custom` of the struct configurator gets the whole object, for the rules on several fields. `h.ErrorT` takes the pointer to the field of the *validated object* to build its location:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
)

type Booking struct {
	From time.Time
	To   time.Time
}

func main() {
	v := valigo.New()
	valigo.Configure[Booking](v, func(b valigo.Configurator[Booking], _ *Booking) {
		b.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *Booking) []shared.Error {
			if !obj.To.After(obj.From) {
				return []shared.Error{h.ErrorT(ctx, &obj.To, obj.To.Format(time.DateOnly), "Should be after %s", obj.From.Format(time.DateOnly))}
			}
			return nil
		})
	})

	from := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	fmt.Println(v.Validate(context.Background(), &Booking{From: from, To: from.AddDate(0, 0, -1)}))
}
```

Output:

```text
[Should be after 2026-05-10 (To: 2026-05-09)]
```

- Struct `Custom` can be combined with [`When`](conditions.md): `b.When(cond).Custom(fn)`.
- The helper is valid only during the call of the function, don't keep it.
- A pointer that doesn't point to a field of the validated object panics.

## Rules from data

The rules are configured once, but they don't have to be constant: a `Custom` rule is a closure, it can read its limits from a database, a cache or the settings of a tenant on every validation, and it gets the context of the request. A `When` condition can do the same. Use the built-in message keys and codes, and the errors are translated and coded like the built-in ones:

```go
package main

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
)

type tenantKey struct{}

// Limits are the settings of the tenants, a database or a cache in a real service.
type Limits map[string]int

func (l Limits) MaxTitle(ctx context.Context) int {
	tenant, _ := ctx.Value(tenantKey{}).(string)
	return l[tenant]
}

type Post struct {
	Title string
}

func main() {
	limits := Limits{"free": 10, "pro": 100}

	v := valigo.New()
	valigo.Configure[Post](v, func(b valigo.Configurator[Post], p *Post) {
		b.Field(&p.Title).As(str.Of).Required().Custom(func(ctx context.Context, h *shared.FieldCustomHelper, title string) []shared.Error {
			if maxLen := limits.MaxTitle(ctx); utf8.RuneCountInString(title) > maxLen {
				msg := shared.M("string.max_runes", "validation:string:Cannot be longer than %d characters", shared.P("max", maxLen))
				return []shared.Error{h.ErrorMsg(ctx, title, msg)}
			}
			return nil
		})
	})

	post := &Post{Title: "Rules from the data"}
	free := context.WithValue(context.Background(), tenantKey{}, "free")
	pro := context.WithValue(context.Background(), tenantKey{}, "pro")
	fmt.Println(v.Validate(free, post))
	fmt.Println(v.Validate(pro, post))
}
```

Output:

```text
[Cannot be longer than 10 characters (Title: Rules from the data)]
[]
```

The message of the error is translated by the plural forms of the built-in key, and the client gets the code `string.max_runes` with the `max` of the tenant.

