package valigo

import (
	"context"
	"reflect"
	"testing"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

func TestStorageNewOnStruct(t *testing.T) {
	s := newStorage()
	temp := struct{}{}
	enabler := func(ctx context.Context, obj any) bool {
		return true
	}
	fn := func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
		return nil
	}
	s.newOnStruct(temp, enabler, fn)
	if len(s.validators[reflect.TypeOf(temp)]) != 1 {
		t.Errorf("expected 1 validator, got %d", len(s.validators[reflect.TypeOf(temp)]))
	}
}

func TestStorageNewOnStructWithNilEnabler(t *testing.T) {
	s := newStorage()
	temp := struct{}{}
	var enabler func(context.Context, any) bool
	fn := func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
		return nil
	}
	s.newOnStruct(temp, enabler, fn)
	if len(s.validators[reflect.TypeOf(temp)]) != 1 {
		t.Errorf("expected 1 validator, got %d", len(s.validators[reflect.TypeOf(temp)]))
	}
}

type user struct {
	Name string
	Age  int
}

func TestStorageNewOnField(t *testing.T) {
	s := newStorage()
	temp := struct{}{}
	enabler := func(ctx context.Context, obj any) bool {
		return true
	}
	testUser := user{
		Name: "Alex",
		Age:  25,
	}
	strg, _ := fmap.GetFrom(testUser)
	field := strg.MustFind("Age")
	fn := func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
		return nil
	}
	s.newOnField(temp, enabler)(field, fn)
	if len(s.validators[reflect.TypeOf(temp)]) != 1 {
		t.Errorf("expected 1 validator, got %d", len(s.validators[reflect.TypeOf(temp)]))
	}
}

func TestStorageNewOnFieldWithNilEnabler(t *testing.T) {
	s := newStorage()
	temp := struct{}{}
	var enabler func(context.Context, any) bool
	testUser := user{
		Name: "Alex",
		Age:  25,
	}
	strg, _ := fmap.GetFrom(testUser)
	field := strg.MustFind("Name")
	fn := func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
		return nil
	}
	s.newOnField(temp, enabler)(field, fn)
	if len(s.validators[reflect.TypeOf(temp)]) != 1 {
		t.Errorf("expected 1 validator, got %d", len(s.validators[reflect.TypeOf(temp)]))
	}
}
