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
	testUser := user{Name: "   Rebecca "}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	field := storage.MustFind("Name")
	builder := stringBuilder[string]{
		field: field,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			fn(context.Background(), &helper, &testUser.Name)
		},
	}
	builder.Trim()
	trimmedValue := field.Get(&testUser)
	if trimmedValue != "Rebecca" {
		t.Errorf("expected 'test value', got '%v'", trimmedValue)
	}
}
func TestStringPtrBuilderTrim(t *testing.T) {
	name := "   Rebecca "
	testUser := user{NamePtr: &name}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}

	field := storage.MustFind("NamePtr")
	builder := stringBuilder[*string]{
		field: field,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			fn(context.Background(), &helper, &testUser.NamePtr)
		},
	}
	builder.Trim()
	trimmedValue := field.Get(&testUser).(*string)
	if *trimmedValue != "Rebecca" {
		t.Errorf("expected 'test value', got '%v'", trimmedValue)
	}
}

func TestStringBuilderMaxLen(t *testing.T) {
	testUser := user{
		PhoneNumber1: "987767213987377283",
		PhoneNumber2: "987767272",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("PhoneNumber1")
	builder1 := stringBuilder[string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.PhoneNumber1)
		},
	}
	builder1.MaxLen(12)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumber2")
	builder2 := stringBuilder[string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.PhoneNumber2)
		},
	}
	builder2.MaxLen(12)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
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
	var errs []shared.Error

	field1 := storage.MustFind("PhoneNumberPtr1")
	builder1 := stringBuilder[*string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.PhoneNumberPtr1)
		},
	}
	builder1.MaxLen(12)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumberPtr2")
	builder2 := stringBuilder[*string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.PhoneNumberPtr2)
		},
	}
	builder2.MaxLen(12)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringBuilderMinLen(t *testing.T) {
	testUser := user{
		PhoneNumber1: "911",
		PhoneNumber2: "987767272",
	}
	storage, _ := fmap.GetFrom(testUser)

	helper := helperStringImpl{}
	var errs []shared.Error
	field1 := storage.MustFind("PhoneNumber1")

	builder1 := stringBuilder[string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.PhoneNumber1)
		},
	}
	builder1.MinLen(9)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumber2")
	builder2 := stringBuilder[string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.PhoneNumber2)
		},
	}
	builder2.MaxLen(9)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
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
	var errs []shared.Error
	field1 := storage.MustFind("PhoneNumberPtr1")

	builder1 := stringBuilder[*string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.PhoneNumberPtr1)
		},
	}
	builder1.MinLen(9)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("PhoneNumberPtr2")
	builder2 := stringBuilder[*string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.PhoneNumberPtr2)
		},
	}
	builder2.MaxLen(9)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringBuilderRequired(t *testing.T) {
	testUser := user{
		Name: "Rebecca",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("LastName")
	builder1 := stringBuilder[string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.LastName)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("Name")
	builder2 := stringBuilder[string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Name)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringPtrBuilderRequired(t *testing.T) {
	name := "Rebecca"
	testUser := user{
		NamePtr: &name,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("LastNamePtr")
	builder1 := stringBuilder[*string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.LastNamePtr)
		},
	}
	builder1.Required()
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("NamePtr")
	builder2 := stringBuilder[*string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.NamePtr)
		},
	}
	builder2.Required()
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringBuilderRegexp(t *testing.T) {
	testUser := user{
		Name:     "R",
		LastName: "Smith",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}
	var errs []shared.Error
	re := regexp.MustCompile(`^[a-zA-Z]{3,25}`)

	field1 := storage.MustFind("Name")
	builder1 := stringBuilder[string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Name)
		},
	}
	builder1.Regexp(re)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("LastName")
	builder2 := stringBuilder[string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.LastName)
		},
	}
	builder2.Regexp(re)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
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
	var errs []shared.Error
	re := regexp.MustCompile(`^[a-zA-Z]{3,25}`)

	field1 := storage.MustFind("NamePtr")
	builder1 := stringBuilder[*string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.NamePtr)
		},
	}
	builder1.Regexp(re)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("LastNamePtr")
	builder2 := stringBuilder[*string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.LastNamePtr)
		},
	}
	builder2.Regexp(re)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringBuilderAnyOf(t *testing.T) {
	testUser := user{
		Name:     "Martin",
		LastName: "Smith",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}
	var errs []shared.Error
	allowedNames := []string{"Rebecca", "John", "Alex"}
	allowedLastNames := []string{"Smith", "Doe", "Williams"}

	field1 := storage.MustFind("Name")
	builder1 := stringBuilder[string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Name)
		},
	}
	builder1.AnyOf(allowedNames...)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("LastName")
	builder2 := stringBuilder[string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.LastName)
		},
	}
	builder2.AnyOf(allowedLastNames...)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
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
	var errs []shared.Error
	allowedNames := []string{"Rebecca", "John", "Alex"}
	allowedLastNames := []string{"Smith", "Doe", "Williams"}

	field1 := storage.MustFind("NamePtr")
	builder1 := stringBuilder[*string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.NamePtr)
		},
	}
	builder1.AnyOf(allowedNames...)
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("LastNamePtr")
	builder2 := stringBuilder[*string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.LastNamePtr)
		},
	}
	builder2.AnyOf(allowedLastNames...)
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringBuilderCustom(t *testing.T) {
	testUser := user{
		LastName: "Smith",
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("Name")
	builder1 := stringBuilder[string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.Name)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *string) []shared.Error {
		if value == nil || *value == "" {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("LastName")
	builder2 := stringBuilder[string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.LastName)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *string) []shared.Error {
		if value == nil || *value == "" {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}

func TestStringPtrBuilderCustom(t *testing.T) {
	lastName := "Smith"
	testUser := user{
		LastNamePtr: &lastName,
	}
	storage, _ := fmap.GetFrom(testUser)
	helper := helperStringImpl{}
	var errs []shared.Error

	field1 := storage.MustFind("NamePtr")
	builder1 := stringBuilder[*string]{
		field: field1,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.NamePtr)
		},
	}
	builder1.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **string) []shared.Error {
		if value == nil || *value == nil || **value == "" {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) == 0 {
		t.Errorf("expected error, got nil")
	}

	field2 := storage.MustFind("LastNamePtr")
	builder2 := stringBuilder[*string]{
		field: field2,
		h:     &helper,
		appendFn: func(field fmap.Field, fn shared.FieldValidationFn) {
			errs = fn(context.Background(), &helper, &testUser.LastNamePtr)
		},
	}
	builder2.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **string) []shared.Error {
		if value == nil || *value == nil || **value == "" {
			return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
		}
		return nil
	})
	if len(errs) > 0 {
		t.Errorf("expected nil, got %v", errs)
	}
}
