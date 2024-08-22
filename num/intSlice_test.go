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

	testCases := []struct {
		name          string
		fieldName     string
		max           int
		value         any
		expectedError int
	}{
		{
			name:          "RoleIDs slice max check",
			fieldName:     "RoleIDs",
			max:           3,
			value:         &testAdmin.RoleIDs,
			expectedError: 1,
		},
		{
			name:          "RoleIDs slice max check",
			fieldName:     "RoleIDs",
			max:           7,
			value:         &testAdmin.RoleIDs,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := intSliceBuilder[[]int]{
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

func TestIntSlicePtrBuilderMax(t *testing.T) {
	roleIDs := []int{1, 2, 3, 4, 5}
	testAdmin := admin{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		max           int
		value         any
		expectedError int
	}{
		{
			name:          "RoleIDsPtr slice max check",
			fieldName:     "RoleIDsPtr",
			max:           3,
			value:         &testAdmin.RoleIDsPtr,
			expectedError: 1,
		},
		{
			name:          "RoleIDsPtr slice max check",
			fieldName:     "RoleIDsPtr",
			max:           7,
			value:         &testAdmin.RoleIDsPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := intSliceBuilder[*[]int]{
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

func TestIntSliceBuilderMin(t *testing.T) {
	testAdmin := admin{
		RoleIDs: []int{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		min           int
		value         any
		expectedError int
	}{
		{
			name:          "RoleIDs slice min check",
			fieldName:     "RoleIDs",
			min:           2,
			value:         &testAdmin.RoleIDs,
			expectedError: 1,
		},
		{
			name:          "RoleIDs slice min check",
			fieldName:     "RoleIDs",
			min:           1,
			value:         &testAdmin.RoleIDs,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := intSliceBuilder[[]int]{
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

func TestIntSlicePtrBuilderMin(t *testing.T) {
	roleIDs := []int{1, 2, 3, 4, 5}
	testAdmin := admin{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		min           int
		value         any
		expectedError int
	}{
		{
			name:          "RoleIDsPtr slice min check",
			fieldName:     "RoleIDsPtr",
			min:           2,
			value:         &testAdmin.RoleIDsPtr,
			expectedError: 1,
		},
		{
			name:          "RoleIDsPtr slice min check",
			fieldName:     "RoleIDsPtr",
			min:           1,
			value:         &testAdmin.RoleIDsPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := intSliceBuilder[*[]int]{
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

func TestIntSliceBuilderRequired(t *testing.T) {
	testAdmin := admin{
		RoleIDs: []int{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "ChatIDs slice required check",
			fieldName:     "ChatIDs",
			value:         &testAdmin.ChatIDs,
			expectedError: 1,
		},
		{
			name:          "RoleIDs slice required check",
			fieldName:     "RoleIDs",
			value:         &testAdmin.RoleIDs,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := intSliceBuilder[[]int]{
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

func TestIntSlicePtrBuilderRequired(t *testing.T) {
	roleIDs := []int{1, 2, 3, 4, 5}
	testAdmin := admin{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "ChatIDsPtr slice required check",
			fieldName:     "ChatIDsPtr",
			value:         &testAdmin.ChatIDsPtr,
			expectedError: 1,
		},
		{
			name:          "RoleIDsPtr slice required check",
			fieldName:     "RoleIDsPtr",
			value:         &testAdmin.RoleIDsPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := intSliceBuilder[*[]int]{
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

func TestIntSliceBuilderCustom(t *testing.T) {
	testAdmin := admin{
		RoleIDs: []int{1, 2, 3, 4, 5},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "ChatIDs slice custom function check",
			fieldName:     "ChatIDs",
			value:         &testAdmin.ChatIDs,
			expectedError: 1,
		},
		{
			name:          "RoleIDs slice custom function check",
			fieldName:     "RoleIDs",
			value:         &testAdmin.RoleIDs,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := intSliceBuilder[[]int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]int) []shared.Error {
				if value == nil || *value == nil || len(*value) == 0 {
					return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
				}
				return nil
			})
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestIntSlicePtrBuilderCustom(t *testing.T) {
	roleIDs := []int{1, 2, 3, 4, 5}
	testAdmin := admin{
		RoleIDsPtr: &roleIDs,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperIntSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "ChatIDsPtr slice custom function check",
			fieldName:     "ChatIDsPtr",
			value:         &testAdmin.ChatIDsPtr,
			expectedError: 1,
		},
		{
			name:          "RoleIDsPtr slice custom function check",
			fieldName:     "RoleIDsPtr",
			value:         &testAdmin.RoleIDsPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := intSliceBuilder[*[]int]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]int) []shared.Error {
				if value == nil || *value == nil || len(**value) == 0 {
					return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
				}
				return nil
			})
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}
