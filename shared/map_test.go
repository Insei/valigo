package shared

import (
	"context"
	"testing"

	"github.com/insei/fmap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mapsSample struct {
	Labels map[string]string
	Counts map[string]int
	Inner  *struct {
		Tags map[string]bool
	}
	Limit int
}

var mapsModel = &mapsSample{Inner: &struct{ Tags map[string]bool }{}}

func mapsFields(fieldPtr any) fmap.Field {
	fields, err := fmap.GetFrom(mapsModel)
	if err != nil {
		panic(err)
	}
	field, err := fields.GetFieldByPtr(mapsModel, fieldPtr)
	if err != nil {
		panic(err)
	}
	return field
}

func newLabels(r *rules) *MapFieldConfigurator[string, string] {
	return NewMapFieldConfigurator[string, string](mapsFields(&mapsModel.Labels), pathHelper{}, mapsFields, r.add, nil)
}

func newCounts(r *rules) *MapFieldConfigurator[string, int] {
	return NewMapFieldConfigurator[string, int](mapsFields(&mapsModel.Counts), pathHelper{}, mapsFields, r.add, nil)
}

func TestMapFieldConfigurator(t *testing.T) {
	notEmptyValue := func(v string) bool { return v != "" }
	tests := []struct {
		name  string
		setup func(r *rules)
		obj   *mapsSample
		want  []Error
	}{
		{"required nil", func(r *rules) { newLabels(r).Required() }, &mapsSample{}, []Error{{Location: "Labels", Message: "validation:map:Should be fulfilled"}}},
		{"required empty", func(r *rules) { newLabels(r).Required() }, &mapsSample{Labels: map[string]string{}}, []Error{{Location: "Labels", Message: "validation:map:Should be fulfilled"}}},
		{"required ok", func(r *rules) { newLabels(r).Required() }, &mapsSample{Labels: map[string]string{"a": "b"}}, nil},
		{"max len", func(r *rules) { newLabels(r).MaxLen(1) }, &mapsSample{Labels: map[string]string{"a": "1", "b": "2"}},
			[]Error{{Location: "Labels", Message: "validation:map:Cannot contain more than 1 elements", Value: 2}}},
		{"min len", func(r *rules) { newLabels(r).MinLen(2) }, &mapsSample{Labels: map[string]string{"a": "1"}},
			[]Error{{Location: "Labels", Message: "validation:map:Cannot contain less than 2 elements", Value: 1}}},
		{"each key", func(r *rules) {
			newLabels(r).EachKey(func(k *Field[string]) { k.Check(func(v string) bool { return len(v) == 1 }, "short") })
		}, &mapsSample{Labels: map[string]string{"bb": "1", "a": "2", "cc": "3"}}, []Error{
			{Location: "Labels[bb]", Message: "short", Value: "bb"},
			{Location: "Labels[cc]", Message: "short", Value: "cc"},
		}},
		{"each value", func(r *rules) {
			newLabels(r).EachValue(func(v *Field[string]) { v.Required(notEmptyValue, "required") })
		}, &mapsSample{Labels: map[string]string{"a": "", "b": "x"}}, []Error{{Location: "Labels[a]", Message: "required", Value: ""}}},
		{"each value with other field", func(r *rules) {
			newCounts(r).EachValue(func(v *Field[int]) {
				v.CheckWith(&mapsModel.Limit, func(v, l int) bool { return v <= l }, "Cannot be greater than %s")
			})
		}, &mapsSample{Counts: map[string]int{"a": 1, "b": 5}, Limit: 3}, []Error{{Location: "Counts[b]", Message: "Cannot be greater than Limit", Value: 5}}},
		{"modify of values doesn't change the map", func(r *rules) {
			newLabels(r).EachValue(func(v *Field[string]) {
				v.Modify(func(s *string) { *s = "changed" })
				v.Check(func(s string) bool { return s != "changed" }, "sees the modified copy")
			})
		}, &mapsSample{Labels: map[string]string{"a": "x"}}, []Error{{Location: "Labels[a]", Message: "sees the modified copy", Value: "changed"}}},
		{"rules of absent map are skipped", func(r *rules) {
			NewMapFieldConfigurator[string, bool](mapsFields(&mapsModel.Inner.Tags), pathHelper{}, mapsFields, r.add, nil).
				MaxLen(0).EachKey(func(k *Field[string]) { k.Check(func(string) bool { return false }, "bad") })
		}, &mapsSample{}, nil},
		{"custom", func(r *rules) {
			newCounts(r).Custom(func(ctx context.Context, h *FieldCustomHelper, m map[string]int) []Error {
				sum := 0
				for _, v := range m {
					sum += v
				}
				if sum > 10 {
					return []Error{h.ErrorT(ctx, sum, "Sum is too big")}
				}
				return nil
			})
		}, &mapsSample{Counts: map[string]int{"a": 6, "b": 6}}, []Error{{Location: "Counts", Message: "Sum is too big", Value: 12}}},
		{"custom skips absent", func(r *rules) {
			NewMapFieldConfigurator[string, bool](mapsFields(&mapsModel.Inner.Tags), pathHelper{}, mapsFields, r.add, nil).
				Custom(func(context.Context, *FieldCustomHelper, map[string]bool) []Error { panic("must not be called") })
		}, &mapsSample{}, nil},
		{"when false", func(r *rules) {
			newLabels(r).When(func(_ context.Context, m map[string]string) bool { return m != nil }).Required()
		}, &mapsSample{}, nil},
		{"when true", func(r *rules) {
			newLabels(r).When(func(_ context.Context, m map[string]string) bool { return len(m) > 1 }).MaxLen(1)
		}, &mapsSample{Labels: map[string]string{"a": "1", "b": "2"}}, []Error{{Location: "Labels", Message: "validation:map:Cannot contain more than 1 elements", Value: 2}}},
		{"nil when", func(r *rules) { newLabels(r).When(nil).Required() }, &mapsSample{}, []Error{{Location: "Labels", Message: "validation:map:Should be fulfilled"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rules
			tt.setup(&r)
			assert.Equal(t, tt.want, r.validate(tt.obj))
		})
	}
}

func TestMapAs(t *testing.T) {
	var r rules
	count := newCounts(&r).As(func(m *MapFieldConfigurator[string, int]) int { return 1 })
	require.Equal(t, 1, count)
}

func TestMapSourceRegistersOneValidationFn(t *testing.T) {
	var r rules
	newLabels(&r).EachKey(func(k *Field[string]) {
		k.Check(func(string) bool { return true }, "a")
		k.Check(func(string) bool { return true }, "b")
	})
	assert.Len(t, r, 1)
}
