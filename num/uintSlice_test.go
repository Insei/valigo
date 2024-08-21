package num

import (
	"context"
	"fmt"
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type helperUintSliceImpl struct{}

func (h *helperUintSliceImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

type test struct {
	RoleIDs    []uint
	ChatIDs    []uint16
	RoleIDsPtr *[]uint
	ChatIDsPtr *[]uint16
}

func TestUintSliceBuilderMax(t *testing.T) {
	tt := test{
		RoleIDs: []uint{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(tt)
	helper := helperUintSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RoleIDs")
	builder1 := uintSliceBuilder[[]uint]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDs)
		},
	}
	builder1.Max(3)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := uintSliceBuilder[[]uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDs)
		},
	}
	builder2.Max(7)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintSlicePtrBuilderMax(t *testing.T) {
	roleIDs := []uint{1, 2, 3, 4, 5}
	tt := test{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(tt)
	helper := helperUintSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RoleIDsPtr")
	builder1 := uintSliceBuilder[*[]uint]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDsPtr)
		},
	}
	builder1.Max(3)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDsPtr")
	builder2 := uintSliceBuilder[*[]uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDsPtr)
		},
	}
	builder2.Max(7)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintSliceBuilderMin(t *testing.T) {
	tt := test{
		RoleIDs: []uint{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(tt)
	helper := helperUintSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RoleIDs")
	builder1 := uintSliceBuilder[[]uint]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDs)
		},
	}
	builder1.Min(2)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := uintSliceBuilder[[]uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDs)
		},
	}
	builder2.Min(1)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintSlicePtrBuilderMin(t *testing.T) {
	roleIDs := []uint{1, 2, 3, 4, 5}
	tt := test{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(tt)
	helper := helperUintSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RoleIDsPtr")
	builder1 := uintSliceBuilder[*[]uint]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDsPtr)
		},
	}
	builder1.Min(2)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDsPtr")
	builder2 := uintSliceBuilder[*[]uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDsPtr)
		},
	}
	builder2.Min(1)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintSliceBuilderRequired(t *testing.T) {
	tt := test{
		RoleIDs: []uint{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(tt)
	helper := helperUintSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("ChatIDs")
	builder1 := uintSliceBuilder[[]uint16]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.ChatIDs)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := uintSliceBuilder[[]uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDs)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintSlicePtrBuilderRequired(t *testing.T) {
	roleIDs := []uint{1, 2, 3, 4, 5}
	tt := test{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(tt)
	helper := helperUintSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("ChatIDsPtr")
	builder1 := uintSliceBuilder[*[]uint16]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.ChatIDsPtr)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := uintSliceBuilder[*[]uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDsPtr)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintSliceBuilderCustom(t *testing.T) {
	tt := test{
		RoleIDs: []uint{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(tt)
	helper := helperUintSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("ChatIDs")
	builder1 := uintSliceBuilder[[]uint16]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.ChatIDs)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]uint16) []shared.Error {
		if len(*value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := uintSliceBuilder[[]uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDs)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]uint) []shared.Error {
		if len(*value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintSlicePtrBuilderCustom(t *testing.T) {
	roleIDs := []uint{1, 2, 3, 4, 5}
	tt := test{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(tt)
	helper := helperUintSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("ChatIDsPtr")
	builder1 := uintSliceBuilder[*[]uint16]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.ChatIDsPtr)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]uint16) []shared.Error {
		if value == nil || *value == nil || len(**value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDsPtr")
	builder2 := uintSliceBuilder[*[]uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &tt.RoleIDsPtr)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]uint) []shared.Error {
		if value == nil || *value == nil || len(**value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}
