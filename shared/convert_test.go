package shared

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func notEmpty(v string) bool { return v != "" }

// typedSlice is a typed slice configurator built like the ones of the str and uuid packages.
type typedSlice struct {
	SliceRules[*typedSlice, string]
}

func newTypedSlice(s *SliceFieldConfigurator[string]) *typedSlice {
	t := &typedSlice{}
	t.SliceRules = NewSliceRules(s, t, newTypedSlice)
	return t
}

func toTypedSlice[E ~string | *string](s *SliceFieldConfigurator[E]) *typedSlice {
	return newTypedSlice(RetypeSlice[E, string](s))
}

func (t *typedSlice) NotEmpty() *typedSlice {
	t.Slice().Each(func(e *Field[string]) { e.Required(notEmpty, "required") })
	return t
}

func TestFieldAs(t *testing.T) {
	var r rules
	f := NewValueField[string](sampleField(t, "Text"), pathHelper{}, sampleFields, r.add)
	got := f.As(func(f *Field[string]) *Field[string] { return f.Check(notEmpty, "required") })
	assert.Same(t, f, got)
	assert.Len(t, r.validate(&sample{}), 1)
}

func TestDeref(t *testing.T) {
	var r rules
	NewValueField[*string](sampleField(t, "TextPtr"), pathHelper{}, sampleFields, r.add).As(Deref).Required(notEmpty, "required")
	assert.Equal(t, []Error{{Location: "TextPtr", Message: "required"}}, r.validate(&sample{}), "nil pointer is absent")
	assert.Equal(t, []Error{{Location: "TextPtr", Message: "required", Value: ""}}, r.validate(&sample{TextPtr: ptr("")}))
	assert.Empty(t, r.validate(&sample{TextPtr: ptr("a")}))
}

func TestDerefOfNilNestedStruct(t *testing.T) {
	var r rules
	NewValueField[*string](sampleField(t, "Inner.TextPtr"), pathHelper{}, sampleFields, r.add).As(Deref).Check(notEmpty, "required")
	assert.Empty(t, r.validate(&sample{}))
}

func TestRetype(t *testing.T) {
	var r rules
	named := NewValueField[label](sampleField(t, "Label"), pathHelper{}, sampleFields, r.add)
	Retype[label, string](named).Check(func(v string) bool { return v == "ok" }, "bad")
	assert.Equal(t, []Error{{Location: "Label", Message: "bad", Value: "no"}}, r.validate(&sample{Label: "no"}))

	text := NewValueField[string](sampleField(t, "Text"), pathHelper{}, sampleFields, r.add)
	assert.Same(t, text, Retype[string, string](text), "the same type isn't wrapped")
	assert.PanicsWithValue(t, "valigo: field Num has type int, string expected", func() {
		Retype[int, string](NewValueField[int](sampleField(t, "Num"), pathHelper{}, sampleFields, r.add))
	})
}

func TestRetypeSlice(t *testing.T) {
	var r rules
	NewSliceFieldConfigurator[*string](sampleField(t, "TextPtrs"), pathHelper{}, sampleFields, r.add).As(toTypedSlice).Required().NotEmpty()
	assert.Equal(t, []Error{
		{Location: "TextPtrs[1]", Message: "required"},
		{Location: "TextPtrs[2]", Message: "required", Value: ""},
	}, r.validate(&sample{TextPtrs: []*string{ptr("a"), nil, ptr("")}}))

	texts := NewSliceFieldConfigurator[string](sampleField(t, "Texts"), pathHelper{}, sampleFields, r.add)
	assert.Same(t, texts, RetypeSlice[string, string](texts), "the same type isn't wrapped")
}

func TestSlicePtr(t *testing.T) {
	var r rules
	NewValueField[*[]string](sampleField(t, "TextsPtr"), pathHelper{}, sampleFields, r.add).As(SlicePtr).As(toTypedSlice).MaxLen(1).NotEmpty()
	assert.Empty(t, r.validate(&sample{}), "nil pointer to slice is absent")
	assert.Equal(t, []Error{
		{Location: "TextsPtr", Message: "validation:slice:Cannot contain more than 1 elements", Value: 2},
		{Location: "TextsPtr[1]", Message: "required", Value: ""},
	}, r.validate(&sample{TextsPtr: &[]string{"a", ""}}))
}

func TestSlicePtrAfterWhenPanics(t *testing.T) {
	var r rules
	f := NewValueField[*[]string](sampleField(t, "TextsPtr"), pathHelper{}, sampleFields, r.add)
	assert.PanicsWithValue(t, "valigo: SlicePtr of field TextsPtr must be called before When and other conversions", func() {
		SlicePtr(f.When(func(context.Context, *[]string) bool { return true }))
	})
}

func TestSliceRules(t *testing.T) {
	tests := []struct {
		name  string
		setup func(s *typedSlice)
		obj   *sample
		want  []Error
	}{
		{"required", func(s *typedSlice) { s.Required() }, &sample{}, []Error{{Location: "Texts", Message: "validation:slice:Should be fulfilled"}}},
		{"min len", func(s *typedSlice) { s.MinLen(2) }, &sample{Texts: []string{"a"}}, []Error{{Location: "Texts", Message: "validation:slice:Cannot contain less than 2 elements", Value: 1}}},
		{"max len", func(s *typedSlice) { s.MaxLen(0) }, &sample{Texts: []string{"a"}}, []Error{{Location: "Texts", Message: "validation:slice:Cannot contain more than 0 elements", Value: 1}}},
		{"custom", func(s *typedSlice) {
			s.Custom(func(ctx context.Context, h *FieldCustomHelper, v Elems[string]) []Error {
				return []Error{h.ErrorT(ctx, v.Len(), "custom")}
			})
		}, &sample{Texts: []string{"a"}}, []Error{{Location: "Texts", Message: "custom", Value: 1}}},
		{"when keeps the typed configurator", func(s *typedSlice) {
			s.When(func(_ context.Context, v Elems[string]) bool { return v.Len() > 1 }).NotEmpty()
		}, &sample{Texts: []string{""}}, nil},
		{"nil when", func(s *typedSlice) { s.When(nil).NotEmpty() }, &sample{Texts: []string{""}}, []Error{{Location: "Texts[0]", Message: "required", Value: ""}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rules
			tt.setup(NewSliceFieldConfigurator[string](sampleField(t, "Texts"), pathHelper{}, sampleFields, r.add).As(toTypedSlice))
			assert.Equal(t, tt.want, r.validate(tt.obj))
		})
	}
}
