# Comparable values

`eq.Of` converts a field of any `comparable` type — an enum of a named type, a `bool`, a number, a string — to the configurator of equality rules. `eq.Ptr` converts a pointer to such a field, `eq.Slice` — a slice of them:

```go
type Status string

const (
    Active  Status = "active"
    Blocked Status = "blocked"
)

b.Field(&u.Status).As(eq.Of).Required().OneOf(Active, Blocked) // the arguments are Status
b.Field(&u.Agreed).As(eq.Of).Equal(true)                        // a bool
b.Field(&u.Role).As(eq.Ptr).OneOf(Admin, Editor)                // a nil Role is skipped
b.Slice(&u.Roles).As(eq.Slice).Unique()
```

For strings and numbers [`str.Of`](https://insei.github.io/valigo/latest/features/strings/index.md) and [`num.Of`](https://insei.github.io/valigo/latest/features/numbers/index.md) have more rules, `eq.Of` is for the types without their own configurator. `eq.Of` of a pointer field compares the pointers, use `eq.Ptr` to compare the values.

## Rules

| Rule                  | Fails when                                                                                                                                       | Message                                                                                                                                    |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `Required()`          | `eq.Of`: the value is zero; `eq.Ptr`: the value is absent                                                                                        | Should be fulfilled                                                                                                                        |
| `OneOf(values...)`    | the value isn't one of `values`                                                                                                                  | Only %v values is allowed                                                                                                                  |
| `NotOneOf(values...)` | the value is one of `values`                                                                                                                     | Values %v are not allowed                                                                                                                  |
| `Equal(v)`            | the value isn't `v`                                                                                                                              | Should be %v                                                                                                                               |
| `NotEqual(v)`         | the value is `v`                                                                                                                                 | Cannot be %v                                                                                                                               |
| `EqualField(&m.X)`    | the value isn't equal to the value of the field `X`, see [Rules on other fields](https://insei.github.io/valigo/latest/features/fields/index.md) | Should be equal to %s                                                                                                                      |
| `Custom(fn)`          | `fn` returns errors                                                                                                                              | your messages, see [Custom rules](https://insei.github.io/valigo/latest/features/custom/index.md)                                          |
| `When(fn)`            | —                                                                                                                                                | the next rules are applied only if `fn` returns true, see [Conditions](https://insei.github.io/valigo/latest/features/conditions/index.md) |

`eq.Slice` has the rules of the [slice itself](https://insei.github.io/valigo/latest/features/slices/index.md), `Each` with the rules above for the elements, `OneOf` for every element and `Unique`:

| Rule       | Fails when                                                               | Message                      |
| ---------- | ------------------------------------------------------------------------ | ---------------------------- |
| `Unique()` | the slice has duplicate elements, the error value is the first duplicate | Should contain unique values |

`Unique` compares the elements pairwise without allocations for slices up to 32 elements, and by a map for longer ones.

## Example

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/eq"
)

type Plan string

const (
    Free Plan = "free"
    Pro  Plan = "pro"
)

type SignUp struct {
    Plan     Plan
    Terms    bool
    Features []string
}

func main() {
    v := valigo.New()
    valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
        b.Field(&s.Plan).As(eq.Of).Required().OneOf(Free, Pro)
        b.Field(&s.Terms).As(eq.Of).Equal(true)
        b.Slice(&s.Features).As(eq.Slice).MaxLen(5).Unique()
    })

    for _, err := range v.Validate(context.Background(), &SignUp{Plan: "gold", Features: []string{"sso", "api", "sso"}}) {
        fmt.Println(err)
    }
}
```

Output:

```text
Only [free pro] values is allowed (Plan: gold)
Should be true (Terms: false)
Should contain unique values (Features: sso)
```
