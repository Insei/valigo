package valigo

import (
	"context"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo/rules/eq"
	"github.com/insei/valigo/rules/num"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/shared"
)

type constraintsAddress struct {
	City string
}

type constraintsModel struct {
	Name    string
	Code    string
	Email   *string
	Age     int
	Rate    *float64
	Role    string
	Kind    string
	Day     string
	Tags    []string
	Roles   []string
	Labels  map[string]int
	Note    string
	Address *constraintsAddress
}

type described struct {
	Field   string
	Path    []string
	Keyword string
	Value   any
}

func describedOf(constraints []shared.Constraint) []described {
	result := make([]described, len(constraints))
	for i, c := range constraints {
		result[i] = described{Field: c.Field.GetStructPath(), Path: c.Path, Keyword: c.Keyword, Value: c.Value}
	}
	return result
}

func TestConstraints(t *testing.T) {
	v := New()
	Configure[constraintsModel](v, func(b Configurator[constraintsModel], m *constraintsModel) {
		b.Field(&m.Name).As(str.Of).Trim().Required().MaxRunes(64)
		b.Field(&m.Code).As(str.Of).MinRunes(2).MaxLen(8).Regexp(regexp.MustCompile(`^[A-Z]+$`)).Prefix("A.")
		b.Field(&m.Email).As(str.Of).Email()
		b.Field(&m.Age).As(num.Of).Required().AnyOfInterval(18, 120)
		b.Field(&m.Rate).As(num.Ptr).Min(0).Max(1)
		b.Field(&m.Role).As(eq.Of).OneOf("admin", "user")
		b.Field(&m.Kind).As(eq.Of).Equal("person")
		b.Field(&m.Day).As(str.Of).Datetime(time.DateOnly).Datetime(time.Kitchen)
		b.Slice(&m.Tags).As(str.Slice).Required().MaxLen(5).Each(func(tag *str.BaseConfigurator) { tag.MaxRunes(16) })
		b.Slice(&m.Roles).As(eq.Slice).MinLen(1).Unique()
		b.Map(&m.Labels).Required().MaxLen(3).
			EachKey(func(k *shared.Field[string]) { k.As(str.Of).Hostname() }).
			EachValue(func(v *shared.Field[int]) { v.As(num.Of).Min(1) })
		b.Field(&m.Note).As(str.Of).When(func(context.Context, string) bool { return true }).MaxLen(10)
		b.Field(&m.Address.City).As(str.Of).Required()
		b.When(func(context.Context, *constraintsModel) bool { return true }).Field(&m.Note).As(str.Of).MaxLen(20)
	})
	items, props, names := []string{"items"}, []string{"additionalProperties"}, []string{"propertyNames"}
	assert.Equal(t, []described{
		{"Name", nil, "required", true},
		{"Code", nil, "minLength", 2},
		{"Code", nil, "maxLength", 8},
		{"Code", nil, "pattern", "^[A-Z]+$"},
		{"Code", nil, "pattern", `^A\.`},
		{"Email", nil, "format", "email"},
		{"Age", nil, "required", true},
		{"Age", nil, "minimum", 18},
		{"Age", nil, "maximum", 120},
		{"Rate", nil, "minimum", 0.0},
		{"Rate", nil, "maximum", 1.0},
		{"Role", nil, "enum", []any{"admin", "user"}},
		{"Kind", nil, "const", "person"},
		{"Day", nil, "format", "date"},
		{"Tags", nil, "required", true},
		{"Tags", nil, "minItems", 1},
		{"Tags", nil, "maxItems", 5},
		{"Tags", items, "maxLength", 16},
		{"Roles", nil, "minItems", 1},
		{"Roles", nil, "uniqueItems", true},
		{"Labels", nil, "required", true},
		{"Labels", nil, "minProperties", 1},
		{"Labels", nil, "maxProperties", 3},
		{"Labels", names, "format", "hostname"},
		{"Labels", props, "minimum", 1},
		{"Address.City", nil, "required", true},
		{"Address.City", nil, "minLength", 1},
	}, describedOf(v.Constraints()))
	for _, c := range v.Constraints() {
		assert.Equal(t, reflect.TypeFor[constraintsModel](), c.Struct)
	}

	Configure[constraintsModel](v, func(b Configurator[constraintsModel], m *constraintsModel) {
		b.Field(&m.Email).As(str.Of).Email()
	})
	assert.Len(t, v.Constraints(), 1, "Configure replaces the constraints")
}

func TestConstraintsOrder(t *testing.T) {
	v := New()
	Configure[constraintsModel](v, func(b Configurator[constraintsModel], m *constraintsModel) {
		b.Field(&m.Name).As(str.Of).Required()
	})
	Configure[constraintsAddress](v, func(b Configurator[constraintsAddress], a *constraintsAddress) {
		b.Field(&a.City).As(str.Of).Required()
	})
	assert.Equal(t, "City", v.Constraints()[0].Field.GetStructPath(), "the structs are sorted by name")
	assert.Equal(t, "Name", v.Constraints()[2].Field.GetStructPath())
}
