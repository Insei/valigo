package timex

import (
	"time"

	"github.com/insei/valigo/shared"
)

// SliceConfigurator configures rules for a slice of times.
// Required, MinLen, MaxLen, Custom and When are applied to the slice itself, use Each for rules on elements.
type SliceConfigurator struct {
	shared.SliceRules[*SliceConfigurator, time.Time]
}

// Slice converts a slice configurator of time.Time or *time.Time to the time slice configurator:
// b.Slice(&m.Dates).As(timex.Slice). Nil elements of []*time.Time are absent.
func Slice[E time.Time | *time.Time](s *shared.SliceFieldConfigurator[E]) *SliceConfigurator {
	return newSliceConfigurator(shared.RetypeSlice[E, time.Time](s))
}

func newSliceConfigurator(s *shared.SliceFieldConfigurator[time.Time]) *SliceConfigurator {
	c := &SliceConfigurator{}
	c.SliceRules = shared.NewSliceRules(s, c, newSliceConfigurator)
	return c
}

// Each configures rules applied to every time of the slice.
func (s *SliceConfigurator) Each(fn func(elem *BaseConfigurator)) *SliceConfigurator {
	fn(NewBaseConfigurator(s.Slice().EachField()))
	return s
}
