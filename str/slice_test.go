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
	roles := []string{"   root ", "read  "}
	testAdmin := admin{
		Roles:    roles,
		RolesPtr: &roles,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringSliceImpl{}

	field := storage.MustFind("Roles")
	builder := stringSliceBuilder[[]string]{
		field: field,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			fn(context.Background(), &helper, &testAdmin.Roles)
		},
	}
	builder.Trim()
	trimmedValue := field.Get(&testAdmin)
	trimmedSliceValue := trimmedValue.([]string)

	if trimmedSliceValue[0] != "root" || trimmedSliceValue[1] != "read" {
		t.Errorf("expected 'test value', got '%v'", trimmedValue)
	}
}

func TestStringSlicePtrBuilderTrim(t *testing.T) {
	roles := []string{"   root ", "read  "}
	testAdmin := admin{
		RolesPtr: &roles,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringSliceImpl{}

	field := storage.MustFind("RolesPtr")
	builder := stringSliceBuilder[*[]string]{
		field: field,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			fn(context.Background(), &helper, &testAdmin.RolesPtr)
		},
	}
	builder.Trim()
	trimmedValue := field.Get(&testAdmin)
	if trimmedValue != nil {
		trimmedSliceValue := trimmedValue.(*[]string)
		if (*trimmedSliceValue)[0] != "root" || (*trimmedSliceValue)[1] != "read" {
			t.Errorf("expected 'test value', got '%v'", trimmedValue)
		}
	} else {
		t.Errorf("trimmedValue is nil")
	}
}

func TestStringSliceBuilderMaxLen(t *testing.T) {
	phoneNumbers := []string{"911", "8982168233", "9876543210923"}
	testAdmin := admin{
		PhoneNumbers:    phoneNumbers,
		PhoneNumbersPtr: &phoneNumbers,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("PhoneNumbers")
	builder1 := stringSliceBuilder[[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbers)
		},
	}
	builder1.MaxLen(12)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumbers")
	builder2 := stringSliceBuilder[[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbers)
		},
	}
	builder2.MaxLen(15)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringSlicePtrBuilderMaxLen(t *testing.T) {
	phoneNumbers := []string{"911", "8982168233", "9876543210923"}
	testAdmin := admin{
		PhoneNumbersPtr: &phoneNumbers,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("PhoneNumbersPtr")
	builder1 := stringSliceBuilder[*[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbersPtr)
		},
	}
	builder1.MaxLen(12)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumbersPtr")
	builder2 := stringSliceBuilder[*[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbersPtr)
		},
	}
	builder2.MaxLen(15)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringSliceBuilderMinLen(t *testing.T) {
	testAdmin := admin{
		PhoneNumbers: []string{"911", "8982168233", "9876543210923"},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("PhoneNumbers")
	builder1 := stringSliceBuilder[[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbers)
		},
	}
	builder1.MinLen(8)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumbers")
	builder2 := stringSliceBuilder[[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbers)
		},
	}
	builder2.MinLen(3)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringSlicePtrBuilderMinLen(t *testing.T) {
	phoneNumber := []string{"911", "8982168233", "9876543210923"}
	testAdmin := admin{
		PhoneNumbersPtr: &phoneNumber,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("PhoneNumbersPtr")
	builder1 := stringSliceBuilder[*[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbersPtr)
		},
	}
	builder1.MinLen(8)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumbersPtr")
	builder2 := stringSliceBuilder[*[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbers)
		},
	}
	builder2.MinLen(3)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringSliceBuilderRequired(t *testing.T) {
	testAdmin := admin{
		PhoneNumbers: []string{"911", "8982168233", "9876543210923"},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("Roles")
	builder1 := stringSliceBuilder[[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.Roles)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumbers")
	builder2 := stringSliceBuilder[[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbers)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringSlicePtrBuilderRequired(t *testing.T) {
	phoneNumbers := []string{"911", "8982168233", "9876543210923"}
	testAdmin := admin{
		PhoneNumbersPtr: &phoneNumbers,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("RolesPtr")
	builder1 := stringSliceBuilder[*[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RolesPtr)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumbersPtr")
	builder2 := stringSliceBuilder[*[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbersPtr)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringSliceBuilderRegexp(t *testing.T) {
	testAdmin := admin{
		Roles:        []string{"root", "read"},
		PhoneNumbers: []string{"911", "8982168233", "9876543210923"},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error
	re := regexp.MustCompile(`^[0-9]{3,25}`)

	field1 := storage.MustFind("Roles")
	builder1 := stringSliceBuilder[[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.Roles)
		},
	}
	builder1.Regexp(re)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumbers")
	builder2 := stringSliceBuilder[[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbers)
		},
	}
	builder2.Regexp(re)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
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
	var errs []shared.Error
	re := regexp.MustCompile(`^[0-9]{3,25}`)

	field1 := storage.MustFind("RolesPtr")
	builder1 := stringSliceBuilder[*[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RolesPtr)
		},
	}
	builder1.Regexp(re)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumbersPtr")
	builder2 := stringSliceBuilder[*[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbersPtr)
		},
	}
	builder2.Regexp(re)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringSliceBuilderCustom(t *testing.T) {
	testAdmin := admin{
		Roles: []string{"root", "read"},
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("PhoneNumbers")
	builder1 := stringSliceBuilder[[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbers)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]string) []shared.Error {
		if len(*value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Roles")
	builder2 := stringSliceBuilder[[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.Roles)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *[]string) []shared.Error {
		if len(*value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringSlicePtrBuilderCustom(t *testing.T) {
	roles := []string{"root", "read"}
	testAdmin := admin{
		RolesPtr: &roles,
	}
	storage, _ := fmap.GetFrom(testAdmin)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("PhoneNumbersPtr")
	builder1 := stringSliceBuilder[*[]string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.PhoneNumbersPtr)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]string) []shared.Error {
		if *value == nil || len(**value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("RolesPtr")
	builder2 := stringSliceBuilder[*[]string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testAdmin.RolesPtr)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **[]string) []shared.Error {
		if *value == nil || len(**value) == 0 {
			return []shared.Error{h.ErrorT(ctx, value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}
