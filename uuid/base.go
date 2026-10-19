package uuid

import (
	"context"
	"slices"

	"github.com/google/uuid"

	"github.com/insei/valigo/shared"
)

const (
	requiredLocaleKey = "validation:uuid:Should be fulfilled"
	anyOfLocaleKey    = "validation:uuid:Only %s values is allowed"
)

// BaseConfigurator configures rules for an uuid value.
type BaseConfigurator struct {
	f *shared.Field[uuid.UUID]
}

// NewBaseConfigurator creates a new BaseConfigurator.
func NewBaseConfigurator(f *shared.Field[uuid.UUID]) *BaseConfigurator {
	return &BaseConfigurator{f: f}
}

// Required checks if the uuid is present and not uuid.Nil.
func (i *BaseConfigurator) Required() *BaseConfigurator {
	i.f.Required(func(v uuid.UUID) bool { return v != uuid.Nil }, requiredLocaleKey)
	return i
}

// AnyOf checks if the uuid value is one of the allowed values.
func (i *BaseConfigurator) AnyOf(allowed ...uuid.UUID) *BaseConfigurator {
	i.f.Check(func(v uuid.UUID) bool { return slices.Contains(allowed, v) }, anyOfLocaleKey, allowed)
	return i
}

// Custom allows for custom validation logic to be applied to the uuid value.
func (i *BaseConfigurator) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value uuid.UUID) []shared.Error) *BaseConfigurator {
	i.f.Custom(f)
	return i
}

// When allows for conditional validation logic to be applied to the uuid value.
func (i *BaseConfigurator) When(whenFn func(ctx context.Context, value uuid.UUID) bool) *BaseConfigurator {
	return NewBaseConfigurator(i.f.When(whenFn))
}
