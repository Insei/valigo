package uuid

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
	"testing"
)

type helperUUIDSliceImpl struct{}

func (h *helperUUIDSliceImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

type entity struct {
	ClientIds    []uuid.UUID
	AppIds       []uuid.UUID
	ClientIdsPtr *[]uuid.UUID
	AppIdsPtr    *[]uuid.UUID
}

func TestSliceBuilderRequired(t *testing.T) {
	testEntity := entity{
		ClientIds: []uuid.UUID{uuid.New(), uuid.New()},
	}
	storage, _ := fmap.GetFrom(testEntity)
	helper := helperUUIDSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "AppIds slice required check",
			fieldName:     "AppIds",
			value:         &testEntity.AppIds,
			expectedError: 1,
		},
		{
			name:          "ClientIds slice required check",
			fieldName:     "ClientIds",
			value:         &testEntity.ClientIds,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidSliceBuilder[[]uuid.UUID]{
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

func TestSlicePtrBuilderRequired(t *testing.T) {
	AppIds := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	testEntity := entity{
		AppIdsPtr: &AppIds,
	}
	storage, _ := fmap.GetFrom(testEntity)
	helper := helperUUIDSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "ClientIdsPtr slice required check",
			fieldName:     "ClientIdsPtr",
			value:         &testEntity.ClientIdsPtr,
			expectedError: 1,
		},
		{
			name:          "AppIdsPtr slice required check",
			fieldName:     "AppIdsPtr",
			value:         &testEntity.AppIdsPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidSliceBuilder[*[]uuid.UUID]{
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

func TestSliceBuilderCustom(t *testing.T) {
	testEntity := entity{
		AppIds: []uuid.UUID{uuid.New(), uuid.New()},
	}
	storage, _ := fmap.GetFrom(testEntity)
	helper := helperUUIDSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "ClientIds slice custom function check",
			fieldName:     "ClientIds",
			value:         &testEntity.ClientIds,
			expectedError: 1,
		},
		{
			name:          "AppIds slice custom function check",
			fieldName:     "AppIds",
			value:         &testEntity.AppIds,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidSliceBuilder[[]uuid.UUID]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]uuid.UUID) []shared.Error {
				if len(*value) == 0 {
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

func TestStringSlicePtrBuilderCustom(t *testing.T) {
	appIds := []uuid.UUID{uuid.New(), uuid.New()}
	testEntity := entity{
		AppIdsPtr: &appIds,
	}
	storage, _ := fmap.GetFrom(testEntity)
	helper := helperUUIDSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "ClientIdsPtr slice custom function check",
			fieldName:     "ClientIdsPtr",
			value:         &testEntity.ClientIdsPtr,
			expectedError: 1,
		},
		{
			name:          "AppIdsPtr slice custom function check",
			fieldName:     "AppIdsPtr",
			value:         &testEntity.AppIdsPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidSliceBuilder[*[]uuid.UUID]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]uuid.UUID) []shared.Error {
				if *value == nil || len(**value) == 0 {
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
