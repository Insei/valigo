package uuid

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

const (
	requiredLocaleKey = "validation:uuid:Should be fulfilled"
	anyOfLocaleKey    = "validation:uuid:Only %s values is allowed"
)

type baseConfigurator struct {
	c     *shared.FieldConfigurator[uuid.UUID]
	field fmap.Field
	h     shared.Helper
}

// Required checks if the uuid is not empty.
func (i *baseConfigurator) Required() BaseConfigurator {
	i.c.Append(func(v uuid.UUID) bool {
		return v != uuid.Nil
	}, requiredLocaleKey)

	return i
}

// AnyOf checks if the uuid value is one of the allowed values.
func (i *baseConfigurator) AnyOf(allowed ...uuid.UUID) BaseConfigurator {
	i.c.Append(func(v uuid.UUID) bool {
		return slices.Contains(allowed, v)
	}, anyOfLocaleKey)
	return i
}

// Custom allows for custom validation logic to be applied to the uuid value.
func (i *baseConfigurator) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error) BaseConfigurator {
	customHelper := shared.NewFieldCustomHelper(i.field, i.h)
	i.c.CustomAppend(func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value)
	})
	return i
}

// When allows for conditional validation logic to be applied to the uuid value.
func (i *baseConfigurator) When(whenFn func(ctx context.Context, value any) bool) BaseConfigurator {
	if whenFn == nil {
		return i
	}
	base := i.c.NewWithWhen(whenFn)
	return &baseConfigurator{
		c:     base,
		field: i.field,
		h:     i.h,
	}
}
