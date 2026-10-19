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
		b.String(&o.Comment).MaxLen(200)
		b.NumberPtr(&o.Discount).AnyOfInterval(1, 50)
		b.StringSlice(&o.Coupons).MaxLen(3)
		b.String(&o.Delivery.City).Required()
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

| Configurator | Fails on |
|---|---|
| `String` | absent value, empty string |
| `Number` | zero |
| `NumberPtr` | absent value; zero is a valid value of a pointer |
| `UUID` | absent value, `uuid.Nil` |
| `Slice`, `StringSlice`, `UUIDSlice` | absent slice, empty slice |

Nil slices of value fields (`[]string(nil)`) are present and empty: length rules check them, `Each` has nothing to check.

## Conditions

A field `When` doesn't call its condition for an absent value and skips its rules, including `Required`: the condition is about the value, and there is no value. Put `Required` before `When`, or use `When` of the struct, if the field is required under a condition. See [Conditions](conditions.md).
