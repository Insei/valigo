package num

import (
	"context"
	"fmt"
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type helperIntImpl struct{}

func (h *helperIntImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

type user struct {
	Age       int
	Height    int
	AgePtr    *int
	HeightPtr *int
}

func TestIntBuilderMax(t *testing.T) {
	testUser := user{
		Age:    40,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("Age")
	builder1 := intBuilder[int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Age)
		},
	}
	builder1.Max(35)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Height")
	builder2 := intBuilder[int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Height)
		},
	}
	builder2.Max(190)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntPtrBuilderMax(t *testing.T) {
	age := 40
	height := 185
	testUser := user{
		AgePtr:    &age,
		HeightPtr: &height,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("AgePtr")
	builder1 := intBuilder[*int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.AgePtr)
		},
	}
	builder1.Max(35)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("HeightPtr")
	builder2 := intBuilder[*int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.HeightPtr)
		},
	}
	builder2.Max(190)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntBuilderMin(t *testing.T) {
	testUser := user{
		Age:    20,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testUser)

	helper := helperIntImpl{}
	var errs []shared.Error
	field1 := storage.MustFind("Age")

	builder1 := intBuilder[int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Age)
		},
	}
	builder1.Min(25)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Height")
	builder2 := intBuilder[int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Height)
		},
	}
	builder2.Min(165)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntPtrBuilderMin(t *testing.T) {
	age := 20
	height := 185
	testUser := user{
		AgePtr:    &age,
		HeightPtr: &height,
	}
	storage, _ := fmap.GetFrom(testUser)

	helper := helperIntImpl{}
	var errs []shared.Error
	field1 := storage.MustFind("AgePtr")

	builder1 := intBuilder[*int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.AgePtr)
		},
	}
	builder1.Min(25)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("HeightPtr")
	builder2 := intBuilder[*int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.HeightPtr)
		},
	}
	builder2.Min(165)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntBuilderRequired(t *testing.T) {
	testUser := user{
		Age: 20,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("Height")
	builder1 := intBuilder[*int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Height)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Age")
	builder2 := intBuilder[int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Age)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntPtrBuilderRequired(t *testing.T) {
	age := 20
	testUser := user{
		AgePtr: &age,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("HeightPtr")
	builder1 := intBuilder[*int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.HeightPtr)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("AgePtr")
	builder2 := intBuilder[*int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.AgePtr)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntBuilderAnyOf(t *testing.T) {
	testUser := user{
		Age:    18,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error
	allowedAges := []int{20, 30, 40}
	allowedHeights := []int{180, 185, 190}

	field1 := storage.MustFind("Age")
	builder1 := intBuilder[int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Age)
		},
	}
	builder1.AnyOf(allowedAges...)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Height")
	builder2 := intBuilder[int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Height)
		},
	}
	builder2.AnyOf(allowedHeights...)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntPtrBuilderAnyOf(t *testing.T) {
	age := 18
	height := 185
	testUser := user{
		AgePtr:    &age,
		HeightPtr: &height,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error
	allowedAges := []int{20, 30, 40}
	allowedHeights := []int{180, 185, 190}

	field1 := storage.MustFind("AgePtr")
	builder1 := intBuilder[*int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.AgePtr)
		},
	}
	builder1.AnyOf(allowedAges...)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("HeightPtr")
	builder2 := intBuilder[*int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.HeightPtr)
		},
	}
	builder2.AnyOf(allowedHeights...)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntBuilderAnyOfInterval(t *testing.T) {
	testUser := user{
		Age:    18,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error
	beginAgeInterval := 20
	endAgeInterval := 50
	beginHeightInterval := 160
	endHeightInterval := 190

	field1 := storage.MustFind("Age")
	builder1 := intBuilder[int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Age)
		},
	}
	builder1.AnyOfInterval(beginAgeInterval, endAgeInterval)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Height")
	builder2 := intBuilder[int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Height)
		},
	}
	builder2.AnyOfInterval(beginHeightInterval, endHeightInterval)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntPtrBuilderAnyOfInterval(t *testing.T) {
	age := 18
	height := 185
	testUser := user{
		AgePtr:    &age,
		HeightPtr: &height,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error
	beginAgeInterval := 20
	endAgeInterval := 50
	beginHeightInterval := 160
	endHeightInterval := 190

	field1 := storage.MustFind("AgePtr")
	builder1 := intBuilder[*int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.AgePtr)
		},
	}
	builder1.AnyOfInterval(beginAgeInterval, endAgeInterval)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("HeightPtr")
	builder2 := intBuilder[*int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.HeightPtr)
		},
	}
	builder2.AnyOfInterval(beginHeightInterval, endHeightInterval)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntBuilderCustom(t *testing.T) {
	testUser := user{
		Age: 18,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("Height")
	builder1 := intBuilder[int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Height)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *int) []shared.Error {
		if value == nil || *value == 0 {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Age")
	builder2 := intBuilder[int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Age)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *int) []shared.Error {
		if value == nil || *value == 0 {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestIntPtrBuilderCustom(t *testing.T) {
	age := 18
	testUser := user{
		AgePtr: &age,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("HeightPtr")
	builder1 := intBuilder[*int]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.HeightPtr)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **int) []shared.Error {
		if value == nil || *value == nil || **value == 0 {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("AgePtr")
	builder2 := intBuilder[*int]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.AgePtr)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **int) []shared.Error {
		if value == nil || *value == nil || **value == 0 {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}
