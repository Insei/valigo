package shared

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestElemsGetterLayouts(t *testing.T) {
	a, b := "a", "b"
	obj := &sample{
		Texts:     []string{a, b},
		TextsPtr:  &[]string{a, b},
		TextPtrs:  []*string{&a, nil, &b},
		TextsPtrs: &[]*string{&a, nil},
	}
	tests := []struct {
		path string
		want []*string
	}{
		{"Texts", []*string{&obj.Texts[0], &obj.Texts[1]}},
		{"TextsPtr", []*string{&(*obj.TextsPtr)[0], &(*obj.TextsPtr)[1]}},
		{"TextPtrs", []*string{&a, nil, &b}},
		{"TextsPtrs", []*string{&a, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			elems, ok := elemsGetter[string](sampleField(t, tt.path))(obj)
			assert.True(t, ok)
			assert.Equal(t, len(tt.want), elems.Len())
			for i, want := range tt.want {
				assert.Same(t, want, elems.At(i), "element %d", i)
			}
		})
	}
}

func TestElemsGetterAbsent(t *testing.T) {
	_, ok := elemsGetter[string](sampleField(t, "TextsPtr"))(&sample{})
	assert.False(t, ok, "nil pointer to slice is absent")

	elems, ok := elemsGetter[string](sampleField(t, "Texts"))(&sample{})
	assert.True(t, ok, "nil slice is present and empty")
	assert.Zero(t, elems.Len())
}

func TestElemsGetterPointerElements(t *testing.T) {
	one := 1
	elems, ok := elemsGetter[*int](sampleField(t, "NumPtrs"))(&sample{NumPtrs: []*int{&one, nil}})
	assert.True(t, ok)
	assert.Same(t, &one, *elems.At(0))
	assert.Nil(t, *elems.At(1))
}

func TestElemsGetterPanics(t *testing.T) {
	assert.PanicsWithValue(t, "valigo: field Text has type string, slice expected", func() {
		elemsGetter[string](sampleField(t, "Text"))
	})
	assert.PanicsWithValue(t, "valigo: field NumPtrs has type []*int, string expected", func() {
		elemsGetter[string](sampleField(t, "NumPtrs"))
	})
}

func TestElemsAll(t *testing.T) {
	elems, _ := elemsGetter[string](sampleField(t, "Texts"))(&sample{Texts: []string{"a", "b", "c"}})
	var got []string
	for i, v := range elems.All() {
		got = append(got, *v)
		if i == 1 {
			break
		}
	}
	assert.Equal(t, []string{"a", "b"}, got)
}

func TestSliceFieldConfigurator(t *testing.T) {
	alwaysFalse := func(string) bool { return false }
	tests := []struct {
		name  string
		path  string
		setup func(s *SliceFieldConfigurator[string])
		obj   *sample
		want  []Error
	}{
		{
			name:  "max len",
			path:  "Texts",
			setup: func(s *SliceFieldConfigurator[string]) { s.MaxLen(1) },
			obj:   &sample{Texts: []string{"a", "b"}},
			want:  []Error{{Location: "Texts", Message: "validation:slice:Cannot contain more than 1 elements", Value: 2}},
		},
		{
			name:  "min len",
			path:  "Texts",
			setup: func(s *SliceFieldConfigurator[string]) { s.MinLen(3) },
			obj:   &sample{Texts: []string{"a", "b"}},
			want:  []Error{{Location: "Texts", Message: "validation:slice:Cannot contain less than 3 elements", Value: 2}},
		},
		{
			name:  "length rules skip absent",
			path:  "TextsPtr",
			setup: func(s *SliceFieldConfigurator[string]) { s.MinLen(3).MaxLen(0) },
			obj:   &sample{},
		},
		{
			name:  "required absent",
			path:  "TextsPtr",
			setup: func(s *SliceFieldConfigurator[string]) { s.Required() },
			obj:   &sample{},
			want:  []Error{{Location: "TextsPtr", Message: "validation:slice:Should be fulfilled"}},
		},
		{
			name:  "required present",
			path:  "Texts",
			setup: func(s *SliceFieldConfigurator[string]) { s.Required() },
			obj:   &sample{Texts: []string{"a"}},
		},
		{
			name:  "required empty",
			path:  "Texts",
			setup: func(s *SliceFieldConfigurator[string]) { s.Required() },
			obj:   &sample{Texts: []string{}},
			want:  []Error{{Location: "Texts", Message: "validation:slice:Should be fulfilled"}},
		},
		{
			name: "custom",
			path: "Texts",
			setup: func(s *SliceFieldConfigurator[string]) {
				s.Custom(func(ctx context.Context, h *FieldCustomHelper, value Elems[string]) []Error {
					return []Error{h.ErrorT(ctx, value.Len(), "custom")}
				})
			},
			obj:  &sample{Texts: []string{"a"}},
			want: []Error{{Location: "Texts", Message: "custom", Value: 1}},
		},
		{
			name: "custom skips absent",
			path: "TextsPtr",
			setup: func(s *SliceFieldConfigurator[string]) {
				s.Custom(func(context.Context, *FieldCustomHelper, Elems[string]) []Error { panic("must not be called") })
			},
			obj: &sample{},
		},
		{
			name: "when false",
			path: "Texts",
			setup: func(s *SliceFieldConfigurator[string]) {
				s.When(func(_ context.Context, v Elems[string]) bool { return v.Len() > 5 }).MaxLen(0)
			},
			obj: &sample{Texts: []string{"a"}},
		},
		{
			name: "when true",
			path: "Texts",
			setup: func(s *SliceFieldConfigurator[string]) {
				s.When(func(_ context.Context, v Elems[string]) bool { return v.Len() > 0 }).MaxLen(0)
			},
			obj:  &sample{Texts: []string{"a"}},
			want: []Error{{Location: "Texts", Message: "validation:slice:Cannot contain more than 0 elements", Value: 1}},
		},
		{
			name:  "nil when",
			path:  "Texts",
			setup: func(s *SliceFieldConfigurator[string]) { s.When(nil).MaxLen(0) },
			obj:   &sample{Texts: []string{"a"}},
			want:  []Error{{Location: "Texts", Message: "validation:slice:Cannot contain more than 0 elements", Value: 1}},
		},
		{
			name: "each with index",
			path: "TextPtrs",
			setup: func(s *SliceFieldConfigurator[string]) {
				s.Each(func(e *Field[string]) {
					e.Required(func(v string) bool { return v != "" }, "required")
					e.Check(alwaysFalse, "bad")
				})
			},
			obj: &sample{TextPtrs: []*string{ptr("a"), nil}},
			want: []Error{
				{Location: "TextPtrs[0]", Message: "bad", Value: "a"},
				{Location: "TextPtrs[1]", Message: "required"},
			},
		},
		{
			name: "each skips absent slice",
			path: "TextsPtrs",
			setup: func(s *SliceFieldConfigurator[string]) {
				s.Each(func(e *Field[string]) { e.Required(func(string) bool { return true }, "required") })
			},
			obj: &sample{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rules
			tt.setup(NewSliceFieldConfigurator[string](sampleField(t, tt.path), pathHelper{}, r.add))
			assert.Equal(t, tt.want, r.validate(tt.obj))
		})
	}
}

func TestEachRegistersOneValidationFn(t *testing.T) {
	var r rules
	NewSliceFieldConfigurator[string](sampleField(t, "Texts"), pathHelper{}, r.add).Each(func(e *Field[string]) {
		e.Check(func(string) bool { return true }, "a")
		e.Check(func(string) bool { return true }, "b")
	})
	assert.Len(t, r, 1, "all element rules must share one read of the slice")
}
