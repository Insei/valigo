# Rules on other fields

A rule can compare the value with the value of another field of the same object. The other field is passed by its pointer to the model, like the field itself, so its type is checked by the compiler and a pointer that isn't a field of the model panics in `Configure`:

```go
b.Field(&o.EndsAt).As(timex.Of).AfterField(&o.StartsAt)
b.Field(&o.Price).As(num.Of).MinField(&o.MinPrice).MaxField(&o.MaxPrice)
b.Field(&u.Confirmation).As(str.Of).EqualField(&u.Password)
b.Field(&u.Plan).As(eq.Of).EqualField(&u.CurrentPlan)
```

| Configurator | Rule | Fails when | Message |
|---|---|---|---|
| [strings](strings.md) | `EqualField(&m.X)` | the value isn't equal to `X` | Should be equal to %s |
| [numbers](numbers.md) | `MinField(&m.X)` | the value is less than `X` | Cannot be less than %v |
| [numbers](numbers.md) | `MaxField(&m.X)` | the value is greater than `X` | Cannot be greater than %v |
| [time](time.md) | `BeforeField(&m.X)` | the value isn't before `X` | Should be before %s |
| [time](time.md) | `AfterField(&m.X)` | the value isn't after `X` | Should be after %s |
| [comparable values](values.md) | `EqualField(&m.X)` | the value isn't equal to `X` | Should be equal to %s |

- The other field has the type of the configurator: `string` for `str.Of`, the number type for `num.Of`, `time.Time` for `timex.Of`.
- The message has the location of the other field, by the [location function](errors.md#location-of-fields) of the validator: `Should be after StartsAt`, or `Should be after starts_at` with JSON paths.
- The rule is skipped when the value or the value of the other field is [absent](absent.md), e.g. a field of a nil nested struct.
- In `Each` the elements are compared with a field of the object: `e.As(timex.Of).AfterField(&m.StartsAt)`.

## Example

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/insei/valigo"
	"github.com/insei/valigo/rules/num"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/rules/timex"
)

type Auction struct {
	StartsAt time.Time
	EndsAt   time.Time
	MinBid   int
	Bids     []int
	Password string
	Repeat   string
}

func main() {
	v := valigo.New()
	valigo.Configure[Auction](v, func(b valigo.Configurator[Auction], a *Auction) {
		b.Field(&a.EndsAt).As(timex.Of).AfterField(&a.StartsAt)
		b.Slice(&a.Bids).As(num.Slice).Each(func(bid *num.BaseConfigurator[int]) { bid.MinField(&a.MinBid) })
		b.Field(&a.Repeat).As(str.Of).EqualField(&a.Password)
	})

	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	auction := &Auction{
		StartsAt: start,
		EndsAt:   start.Add(-time.Hour),
		MinBid:   100,
		Bids:     []int{150, 90},
		Password: "secret",
		Repeat:   "Secret",
	}
	for _, err := range v.Validate(context.Background(), auction) {
		fmt.Println(err)
	}
}
```

Output:

```text
Should be after StartsAt (EndsAt: 2026-04-30 23:00:00 +0000 UTC)
Cannot be less than MinBid (Bids[1]: 90)
Should be equal to Password (Repeat: Secret)
```

For a comparison of several fields that doesn't fit one rule, use a [struct custom rule](custom.md#struct-rules).
