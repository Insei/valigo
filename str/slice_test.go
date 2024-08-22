package str

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type helperStringSliceImpl struct{}

func (h *helperStringSliceImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

type admin struct {
	Roles           []string
	PhoneNumbers    []string
	RolesPtr        *[]string
	PhoneNumbersPtr *[]string
}

func TestStringSliceBuilderTrim(t *testing.T) {
	testAdmin := admin{
		Roles: []string{"   root "},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		trimmedValue  string
		value         any
		expectedError bool
	}{
		{
			name:          "Roles slice trim check",
			fieldName:     "Roles",
			trimmedValue:  " root",
			value:         &testAdmin.Roles,
			expectedError: true,
		},
		{
			name:          "Roles slice trim check",
			fieldName:     "Roles",
			trimmedValue:  "root",
			value:         &testAdmin.Roles,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Trim()
			trimmedValue := field.Get(&testAdmin)
			if trimmedValue != nil {
				trimmedSliceValue := trimmedValue.([]string)
				if trimmedSliceValue[0] != tc.trimmedValue && !tc.expectedError {
					t.Errorf("expected 'test value', got '%v'", trimmedValue)
				}
			} else {
				t.Errorf("trimmedValue is nil")
			}
		})
	}
}

func TestStringSlicePtrBuilderTrim(t *testing.T) {
	roles := []string{"   root "}
	testAdmin := admin{
		RolesPtr: &roles,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringSliceImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		trimmedValue  string
		value         any
		expectedError bool
	}{
		{
			name:          "RolesPtr slice trim check",
			fieldName:     "RolesPtr",
			trimmedValue:  " root",
			value:         &testAdmin.RolesPtr,
			expectedError: true,
		},
		{
			name:          "RolesPtr slice trim check",
			fieldName:     "RolesPtr",
			trimmedValue:  "root",
			value:         &testAdmin.RolesPtr,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[*[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Trim()
			trimmedValue := field.Get(&testAdmin)
			if trimmedValue != nil {
				trimmedSliceValue := trimmedValue.(*[]string)
				if (*trimmedSliceValue)[0] != tc.trimmedValue && !tc.expectedError {
					t.Errorf("expected 'test value', got '%v'", trimmedValue)
				}
			} else {
				t.Errorf("trimmedValue is nil")
			}
		})
	}
}

func TestStringSliceBuilderMax(t *testing.T) {
	phoneNumbers := []string{"911", "8982168233", "9876543210923"}
	testAdmin := admin{
		PhoneNumbers:    phoneNumbers,
		PhoneNumbersPtr: &phoneNumbers,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		maxLen        int
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumbers slice max length check",
			fieldName:     "PhoneNumbers",
			maxLen:        12,
			value:         &testAdmin.PhoneNumbers,
			expectedError: 1,
		},
		{
			name:          "PhoneNumbers slice max length check",
			fieldName:     "PhoneNumbers",
			maxLen:        15,
			value:         &testAdmin.PhoneNumbers,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.MaxLen(tc.maxLen)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringSlicePtrBuilderMaxLen(t *testing.T) {
	phoneNumbers := []string{"911", "8982168233", "9876543210923"}
	testAdmin := admin{
		PhoneNumbersPtr: &phoneNumbers,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		maxLen        int
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumbersPtr slice max length check",
			fieldName:     "PhoneNumbers",
			maxLen:        12,
			value:         &testAdmin.PhoneNumbersPtr,
			expectedError: 1,
		},
		{
			name:          "PhoneNumbersPtr slice max length check",
			fieldName:     "PhoneNumbersPtr",
			maxLen:        15,
			value:         &testAdmin.PhoneNumbersPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[*[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.MaxLen(tc.maxLen)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringSliceBuilderMinlen(t *testing.T) {
	testAdmin := admin{
		PhoneNumbers: []string{"911", "8982168233", "9876543210923"},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		minLength     int
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumbers slice min length check",
			fieldName:     "PhoneNumbers",
			minLength:     8,
			value:         &testAdmin.PhoneNumbers,
			expectedError: 1,
		},
		{
			name:          "PhoneNumbers slice min length check",
			fieldName:     "PhoneNumbers",
			minLength:     3,
			value:         &testAdmin.PhoneNumbers,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[*[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.MinLen(tc.minLength)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringSlicePtrBuilderMinLen(t *testing.T) {
	phoneNumber := []string{"911", "8982168233", "9876543210923"}
	testAdmin := admin{
		PhoneNumbersPtr: &phoneNumber,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		min           int
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumbersPtr slice min length check",
			fieldName:     "PhoneNumbersPtr",
			min:           8,
			value:         &testAdmin.PhoneNumbersPtr,
			expectedError: 1,
		},
		{
			name:          "PhoneNumbersPtr slice min length check",
			fieldName:     "PhoneNumbersPtr",
			min:           3,
			value:         &testAdmin.PhoneNumbersPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[*[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.MinLen(tc.min)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringSliceBuilderRequired(t *testing.T) {
	testAdmin := admin{
		PhoneNumbers: []string{"911", "8982168233", "9876543210923"},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "Roles slice required check",
			fieldName:     "Roles",
			value:         &testAdmin.Roles,
			expectedError: 1,
		},
		{
			name:          "PhoneNumbers slice required check",
			fieldName:     "PhoneNumbers",
			value:         &testAdmin.PhoneNumbers,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[[]string]{
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

func TestStringSlicePtrBuilderRequired(t *testing.T) {
	phoneNumbers := []string{"911", "8982168233", "9876543210923"}
	testAdmin := admin{
		PhoneNumbersPtr: &phoneNumbers,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "RolesPtr slice required check",
			fieldName:     "RolesPtr",
			value:         &testAdmin.RolesPtr,
			expectedError: 1,
		},
		{
			name:          "PhoneNumbersPtr slice required check",
			fieldName:     "PhoneNumbersPtr",
			value:         &testAdmin.PhoneNumbersPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[*[]string]{
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

func TestStringSliceBuilderRegexp(t *testing.T) {
	testAdmin := admin{
		Roles:        []string{"root", "read"},
		PhoneNumbers: []string{"911", "8982168233", "9876543210923"},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		re            *regexp.Regexp
		value         any
		expectedError int
	}{
		{
			name:          "Roles slice regex check",
			fieldName:     "Roles",
			re:            regexp.MustCompile(`^[0-9]{3,25}`),
			value:         &testAdmin.Roles,
			expectedError: 1,
		},
		{
			name:          "PhoneNumbers slice regex check",
			fieldName:     "PhoneNumbers",
			re:            regexp.MustCompile(`^[0-9]{3,25}`),
			value:         &testAdmin.PhoneNumbers,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Regexp(tc.re)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringSlicePtrBuilderRegexp(t *testing.T) {
	phoneNumbers := []string{"911", "8982168233", "9876543210923"}
	roles := []string{"root", "read"}
	testAdmin := admin{
		RolesPtr:        &roles,
		PhoneNumbersPtr: &phoneNumbers,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		re            *regexp.Regexp
		value         any
		expectedError int
	}{
		{
			name:          "RolesPtr slice regex check",
			fieldName:     "RolesPtr",
			re:            regexp.MustCompile(`^[0-9]{3,25}`),
			value:         &testAdmin.RolesPtr,
			expectedError: 1,
		},
		{
			name:          "PhoneNumbersPtr slice regex check",
			fieldName:     "PhoneNumbersPtr",
			re:            regexp.MustCompile(`^[0-9]{3,25}`),
			value:         &testAdmin.PhoneNumbersPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[*[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Regexp(tc.re)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringSliceBuilderCustom(t *testing.T) {
	testAdmin := admin{
		Roles: []string{"root", "read"},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumbers slice custom function check",
			fieldName:     "PhoneNumbers",
			value:         &testAdmin.PhoneNumbers,
			expectedError: 1,
		},
		{
			name:          "Roles slice custom function check",
			fieldName:     "Roles",
			value:         &testAdmin.Roles,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]string) []shared.Error {
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
	roles := []string{"root", "read"}
	testAdmin := admin{
		RolesPtr: &roles,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumbersPtr slice custom function check",
			fieldName:     "PhoneNumbersPtr",
			value:         &testAdmin.PhoneNumbersPtr,
			expectedError: 1,
		},
		{
			name:          "RolesPtr slice custom function check",
			fieldName:     "RolesPtr",
			value:         &testAdmin.RolesPtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringSliceBuilder[*[]string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]string) []shared.Error {
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
