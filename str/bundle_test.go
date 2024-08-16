package str

import (
	"context"
	"fmt"
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type testStruct struct {
	Str         string
	StrPtr      *string
	StrSlice    []string
	StrSlicePtr *[]string
}

type helperImpl struct{}

func (h *helperImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
	return shared.Error{
		Message: fmt.Sprintf(localeKey, value),
	}
}

func TestNewStringBundle(t *testing.T) {
	// Test case 1: Valid input
	test := testStruct{
		Str: "test",
	}
	storage, _ := fmap.GetFrom(test)
	helper := helperImpl{}
	deps := shared.BundleDependencies{
		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
		Fields:   storage,
		Object:   struct{}{},
		Helper:   &helper,
	}
	sb := NewStringBundle(deps)

	if sb.storage != deps.Fields {
		t.Errorf("storage not set correctly")
	}
	if sb.obj != deps.Object {
		t.Errorf("obj not set correctly")
	}
	if sb.h != deps.Helper {
		t.Errorf("h not set correctly")
	}

	// Test case 2: Nil storage
	deps.Fields = nil
	sb = NewStringBundle(deps)
	if sb.storage != nil {
		t.Errorf("storage should be nil")
	}

	// Test case 3: Nil helper
	deps.Fields = storage
	deps.Helper = nil
	sb = NewStringBundle(deps)
	if sb.h != nil {
		t.Errorf("helper should be nil")
	}

	// Test case 4: Invalid AppendFn
	deps.Helper = &helper
	deps.AppendFn = func(field fmap.Field, fn shared.FieldValidationFn) { panic("invalid AppendFn") }
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic")
		}
	}()
	sb = NewStringBundle(deps)
	sb.appendFn(nil, nil)

}

func TestStringBundleString(t *testing.T) {
	test := testStruct{
		Str: "test",
	}
	storage, err := fmap.GetFrom(test)
	if err != nil {
		t.Fatal(err)
	}
	helper := helperImpl{}
	deps := shared.BundleDependencies{
		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
		Fields:   storage,
		Object:   &test,
		Helper:   &helper,
	}

	testCases := []struct {
		name          string
		field         *string
		expectedError bool
	}{
		{
			name:          "Valid field",
			field:         &test.Str,
			expectedError: false,
		},
		{
			name:          "Invalid field",
			field:         (*string)(nil),
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tc.expectedError {
						t.Logf("Expected error: %v", r)
					} else {
						t.Errorf("Unexpected error: %v", r)
					}
				}
			}()

			sb := NewStringBundle(deps)
			sbStr := sb.String(tc.field)
			if sbStr == nil && !tc.expectedError {
				t.Errorf("sbStr should not be nil")
			}
		})
	}
}

func TestStringBundleStringPtr(t *testing.T) {
	temp := "test"
	test := testStruct{
		StrPtr: &temp,
	}
	storage, err := fmap.GetFrom(test)
	if err != nil {
		t.Fatal(err)
	}
	helper := helperImpl{}
	deps := shared.BundleDependencies{
		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
		Fields:   storage,
		Object:   &test,
		Helper:   &helper,
	}
	sb := NewStringBundle(deps)
	sbStr := sb.StringPtr(&test.StrPtr)
	if sbStr == nil {
		t.Errorf("field not set correctly")
	}

	// Test case with an invalid field pointer
	invalidField := (**string)(nil)
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); !ok {
				t.Errorf("expected error, got %v", r)
			} else if err.Error() != "field not found" {
				t.Errorf("expected error 'field not found', got %v", err)
			}
		} else {
			t.Errorf("expected panic, got none")
		}
	}()
	sb.StringPtr(invalidField)
}
func TestStringBundleStringSlice(t *testing.T) {
	test := testStruct{
		StrSlice: []string{"test1", "test2"},
	}
	storage, err := fmap.GetFrom(test)
	if err != nil {
		t.Fatal(err)
	}
	helper := helperImpl{}
	deps := shared.BundleDependencies{
		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
		Fields:   storage,
		Object:   &test,
		Helper:   &helper,
	}
	sb := NewStringBundle(deps)
	modelsPtr := &test.StrSlice
	sbStr := sb.StringSlice(modelsPtr)
	if sbStr == nil {
		t.Errorf("field not set correctly")
	}

	// Test case with an invalid field pointer
	invalidField := (*[]string)(nil)
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); !ok {
				t.Errorf("expected error, got %v", r)
			} else if err.Error() != "field not found" {
				t.Errorf("expected error 'field not found', got %v", err)
			}
		} else {
			t.Errorf("expected panic, got none")
		}
	}()
	sb.StringSlice(invalidField)
}

func TestStringBundleStringSlicePtr(t *testing.T) {
	temp := []string{"test1", "test2"}
	test := testStruct{
		StrSlicePtr: &temp,
	}
	storage, err := fmap.GetFrom(test)
	if err != nil {
		t.Fatal(err)
	}
	helper := helperImpl{}
	deps := shared.BundleDependencies{
		AppendFn: func(field fmap.Field, fn shared.FieldValidationFn) {},
		Fields:   storage,
		Object:   &test,
		Helper:   &helper,
	}
	sb := NewStringBundle(deps)
	sbStr := sb.StringSlicePtr(&test.StrSlicePtr)
	if sbStr == nil {
		t.Errorf("field not set correctly")
	}

	// Test case with an invalid field pointer
	invalidField := (**[]string)(nil)
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); !ok {
				t.Errorf("expected error, got %v", r)
			} else if err.Error() != "field not found" {
				t.Errorf("expected error 'field not found', got %v", err)
			}
		} else {
			t.Errorf("expected panic, got none")
		}
	}()
	sb.StringSlicePtr(invalidField)
}
