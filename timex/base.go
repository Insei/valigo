// Package timex configures the rules of time.Time values: b.Field(&m.StartsAt).As(timex.Of).After(start).
package timex

import (
	"context"
	"time"

	"github.com/insei/valigo/shared"
)

const (
	requiredLocaleKey = "validation:time:Should be fulfilled"
	beforeLocaleKey   = "validation:time:Should be before %s"
	afterLocaleKey    = "validation:time:Should be after %s"
	betweenLocaleKey  = "validation:time:Should be between %s and %s"
	pastLocaleKey     = "validation:time:Should be in the past"
	futureLocaleKey   = "validation:time:Should be in the future"
)

// now returns the current time for InPast and InFuture, tests replace it.
var now = time.Now

// BaseConfigurator configures rules for a time.Time value.
type BaseConfigurator struct {
	f *shared.Field[time.Time]
}

// NewBaseConfigurator creates a new BaseConfigurator.
func NewBaseConfigurator(f *shared.Field[time.Time]) *BaseConfigurator {
	return &BaseConfigurator{f: f}
}

// Of converts a field of time.Time or *time.Time to the time configurator: b.Field(&m.At).As(timex.Of).
// A nil *time.Time is absent.
func Of[T time.Time | *time.Time](f *shared.Field[T]) *BaseConfigurator {
	if p, ok := any(f).(*shared.Field[*time.Time]); ok {
		return NewBaseConfigurator(shared.Deref(p))
	}
	return NewBaseConfigurator(shared.Retype[T, time.Time](f))
}

// Required checks if the time is present and not zero.
func (i *BaseConfigurator) Required() *BaseConfigurator {
	i.f.Required(func(v time.Time) bool { return !v.IsZero() }, requiredLocaleKey)
	return i
}

// Before checks if the time is before t.
func (i *BaseConfigurator) Before(t time.Time) *BaseConfigurator {
	i.f.Check(func(v time.Time) bool { return v.Before(t) }, beforeLocaleKey, t.Format(time.RFC3339))
	return i
}

// After checks if the time is after t.
func (i *BaseConfigurator) After(t time.Time) *BaseConfigurator {
	i.f.Check(func(v time.Time) bool { return v.After(t) }, afterLocaleKey, t.Format(time.RFC3339))
	return i
}

// Between checks if the time is within [from, to], both ends are included.
func (i *BaseConfigurator) Between(from, to time.Time) *BaseConfigurator {
	i.f.Check(func(v time.Time) bool { return !v.Before(from) && !v.After(to) }, betweenLocaleKey,
		from.Format(time.RFC3339), to.Format(time.RFC3339))
	return i
}

// InPast checks if the time is before the moment of the validation.
func (i *BaseConfigurator) InPast() *BaseConfigurator {
	i.f.Check(func(v time.Time) bool { return v.Before(now()) }, pastLocaleKey)
	return i
}

// InFuture checks if the time is after the moment of the validation.
func (i *BaseConfigurator) InFuture() *BaseConfigurator {
	i.f.Check(func(v time.Time) bool { return v.After(now()) }, futureLocaleKey)
	return i
}

// Custom allows for custom validation logic to be applied to the time value.
func (i *BaseConfigurator) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value time.Time) []shared.Error) *BaseConfigurator {
	i.f.Custom(f)
	return i
}

// When allows for conditional validation logic to be applied to the time value.
func (i *BaseConfigurator) When(whenFn func(ctx context.Context, value time.Time) bool) *BaseConfigurator {
	return NewBaseConfigurator(i.f.When(whenFn))
}
