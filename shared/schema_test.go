package shared

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// schema collects the constraints of the configurators.
type schema []Constraint

func (s *schema) describe(c Constraint) {
	*s = append(*s, c)
}

func (s schema) keywords() []string {
	result := make([]string, len(s))
	for i, c := range s {
		result[i] = c.Field.GetStructPath() + " " + c.Keyword
		for _, p := range c.Path {
			result[i] += " " + p
		}
	}
	return result
}

func TestSchema(t *testing.T) {
	var r rules
	var s schema
	f := NewValueField[string](sampleField(t, "Text"), pathHelper{}, sampleFields, r.add, s.describe)
	f.Schema("maxLength", 3).As(Retype[string, string]).Schema("pattern", "^a")
	f.When(func(context.Context, string) bool { return true }).Schema("minLength", 1)
	f.Modify(func(*string) {}).Schema("maxLength", 2).Schema(KeywordRequired, true)

	NewValueField[*string](sampleField(t, "TextPtr"), pathHelper{}, sampleFields, r.add, s.describe).As(Deref).Schema("format", "email")
	NewValueField[*[]string](sampleField(t, "TextsPtr"), pathHelper{}, sampleFields, r.add, s.describe).As(SlicePtr).MaxLen(2)

	texts := NewSliceFieldConfigurator[string](sampleField(t, "Texts"), pathHelper{}, sampleFields, r.add, s.describe)
	texts.Schema("uniqueItems", true).Each(func(e *Field[string]) {
		e.Modify(func(*string) {}).Schema("maxLength", 1)
	}).EachField().Schema("minLength", 1)
	texts.When(func(context.Context, Elems[string]) bool { return true }).MaxLen(1).EachField().Schema("maxLength", 1)

	assert.Equal(t, []string{
		"Text maxLength", "Text pattern", "Text required",
		"TextPtr format",
		"TextsPtr maxItems",
		"Texts uniqueItems", "Texts minLength items",
	}, s.keywords())

	// without describe there is nothing to register
	NewValueField[string](sampleField(t, "Text"), pathHelper{}, sampleFields, r.add, nil).Modify(func(*string) {}).Schema("maxLength", 1)
	NewSliceFieldConfigurator[string](sampleField(t, "Texts"), pathHelper{}, sampleFields, r.add, nil).EachField().Schema("maxLength", 1)
}

func TestValues(t *testing.T) {
	assert.Equal(t, []any{1, 2}, Values([]int{1, 2}))
	assert.Equal(t, []any{}, Values([]string{}))
}
