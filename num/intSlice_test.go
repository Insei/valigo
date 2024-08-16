package num

import (
	"context"
	"fmt"
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type helperIntSliceImpl struct{}

func (h *helperIntSliceImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

type admin struct {
	RoleIDs    []int
	ChatIDs    []int
	RoleIDsPtr *[]int
	ChatIDsPtr *[]int
}

func TestIntSliceBuilderMax(t *testing.T) {
	testAdmin := admin{
		RoleIDs: []int{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RoleIDs")
	builder1 := intSliceBuilder[[]int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDs)
		},
	}
	builder1.Max(3)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := intSliceBuilder[[]int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDs)
		},
	}
	builder2.Max(7)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntSlicePtrBuilderMax(t *testing.T) {
	roleIDs := []int{1, 2, 3, 4, 5}
	testAdmin := admin{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RoleIDsPtr")
	builder1 := intSliceBuilder[*[]int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDsPtr)
		},
	}
	builder1.Max(3)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDsPtr")
	builder2 := intSliceBuilder[*[]int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDsPtr)
		},
	}
	builder2.Max(7)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntSliceBuilderMin(t *testing.T) {
	testAdmin := admin{
		RoleIDs: []int{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RoleIDs")
	builder1 := intSliceBuilder[[]int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDs)
		},
	}
	builder1.Min(2)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := intSliceBuilder[[]int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDs)
		},
	}
	builder2.Min(1)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntSlicePtrBuilderMin(t *testing.T) {
	roleIDs := []int{1, 2, 3, 4, 5}
	testAdmin := admin{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RoleIDsPtr")
	builder1 := intSliceBuilder[*[]int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDsPtr)
		},
	}
	builder1.Min(2)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDsPtr")
	builder2 := intSliceBuilder[*[]int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDsPtr)
		},
	}
	builder2.Min(1)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntSliceBuilderRequired(t *testing.T) {
	testAdmin := admin{
		RoleIDs: []int{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("ChatIDs")
	builder1 := intSliceBuilder[[]int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.ChatIDs)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := intSliceBuilder[[]int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDs)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntSlicePtrBuilderRequired(t *testing.T) {
	roleIDs := []int{1, 2, 3, 4, 5}
	testAdmin := admin{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("ChatIDsPtr")
	builder1 := intSliceBuilder[*[]int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.ChatIDsPtr)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := intSliceBuilder[*[]int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDsPtr)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntSliceBuilderCustom(t *testing.T) {
	testAdmin := admin{
		RoleIDs: []int{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("ChatIDs")
	builder1 := intSliceBuilder[[]int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.ChatIDs)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]int) []shared.Error {
		if len(*value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDs")
	builder2 := intSliceBuilder[[]int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDs)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]int) []shared.Error {
		if len(*value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntSlicePtrBuilderCustom(t *testing.T) {
	roleIDs := []int{1, 2, 3, 4, 5}
	testAdmin := admin{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("ChatIDsPtr")
	builder1 := intSliceBuilder[*[]int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.ChatIDsPtr)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]int) []shared.Error {
		if value == nil || *value == nil || len(**value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RoleIDsPtr")
	builder2 := intSliceBuilder[*[]int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RoleIDsPtr)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]int) []shared.Error {
		if value == nil || *value == nil || len(**value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}
