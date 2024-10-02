package valigo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo/shared"
)

func TestBuilderWhen(t *testing.T) {
	type TestStruct struct {
		Field string
	}
	obj := &TestStruct{
		Field: "test",
	}
	validator = New()
	bld := configure[TestStruct](validator, obj, func(ctx context.Context, obj any) bool {
		return true
	})
	bld.When(func(ctx context.Context, obj *TestStruct) bool {
		return (*obj).Field == "test"
	}).String(&obj.Field).Required()

	assert.Equal(t, 1, len(validator.storage.validators))
}

func TestBuilderCustom(t *testing.T) {
	type TestStruct struct {
		Field string
	}
	obj := &TestStruct{
		Field: "test",
	}
	vld := New()
	bld := configure[*TestStruct](vld, obj, func(ctx context.Context, obj any) bool {
		return true
	})
	bld.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj **TestStruct) []shared.Error {
		if (*obj).Field != "test" {
			return []shared.Error{h.ErrorT(ctx, &(*obj).Field, (*obj).Field, "validation:string:Field is not 'test'")}
		}
		return nil
	})

	assert.Equal(t, 1, len(vld.storage.validators))
}

func TestBuilderWhenAndCustomMultipleConditions(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 string
	}
	obj := &TestStruct{
		Field1: "test1",
		Field2: "test2",
	}
	vld := New()
	bld := configure[TestStruct](vld, obj, func(ctx context.Context, obj any) bool {
		return true
	})
	bld.When(func(ctx context.Context, obj *TestStruct) bool {
		return (*obj).Field1 == "test1"
	}).When(func(ctx context.Context, obj *TestStruct) bool {
		return (*obj).Field2 == "test2"
	}).Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *TestStruct) []shared.Error {
		if (*obj).Field1 != "test1" {
			return []shared.Error{h.ErrorT(ctx, &(*obj).Field1, (*obj).Field1, "validation:string:Field1 is not 'test1'")}
		}
		if (*obj).Field2 != "test2" {
			return []shared.Error{h.ErrorT(ctx, &(*obj).Field2, (*obj).Field2, "validation:string:Field2 is not 'test2'")}
		}
		return nil
	})

	assert.Equal(t, 1, len(vld.storage.validators))
}
func TestBuilderCustomMultipleErrors(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 string
	}
	obj := &TestStruct{
		Field1: "test1",
		Field2: "test2",
	}
	vld := New()
	bld := configure[*TestStruct](vld, obj, func(ctx context.Context, obj any) bool {
		return true
	})
	bld.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj **TestStruct) []shared.Error {
		errs := make([]shared.Error, 0)
		if (*obj).Field1 != "test1" {
			errs = append(errs, h.ErrorT(ctx, &(*obj).Field1, (*obj).Field1, "validation:string:Field1 is not 'test1'"))
		}
		if (*obj).Field2 != "test2" {
			errs = append(errs, h.ErrorT(ctx, &(*obj).Field2, (*obj).Field2, "validation:string:Field2 is not 'test2'"))
		}
		return errs
	})

	assert.Equal(t, 1, len(vld.storage.validators))
}

func TestConfigurePanic(t *testing.T) {
	type TestStruct struct {
		Field1 string
	}
	obj := &TestStruct{
		Field1: "test1",
	}
	vld := New()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("configure did not panic when GetFrom returned an error")
		}
	}()
	configure[TestStruct](vld, &obj, nil)
}

func TestBuilder_Slice(t *testing.T) {
	type TestStruct struct {
		Slice  []string
		PtrInt float32
	}
	obj := &TestStruct{
		Slice:  []string{"  test1   ", "  test2  "},
		PtrInt: 44,
	}
	vld := New()
	bld := configure[*TestStruct](vld, obj, nil)
	bld.Slice(&obj.Slice).MaxLen(1)
	bld.Number(&obj.PtrInt).Max(float32(22.22))
	err := vld.Validate(context.Background(), obj)
	_ = err
}
