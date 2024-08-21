package shared

import (
	"context"
	"reflect"
	"testing"

	"github.com/insei/fmap/v3"
)

type user struct {
	Name string
}

func TestFieldCustomHelperErrorT(t *testing.T) {
	h := &mockHelper{}
	testUser := user{
		Name: "Alex",
	}
	strg, _ := fmap.GetFrom(testUser)
	field := strg.MustFind("Name")
	fch := NewFieldCustomHelper(field, h)

	ctx := context.Background()
	value := "Alex"
	localeKey := "testLocaleKey"
	args := []any{"arg1", "arg2"}

	h.ctx = ctx
	h.value = value
	h.localeKey = localeKey
	h.args = args
	_ = fch.ErrorT(ctx, testUser.Name, localeKey, args...)

	if !h.errorTCalled {
		t.Errorf("mockHelper.ErrorT was not called")
	}
	if h.ctx != ctx {
		t.Errorf("mockHelper.ctx = %v, want %v", h.ctx, ctx)
	}
	if h.field != fch.field {
		t.Errorf("mockHelper.field = %v, want %v", h.field, fch.field)
	}
	if h.value != value {
		t.Errorf("mockHelper.value = %v, want %v", h.value, value)
	}
	if h.localeKey != localeKey {
		t.Errorf("mockHelper.localeKey = %v, want %v", h.localeKey, localeKey)
	}
	if !reflect.DeepEqual(h.args, args) {
		t.Errorf("mockHelper.args = %v, want %v", h.args, args)
	}
}

func TestNewFieldCustomHelper(t *testing.T) {
	h := &mockHelper{}
	testUser := user{
		Name: "Alex",
	}
	strg, _ := fmap.GetFrom(testUser)
	field := strg.MustFind("Name")
	fch := NewFieldCustomHelper(field, h)

	if fch.field != field {
		t.Errorf("fch.field = %v, want %v", fch.field, field)
	}
	if fch.h != h {
		t.Errorf("fch.h = %v, want %v", fch.h, h)
	}
}

type mockHelper struct {
	errorTCalled bool
	ctx          context.Context
	field        fmap.Field
	value        any
	localeKey    string
	args         []any
}

func (m *mockHelper) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) Error {
	m.errorTCalled = true
	m.ctx = ctx
	m.field = field
	m.value = value
	m.localeKey = localeKey
	m.args = args
	return Error{}
}
