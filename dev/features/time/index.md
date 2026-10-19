# Time

`timex.Of` converts a `time.Time` field or a `*time.Time` field to the time configurator, `timex.Slice` — a slice of them:

```go
type Booking struct {
    CheckIn   time.Time
    CheckOut  *time.Time
    Reminders []time.Time
}

b.Field(&bk.CheckIn).As(timex.Of).Required().InFuture()
b.Field(&bk.CheckOut).As(timex.Of).After(start) // a nil CheckOut is skipped
b.Slice(&bk.Reminders).As(timex.Slice).MaxLen(3).Each(func(r *timex.BaseConfigurator) { r.InFuture() })
```

The package is named `timex` to not shadow the standard `time`.

## Rules

| Rule                                    | Fails when                                                                                                                                             | Message                                                                                                                                    |
| --------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------ |
| `Required()`                            | the value is absent or zero                                                                                                                            | Should be fulfilled                                                                                                                        |
| `Before(t)`                             | the value isn't before `t`                                                                                                                             | Should be before %s                                                                                                                        |
| `After(t)`                              | the value isn't after `t`                                                                                                                              | Should be after %s                                                                                                                         |
| `Between(from, to)`                     | the value is outside `[from, to]`, both ends are included                                                                                              | Should be between %s and %s                                                                                                                |
| `BeforeField(&m.X)`, `AfterField(&m.X)` | the value isn't before or after the time of the field `X`, see [Rules on other fields](https://insei.github.io/valigo/latest/features/fields/index.md) | Should be before %s, Should be after %s                                                                                                    |
| `InPast()`                              | the value isn't before the moment of the validation                                                                                                    | Should be in the past                                                                                                                      |
| `InFuture()`                            | the value isn't after the moment of the validation                                                                                                     | Should be in the future                                                                                                                    |
| `Custom(fn)`                            | `fn` returns errors                                                                                                                                    | your messages, see [Custom rules](https://insei.github.io/valigo/latest/features/custom/index.md)                                          |
| `When(fn)`                              | —                                                                                                                                                      | the next rules are applied only if `fn` returns true, see [Conditions](https://insei.github.io/valigo/latest/features/conditions/index.md) |

The times in the messages are formatted by RFC 3339. `Before`, `After` and `Between` take the bounds of the configuration, `BeforeField` and `AfterField` — of another field of the object.

## Example

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/timex"
)

type Report struct {
    From time.Time
    To   *time.Time
}

func main() {
    start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
    end := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)

    v := valigo.New()
    valigo.Configure[Report](v, func(b valigo.Configurator[Report], r *Report) {
        b.Field(&r.From).As(timex.Of).Required().Between(start, end)
        b.Field(&r.To).As(timex.Of).Before(end)
    })

    late := end.AddDate(0, 1, 0)
    fmt.Println(v.Validate(context.Background(), &Report{}))
    fmt.Println(v.Validate(context.Background(), &Report{From: start, To: &late}))
}
```

Output:

```text
[Should be fulfilled (From: 0001-01-01 00:00:00 +0000 UTC) Should be between 2026-01-01T00:00:00Z and 2026-12-31T00:00:00Z (From: 0001-01-01 00:00:00 +0000 UTC)]
[Should be before 2026-12-31T00:00:00Z (To: 2027-01-31 00:00:00 +0000 UTC)]
```
