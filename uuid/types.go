package uuid

import (
	"context"

	"github.com/google/uuid"

	"github.com/insei/valigo/shared"
)

// UUIDBuilder is a builder interface for uuid fields.
type UUIDBuilder[T uuid.UUID | *uuid.UUID] interface {
	// Required checks if the uuid.UUID value is not empty.
	Required() UUIDBuilder[T]

	// AnyOf checks if the uuid.UUID value is one of the allowed values.
	AnyOf(allowed ...uuid.UUID) UUIDBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) UUIDBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(whenFn func(ctx context.Context, value *T) bool) UUIDBuilder[T]
}

// UUIDSliceBuilder is a builder interface for uuid slice fields.
type UUIDSliceBuilder[T []uuid.UUID | *[]uuid.UUID] interface {
	// Required checks if the slice is not empty.
	Required() UUIDSliceBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) UUIDSliceBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(f func(ctx context.Context, value *T) bool) UUIDSliceBuilder[T]
}

// UuidBundleBuilder is a builder interface for a bundle of uuid fields.
type UuidBundleBuilder interface {
	// UUID adds an uuid field to the bundle.
	UUID(filed *uuid.UUID) UUIDBuilder[uuid.UUID]

	// UUIDPtr adds a pointer to an uuid field to the bundle.
	UUIDPtr(filed **uuid.UUID) UUIDBuilder[*uuid.UUID]

	// UUIDSlice adds an uuid slice field to the bundle.
	UUIDSlice(filed *[]uuid.UUID) UUIDSliceBuilder[[]uuid.UUID]

	// UUIDSlicePtr adds a pointer to an uuid slice field to the bundle.
	UUIDSlicePtr(filed **[]uuid.UUID) UUIDSliceBuilder[*[]uuid.UUID]
}
