package huma_test

import (
	"context"
	"testing"

	dhuma "github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insei/valigo"
	"github.com/insei/valigo/eq"
	vhuma "github.com/insei/valigo/huma"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
)

type Item struct {
	Title string `json:"title"`
}

type Base struct {
	Version int `json:"version"`
}

type CreateItemInput struct {
	Limit int `query:"limit"`
	Body  struct {
		Base
		Name    string            `json:"name" maxLength:"10"`
		Code    string            `json:"code,omitempty" pattern:"^[A-Z]+$"`
		Kind    string            `json:"kind,omitempty"`
		Tags    []string          `json:"tags"`
		Labels  map[string]int    `json:"labels,omitempty"`
		Flags   map[string]string `json:"flags,omitempty"`
		Item    *Item             `json:"item,omitempty"`
		Secret  string            `json:"-"`
		Hidden  string            `json:"hidden,omitempty" hidden:"true"`
		Comment string
	}
}

func newOpenAPI(t *testing.T) map[string]any {
	t.Helper()
	v := valigo.New()
	valigo.Configure[CreateItemInput](v, func(b valigo.Configurator[CreateItemInput], i *CreateItemInput) {
		b.Field(&i.Limit).As(num.Of).Max(100)
		b.Field(&i.Body.Version).As(num.Of).Min(1)
		b.Field(&i.Body.Name).As(str.Of).Required().MaxRunes(10).MaxLen(40)
		b.Field(&i.Body.Code).As(str.Of).Required().Prefix("A")
		b.Field(&i.Body.Kind).As(eq.Of).OneOf("a", "b")
		b.Slice(&i.Body.Tags).As(str.Slice).MaxLen(3).Each(func(tag *str.BaseConfigurator) { tag.Required().MaxRunes(16) })
		b.Map(&i.Body.Labels).EachValue(func(v *shared.Field[int]) { v.As(num.Of).Min(1) })
		b.Map(&i.Body.Flags).EachKey(func(k *shared.Field[string]) { k.As(str.Of).Hostname() })
		b.Field(&i.Body.Secret).As(str.Of).Required()
		b.Field(&i.Body.Hidden).As(str.Of).MaxRunes(1)
		b.Field(&i.Body.Comment).As(str.Of).MaxRunes(200)
	})
	valigo.Configure[Item](v, func(b valigo.Configurator[Item], i *Item) {
		b.Field(&i.Title).As(str.Of).MinRunes(2)
	})
	_, api := humatest.New(t)
	dhuma.Post(api, "/items", func(context.Context, *CreateItemInput) (*struct{}, error) { return nil, nil })
	doc, err := vhuma.OpenAPI(api, v)
	require.NoError(t, err)
	return doc
}

func schemaOf(t *testing.T, doc map[string]any, name string) map[string]any {
	t.Helper()
	s, ok := doc["components"].(map[string]any)["schemas"].(map[string]any)[name].(map[string]any)
	require.True(t, ok, name)
	return s
}

func TestOpenAPI(t *testing.T) {
	doc := newOpenAPI(t)
	body := schemaOf(t, doc, "CreateItemInputBody")
	props := body["properties"].(map[string]any)

	assert.Equal(t, map[string]any{"type": "integer", "format": "int64", "minimum": 1.0}, props["version"], "the embedded struct is flattened")
	assert.Equal(t, map[string]any{"type": "string", "minLength": 1.0, "maxLength": 10.0, "allOf": []any{map[string]any{"maxLength": 40.0}}}, props["name"],
		"the same value is kept once, another one is in allOf")
	assert.Equal(t, map[string]any{"type": "string", "minLength": 1.0, "pattern": "^[A-Z]+$", "allOf": []any{map[string]any{"pattern": "^A"}}}, props["code"])
	assert.Equal(t, []any{"a", "b"}, props["kind"].(map[string]any)["enum"])
	tags := props["tags"].(map[string]any)
	assert.Equal(t, 3.0, tags["maxItems"])
	assert.Equal(t, map[string]any{"type": "string", "minLength": 1.0, "maxLength": 16.0}, tags["items"], "elements aren't required")
	assert.Equal(t, 1.0, props["labels"].(map[string]any)["additionalProperties"].(map[string]any)["minimum"])
	assert.Equal(t, map[string]any{"format": "hostname"}, props["flags"].(map[string]any)["propertyNames"])
	assert.Equal(t, 200.0, props["Comment"].(map[string]any)["maxLength"])
	assert.NotContains(t, props, "Secret")
	assert.NotContains(t, props, "hidden", "a hidden property isn't added")

	required := body["required"].([]any)
	assert.Contains(t, required, "code", "Required adds the omitempty field to the required properties")
	assert.Equal(t, 1, countOf(required, "name"), "a required property isn't repeated")
	assert.NotContains(t, required, "Secret")

	assert.Equal(t, 2.0, schemaOf(t, doc, "Item")["properties"].(map[string]any)["title"].(map[string]any)["minLength"], "the rules of a nested struct type")

	params := doc["paths"].(map[string]any)["/items"].(map[string]any)["post"].(map[string]any)["parameters"].([]any)
	assert.NotContains(t, params[0].(map[string]any)["schema"], "maximum", "the parameters keep the schemas of huma")
}

func TestOpenAPIDoesntChangeHuma(t *testing.T) {
	v := valigo.New()
	valigo.Configure[Item](v, func(b valigo.Configurator[Item], i *Item) { b.Field(&i.Title).As(str.Of).MaxRunes(2) })
	_, api := humatest.New(t)
	dhuma.Post(api, "/items", func(context.Context, *struct{ Body Item }) (*struct{}, error) { return nil, nil })
	_, err := vhuma.OpenAPI(api, v)
	require.NoError(t, err)
	assert.Nil(t, api.OpenAPI().Components.Schemas.Map()["Item"].Properties["title"].MaxLength)
}

func TestOpenAPIError(t *testing.T) {
	_, api := humatest.New(t)
	api.OpenAPI().Extensions = map[string]any{"x-bad": func() {}}
	_, err := vhuma.OpenAPI(api, valigo.New())
	assert.Error(t, err)
}

func countOf(values []any, v any) int {
	n := 0
	for _, x := range values {
		if x == v {
			n++
		}
	}
	return n
}
