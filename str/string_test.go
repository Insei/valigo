package str

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type helperStringImpl struct{}

func (h *helperStringImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

type user struct {
	Name            string
	LastName        string
	PhoneNumber1    string
	PhoneNumber2    string
	NamePtr         *string
	LastNamePtr     *string
	PhoneNumberPtr1 *string
	PhoneNumberPtr2 *string
}

func TestStringBuilderTrim(t *testing.T) {
	testUser := user{
		Name:     "   Rebecca ",
		LastName: "Smith  ",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		trimmedValue  string
		value         any
		expectedError bool
	}{
		{
			name:          "Name trim check",
			fieldName:     "Name",
			trimmedValue:  " Rebecca",
			value:         &testUser.Name,
			expectedError: true,
		},
		{
			name:          "LastName trim check",
			fieldName:     "LastName",
			trimmedValue:  "Smith",
			value:         &testUser.LastName,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Trim()
			trimmedValue := field.Get(&testUser)
			if trimmedValue != tc.trimmedValue && !tc.expectedError {
				t.Errorf("expected '%v', got '%v'", tc.trimmedValue, trimmedValue)
			}
		})
	}
}

func TestStringPtrBuilderTrim(t *testing.T) {
	name := "   Rebecca "
	lastName := "Smith  "
	testUser := user{
		NamePtr:     &name,
		LastNamePtr: &lastName,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		trimmedValue  string
		value         any
		expectedError bool
	}{
		{
			name:          "NamePtr trim check",
			fieldName:     "NamePtr",
			trimmedValue:  " Rebecca",
			value:         &testUser.NamePtr,
			expectedError: true,
		},
		{
			name:          "LastNamePtr trim check",
			fieldName:     "LastNamePtr",
			trimmedValue:  "Smith",
			value:         &testUser.LastNamePtr,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[*string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Trim()
			trimmedValue := field.Get(&testUser).(*string)
			if *trimmedValue != tc.trimmedValue && !tc.expectedError {
				t.Errorf("expected '%v', got '%v'", tc.trimmedValue, *trimmedValue)
			}
		})
	}
}

