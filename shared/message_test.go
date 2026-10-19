package shared

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestM(t *testing.T) {
	m := M("number.interval", "between %v and %v", P("min", 1), P("max", 5))
	assert.Equal(t, Msg{Code: "number.interval", Key: "between %v and %v", Args: []any{1, 5}, Params: map[string]any{"min": 1, "max": 5}}, m)
	assert.Equal(t, Msg{Code: "string.email", Key: "email"}, M("string.email", "email"))
}

func TestMsgWith(t *testing.T) {
	m := M("time.after_field", "after %s")
	with := m.with("field", "StartsAt")
	assert.Equal(t, Msg{Code: "time.after_field", Key: "after %s", Args: []any{"StartsAt"}, Params: map[string]any{"field": "StartsAt"}}, with)
	assert.Nil(t, m.Params, "with doesn't change the message")

	base := M("x", "%v %v", P("a", 1))
	base.with("field", "F")
	assert.Equal(t, []any{1}, base.Args, "with doesn't change the arguments")
}

func TestMessages(t *testing.T) {
	tests := []struct {
		name  string
		setup func(r *rules)
		obj   *sample
		want  []Error
	}{
		{
			name: "check by a key",
			setup: func(r *rules) {
				NewValueField[string](sampleField(t, "Text"), pathHelper{}, sampleFields, r.add).Check(notEmpty, "required %s", "x")
			},
			obj:  &sample{},
			want: []Error{{Location: "Text", Message: "required x", Value: "", Code: "required %s"}},
		},
		{
			name: "check by a message",
			setup: func(r *rules) {
				NewValueField[string](sampleField(t, "Text"), pathHelper{}, sampleFields, r.add).CheckMsg(notEmpty, M("string.required", "required"))
			},
			obj:  &sample{},
			want: []Error{{Location: "Text", Message: "required", Value: "", Code: "string.required"}},
		},
		{
			name: "required by a message",
			setup: func(r *rules) {
				NewValueField[*string](sampleField(t, "TextPtr"), pathHelper{}, sampleFields, r.add).As(Deref).RequiredMsg(notEmpty, M("string.required", "required"))
			},
			obj:  &sample{},
			want: []Error{{Location: "TextPtr", Message: "required", Code: "string.required"}},
		},
		{
			name: "check with another field",
			setup: func(r *rules) {
				NewValueField[string](sampleField(t, "Text"), pathHelper{}, sampleFields, r.add).
					CheckWithMsg(&sampleModel.Label2, func(v, o string) bool { return v == o }, M("string.equal_field", "equal to %s"))
			},
			obj:  &sample{Text: "a", Label2: "b"},
			want: []Error{{Location: "Text", Message: "equal to Label2", Value: "a", Code: "string.equal_field", Params: map[string]any{"field": "Label2"}}},
		},
		{
			name: "slice length",
			setup: func(r *rules) {
				NewSliceFieldConfigurator[string](sampleField(t, "Texts"), pathHelper{}, sampleFields, r.add).MinLen(2).MaxLen(0)
			},
			obj: &sample{Texts: []string{"a"}},
			want: []Error{
				{Location: "Texts", Message: "validation:slice:Cannot contain less than 2 elements", Value: 1, Code: "slice.min_len", Params: map[string]any{"min": 2}},
				{Location: "Texts", Message: "validation:slice:Cannot contain more than 0 elements", Value: 1, Code: "slice.max_len", Params: map[string]any{"max": 0}},
			},
		},
		{
			name: "slice required",
			setup: func(r *rules) {
				NewSliceFieldConfigurator[string](sampleField(t, "Texts"), pathHelper{}, sampleFields, r.add).Required()
			},
			obj:  &sample{},
			want: []Error{{Location: "Texts", Message: "validation:slice:Should be fulfilled", Code: "slice.required"}},
		},
		{
			name: "custom of an element by a message",
			setup: func(r *rules) {
				NewSliceFieldConfigurator[string](sampleField(t, "Texts"), pathHelper{}, sampleFields, r.add).Each(func(e *Field[string]) {
					e.Custom(func(ctx context.Context, h *FieldCustomHelper, v string) []Error {
						return []Error{h.ErrorMsg(ctx, v, M("tag.reserved", "reserved %v", P("tag", v)))}
					})
				})
			},
			obj:  &sample{Texts: []string{"go"}},
			want: []Error{{Location: "Texts[0]", Message: "reserved go", Value: "go", Code: "tag.reserved", Params: map[string]any{"tag": "go"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rules
			tt.setup(&r)
			assert.Equal(t, tt.want, r.validateFull(tt.obj))
		})
	}
}

func TestMapMessages(t *testing.T) {
	var r rules
	newLabels(&r).Required().MinLen(2).MaxLen(0)
	assert.Equal(t, []Error{
		{Location: "Labels", Message: "validation:map:Cannot contain less than 2 elements", Value: 1, Code: "map.min_len", Params: map[string]any{"min": 2}},
		{Location: "Labels", Message: "validation:map:Cannot contain more than 0 elements", Value: 1, Code: "map.max_len", Params: map[string]any{"max": 0}},
	}, r.validateFull(&mapsSample{Labels: map[string]string{"a": "b"}}))

	var required rules
	newLabels(&required).Required()
	assert.Equal(t, []Error{{Location: "Labels", Message: "validation:map:Should be fulfilled", Code: "map.required"}}, required.validateFull(&mapsSample{}))
}
