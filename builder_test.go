package valigo

import (
	"context"
	"regexp"
	"sync"
	"testing"

	"github.com/insei/fmap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insei/valigo/internal/testenv"
	"github.com/insei/valigo/rules/num"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/shared"
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
	Tags    []string
	TagsPtr *[]string
	TagPtrs []*string
	Ints    []int
	IntPtrs []*int
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
	tests := []struct {
		name     string
		cfg      func(c Configurator[model], m *model)
		obj      *model
		wantErrs int
	}{
		{"string max len ok", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).MaxLen(3) }, &model{Name: "abc"}, 0},
		{"string max len fail", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).MaxLen(3) }, &model{Name: "abcd"}, 1},
		{"string min len fail", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).MinLen(3) }, &model{Name: "ab"}, 1},
		{"string required fail", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).Required() }, &model{}, 1},
		{"string any of", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).AnyOf("a", "b") }, &model{Name: "c"}, 1},
		{"string regexp", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).Regexp(regexp.MustCompile(`^\d+$`)) }, &model{Name: "12a"}, 1},
		{"string email ok", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).Email() }, &model{Name: "a@b.cc"}, 0},
		{"string email fail", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).Email() }, &model{Name: "a.b.cc"}, 1},
		{"string chain errors", func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).MinLen(5).Email() }, &model{Name: "ab"}, 2},
		{"named string type", func(c Configurator[model], m *model) { c.Field(&m.Status).As(str.Of).AnyOf("on") }, &model{Status: "off"}, 1},
		{"nil *string skipped", func(c Configurator[model], m *model) { c.Field(&m.NamePtr).As(str.Of).MaxLen(1).MinLen(5) }, &model{}, 0},
		{"nil *string required", func(c Configurator[model], m *model) { c.Field(&m.NamePtr).As(str.Of).Required() }, &model{}, 1},
		{"*string checked", func(c Configurator[model], m *model) { c.Field(&m.NamePtr).As(str.Of).MaxLen(1) }, &model{NamePtr: ptr("ab")}, 1},
		{"nil nested struct skipped", func(c Configurator[model], m *model) { c.Field(&m.Addr.City).As(str.Of).MaxLen(1) }, &model{}, 0},
		{"nil nested struct required", func(c Configurator[model], m *model) { c.Field(&m.Addr.City).As(str.Of).Required() }, &model{}, 1},
		{"nested struct checked", func(c Configurator[model], m *model) { c.Field(&m.Addr.City).As(str.Of).MaxLen(1) }, &model{Addr: &address{City: "ab"}}, 1},
		{"inner struct checked", func(c Configurator[model], m *model) { c.Field(&m.Inner.City).As(str.Of).MaxLen(1) }, &model{Inner: address{City: "ab"}}, 1},

		{"number max inclusive", func(c Configurator[model], m *model) { c.Field(&m.Age).As(num.Of).Max(10) }, &model{Age: 10}, 0},
		{"number max fail", func(c Configurator[model], m *model) { c.Field(&m.Age).As(num.Of).Max(10) }, &model{Age: 11}, 1},
		{"number min inclusive", func(c Configurator[model], m *model) { c.Field(&m.Age).As(num.Of).Min(10) }, &model{Age: 10}, 0},
		{"number min fail", func(c Configurator[model], m *model) { c.Field(&m.Age).As(num.Of).Min(10) }, &model{Age: 9}, 1},
		{"number int64 untyped const", func(c Configurator[model], m *model) { c.Field(&m.ID).As(num.Of).Max(5) }, &model{ID: 6}, 1},
		{"number float32", func(c Configurator[model], m *model) { c.Field(&m.Score).As(num.Of).Max(1.5) }, &model{Score: 1.6}, 1},
		{"number any of", func(c Configurator[model], m *model) { c.Field(&m.Age).As(num.Of).AnyOf(1, 2) }, &model{Age: 3}, 1},
		{"number interval inclusive", func(c Configurator[model], m *model) { c.Field(&m.Age).As(num.Of).AnyOfInterval(1, 10) }, &model{Age: 10}, 0},
		{"number interval fail", func(c Configurator[model], m *model) { c.Field(&m.Age).As(num.Of).AnyOfInterval(1, 10) }, &model{Age: 11}, 1},
		{"number required zero", func(c Configurator[model], m *model) { c.Field(&m.Age).As(num.Of).Required() }, &model{}, 1},
		{"number ptr nil skipped", func(c Configurator[model], m *model) { c.Field(&m.AgePtr).As(num.Ptr).Max(10).Min(20) }, &model{}, 0},
		{"number ptr nil required", func(c Configurator[model], m *model) { c.Field(&m.AgePtr).As(num.Ptr).Required() }, &model{}, 1},
		{"number ptr zero required", func(c Configurator[model], m *model) { c.Field(&m.AgePtr).As(num.Ptr).Required() }, &model{AgePtr: ptr(0)}, 0},
		{"number ptr checked", func(c Configurator[model], m *model) { c.Field(&m.AgePtr).As(num.Ptr).Max(10) }, &model{AgePtr: ptr(11)}, 1},

		{"slice of ptrs", func(c Configurator[model], m *model) {
			c.Slice(&m.IntPtrs).MaxLen(1).Each(func(e *shared.Field[*int]) { e.Check(func(v *int) bool { return v != nil }, "not nil") })
		}, &model{IntPtrs: []*int{ptr(1), nil}}, 2},
		{"slice each as num ptr", func(c Configurator[model], m *model) {
			c.Slice(&m.IntPtrs).Each(func(e *shared.Field[*int]) { e.As(num.Ptr).Required().Max(5) })
		}, &model{IntPtrs: []*int{ptr(1), nil, ptr(6)}}, 2},
		{"slice each as str", func(c Configurator[model], m *model) {
			c.Slice(&m.TagPtrs).Each(func(e *shared.Field[*string]) { e.As(str.Of).MinLen(2) })
		}, &model{TagPtrs: []*string{ptr("a"), nil, ptr("ab")}}, 1},
		{"base field rules", func(c Configurator[model], m *model) {
			c.Field(&m.Age).Check(func(v int) bool { return v%2 == 0 }, "Should be even").Check(func(v int) bool { return v > 0 }, "Should be positive")
		}, &model{Age: -1}, 2},
		{"slice max len", func(c Configurator[model], m *model) { c.Slice(&m.Ints).MaxLen(2) }, &model{Ints: []int{1, 2, 3}}, 1},
		{"slice min len", func(c Configurator[model], m *model) { c.Slice(&m.Ints).MinLen(2) }, &model{Ints: []int{1}}, 1},
		{"slice required", func(c Configurator[model], m *model) { c.Slice(&m.Ints).Required() }, &model{Ints: []int{}}, 1},
		{"slice each", func(c Configurator[model], m *model) {
			c.Slice(&m.Ints).Each(func(e *shared.Field[int]) { e.Check(func(v int) bool { return v > 0 }, "positive") })
		}, &model{Ints: []int{1, -2, 3, -4}}, 2},
		{"string slice len", func(c Configurator[model], m *model) { c.Slice(&m.Tags).As(str.Slice).MaxLen(1) }, &model{Tags: []string{"a", "b"}}, 1},
		{"string slice each", func(c Configurator[model], m *model) {
			c.Slice(&m.Tags).As(str.Slice).Each(func(e *str.BaseConfigurator) { e.MaxLen(2) })
		}, &model{Tags: []string{"ab", "abc", "abcd"}}, 2},
		{"string slice regexp", func(c Configurator[model], m *model) {
			c.Slice(&m.Tags).As(str.Slice).Required().Regexp(regexp.MustCompile(`^a`))
		}, &model{Tags: []string{"ab", "ba"}}, 1},
		{"string slice email", func(c Configurator[model], m *model) { c.Slice(&m.Tags).As(str.Slice).Email() }, &model{Tags: []string{"a@b.cc", "ab"}}, 1},
		{"string slice of ptrs", func(c Configurator[model], m *model) {
			c.Slice(&m.TagPtrs).As(str.Slice).Each(func(e *str.BaseConfigurator) { e.Required() })
		}, &model{TagPtrs: []*string{ptr("a"), nil, ptr("")}}, 2},
		{"ptr to string slice nil skipped", func(c Configurator[model], m *model) {
			c.Field(&m.TagsPtr).As(shared.SlicePtr).As(str.Slice).MinLen(1).Email()
		}, &model{}, 0},
		{"ptr to string slice nil required", func(c Configurator[model], m *model) {
			c.Field(&m.TagsPtr).As(shared.SlicePtr).As(str.Slice).Required()
		}, &model{}, 1},
		{"ptr to string slice checked", func(c Configurator[model], m *model) { c.Field(&m.TagsPtr).As(shared.SlicePtr).As(str.Slice).Email() }, &model{TagsPtr: &[]string{"a"}}, 1},
		{"string slice when", func(c Configurator[model], m *model) {
			c.Slice(&m.Tags).As(str.Slice).When(func(_ context.Context, v shared.Elems[string]) bool { return v.Len() > 1 }).MaxLen(1)
		}, &model{Tags: []string{"a", "b"}}, 1},

		{"field when false", func(c Configurator[model], m *model) {
			c.Field(&m.Name).As(str.Of).When(func(_ context.Context, v string) bool { return v != "skip" }).MaxLen(1)
		}, &model{Name: "skip"}, 0},
		{"field when true", func(c Configurator[model], m *model) {
			c.Field(&m.Name).As(str.Of).When(func(_ context.Context, v string) bool { return v != "skip" }).MaxLen(1)
		}, &model{Name: "other"}, 1},
		{"field when absent skips required", func(c Configurator[model], m *model) {
			c.Field(&m.NamePtr).As(str.Of).When(func(_ context.Context, v string) bool { return true }).Required()
		}, &model{}, 0},
		{"number ptr when absent", func(c Configurator[model], m *model) {
			c.Field(&m.AgePtr).As(num.Ptr).When(func(_ context.Context, v int) bool { return v == 0 }).Max(-1)
		}, &model{}, 0},
		{"number when", func(c Configurator[model], m *model) {
			c.Field(&m.Age).As(num.Of).When(func(_ context.Context, v int) bool { return v > 100 }).Max(10)
		}, &model{Age: 50}, 0},
		{"struct when", func(c Configurator[model], m *model) {
			c.When(func(_ context.Context, m *model) bool { return m.Age > 0 }).Field(&m.Name).As(str.Of).Required()
		}, &model{}, 0},
		{"struct when chained", func(c Configurator[model], m *model) {
			c.When(func(_ context.Context, m *model) bool { return m.Age > 0 }).
				When(func(_ context.Context, m *model) bool { return m.ID > 0 }).
				Field(&m.Name).As(str.Of).Required()
		}, &model{Age: 1}, 0},
		{"struct when does not affect parent", func(c Configurator[model], m *model) {
			c.When(func(_ context.Context, m *model) bool { return false })
			c.Field(&m.Name).As(str.Of).Required()
		}, &model{}, 1},
		{"field custom typed", func(c Configurator[model], m *model) {
			c.Field(&m.Age).As(num.Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v int) []shared.Error {
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
		c.Field(&m.Inner.City).As(str.Of).MaxLen(3)
		c.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *model) []shared.Error {
			return []shared.Error{h.ErrorT(ctx, &obj.Age, obj.Age, "custom %d", 1)}
		})
	}, &model{Inner: address{City: "abcd"}, Age: 7})
	require.Len(t, errs, 2)
	assert.Equal(t, shared.Error{Location: "Inner.City", Message: "Cannot be longer than 3 characters", Value: "abcd", Code: "string.max_len", Params: map[string]any{"max": 3}}, errs[0])
	assert.Equal(t, shared.Error{Location: "Age", Message: "custom 1", Value: 7, Code: "custom %d"}, errs[1])
}

