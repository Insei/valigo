package str

import (
	"context"
	"regexp"

	"github.com/insei/valigo/shared"
)

// StringSliceFieldConfigurator configures rules for a slice of strings.
// MaxLen, MinLen and Required are applied to the slice itself, use Each for rules on elements.
type StringSliceFieldConfigurator struct {
	s *shared.SliceFieldConfigurator[string]
}

// NewStringSliceFieldConfigurator creates a new StringSliceFieldConfigurator.
func NewStringSliceFieldConfigurator(s *shared.SliceFieldConfigurator[string]) *StringSliceFieldConfigurator {
	return &StringSliceFieldConfigurator{s: s}
}

// Each configures rules applied to every string of the slice.
func (s *StringSliceFieldConfigurator) Each(fn func(elem *BaseConfigurator)) *StringSliceFieldConfigurator {
	fn(NewBaseConfigurator(s.s.EachField()))
	return s
}

// Trim removes leading and trailing whitespace from every string of the slice.
func (s *StringSliceFieldConfigurator) Trim() *StringSliceFieldConfigurator {
	return s.Each(func(elem *BaseConfigurator) { elem.Trim() })
}

// Regexp checks if every string of the slice matches the given regular expression.
func (s *StringSliceFieldConfigurator) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) *StringSliceFieldConfigurator {
	return s.Each(func(elem *BaseConfigurator) { elem.Regexp(regexp, opts...) })
}

// Email checks if every string of the slice is an email address.
func (s *StringSliceFieldConfigurator) Email() *StringSliceFieldConfigurator {
	return s.Each(func(elem *BaseConfigurator) { elem.Email() })
}

// MaxLen checks that the slice contains no more than maxLen elements.
func (s *StringSliceFieldConfigurator) MaxLen(maxLen int) *StringSliceFieldConfigurator {
	s.s.MaxLen(maxLen)
	return s
}

// MinLen checks that the slice contains at least minLen elements.
func (s *StringSliceFieldConfigurator) MinLen(minLen int) *StringSliceFieldConfigurator {
	s.s.MinLen(minLen)
	return s
}

// Required checks that the slice is present and not empty.
func (s *StringSliceFieldConfigurator) Required() *StringSliceFieldConfigurator {
	s.s.Required()
	return s
}

// Custom allows for custom validation logic on the whole slice.
func (s *StringSliceFieldConfigurator) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value shared.Elems[string]) []shared.Error) *StringSliceFieldConfigurator {
	s.s.Custom(f)
	return s
}

// When returns a StringSliceFieldConfigurator whose rules are applied only if whenFn returns true.
func (s *StringSliceFieldConfigurator) When(whenFn func(ctx context.Context, value shared.Elems[string]) bool) *StringSliceFieldConfigurator {
	return NewStringSliceFieldConfigurator(s.s.When(whenFn))
}
