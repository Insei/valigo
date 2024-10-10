package uuid

import (
	"context"
	"github.com/google/uuid"

	"github.com/insei/valigo/shared"
)

type BaseConfigurator interface {
	// Required checks if the uuid.UUID value is not empty.
	Required() BaseConfigurator

	// AnyOf checks if the uuid.UUID value is one of the allowed values.
	AnyOf(allowed ...uuid.UUID) BaseConfigurator

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error) BaseConfigurator

	// When allows for conditional validation based on a given condition.
	When(whenFn func(ctx context.Context, value any) bool) BaseConfigurator
}

// UUIDBundleConfigurator is a builder interface for a bundle of uuid fields.
// It provides methods for adding uuid fields to the bundle.
type UUIDBundleConfigurator interface {
	UUID(fieldPtr any) BaseConfigurator
}
