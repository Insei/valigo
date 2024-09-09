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

	testCases := []struct {
		name          string
		fieldName     string
		max           int
		value         any
		expectedError int
	}{
		{
			name:          "Age maxT check",
			fieldName:     "Age",
			max:           35,
			value:         &testUser.Age,
			expectedError: 1,
		},
		{
			name:          "Height maxT check",
			fieldName:     "Height",
			max:           190,
			value:         &testUser.Height,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Max(tc.max)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
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

	testCases := []struct {
		name          string
		fieldName     string
		max           int
		value         any
		expectedError int
	}{
		{
			name:          "AgePtr maxT check",
			fieldName:     "AgePtr",
			max:           35,
			value:         &testUser.AgePtr,
			expectedError: 1,
		},
		{
			name:          "HeightPtr maxT check",
			fieldName:     "HeightPtr",
			max:           190,
			value:         &testUser.HeightPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[*int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Max(tc.max)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestIntBuilderMin(t *testing.T) {
	testUser := user{
		Age:    20,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		min           int
		value         any
		expectedError int
	}{
		{
			name:          "Age min check",
			fieldName:     "Age",
			min:           25,
			value:         &testUser.Age,
			expectedError: 1,
		},
		{
			name:          "Height min check",
			fieldName:     "Height",
			min:           165,
			value:         &testUser.Height,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Min(tc.min)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
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

	testCases := []struct {
		name          string
		fieldName     string
		min           int
		value         any
		expectedError int
	}{
		{
			name:          "AgePtr min check",
			fieldName:     "AgePtr",
			min:           25,
			value:         &testUser.AgePtr,
			expectedError: 1,
		},
		{
			name:          "HeightPtr min check",
			fieldName:     "HeightPtr",
			min:           165,
			value:         &testUser.HeightPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[*int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Min(tc.min)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestIntBuilderRequired(t *testing.T) {
	testUser := user{
		Age: 20,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "Height required check",
			fieldName:     "Height",
			value:         &testUser.Height,
			expectedError: 1,
		},
		{
			name:          "Age required check",
			fieldName:     "Age",
			value:         &testUser.Age,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Required()
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestIntPtrBuilderRequired(t *testing.T) {
	age := 20
	testUser := user{
		AgePtr: &age,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "HeightPtr required check",
			fieldName:     "HeightPtr",
			value:         &testUser.HeightPtr,
			expectedError: 1,
		},
		{
			name:          "AgePtr required check",
			fieldName:     "AgePtr",
			value:         &testUser.AgePtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[*int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Required()
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestIntBuilderAnyOf(t *testing.T) {
	testUser := user{
		Age:    18,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		allowed       []int
		value         any
		expectedError int
	}{
		{
			name:          "Age any of check",
			fieldName:     "Age",
			allowed:       []int{20, 30, 40},
			value:         &testUser.Age,
			expectedError: 1,
		},
		{
			name:          "Height any of check",
			fieldName:     "Height",
			allowed:       []int{180, 185, 190},
			value:         &testUser.Height,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.AnyOf(tc.allowed...)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
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

	testCases := []struct {
		name          string
		fieldName     string
		allowed       []int
		value         any
		expectedError int
	}{
		{
			name:          "AgePtr any of check",
			fieldName:     "AgePtr",
			allowed:       []int{20, 30, 40},
			value:         &testUser.AgePtr,
			expectedError: 1,
		},
		{
			name:          "HeightPtr any of check",
			fieldName:     "HeightPtr",
			allowed:       []int{180, 185, 190},
			value:         &testUser.HeightPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[*int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.AnyOf(tc.allowed...)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestIntBuilderAnyOfInterval(t *testing.T) {
	testUser := user{
		Age:    18,
		Height: 185,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		begin         int
		end           int
		value         any
		expectedError int
	}{
		{
			name:          "Age any of interval check",
			fieldName:     "Age",
			begin:         20,
			end:           50,
			value:         &testUser.Age,
			expectedError: 1,
		},
		{
			name:          "Height any of interval  check",
			fieldName:     "Height",
			begin:         160,
			end:           190,
			value:         &testUser.Height,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.AnyOfInterval(tc.begin, tc.end)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
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

	testCases := []struct {
		name          string
		fieldName     string
		begin         int
		end           int
		value         any
		expectedError int
	}{
		{
			name:          "AgePtr any of interval check",
			fieldName:     "AgePtr",
			begin:         20,
			end:           50,
			value:         &testUser.AgePtr,
			expectedError: 1,
		},
		{
			name:          "HeightPtr any of interval  check",
			fieldName:     "HeightPtr",
			begin:         160,
			end:           190,
			value:         &testUser.HeightPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[*int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.AnyOfInterval(tc.begin, tc.end)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestIntBuilderCustom(t *testing.T) {
	testUser := user{
		Age: 18,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "Height custom function check",
			fieldName:     "Height",
			value:         &testUser.Height,
			expectedError: 1,
		},
		{
			name:          "Age custom function check",
			fieldName:     "Age",
			value:         &testUser.Age,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *int) []shared.Error {
				if value == nil || *value == 0 {
					return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
				}
				return nil
			})
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestIntPtrBuilderCustom(t *testing.T) {
	age := 18
	testUser := user{
		AgePtr: &age,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperIntImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "HeightPtr custom function check",
			fieldName:     "HeightPtr",
			value:         &testUser.HeightPtr,
			expectedError: 1,
		},
		{
			name:          "AgePtr custom function check",
			fieldName:     "AgePtr",
			value:         &testUser.AgePtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := numBuilder[*int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **int) []shared.Error {
				if value == nil || *value == nil || **value == 0 {
					return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
				}
				return nil
			})
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}
