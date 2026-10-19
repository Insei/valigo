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
    "github.com/insei/valigo/shared"
)

type Account struct {
    Login string
    Age   int
}

func main() {
    v := valigo.New()
    valigo.Configure[Account](v, func(b valigo.Configurator[Account], a *Account) {
        b.String(&a.Login).Required().Custom(func(ctx context.Context, h *shared.FieldCustomHelper, login string) []shared.Error {
            if strings.EqualFold(login, "admin") {
                return []shared.Error{h.ErrorT(ctx, login, "Login %q is reserved", login)}
            }
            return nil
        })
        b.Number(&a.Age).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, age int) []shared.Error {
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

`h.ErrorT(ctx, value, localeKey, args...)` translates the message: `localeKey` is looked up in the [translations](https://insei.github.io/valigo/latest/features/translations/index.md), a key without a translation is used as is, as a format for `args`.

The function isn't called for an [absent value](https://insei.github.io/valigo/latest/features/absent/index.md).

## Rules on slice elements

The `Each` of `Slice` gives a `*shared.Field[E]` of the elements. Besides `Custom` and `When` it has `Check`, the base of all built-in rules: a predicate and a message for the elements that don't match it.

```go
b.Slice(&o.Quantities).Each(func(q *shared.Field[int]) {
    q.Check(func(v int) bool { return v > 0 }, "Should be positive")
})
```

`Required(ok, localeKey)` fails for absent elements and the elements that don't match `ok`, `Modify(fn)` changes the elements in place. See [Slices](https://insei.github.io/valigo/latest/features/slices/index.md).

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

- Struct `Custom` can be combined with [`When`](https://insei.github.io/valigo/latest/features/conditions/index.md): `b.When(cond).Custom(fn)`.
- The helper is valid only during the call of the function, don't keep it.
- A pointer that doesn't point to a field of the validated object panics.