func TestElementErrorLocation(t *testing.T) {
	errs := validateModel(t, func(c Configurator[model], m *model) {
		c.Slice(&m.Tags).As(str.Slice).MaxLen(2).Email()
	}, &model{Tags: []string{"a@b.cc", "x", "y"}})
	require.Len(t, errs, 3)
	assert.Equal(t, "Tags", errs[0].Location)
	assert.Equal(t, "Tags[1]", errs[1].Location)
	assert.Equal(t, "Tags[2]", errs[2].Location)
}

func TestOtherFieldLocationByNamingFn(t *testing.T) {
	type signUp struct {
		Password string `json:"password"`
		Repeat   string `json:"repeat"`
	}
	v := New(WithFieldLocationNamingFn(func(f fmap.Field) string { return f.GetTagPath("json", false) }))
	Configure[signUp](v, func(c Configurator[signUp], m *signUp) { c.Field(&m.Repeat).As(str.Of).EqualField(&m.Password) })
	errs := v.ValidateTyped(context.Background(), &signUp{Password: "a", Repeat: "b"})
	assert.Equal(t, []shared.Error{{Location: "repeat", Message: "Should be equal to password", Value: "b", Code: "string.equal_field", Params: map[string]any{"field": "password"}}}, errs)
}

func TestMap(t *testing.T) {
	type service struct {
		Labels map[string]string
		Ports  map[string]int
	}
	v := New()
	Configure[service](v, func(c Configurator[service], m *service) {
		c.Map(&m.Labels).Required().MaxLen(3).
			EachKey(func(k *shared.Field[string]) { k.As(str.Of).Alphanumeric() }).
			EachValue(func(v *shared.Field[string]) { v.As(str.Of).MaxLen(5) })
		c.Map(&m.Ports).EachValue(func(p *shared.Field[int]) { p.As(num.Of).AnyOfInterval(1, 65535) })
	})
	errs := v.ValidateTyped(context.Background(), &service{
		Labels: map[string]string{"app": "valigo-api", "tier!": "web"},
		Ports:  map[string]int{"http": 80, "admin": 0},
	})
	assert.Equal(t, []shared.Error{
		{Location: "Labels[tier!]", Message: "Should contain only letters and digits", Value: "tier!", Code: "string.alphanumeric"},
		{Location: "Labels[app]", Message: "Cannot be longer than 5 characters", Value: "valigo-api", Code: "string.max_len", Params: map[string]any{"max": 5}},
		{Location: "Ports[admin]", Message: "Only interval[1 - 65535] is allowed", Value: 0, Code: "number.interval", Params: map[string]any{"min": 1, "max": 65535}},
	}, errs)

	valid := &service{Labels: map[string]string{"app": "api"}, Ports: map[string]int{"http": 80}}
	assert.Nil(t, v.Validate(context.Background(), valid))
	if !testenv.RaceEnabled {
		assert.Zero(t, testing.AllocsPerRun(100, func() { v.Validate(context.Background(), valid) }))
	}
}

