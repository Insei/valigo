# Maps

`b.Map(&m.X)` returns the configurator of a `map[K]V` field: the rules of the map itself, of its keys and of its values. The keys and the values are converted to the configurators of their types by `As`, like the elements of [slices](https://insei.github.io/valigo/dev/features/slices/index.md):

```go
b.Map(&s.Labels).Required().MaxLen(10).
    EachKey(func(k *shared.Field[string]) { k.As(str.Of).Alphanumeric().MaxLen(32) }).
    EachValue(func(v *shared.Field[string]) { v.As(str.Of).MaxLen(64) })
b.Map(&s.Ports).EachValue(func(p *shared.Field[int]) { p.As(num.Of).AnyOfInterval(1, 65535) })
```

## Rules

| Rule            | Fails when                                 | Message                                                                                        |
| --------------- | ------------------------------------------ | ---------------------------------------------------------------------------------------------- |
| `Required()`    | the map is absent or empty                 | Should be fulfilled                                                                            |
| `MinLen(n)`     | the map has less than `n` elements         | Cannot contain less than %d elements                                                           |
| `MaxLen(n)`     | the map has more than `n` elements         | Cannot contain more than %d elements                                                           |
| `EachKey(fn)`   | —                                          | `fn` configures the rules of the keys                                                          |
| `EachValue(fn)` | —                                          | `fn` configures the rules of the values                                                        |
| `Custom(fn)`    | `fn` returns errors, it gets the whole map | your messages, see [Custom rules](https://insei.github.io/valigo/dev/features/custom/index.md) |
| `When(fn)`      | —                                          | the next rules are applied only if `fn` returns true for the map, a nil map for an absent one  |

- The location of an error of a key or a value is `Labels[key]`, the error value is the key or the value.
- The order of a map is random, so the errors of the keys and the errors of the values are sorted by location.
- The values of a map aren't addressable: [`Trim`](https://insei.github.io/valigo/dev/features/strings/#trim) and other modifying rules change a copy, which is checked by the next rules, the map keeps the original value.
- A nil map of a value field is present and empty, the map of a nil nested struct is [absent](https://insei.github.io/valigo/dev/features/absent/index.md).
- Valid maps are validated without allocations, the copy of the checked key or value is pooled.

## Example

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/num"
    "github.com/insei/valigo/rules/str"
    "github.com/insei/valigo/shared"
)

type Service struct {
    Labels map[string]string
    Ports  map[string]int
}

func main() {
    v := valigo.New()
    valigo.Configure[Service](v, func(b valigo.Configurator[Service], s *Service) {
        b.Map(&s.Labels).Required().MaxLen(10).
            EachKey(func(k *shared.Field[string]) { k.As(str.Of).Alphanumeric() }).
            EachValue(func(v *shared.Field[string]) { v.As(str.Of).MaxLen(5) })
        b.Map(&s.Ports).EachValue(func(p *shared.Field[int]) { p.As(num.Of).AnyOfInterval(1, 65535) })
    })

    errs := v.Validate(context.Background(), &Service{
        Labels: map[string]string{"app": "valigo-api", "tier!": "web"},
        Ports:  map[string]int{"http": 80, "admin": 0},
    })
    for _, err := range errs {
        fmt.Println(err)
    }
}
```

Output:

```text
Should contain only letters and digits (Labels[tier!]: tier!)
Cannot be longer than 5 characters (Labels[app]: valigo-api)
Only interval[1 - 65535] is allowed (Ports[admin]: 0)
```
