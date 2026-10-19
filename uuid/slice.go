package uuid

import (
	"github.com/google/uuid"

	"github.com/insei/valigo/shared"
)

// SliceConfigurator configures rules for a slice of uuids.
// Required, MinLen, MaxLen, Custom and When are applied to the slice itself, use Each for rules on elements.
type SliceConfigurator struct {
	shared.SliceRules[*SliceConfigurator, uuid.UUID]
}

// Slice converts a slice configurator of uuid.UUID or *uuid.UUID to the uuid slice configurator:
// b.Slice(&m.IDs).As(uuid.Slice). Nil elements of []*uuid.UUID are absent.
func Slice[E uuid.UUID | *uuid.UUID](s *shared.SliceFieldConfigurator[E]) *SliceConfigurator {
	return newSliceConfigurator(shared.RetypeSlice[E, uuid.UUID](s))
}

func newSliceConfigurator(s *shared.SliceFieldConfigurator[uuid.UUID]) *SliceConfigurator {
	c := &SliceConfigurator{}
	c.SliceRules = shared.NewSliceRules(s, c, newSliceConfigurator)
	return c
}

// Each configures rules applied to every uuid of the slice.
func (s *SliceConfigurator) Each(fn func(elem *BaseConfigurator)) *SliceConfigurator {
	fn(NewBaseConfigurator(s.Slice().EachField()))
	return s
}

// AnyOf checks if every uuid of the slice is one of the allowed values.
func (s *SliceConfigurator) AnyOf(allowed ...uuid.UUID) *SliceConfigurator {
	return s.Each(func(elem *BaseConfigurator) { elem.AnyOf(allowed...) })
}
