# Absent values

A value is *absent* when there is nothing to check:

- a nil pointer field: `*string`, `*int`, `*uuid.UUID`;
- a nil pointer to a slice: `*[]string`;
- a field of a nested struct stored by a nil pointer: `City` of a nil `*Address`;
- a nil element of a slice of pointers: `[]*string{nil}`.

Every rule skips absent values, only `Required` fails on them. So an optional field needs no special code: the rules check it only when it is set.

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

type Address struct {
    City string
}

type Order struct {
    Comment  *string
    Discount *int
    Coupons  *[]string
    Delivery *Address
}

func main() {
    v := valigo.New()
    valigo.Configure[Order](v, func(b valigo.Configurator[Order], o *Order) {
        b.Field(&o.Comment).As(str.Of).MaxLen(200)
        b.Field(&o.Discount).As(num.Ptr).AnyOfInterval(1, 50)
        b.Field(&o.Coupons).As(shared.SlicePtr).As(str.Slice).MaxLen(3)
        b.Field(&o.Delivery.City).As(str.Of).Required()
    })

    // nothing is set: only the Required rule fails
    fmt.Println(v.Validate(context.Background(), &Order{}))

    discount := 70
    fmt.Println(v.Validate(context.Background(), &Order{Discount: &discount, Delivery: &Address{City: "Paris"}}))
}
```

Output:

```text
[Should be fulfilled (Delivery.City: <nil>)]
[Only interval[1 - 50] is allowed (Discount: 70)]
```

## Required

What `Required` checks depends on the type:

| Conversion                         | Fails on                                         |
| ---------------------------------- | ------------------------------------------------ |
| `str.Of`                           | absent value, empty string                       |
| `num.Of`                           | zero                                             |
| `num.Ptr`                          | absent value; zero is a valid value of a pointer |
| `uuid.Of`                          | absent value, `uuid.Nil`                         |
| `Slice`, `str.Slice`, `uuid.Slice` | absent slice, empty slice                        |

Nil slices of value fields (`[]string(nil)`) are present and empty: length rules check them, `Each` has nothing to check.

## Conditions

A field `When` doesn't call its condition for an absent value and skips its rules, including `Required`: the condition is about the value, and there is no value. Put `Required` before `When`, or use `When` of the struct, if the field is required under a condition. See [Conditions](https://insei.github.io/valigo/0.2/features/conditions/index.md).
