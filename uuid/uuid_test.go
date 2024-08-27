package uuid

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
	"testing"
)

type helperUUIDImpl struct{}

func (h *helperUUIDImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

type user struct {
	Id       uuid.UUID
	AppId    uuid.UUID
	ClientId uuid.UUID

	IdPtr       *uuid.UUID
	AppIdPtr    *uuid.UUID
	ClientIdPtr *uuid.UUID
}

func TestBuilderRequired(t *testing.T) {
	testUser := user{
		Id:       uuid.New(),
		AppId:    uuid.New(),
		ClientId: uuid.Nil,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperUUIDImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "Id required check",
			fieldName:     "Id",
			value:         &testUser.Id,
			expectedError: 0,
		},
		{
			name:          "AppId required check",
			fieldName:     "AppId",
			value:         &testUser.AppId,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidBuilder[uuid.UUID]{
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

func TestPtrBuilderRequired(t *testing.T) {
	id := uuid.New()
	testUser := user{
		IdPtr: &id,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperUUIDImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "AppIdPtr required check",
			fieldName:     "AppIdPtr",
			value:         &testUser.AppIdPtr,
			expectedError: 1,
		},
		{
			name:          "IdPtr required check",
			fieldName:     "IdPtr",
			value:         &testUser.IdPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidBuilder[*uuid.UUID]{
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

func TestBuilderAnyOf(t *testing.T) {
	allowedUUID := uuid.New()

	testUser := user{
		Id:    uuid.New(),
		AppId: allowedUUID,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperUUIDImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		allowedValues []uuid.UUID
		value         any
		expectedError int
	}{
		{
			name:          "Id any of check",
			fieldName:     "Id",
			allowedValues: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			value:         &testUser.Id,
			expectedError: 1,
		},
		{
			name:          "AppId any of check",
			fieldName:     "AppId",
			allowedValues: []uuid.UUID{uuid.New(), allowedUUID, uuid.New()},
			value:         &testUser.AppId,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidBuilder[uuid.UUID]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.AnyOf(tc.allowedValues...)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringPtrBuilderAnyOf(t *testing.T) {
	allowedUUID := uuid.New()
	notAllowedUUID := uuid.New()
	testUser := user{
		IdPtr:    &notAllowedUUID,
		AppIdPtr: &allowedUUID,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperUUIDImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		allowedValues []uuid.UUID
		value         any
		expectedError int
	}{
		{
			name:          "IdPtr any of check",
			fieldName:     "IdPtr",
			allowedValues: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			value:         &testUser.IdPtr,
			expectedError: 1,
		},
		{
			name:          "AppIdPtr any of check",
			fieldName:     "AppIdPtr",
			allowedValues: []uuid.UUID{uuid.New(), allowedUUID, uuid.New()},
			value:         &testUser.AppIdPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidBuilder[*uuid.UUID]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.AnyOf(tc.allowedValues...)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestBuilderCustom(t *testing.T) {
	testUser := user{
		Id:    uuid.Nil,
		AppId: uuid.New(),
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperUUIDImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "Id custom function check",
			fieldName:     "Id",
			value:         &testUser.Id,
			expectedError: 1,
		},
		{
			name:          "AppId custom function check",
			fieldName:     "AppId",
			value:         &testUser.AppId,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidBuilder[uuid.UUID]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *uuid.UUID) []shared.Error {
				if *value == uuid.Nil {
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

func TestPtrBuilderCustom(t *testing.T) {
	id := uuid.New()
	testUser := user{
		IdPtr: &id,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperUUIDImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "IdPtr custom function check",
			fieldName:     "IdPtr",
			value:         &testUser.IdPtr,
			expectedError: 0,
		},
		{
			name:          "AppIdPtr custom function check",
			fieldName:     "AppIdPtr",
			value:         &testUser.AppIdPtr,
			expectedError: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := uuidBuilder[*uuid.UUID]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **uuid.UUID) []shared.Error {
				if value == nil || *value == nil || **value == uuid.Nil {
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