func TestReconfigureReplacesRules(t *testing.T) {
	v := New()
	for range 3 {
		Configure[model](v, func(c Configurator[model], m *model) { c.Field(&m.Name).As(str.Of).Required() })
	}
	assert.Len(t, v.ValidateTyped(context.Background(), &model{}), 1)
}

func TestTrim(t *testing.T) {
	obj := &model{Name: "  a  ", NamePtr: ptr(" b "), Tags: []string{" c ", "d "}, TagPtrs: []*string{ptr(" e"), nil}}
	errs := validateModel(t, func(c Configurator[model], m *model) {
		c.Field(&m.Name).As(str.Of).Trim().MaxLen(1)
		c.Field(&m.NamePtr).As(str.Of).Trim().MaxLen(1)
		c.Slice(&m.Tags).As(str.Slice).Trim().Each(func(e *str.BaseConfigurator) { e.MaxLen(1) })
		c.Slice(&m.TagPtrs).As(str.Slice).Trim()
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
		c.Field(&m.Name).As(str.Of).Trim().Required()
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
	assert.Nil(t, v.Validate(ctx, address{}), "value of a type without rules")
}

func TestConfigurePanics(t *testing.T) {
	foreign := ""
	assert.PanicsWithValue(t, "valigo: *string is not a pointer to a field of *valigo.model: field not found", func() {
		Configure[model](New(), func(c Configurator[model], m *model) { c.Field(&foreign).As(str.Of) })
	})
}

func TestConfigureWithUnexportedStructField(t *testing.T) {
	type withMutex struct {
		mu   sync.Mutex
		Name string
	}
	v := New()
	Configure[withMutex](v, func(c Configurator[withMutex], m *withMutex) { c.Field(&m.Name).As(str.Of).Required() })
	assert.Len(t, v.ValidateTyped(context.Background(), &withMutex{}), 1)
}

func TestConfigureRecursiveType(t *testing.T) {
	type node struct {
		Name string
		Next *node
	}
	v := New()
	Configure[node](v, func(c Configurator[node], m *node) { c.Field(&m.Name).As(str.Of).Required() })
	assert.Len(t, v.ValidateTyped(context.Background(), &node{Next: &node{}}), 1)
}

func TestNilWhen(t *testing.T) {
	errs := validateModel(t, func(c Configurator[model], m *model) {
		c.When(nil).Field(&m.Name).As(str.Of).Required()
	}, &model{})
	assert.Len(t, errs, 1)
}

func TestStructCustomForeignPointerPanics(t *testing.T) {
	foreign := 0
	v := New()
	Configure[model](v, func(c Configurator[model], m *model) {
		c.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *model) []shared.Error {
			return []shared.Error{h.ErrorT(ctx, &foreign, foreign, "never")}
		})
	})
	assert.PanicsWithValue(t, "valigo: *int is not a pointer to a field of *valigo.model: field not found", func() {
		v.Validate(context.Background(), &model{})
	})
}

