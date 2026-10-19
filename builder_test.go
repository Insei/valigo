package valigo

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
	vuuid "github.com/insei/valigo/uuid"
)

type status string

type address struct {
	City string
}

type model struct {
	Name    string
	NamePtr *string
	Status  status
	Age     int
	AgePtr  *int
	ID      int64
	Score   float32
	UID     uuid.UUID
	UIDPtr  *uuid.UUID
	Tags    []string
	TagsPtr *[]string
	TagPtrs []*string
	Ints    []int
	UIDs    []*uuid.UUID
	Addr    *address
	Inner   address
}

func ptr[T any](v T) *T {
	return &v
}

func validateModel(t *testing.T, cfg func(c Configurator[model], m *model), obj any) []shared.Error {
	t.Helper()
	v := New()
	Configure[model](v, cfg)
	return v.ValidateTyped(context.Background(), obj)
}

func TestConfiguratorRules(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	tests := []struct {
		name     string
		cfg      func(c Configurator[model], m *model)
		obj      *model
		wantErrs int
	}{
		{"string max len ok", func(c Configurator[model], m *model) { c.String(&m.Name).MaxLen(3) }, &model{Name: "abc"}, 0},
		{"string max len fail", func(c Configurator[model], m *model) { c.String(&m.Name).MaxLen(3) }, &model{Name: "abcd"}, 1},
		{"string min len fail", func(c Configurator[model], m *model) { c.String(&m.Name).MinLen(3) }, &model{Name: "ab"}, 1},
		{"string required fail", func(c Configurator[model], m *model) { c.String(&m.Name).Required() }, &model{}, 1},
		{"string any of", func(c Configurator[model], m *model) { c.String(&m.Name).AnyOf("a", "b") }, &model{Name: "c"}, 1},
		{"string regexp", func(c Configurator[model], m *model) { c.String(&m.Name).Regexp(regexp.MustCompile(`^\d+$`)) }, &model{Name: "12a"}, 1},
		{"string email ok", func(c Configurator[model], m *model) { c.String(&m.Name).Email() }, &model{Name: "a@b.cc"}, 0},
		{"string email fail", func(c Configurator[model], m *model) { c.String(&m.Name).Email() }, &model{Name: "a.b.cc"}, 1},
		{"string chain errors", func(c Configurator[model], m *model) { c.String(&m.Name).MinLen(5).Email() }, &model{Name: "ab"}, 2},
		{"named string type", func(c Configurator[model], m *model) { c.String(&m.Status).AnyOf("on") }, &model{Status: "off"}, 1},
		{"nil *string skipped", func(c Configurator[model], m *model) { c.String(&m.NamePtr).MaxLen(1).MinLen(5) }, &model{}, 0},
		{"nil *string required", func(c Configurator[model], m *model) { c.String(&m.NamePtr).Required() }, &model{}, 1},
		{"*string checked", func(c Configurator[model], m *model) { c.String(&m.NamePtr).MaxLen(1) }, &model{NamePtr: ptr("ab")}, 1},
		{"nil nested struct skipped", func(c Configurator[model], m *model) { c.String(&m.Addr.City).MaxLen(1) }, &model{}, 0},
		{"nil nested struct required", func(c Configurator[model], m *model) { c.String(&m.Addr.City).Required() }, &model{}, 1},
		{"nested struct checked", func(c Configurator[model], m *model) { c.String(&m.Addr.City).MaxLen(1) }, &model{Addr: &address{City: "ab"}}, 1},
		{"inner struct checked", func(c Configurator[model], m *model) { c.String(&m.Inner.City).MaxLen(1) }, &model{Inner: address{City: "ab"}}, 1},

		{"number max inclusive", func(c Configurator[model], m *model) { c.Number(&m.Age).Max(10) }, &model{Age: 10}, 0},
		{"number max fail", func(c Configurator[model], m *model) { c.Number(&m.Age).Max(10) }, &model{Age: 11}, 1},
		{"number min inclusive", func(c Configurator[model], m *model) { c.Number(&m.Age).Min(10) }, &model{Age: 10}, 0},
		{"number min fail", func(c Configurator[model], m *model) { c.Number(&m.Age).Min(10) }, &model{Age: 9}, 1},
		{"number int64 untyped const", func(c Configurator[model], m *model) { c.Number(&m.ID).Max(5) }, &model{ID: 6}, 1},
		{"number float32", func(c Configurator[model], m *model) { c.Number(&m.Score).Max(1.5) }, &model{Score: 1.6}, 1},
		{"number any of", func(c Configurator[model], m *model) { c.Number(&m.Age).AnyOf(1, 2) }, &model{Age: 3}, 1},
		{"number interval inclusive", func(c Configurator[model], m *model) { c.Number(&m.Age).AnyOfInterval(1, 10) }, &model{Age: 10}, 0},
		{"number interval fail", func(c Configurator[model], m *model) { c.Number(&m.Age).AnyOfInterval(1, 10) }, &model{Age: 11}, 1},
		{"number required zero", func(c Configurator[model], m *model) { c.Number(&m.Age).Required() }, &model{}, 1},
		{"number ptr nil skipped", func(c Configurator[model], m *model) { c.NumberPtr(&m.AgePtr).Max(10).Min(20) }, &model{}, 0},
		{"number ptr nil required", func(c Configurator[model], m *model) { c.NumberPtr(&m.AgePtr).Required() }, &model{}, 1},
		{"number ptr zero required", func(c Configurator[model], m *model) { c.NumberPtr(&m.AgePtr).Required() }, &model{AgePtr: ptr(0)}, 0},
		{"number ptr checked", func(c Configurator[model], m *model) { c.NumberPtr(&m.AgePtr).Max(10) }, &model{AgePtr: ptr(11)}, 1},

		{"uuid required", func(c Configurator[model], m *model) { c.UUID(&m.UID).Required() }, &model{}, 1},
		{"uuid any of", func(c Configurator[model], m *model) { c.UUID(&m.UID).AnyOf(id1) }, &model{UID: id2}, 1},
		{"uuid ptr nil skipped", func(c Configurator[model], m *model) { c.UUID(&m.UIDPtr).AnyOf(id1) }, &model{}, 0},
		{"uuid ptr nil required", func(c Configurator[model], m *model) { c.UUID(&m.UIDPtr).Required() }, &model{}, 1},

		{"slice max len", func(c Configurator[model], m *model) { c.Slice(&m.Ints).MaxLen(2) }, &model{Ints: []int{1, 2, 3}}, 1},
		{"slice min len", func(c Configurator[model], m *model) { c.Slice(&m.Ints).MinLen(2) }, &model{Ints: []int{1}}, 1},
		{"slice required", func(c Configurator[model], m *model) { c.Slice(&m.Ints).Required() }, &model{Ints: []int{}}, 1},
		{"slice each", func(c Configurator[model], m *model) {
			c.Slice(&m.Ints).Each(func(e *shared.Field[int]) { e.Check(func(v int) bool { return v > 0 }, "positive") })
		}, &model{Ints: []int{1, -2, 3, -4}}, 2},
		{"string slice len", func(c Configurator[model], m *model) { c.StringSlice(&m.Tags).MaxLen(1) }, &model{Tags: []string{"a", "b"}}, 1},
		{"string slice each", func(c Configurator[model], m *model) {
			c.StringSlice(&m.Tags).Each(func(e *str.BaseConfigurator) { e.MaxLen(2) })
		}, &model{Tags: []string{"ab", "abc", "abcd"}}, 2},
		{"string slice regexp", func(c Configurator[model], m *model) {
			c.StringSlice(&m.Tags).Required().Regexp(regexp.MustCompile(`^a`))
		}, &model{Tags: []string{"ab", "ba"}}, 1},
		{"string slice email", func(c Configurator[model], m *model) { c.StringSlice(&m.Tags).Email() }, &model{Tags: []string{"a@b.cc", "ab"}}, 1},
		{"string slice of ptrs", func(c Configurator[model], m *model) {
			c.StringSlice(&m.TagPtrs).Each(func(e *str.BaseConfigurator) { e.Required() })
		}, &model{TagPtrs: []*string{ptr("a"), nil, ptr("")}}, 2},
		{"ptr to string slice nil skipped", func(c Configurator[model], m *model) { c.StringSlice(&m.TagsPtr).MinLen(1).Email() }, &model{}, 0},
		{"ptr to string slice nil required", func(c Configurator[model], m *model) { c.StringSlice(&m.TagsPtr).Required() }, &model{}, 1},
		{"ptr to string slice checked", func(c Configurator[model], m *model) { c.StringSlice(&m.TagsPtr).Email() }, &model{TagsPtr: &[]string{"a"}}, 1},
		{"string slice when", func(c Configurator[model], m *model) {
			c.StringSlice(&m.Tags).When(func(_ context.Context, v shared.Elems[string]) bool { return v.Len() > 1 }).MaxLen(1)
		}, &model{Tags: []string{"a", "b"}}, 1},
		{"uuid slice any of", func(c Configurator[model], m *model) { c.UUIDSlice(&m.UIDs).AnyOf(id1) }, &model{UIDs: []*uuid.UUID{&id1, &id2, nil}}, 1},
		{"uuid slice each", func(c Configurator[model], m *model) {
			c.UUIDSlice(&m.UIDs).Each(func(e *vuuid.BaseConfigurator) { e.Required() })
		}, &model{UIDs: []*uuid.UUID{&id1, nil}}, 1},

		{"field when false", func(c Configurator[model], m *model) {
			c.String(&m.Name).When(func(_ context.Context, v string) bool { return v != "skip" }).MaxLen(1)
		}, &model{Name: "skip"}, 0},
		{"field when true", func(c Configurator[model], m *model) {
			c.String(&m.Name).When(func(_ context.Context, v string) bool { return v != "skip" }).MaxLen(1)
		}, &model{Name: "other"}, 1},
		{"number when", func(c Configurator[model], m *model) {
			c.Number(&m.Age).When(func(_ context.Context, v int) bool { return v > 100 }).Max(10)
		}, &model{Age: 50}, 0},
		{"struct when", func(c Configurator[model], m *model) {
			c.When(func(_ context.Context, m *model) bool { return m.Age > 0 }).String(&m.Name).Required()
		}, &model{}, 0},
		{"struct when chained", func(c Configurator[model], m *model) {
			c.When(func(_ context.Context, m *model) bool { return m.Age > 0 }).
				When(func(_ context.Context, m *model) bool { return m.ID > 0 }).
				String(&m.Name).Required()
		}, &model{Age: 1}, 0},
		{"struct when does not affect parent", func(c Configurator[model], m *model) {
			c.When(func(_ context.Context, m *model) bool { return false })
			c.String(&m.Name).Required()
		}, &model{}, 1},
		{"field custom typed", func(c Configurator[model], m *model) {
			c.Number(&m.Age).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v int) []shared.Error {
				if v%2 != 0 {
					return []shared.Error{h.ErrorT(ctx, v, "should be even")}
				}
				return nil
			})
		}, &model{Age: 3}, 1},
		{"struct custom", func(c Configurator[model], m *model) {
			c.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *model) []shared.Error {
				if obj.Name == obj.Inner.City {
					return []shared.Error{h.ErrorT(ctx, &obj.Inner.City, obj.Inner.City, "should differ")}
				}
				return nil
			})
		}, &model{}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := validateModel(t, tt.cfg, tt.obj)
			assert.Len(t, errs, tt.wantErrs, "errors: %v", errs)
		})
	}
}

