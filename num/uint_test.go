package num

import (
	"context"
	"fmt"
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type helperUuintImpl struct{}

func (h *helperUuintImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

type man struct {
	Age       uint8
	Height    uint
	AgePtr    *uint8
	HeightPtr *uint
}

func TestUintBuilderMax(t *testing.T) {
	testMan := man{
		Age:    40,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("Age")
	builder1 := uintBuilder[uint8]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Age)
		},
	}
	builder1.Max(35)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Height")
	builder2 := uintBuilder[uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Height)
		},
	}
	builder2.Max(190)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintPtrBuilderMax(t *testing.T) {
	age := uint8(40)
	height := uint(185)
	testMan := man{
		AgePtr:    &age,
		HeightPtr: &height,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("AgePtr")
	builder1 := uintBuilder[*uint8]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.AgePtr)
		},
	}
	builder1.Max(35)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("HeightPtr")
	builder2 := uintBuilder[*uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.HeightPtr)
		},
	}
	builder2.Max(190)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintBuilderMin(t *testing.T) {
	testMan := man{
		Age:    20,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testMan)

	helper := helperUuintImpl{}
	var errs []shared.Error
	field1 := storage.MustFind("Age")

	builder1 := uintBuilder[uint8]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Age)
		},
	}
	builder1.Min(25)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Height")
	builder2 := uintBuilder[uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Height)
		},
	}
	builder2.Min(165)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintPtrBuilderMin(t *testing.T) {
	age := uint8(20)
	height := uint(185)
	testMan := man{
		AgePtr:    &age,
		HeightPtr: &height,
	}
	storage, _ := fmap.GetFrom(testMan)

	helper := helperUuintImpl{}
	var errs []shared.Error
	field1 := storage.MustFind("AgePtr")

	builder1 := uintBuilder[*uint8]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.AgePtr)
		},
	}
	builder1.Min(25)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("HeightPtr")
	builder2 := uintBuilder[*uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.HeightPtr)
		},
	}
	builder2.Min(165)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintBuilderRequired(t *testing.T) {
	testMan := man{
		Age: 20,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("Height")
	builder1 := uintBuilder[*uint]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Height)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Age")
	builder2 := uintBuilder[uint8]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Age)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintPtrBuilderRequired(t *testing.T) {
	age := uint8(20)
	testMan := man{
		AgePtr: &age,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("HeightPtr")
	builder1 := uintBuilder[*uint]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.HeightPtr)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("AgePtr")
	builder2 := uintBuilder[*uint8]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.AgePtr)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintBuilderAnyOf(t *testing.T) {
	testMan := man{
		Age:    18,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error
	allowedAges := []uint{20, 30, 40}
	allowedHeights := []uint{180, 185, 190}

	field1 := storage.MustFind("Age")
	builder1 := uintBuilder[uint8]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Age)
		},
	}
	builder1.AnyOf(allowedAges...)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Height")
	builder2 := uintBuilder[uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Height)
		},
	}
	builder2.AnyOf(allowedHeights...)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintPtrBuilderAnyOf(t *testing.T) {
	age := uint8(18)
	height := uint(185)
	testMan := man{
		AgePtr:    &age,
		HeightPtr: &height,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error
	allowedAges := []uint{20, 30, 40}
	allowedHeights := []uint{180, 185, 190}

	field1 := storage.MustFind("AgePtr")
	builder1 := uintBuilder[*uint8]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.AgePtr)
		},
	}
	builder1.AnyOf(allowedAges...)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("HeightPtr")
	builder2 := uintBuilder[*uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.HeightPtr)
		},
	}
	builder2.AnyOf(allowedHeights...)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintBuilderAnyOfInterval(t *testing.T) {
	testMan := man{
		Age:    18,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error
	beginAgeInterval := uint(20)
	endAgeInterval := uint(50)
	beginHeightInterval := uint(160)
	endHeightInterval := uint(190)

	field1 := storage.MustFind("Age")
	builder1 := uintBuilder[uint8]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Age)
		},
	}
	builder1.AnyOfInterval(beginAgeInterval, endAgeInterval)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Height")
	builder2 := uintBuilder[uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Height)
		},
	}
	builder2.AnyOfInterval(beginHeightInterval, endHeightInterval)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintPtrBuilderAnyOfInterval(t *testing.T) {
	age := uint8(18)
	height := uint(185)
	testMan := man{
		AgePtr:    &age,
		HeightPtr: &height,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error
	beginAgeInterval := uint(20)
	endAgeInterval := uint(50)
	beginHeightInterval := uint(160)
	endHeightInterval := uint(190)

	field1 := storage.MustFind("AgePtr")
	builder1 := uintBuilder[*uint8]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.AgePtr)
		},
	}
	builder1.AnyOfInterval(beginAgeInterval, endAgeInterval)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("HeightPtr")
	builder2 := uintBuilder[*uint]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.HeightPtr)
		},
	}
	builder2.AnyOfInterval(beginHeightInterval, endHeightInterval)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintBuilderCustom(t *testing.T) {
	testMan := man{
		Age: 18,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("Height")
	builder1 := uintBuilder[uint]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Height)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *uint) []shared.Error {
		if value == nil || *value == 0 {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Age")
	builder2 := uintBuilder[uint8]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.Age)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *uint8) []shared.Error {
		if value == nil || *value == 0 {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestUintPtrBuilderCustom(t *testing.T) {
	age := uint8(18)
	testMan := man{
		AgePtr: &age,
	}
	storage, _ := fmap.GetFrom(testMan)
	helper := helperUuintImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("HeightPtr")
	builder1 := uintBuilder[*uint]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.HeightPtr)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **uint) []shared.Error {
		if value == nil || *value == nil || **value == 0 {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("AgePtr")
	builder2 := uintBuilder[*uint8]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testMan.AgePtr)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **uint8) []shared.Error {
		if value == nil || *value == nil || **value == 0 {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}