func TestConfigureNonStructPanics(t *testing.T) {
	assert.PanicsWithError(t, "not supported type: *int, only struct and ptr to struct is supported", func() {
		Configure[int](New(), func(Configurator[int], *int) {})
	})
}

func TestZeroAllocationsOnValid(t *testing.T) {
	if testenv.RaceEnabled {
		t.Skip("allocations are not stable under the race detector")
	}
	v := New()
	Configure[model](v, func(c Configurator[model], m *model) {
		c.Field(&m.Name).As(str.Of).Trim().Required().MaxLen(10).Email()
		c.Field(&m.NamePtr).As(str.Of).MaxLen(5)
		c.Field(&m.Addr.City).As(str.Of).Required()
		c.Field(&m.Age).As(num.Of).Min(0).Max(100).AnyOfInterval(0, 100)
		c.Field(&m.AgePtr).As(num.Ptr).Required().Max(10)
		c.Slice(&m.Tags).As(str.Slice).Required().MaxLen(5).Each(func(e *str.BaseConfigurator) { e.Trim().MaxLen(3) })
		c.When(func(_ context.Context, m *model) bool { return m.Age > 0 }).Field(&m.ID).As(num.Of).Min(1)
		c.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *model) []shared.Error {
			if obj.Age < 0 {
				return []shared.Error{h.ErrorT(ctx, &obj.Age, obj.Age, "negative")}
			}
			return nil
		})
	})
	obj := &model{
		Name: "a@b.cc", AgePtr: ptr(1), Age: 5, ID: 1,
		Addr: &address{City: "x"}, Tags: []string{"a", "b"},
	}
	ctx := context.Background()
	require.Nil(t, v.Validate(ctx, obj))
	assert.Zero(t, testing.AllocsPerRun(100, func() { v.Validate(ctx, obj) }))
}

func TestConfiguratorTags(t *testing.T) {
	type tagged struct {
		Name string `validate:"required"`
	}
	v := New()
	Configure[tagged](v, func(b Configurator[tagged], _ *tagged) { b.Tags("validate") })
	errs := v.ValidateTyped(context.Background(), &tagged{})
	if len(errs) != 1 || errs[0].Code != "string.required" {
		t.Errorf("errors of the tags: %v", errs)
	}
}