func TestErrorContent(t *testing.T) {
	errs := validateModel(t, func(c Configurator[model], m *model) {
		c.String(&m.Inner.City).MaxLen(3)
		c.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *model) []shared.Error {
			return []shared.Error{h.ErrorT(ctx, &obj.Age, obj.Age, "custom %d", 1)}
		})
	}, &model{Inner: address{City: "abcd"}, Age: 7})
	require.Len(t, errs, 2)
	assert.Equal(t, shared.Error{Location: "Inner.City", Message: "Cannot be longer than 3 characters", Value: "abcd"}, errs[0])
	assert.Equal(t, shared.Error{Location: "Age", Message: "custom 1", Value: 7}, errs[1])
}

func TestTrim(t *testing.T) {
	obj := &model{Name: "  a  ", NamePtr: ptr(" b "), Tags: []string{" c ", "d "}, TagPtrs: []*string{ptr(" e"), nil}}
	errs := validateModel(t, func(c Configurator[model], m *model) {
		c.String(&m.Name).Trim().MaxLen(1)
		c.String(&m.NamePtr).Trim().MaxLen(1)
		c.StringSlice(&m.Tags).Trim().Each(func(e *str.BaseConfigurator) { e.MaxLen(1) })
		c.StringSlice(&m.TagPtrs).Trim()
	}, obj)
	assert.Empty(t, errs)
	assert.Equal(t, "a", obj.Name)
	assert.Equal(t, "b", *obj.NamePtr)
	assert.Equal(t, []string{"c", "d"}, obj.Tags)
	assert.Equal(t, "e", *obj.TagPtrs[0])
}

