# Types and conversions

The configurator knows nothing about the types of the values. It has two methods for the fields:

- `b.Field(&m.X)` returns `*shared.Field[T]`, the base configurator of a field of type `T`;
- `b.Slice(&m.X)` returns `*shared.SliceFieldConfigurator[E]` of a `[]E` field, with the rules of the slice itself: `Required`, `MinLen`, `MaxLen`, `Custom`, `When` and `Each`.

`As` converts a base configurator to the configurator of its type. A conversion is a function with a type constraint, so a conversion that doesn't fit the field doesn't compile:

```go
b.Field(&u.Name).As(str.Of).Required()       // string
b.Field(&u.Discount).As(num.Ptr).Max(50)     // *int
b.Slice(&u.Tags).MaxLen(5).As(str.Slice).Email()

b.Field(&u.Age).As(str.Of)
// int does not satisfy ~string | *string (int missing in ~string | *string)
```

## Built-in conversions

| Conversion        | Fields                                                | Configurator                                                                                                                  |
| ----------------- | ----------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `str.Of`          | `string`, `*string`, named strings                    | [strings](https://insei.github.io/valigo/latest/features/strings/index.md)                                                    |
| `num.Of`          | numbers and named numbers                             | [numbers](https://insei.github.io/valigo/latest/features/numbers/index.md)                                                    |
| `num.Ptr`         | pointers to numbers                                   | [numbers](https://insei.github.io/valigo/latest/features/numbers/index.md), `Required` checks only that the pointer isn't nil |
| `timex.Of`        | `time.Time`, `*time.Time`                             | [time](https://insei.github.io/valigo/latest/features/time/index.md)                                                          |
| `vuuid.Of`        | `uuid.UUID`, `*uuid.UUID` of `github.com/google/uuid` | [uuid](https://insei.github.io/valigo/latest/features/uuid/index.md), the separate module `github.com/insei/valigo/uuid`      |
| `str.Slice`       | `Slice` of strings, `*string` or named strings        | [slices](https://insei.github.io/valigo/latest/features/slices/index.md) with the string rules of elements                    |
| `timex.Slice`     | `Slice` of `time.Time`, `*time.Time`                  | [slices](https://insei.github.io/valigo/latest/features/slices/index.md) with the time rules of elements                      |
| `vuuid.Slice`     | `Slice` of `uuid.UUID`, `*uuid.UUID`                  | [slices](https://insei.github.io/valigo/latest/features/slices/index.md) with the uuid rules of elements                      |
| `shared.SlicePtr` | `Field` of `*[]E`                                     | the slice configurator, a nil pointer is absent                                                                               |

The elements of a slice are converted in `Each`: `b.Slice(&m.Sizes).Each(func(e *shared.Field[int]) { e.As(num.Of).Min(1) })`.

## The base rules

`shared.Field[T]` has the primitives all configurators are built of. They are enough for a one-off rule without a conversion:

| Rule                      | Behavior                                                                                                                                       |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `Check(ok, key, args...)` | an error with the message `key` when `ok(value)` is false, absent values are skipped                                                           |
| `Required(ok, key)`       | an error when the value is absent or `ok(value)` is false                                                                                      |
| `Modify(fn)`              | changes the value in place, e.g. trims it                                                                                                      |
| `Custom(fn)`              | your function with the typed value, see [Custom rules](https://insei.github.io/valigo/latest/features/custom/index.md)                         |
| `When(cond)`              | the next rules are applied only if `cond(value)` is true, see [Conditions](https://insei.github.io/valigo/latest/features/conditions/index.md) |

```go
b.Field(&o.Paid).Check(func(paid bool) bool { return paid }, "Should be paid")
```

## Your types

A conversion of your type is a function that takes the base configurator and returns your configurator over `shared.Field` of the type. Use the conversions of `shared` for pointers and named types:

- `shared.Deref` converts the rules of `*T` to the rules of `T`, a nil pointer is absent;
- `shared.Retype` converts the rules of `T` to the rules of `U` with the same memory layout, e.g. of a named type;
- `shared.RetypeSlice` does the same for a slice configurator, also for the elements stored by pointer;
- `shared.SliceRules` gives `Required`, `MinLen`, `MaxLen`, `Custom` and `When` of the slice to your slice configurator, see the source of `str.Slice`.

A configurator of `time.Time` (the built-in one is [`timex`](https://insei.github.io/valigo/latest/features/time/index.md), this one shows how it's made):

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/insei/valigo"
    "github.com/insei/valigo/shared"
)

// TimeConfigurator configures the rules of a time.Time value.
type TimeConfigurator struct {
    f *shared.Field[time.Time]
}

// TimeOf converts a time.Time or *time.Time field: b.Field(&m.At).As(TimeOf).
func TimeOf[T time.Time | *time.Time](f *shared.Field[T]) *TimeConfigurator {
    if p, ok := any(f).(*shared.Field[*time.Time]); ok {
        return &TimeConfigurator{f: shared.Deref(p)}
    }
    return &TimeConfigurator{f: shared.Retype[T, time.Time](f)}
}

func (c *TimeConfigurator) Required() *TimeConfigurator {
    c.f.Required(func(v time.Time) bool { return !v.IsZero() }, "Should be fulfilled")
    return c
}

func (c *TimeConfigurator) After(t time.Time) *TimeConfigurator {
    c.f.Check(func(v time.Time) bool { return v.After(t) }, "Should be after %s", t.Format(time.DateOnly))
    return c
}

type Event struct {
    StartsAt time.Time
    EndsAt   *time.Time
}

func main() {
    start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

    v := valigo.New()
    valigo.Configure[Event](v, func(b valigo.Configurator[Event], e *Event) {
        b.Field(&e.StartsAt).As(TimeOf).Required().After(start)
        b.Field(&e.EndsAt).As(TimeOf).After(start) // a nil EndsAt is skipped
    })

    before := start.AddDate(0, 0, -1)
    fmt.Println(v.Validate(context.Background(), &Event{}))
    fmt.Println(v.Validate(context.Background(), &Event{StartsAt: start.AddDate(0, 1, 0), EndsAt: &before}))
}
```

Output:

```text
[Should be fulfilled (StartsAt: 0001-01-01 00:00:00 +0000 UTC) Should be after 2026-01-01 (StartsAt: 0001-01-01 00:00:00 +0000 UTC)]
[Should be after 2026-01-01 (EndsAt: 2025-12-31 00:00:00 +0000 UTC)]
```

The type constraint of `TimeOf` makes `b.Field(&e.Title).As(TimeOf)` a compile error for a `string` field. The rules of another uuid package are added the same way, see the source of [`github.com/insei/valigo/uuid`](https://github.com/Insei/valigo/tree/main/uuid).
