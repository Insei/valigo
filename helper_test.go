package valigo

import (
	"context"
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

	err := h.ErrorT(context.Background(), field, value, localeKey, args...)
	if err.Error() == "" {
		t.Errorf("expected error, got %v", err.Message)
	}
}
