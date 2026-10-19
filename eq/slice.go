package eq

import (
	"context"

	"github.com/insei/valigo/shared"
)

const uniqueLocaleKey = "validation:value:Should contain unique values"

// smallSlice is the length up to which Unique compares the elements pairwise without allocations.
const smallSlice = 32

// SliceConfigurator configures rules for a slice of comparable values.
// Required, MinLen, MaxLen, Custom and When are applied to the slice itself, use Each for rules on elements.
type SliceConfigurator[E comparable] struct {
	shared.SliceRules[*SliceConfigurator[E], E]
}

// Slice converts a slice configurator of comparable values: b.Slice(&m.Roles).As(eq.Slice).
func Slice[E comparable](s *shared.SliceFieldConfigurator[E]) *SliceConfigurator[E] {
	c := &SliceConfigurator[E]{}
	c.SliceRules = shared.NewSliceRules(s, c, Slice[E])
	return c
}

// Each configures rules applied to every element of the slice.
func (s *SliceConfigurator[E]) Each(fn func(elem *BaseConfigurator[E])) *SliceConfigurator[E] {
	fn(NewBaseConfigurator(s.Slice().EachField(), false))
	return s
}

// OneOf checks if every element of the slice is one of the allowed values.
func (s *SliceConfigurator[E]) OneOf(allowed ...E) *SliceConfigurator[E] {
	return s.Each(func(elem *BaseConfigurator[E]) { elem.OneOf(allowed...) })
}

// Unique checks that the slice has no duplicate elements, the error value is the first duplicate.
// Slices up to 32 elements are checked without allocations.
func (s *SliceConfigurator[E]) Unique() *SliceConfigurator[E] {
	msg := shared.M("value.unique", uniqueLocaleKey)
	s.Slice().Custom(func(ctx context.Context, h *shared.FieldCustomHelper, elems shared.Elems[E]) []shared.Error {
		if dup, ok := duplicate(elems); ok {
			return []shared.Error{h.ErrorMsg(ctx, dup, msg)}
		}
		return nil
	})
	return s
}

// duplicate returns the first element that occurs earlier in the slice.
func duplicate[E comparable](elems shared.Elems[E]) (E, bool) {
	n := elems.Len()
	if n <= smallSlice {
		for i := 1; i < n; i++ {
			for j := 0; j < i; j++ {
				if *elems.At(i) == *elems.At(j) {
					return *elems.At(i), true
				}
			}
		}
		var zero E
		return zero, false
	}
	seen := make(map[E]struct{}, n)
	for _, v := range elems.All() {
		if _, ok := seen[*v]; ok {
			return *v, true
		}
		seen[*v] = struct{}{}
	}
	var zero E
	return zero, false
}
