package str

import (
	"regexp"

	"github.com/insei/valigo/shared"
)

// SliceConfigurator configures rules for a slice of strings.
// Required, MinLen, MaxLen, Custom and When are applied to the slice itself, use Each for rules on elements.
type SliceConfigurator struct {
	shared.SliceRules[*SliceConfigurator, string]
}

// Slice converts a slice configurator of strings, *string or a named string type to the string slice configurator:
// b.Slice(&m.Tags).As(str.Slice). Nil elements of []*string are absent.
func Slice[E ~string | *string](s *shared.SliceFieldConfigurator[E]) *SliceConfigurator {
	return newSliceConfigurator(shared.RetypeSlice[E, string](s))
}

func newSliceConfigurator(s *shared.SliceFieldConfigurator[string]) *SliceConfigurator {
	c := &SliceConfigurator{}
	c.SliceRules = shared.NewSliceRules(s, c, newSliceConfigurator)
	return c
}

// Each configures rules applied to every string of the slice.
func (s *SliceConfigurator) Each(fn func(elem *BaseConfigurator)) *SliceConfigurator {
	fn(NewBaseConfigurator(s.Slice().EachField()))
	return s
}

// Trim removes leading and trailing whitespace from every string of the slice.
func (s *SliceConfigurator) Trim() *SliceConfigurator {
	return s.Each(func(elem *BaseConfigurator) { elem.Trim() })
}

// Regexp checks if every string of the slice matches the given regular expression.
func (s *SliceConfigurator) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) *SliceConfigurator {
	return s.Each(func(elem *BaseConfigurator) { elem.Regexp(regexp, opts...) })
}

// Email checks if every string of the slice is an email address.
func (s *SliceConfigurator) Email() *SliceConfigurator {
	return s.Each(func(elem *BaseConfigurator) { elem.Email() })
}
