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

// Of converts a field of uuid.UUID or *uuid.UUID to the uuid configurator: b.Field(&m.ID).As(uuid.Of).
// A nil *uuid.UUID is absent.
func Of[U uuid.UUID | *uuid.UUID](f *shared.Field[U]) *BaseConfigurator {
	if p, ok := any(f).(*shared.Field[*uuid.UUID]); ok {
		return NewBaseConfigurator(shared.Deref(p))
	}
	return NewBaseConfigurator(shared.Retype[U, uuid.UUID](f))
}

// Required checks if the uuid is present and not uuid.Nil.
func (i *BaseConfigurator) Required() *BaseConfigurator {
	i.f.RequiredMsg(func(v uuid.UUID) bool { return v != uuid.Nil }, shared.M("uuid.required", requiredLocaleKey))
	return i
}

// AnyOf checks if the uuid value is one of the allowed values.
func (i *BaseConfigurator) AnyOf(allowed ...uuid.UUID) *BaseConfigurator {
	i.f.CheckMsg(func(v uuid.UUID) bool { return slices.Contains(allowed, v) }, shared.M("uuid.any_of", anyOfLocaleKey, shared.P("values", allowed)))
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