func TestStringBuilderMaxLen(t *testing.T) {
	testUser := user{
		PhoneNumber1: "987767213987377283",
		PhoneNumber2: "987767272",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		maxLen        int
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumber1 max length check",
			fieldName:     "PhoneNumber1",
			maxLen:        12,
			value:         &testUser.PhoneNumber1,
			expectedError: 1,
		},
		{
			name:          "PhoneNumber2 max length check",
			fieldName:     "PhoneNumber2",
			maxLen:        12,
			value:         &testUser.PhoneNumber2,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[string]{
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

func TestStringPtrBuilderMaxLen(t *testing.T) {
	phoneNumber1 := "987767213987377283"
	phoneNumber2 := "987767272"
	testUser := user{
		PhoneNumberPtr1: &phoneNumber1,
		PhoneNumberPtr2: &phoneNumber2,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		maxLen        int
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumberPtr1 max length check",
			fieldName:     "PhoneNumberPtr1",
			maxLen:        12,
			value:         &testUser.PhoneNumberPtr1,
			expectedError: 1,
		},
		{
			name:          "PhoneNumberPtr2 max length check",
			fieldName:     "PhoneNumberPtr2",
			maxLen:        12,
			value:         &testUser.PhoneNumberPtr2,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[*string]{
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

func TestStringBuilderMinLen(t *testing.T) {
	testUser := user{
		PhoneNumber1: "911",
		PhoneNumber2: "987767272",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		minLen        int
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumber1 min length check",
			fieldName:     "PhoneNumber1",
			minLen:        9,
			value:         &testUser.PhoneNumber1,
			expectedError: 1,
		},
		{
			name:          "PhoneNumber2 min length check",
			fieldName:     "PhoneNumber2",
			minLen:        9,
			value:         &testUser.PhoneNumber2,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.MinLen(tc.minLen)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringPtrBuilderMinLen(t *testing.T) {
	phoneNumber1 := "911"
	phoneNumber2 := "987767272"
	testUser := user{
		PhoneNumberPtr1: &phoneNumber1,
		PhoneNumberPtr2: &phoneNumber2,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		minLen        int
		value         any
		expectedError int
	}{
		{
			name:          "PhoneNumberPtr1 min length check",
			fieldName:     "PhoneNumberPtr1",
			minLen:        9,
			value:         &testUser.PhoneNumberPtr1,
			expectedError: 1,
		},
		{
			name:          "PhoneNumberPtr2 min length check",
			fieldName:     "PhoneNumberPtr2",
			minLen:        9,
			value:         &testUser.PhoneNumberPtr2,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[*string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.MinLen(tc.minLen)
			if len(errs) != tc.expectedError {
				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
			}
		})
	}
}

func TestStringBuilderRequired(t *testing.T) {
	testUser := user{
		Name: "Rebecca",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "LastName required check",
			fieldName:     "LastName",
			value:         &testUser.LastName,
			expectedError: 1,
		},
		{
			name:          "Name required check",
			fieldName:     "Name",
			value:         &testUser.Name,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[string]{
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

func TestStringPtrBuilderRequired(t *testing.T) {
	name := "Rebecca"
	testUser := user{
		NamePtr: &name,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "LastNamePtr required check",
			fieldName:     "LastNamePtr",
			value:         &testUser.LastNamePtr,
			expectedError: 1,
		},
		{
			name:          "NamePtr required check",
			fieldName:     "NamePtr",
			value:         &testUser.NamePtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[*string]{
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

func TestStringBuilderRegexp(t *testing.T) {
	testUser := user{
		Name:     "R",
		LastName: "Smith",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		re            *regexp.Regexp
		value         any
		expectedError int
	}{
		{
			name:          "Name regex check",
			fieldName:     "Name",
			re:            regexp.MustCompile(`^[a-zA-Z]{3,25}`),
			value:         &testUser.Name,
			expectedError: 1,
		},
		{
			name:          "LastName regex check",
			fieldName:     "LastName",
			re:            regexp.MustCompile(`^[a-zA-Z]{3,25}`),
			value:         &testUser.LastName,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[string]{
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

func TestStringPtrBuilderRegexp(t *testing.T) {
	name := "R"
	lastName := "Smith"
	testUser := user{
		NamePtr:     &name,
		LastNamePtr: &lastName,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		re            *regexp.Regexp
		value         any
		expectedError int
	}{
		{
			name:          "NamePtr regex check",
			fieldName:     "NamePtr",
			re:            regexp.MustCompile(`^[a-zA-Z]{3,25}`),
			value:         &testUser.NamePtr,
			expectedError: 1,
		},
		{
			name:          "LastNamePtr regex check",
			fieldName:     "LastNamePtr",
			re:            regexp.MustCompile(`^[a-zA-Z]{3,25}`),
			value:         &testUser.LastNamePtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[*string]{
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

func TestStringBuilderAnyOf(t *testing.T) {
	testUser := user{
		Name:     "Martin",
		LastName: "Smith",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		allowedValues []string
		value         any
		expectedError int
	}{
		{
			name:          "Name any of check",
			fieldName:     "Name",
			allowedValues: []string{"Rebecca", "John", "Alex"},
			value:         &testUser.Name,
			expectedError: 1,
		},
		{
			name:          "LastName any of check",
			fieldName:     "LastName",
			allowedValues: []string{"Smith", "Doe", "Williams"},
			value:         &testUser.LastName,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[string]{
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
	name := "Martin"
	lastName := "Smith"
	testUser := user{
		NamePtr:     &name,
		LastNamePtr: &lastName,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		allowedValues []string
		value         any
		expectedError int
	}{
		{
			name:          "NamePtr any of check",
			fieldName:     "NamePtr",
			allowedValues: []string{"Rebecca", "John", "Alex"},
			value:         &testUser.NamePtr,
			expectedError: 1,
		},
		{
			name:          "LastNamePtr any of check",
			fieldName:     "LastNamePtr",
			allowedValues: []string{"Smith", "Doe", "Williams"},
			value:         &testUser.LastNamePtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[*string]{
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

func TestStringBuilderCustom(t *testing.T) {
	testUser := user{
		LastName: "Smith",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "Name custom function check",
			fieldName:     "Name",
			value:         &testUser.Name,
			expectedError: 1,
		},
		{
			name:          "LastName custom function check",
			fieldName:     "LastName",
			value:         &testUser.LastName,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *string) []shared.Error {
				if value == nil || *value == "" {
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

func TestStringPtrBuilderCustom(t *testing.T) {
	lastName := "Smith"
	testUser := user{
		LastNamePtr: &lastName,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	testCases := []struct {
		name          string
		fieldName     string
		value         any
		expectedError int
	}{
		{
			name:          "NamePtr custom function check",
			fieldName:     "NamePtr",
			value:         &testUser.NamePtr,
			expectedError: 1,
		},
		{
			name:          "LastNamePtr custom function check",
			fieldName:     "LastNamePtr",
			value:         &testUser.LastNamePtr,
			expectedError: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var errs []shared.Error
			field := storage.MustFind(tc.fieldName)
			builder := stringBuilder[*string]{
				field: field,
				h:     &helper,
				appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
					errs = fn(context.Background(), &helper, tc.value)
				},
			}
			builder.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **string) []shared.Error {
				if value == nil || *value == nil || **value == "" {
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
