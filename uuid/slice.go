package uuid

import (
	"context"
	"github.com/google/uuid"
	"github.com/insei/valigo/shared"
	"slices"
)

type UUIDSliceFieldConfigurator struct {
	*shared.SliceFieldConfigurator
}

func NewUUIDSliceFieldConfigurator(p shared.SliceFieldConfiguratorParams) *UUIDSliceFieldConfigurator {
	return &UUIDSliceFieldConfigurator{
		shared.NewSliceFieldConfigurator(p),
	}
}

func (s *UUIDSliceFieldConfigurator) AnyOf(allowed ...uuid.UUID) *UUIDSliceFieldConfigurator {
	s.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v []*any) []shared.Error {
		values := shared.UnsafeValigoSliceCast[uuid.UUID](v)
		var errs []shared.Error
		for _, val := range values {
			if !slices.Contains(allowed, *val) {
				errs = append(errs, h.ErrorT(ctx, (*val).String(), anyOfLocaleKey, allowed))
			}
		}
		return errs
	})
	return s
}
