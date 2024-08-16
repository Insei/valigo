package valigo

import (
	"context"
	"errors"
	"testing"

	"github.com/insei/fmap/v3"
)

type admin struct {
	Name string
	Age  int
}

func TestHelperErrorT(t *testing.T) {
	h := newHelper()
	testAdmin := admin{
		Name: "John",
		Age:  33,
	}
	strg, _ := fmap.GetFrom(testAdmin)
	field := strg.MustFind("Name")
	value := "test value"
	localeKey := "test.key"
	args := []any{"arg1", "arg2"}

	expectedError := errors.New("test.key%!(EXTRA string=arg1, string=arg2) (Name: test value)")
	err := h.ErrorT(context.Background(), field, value, localeKey, args...)
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error to be %v, got %v", expectedError.Error(), err.Message)
	}
}
