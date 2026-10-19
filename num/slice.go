package num

import (
	"github.com/insei/valigo/shared"
)

// SliceConfigurator configures rules for a slice of numbers.
// Required, MinLen, MaxLen, Custom and When are applied to the slice itself, use Each for rules on elements.
type SliceConfigurator[N Number] struct {
	shared.SliceRules[*SliceConfigurator[N], N]
}

// Slice converts a slice configurator of numbers: b.Slice(&m.Sizes).As(num.Slice).
// For a slice of pointers to numbers convert the elements: b.Slice(&m.X).Each(func(e *shared.Field[*int]) { e.As(num.Ptr) }).
func Slice[N Number](s *shared.SliceFieldConfigurator[N]) *SliceConfigurator[N] {
	c := &SliceConfigurator[N]{}
	c.SliceRules = shared.NewSliceRules(s, c, Slice[N])
	return c
}

// Each configures rules applied to every number of the slice.
func (s *SliceConfigurator[N]) Each(fn func(elem *BaseConfigurator[N])) *SliceConfigurator[N] {
	fn(NewBaseConfigurator(s.Slice().EachField(), false))
	return s
}
