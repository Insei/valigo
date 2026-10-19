package helper

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
	h := New()
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

func TestHelperFieldLocation(t *testing.T) {
	strg, _ := fmap.GetFrom(admin{})
	if got := New().FieldLocation(strg.MustFind("Name")); got != "Name" {
		t.Errorf("FieldLocation() = %q, want Name", got)
	}
}

// plainTranslator is a translator without the named parameters.
type plainTranslator struct{}

func (plainTranslator) ErrorT(_ context.Context, format string, _ ...any) error { return nil }
func (plainTranslator) T(_ context.Context, format string, _ ...any) string     { return "plain " + format }

func TestHelperErrorParams(t *testing.T) {
	strg, _ := fmap.GetFrom(admin{})
	field := strg.MustFind("Name")
	h := New()
	err := h.ErrorParams(context.Background(), field, "v", "validation:string:Cannot be longer than %d characters", []any{3}, map[string]any{"max": 3})
	if err.Message != "Cannot be longer than 3 characters" || err.Location != "Name" || err.Code != "validation:string:Cannot be longer than %d characters" {
		t.Errorf("ErrorParams: %+v", err)
	}
	h.Translator = plainTranslator{}
	if err := h.ErrorParams(context.Background(), field, "v", "k", nil, nil); err.Message != "plain k" {
		t.Errorf("ErrorParams of a translator without the parameters: %+v", err)
	}
}
