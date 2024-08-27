package uuid

import (
	"context"
	"github.com/google/uuid"
	"github.com/insei/valigo/shared"
)

type Builder[T uuid.UUID | *uuid.UUID] interface {
	Required() Builder[T]

	AnyOf(allowed ...uuid.UUID) Builder[T]

	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) Builder[T]

	When(whenFn func(ctx context.Context, value *T) bool) Builder[T]
}

type SliceBuilder[T []uuid.UUID | *[]uuid.UUID] interface {
	Required() SliceBuilder[T]

	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) SliceBuilder[T]

	When(f func(ctx context.Context, value *T) bool) SliceBuilder[T]
}
