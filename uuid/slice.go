package uuid

import (
	"context"

	"github.com/google/uuid"

	"github.com/insei/valigo/shared"
)

// UUIDSliceFieldConfigurator configures rules for a slice of uuids.
// MaxLen, MinLen and Required are applied to the slice itself, use Each for rules on elements.
type UUIDSliceFieldConfigurator struct {
	s *shared.SliceFieldConfigurator[uuid.UUID]
}

// NewUUIDSliceFieldConfigurator creates a new UUIDSliceFieldConfigurator.
func NewUUIDSliceFieldConfigurator(s *shared.SliceFieldConfigurator[uuid.UUID]) *UUIDSliceFieldConfigurator {
	return &UUIDSliceFieldConfigurator{s: s}
}

// Each configures rules applied to every uuid of the slice.
func (s *UUIDSliceFieldConfigurator) Each(fn func(elem *BaseConfigurator)) *UUIDSliceFieldConfigurator {
	fn(NewBaseConfigurator(s.s.EachField()))
	return s
}

// AnyOf checks if every uuid of the slice is one of the allowed values.
func (s *UUIDSliceFieldConfigurator) AnyOf(allowed ...uuid.UUID) *UUIDSliceFieldConfigurator {
	return s.Each(func(elem *BaseConfigurator) { elem.AnyOf(allowed...) })
}

// MaxLen checks that the slice contains no more than maxLen elements.
func (s *UUIDSliceFieldConfigurator) MaxLen(maxLen int) *UUIDSliceFieldConfigurator {
	s.s.MaxLen(maxLen)
	return s
}

// MinLen checks that the slice contains at least minLen elements.
func (s *UUIDSliceFieldConfigurator) MinLen(minLen int) *UUIDSliceFieldConfigurator {
	s.s.MinLen(minLen)
	return s
}

// Required checks that the slice is present and not empty.
func (s *UUIDSliceFieldConfigurator) Required() *UUIDSliceFieldConfigurator {
	s.s.Required()
	return s
}

// Custom allows for custom validation logic on the whole slice.
func (s *UUIDSliceFieldConfigurator) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value shared.Elems[uuid.UUID]) []shared.Error) *UUIDSliceFieldConfigurator {
	s.s.Custom(f)
	return s
}

// When returns an UUIDSliceFieldConfigurator whose rules are applied only if whenFn returns true.
func (s *UUIDSliceFieldConfigurator) When(whenFn func(ctx context.Context, value shared.Elems[uuid.UUID]) bool) *UUIDSliceFieldConfigurator {
	return NewUUIDSliceFieldConfigurator(s.s.When(whenFn))
}