func TestSliceElemsValues(t *testing.T) {
	var got []int
	validateModel(t, func(c Configurator[model], m *model) {
		c.Slice(&m.Ints).Custom(func(_ context.Context, _ *shared.FieldCustomHelper, v shared.Elems[int]) []shared.Error {
			for _, e := range v.All() {
				got = append(got, *e)
			}
			return nil
		})
	}, &model{Ints: []int{1, 2, 3, 4}})
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

func TestValidate(t *testing.T) {
	v := New()
	Configure[model](v, func(c Configurator[model], m *model) {
		c.String(&m.Name).Trim().Required()
	})
	ctx := context.Background()

	assert.Nil(t, v.Validate(ctx, &model{Name: "a"}), "valid object must give nil")
	assert.Len(t, v.Validate(ctx, &model{}), 1)

	byValue := model{Name: "  "}
	assert.Len(t, v.Validate(ctx, byValue), 1, "struct value must be validated")
	assert.Equal(t, "  ", byValue.Name, "struct value must not be modified")

	assert.Nil(t, v.Validate(ctx, (*model)(nil)))
	assert.Nil(t, v.Validate(ctx, nil))
	assert.Nil(t, v.Validate(ctx, &address{}))
}

func TestConfigurePanics(t *testing.T) {
	foreign := ""
	assert.PanicsWithValue(t, "valigo: *string is not a pointer to a field of *valigo.model: field not found", func() {
		Configure[model](New(), func(c Configurator[model], m *model) { c.String(&foreign) })
	})
}

func TestZeroAllocationsOnValid(t *testing.T) {
	v := New()
	Configure[model](v, func(c Configurator[model], m *model) {
		c.String(&m.Name).Trim().Required().MaxLen(10).Email()
		c.String(&m.NamePtr).MaxLen(5)
		c.String(&m.Addr.City).Required()
		c.Number(&m.Age).Min(0).Max(100).AnyOfInterval(0, 100)
		c.NumberPtr(&m.AgePtr).Required().Max(10)
		c.UUID(&m.UID).Required()
		c.StringSlice(&m.Tags).Required().MaxLen(5).Each(func(e *str.BaseConfigurator) { e.Trim().MaxLen(3) })
		c.UUIDSlice(&m.UIDs).Each(func(e *vuuid.BaseConfigurator) { e.Required() })
		c.When(func(_ context.Context, m *model) bool { return m.Age > 0 }).Number(&m.ID).Min(1)
	})
	id := uuid.New()
	obj := &model{
		Name: "a@b.cc", AgePtr: ptr(1), Age: 5, ID: 1, UID: id,
		Addr: &address{City: "x"}, Tags: []string{"a", "b"}, UIDs: []*uuid.UUID{&id},
	}
	ctx := context.Background()
	require.Nil(t, v.Validate(ctx, obj))
	assert.Zero(t, testing.AllocsPerRun(100, func() { v.Validate(ctx, obj) }))
}
