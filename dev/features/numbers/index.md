# Numbers

`Number` configures the rules of a numeric field, `NumberPtr` — of a pointer to a numeric field. Supported types are `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64` and named types based on them.

The arguments of the rules have the type of the field. Untyped constants fit any numeric type, a value of another type doesn't compile:

```go
type Price int64

type Product struct {
    Stock    int
    Price    Price
    Rating   float32
    Discount *uint8
}

b.Number(&p.Stock).Min(0)
b.Number(&p.Price).AnyOfInterval(1, 1_000_000) // the arguments are Price
b.Number(&p.Rating).Max(5)
b.NumberPtr(&p.Discount).Max(90)               // a nil Discount is skipped

var limit int32 = 10
b.Number(&p.Stock).Max(limit) // doesn't compile: int32 isn't int
```

## Rules

| Rule                        | Fails when                                                    | Message                                                                                                                                    |
| --------------------------- | ------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `Required()`                | `Number`: the value is zero; `NumberPtr`: the value is absent | Should be fulfilled                                                                                                                        |
| `Min(n)`                    | the value is less than `n`                                    | Cannot be less than %v                                                                                                                     |
| `Max(n)`                    | the value is greater than `n`                                 | Cannot be greater than %v                                                                                                                  |
| `AnyOf(values...)`          | the value isn't one of `values`                               | Only %v values is allowed                                                                                                                  |
| `AnyOfInterval(begin, end)` | the value is outside `[begin, end]`, both ends are included   | Only interval[%v - %v] is allowed                                                                                                          |
| `Custom(fn)`                | `fn` returns errors                                           | your messages, see [Custom rules](https://insei.github.io/valigo/latest/features/custom/index.md)                                          |
| `When(fn)`                  | —                                                             | the next rules are applied only if `fn` returns true, see [Conditions](https://insei.github.io/valigo/latest/features/conditions/index.md) |

## Zero and absent

Zero is a valid value of a pointer: `NumberPtr(&p.Discount).Required()` passes for a pointer to `0` and fails for `nil`. Use a pointer field when zero must be distinguished from "not set", e.g. in PATCH requests:

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
)

type UpdateProduct struct {
    Stock *int
}

func main() {
    v := valigo.New()
    valigo.Configure[UpdateProduct](v, func(b valigo.Configurator[UpdateProduct], p *UpdateProduct) {
        b.NumberPtr(&p.Stock).Min(0)
    })

    zero, negative := 0, -1
    fmt.Println(v.Validate(context.Background(), &UpdateProduct{}))
    fmt.Println(v.Validate(context.Background(), &UpdateProduct{Stock: &zero}))
    fmt.Println(v.Validate(context.Background(), &UpdateProduct{Stock: &negative}))
}
```

Output:

```text
[]
[]
[Cannot be less than 0 (Stock: -1)]
```
