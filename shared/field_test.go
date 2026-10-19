package shared

import (
	"context"
	"fmt"
	"testing"

	"github.com/insei/fmap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type label string

type inner struct {
	Text    string
	TextPtr *string
}

type sample struct {
	Text      string
	TextPtr   *string
	Label     label
	Num       int
	Inner     *inner
	Texts     []string
	TextsPtr  *[]string
	TextPtrs  []*string
	TextsPtrs *[]*string
	NumPtrs   []*int
}

// pathHelper builds errors with the struct path location and the formatted message.
type pathHelper struct{}

// format is a variable, so go vet doesn't treat Helper.ErrorT as a printf wrapper because of pathHelper.
var format = fmt.Sprintf

func (pathHelper) ErrorT(_ context.Context, field fmap.Field, value any, localeKey string, args ...any) Error {
	return Error{Location: field.GetStructPath(), Message: format(localeKey, args...), Value: value}
}

// rules collects validation functions registered by the configurators.
type rules []ValidationFn

func (r *rules) add(fn ValidationFn) {
	*r = append(*r, fn)
}

func (r *rules) validate(obj any) []Error {
	var errs []Error
	for _, fn := range *r {
		errs = append(errs, fn(context.Background(), pathHelper{}, obj)...)
	}
	return errs
}

func sampleField(t *testing.T, path string) fmap.Field {
	t.Helper()
	fields, err := fmap.Get[sample]()
	require.NoError(t, err)
	return fields.MustFind(path)
}

func ptr[T any](v T) *T {
	return &v
}

func TestValueGetter(t *testing.T) {
	text := valueGetter[string](sampleField(t, "Text"))
	textPtr := valueGetter[string](sampleField(t, "TextPtr"))
	named := valueGetter[string](sampleField(t, "Label"))
	nested := valueGetter[string](sampleField(t, "Inner.Text"))
	nestedPtr := valueGetter[string](sampleField(t, "Inner.TextPtr"))

	obj := &sample{Text: "a", TextPtr: ptr("b"), Label: "c", Inner: &inner{Text: "d"}}
	assert.Equal(t, "a", *text(obj))
	assert.Equal(t, "b", *textPtr(obj))
	assert.Equal(t, "c", *named(obj))
	assert.Equal(t, "d", *nested(obj))

	*text(obj) = "changed"
	assert.Equal(t, "changed", obj.Text, "getter must point into the object")

	empty := &sample{}
	assert.Nil(t, textPtr(empty), "nil pointer field is absent")
	assert.Nil(t, nested(empty), "field of nil nested struct is absent")
	assert.Nil(t, nestedPtr(empty), "pointer field of nil nested struct is absent")
}

func TestValueGetterLayoutMismatch(t *testing.T) {
	assert.PanicsWithValue(t, "valigo: field Num has type int, string expected", func() {
		valueGetter[string](sampleField(t, "Num"))
	})
}

func TestFieldRules(t *testing.T) {
	tests := []struct {
		name  string
		setup func(f *Field[string])
		obj   *sample
		want  []Error
	}{
		{
			name:  "check passes",
			setup: func(f *Field[string]) { f.Check(func(v string) bool { return v == "ok" }, "bad") },
			obj:   &sample{TextPtr: ptr("ok")},
		},
		{
			name:  "check fails",
			setup: func(f *Field[string]) { f.Check(func(v string) bool { return v == "ok" }, "bad %d", 1) },
			obj:   &sample{TextPtr: ptr("no")},
			want:  []Error{{Location: "TextPtr", Message: "bad 1", Value: "no"}},
		},
		{
			name:  "check skips absent",
			setup: func(f *Field[string]) { f.Check(func(string) bool { return false }, "bad") },
			obj:   &sample{},
		},
		{
			name:  "required absent",
			setup: func(f *Field[string]) { f.Required(func(string) bool { return true }, "required") },
			obj:   &sample{},
			want:  []Error{{Location: "TextPtr", Message: "required"}},
		},
		{
			name:  "required present but empty",
			setup: func(f *Field[string]) { f.Required(func(v string) bool { return v != "" }, "required") },
			obj:   &sample{TextPtr: ptr("")},
			want:  []Error{{Location: "TextPtr", Message: "required", Value: ""}},
		},
		{
			name:  "custom gets value",
			setup: customEcho,
			obj:   &sample{TextPtr: ptr("v")},
			want:  []Error{{Location: "TextPtr", Message: "custom", Value: "v"}},
		},
		{
			name:  "custom skips absent",
			setup: customEcho,
			obj:   &sample{},
		},
		{
			name: "when true",
			setup: func(f *Field[string]) {
				f.When(func(_ context.Context, v string) bool { return v == "x" }).Check(func(string) bool { return false }, "bad")
			},
			obj:  &sample{TextPtr: ptr("x")},
			want: []Error{{Location: "TextPtr", Message: "bad", Value: "x"}},
		},
		{
			name: "when false",
			setup: func(f *Field[string]) {
				f.When(func(_ context.Context, v string) bool { return v == "x" }).Check(func(string) bool { return false }, "bad")
			},
			obj: &sample{TextPtr: ptr("y")},
		},
		{
			name: "when skips absent without calling cond",
			setup: func(f *Field[string]) {
				f.When(func(context.Context, string) bool { panic("must not be called") }).Required(func(string) bool { return true }, "required")
			},
			obj: &sample{},
		},
		{
			name: "nil when returns the same field",
			setup: func(f *Field[string]) {
				f.When(nil).Check(func(string) bool { return false }, "bad")
			},
			obj:  &sample{TextPtr: ptr("x")},
			want: []Error{{Location: "TextPtr", Message: "bad", Value: "x"}},
		},
		{
			name: "rules run in order",
			setup: func(f *Field[string]) {
				f.Modify(func(v *string) { *v += "!" })
				f.Check(func(v string) bool { return v == "x" }, "first")
				f.Check(func(v string) bool { return v == "x" }, "second")
			},
			obj: &sample{TextPtr: ptr("x")},
			want: []Error{
				{Location: "TextPtr", Message: "first", Value: "x!"},
				{Location: "TextPtr", Message: "second", Value: "x!"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rules
			tt.setup(NewValueField[string](sampleField(t, "TextPtr"), pathHelper{}, r.add))
			assert.Equal(t, tt.want, r.validate(tt.obj))
		})
	}
}

func customEcho(f *Field[string]) {
	f.Custom(func(ctx context.Context, h *FieldCustomHelper, value string) []Error {
		return []Error{h.ErrorT(ctx, value, "custom")}
	})
}

func TestFieldRegistersOneValidationFn(t *testing.T) {
	var r rules
	f := NewValueField[string](sampleField(t, "Text"), pathHelper{}, r.add)
	f.Check(func(string) bool { return true }, "a")
	f.When(func(context.Context, string) bool { return true }).Check(func(string) bool { return true }, "b")
	f.Modify(func(*string) {})
	assert.Len(t, r, 1, "all rules of a field must share one read of the value")
}

func TestModifySkipsAbsent(t *testing.T) {
	var r rules
	NewValueField[string](sampleField(t, "TextPtr"), pathHelper{}, r.add).Modify(func(*string) { panic("must not be called") })
	assert.Empty(t, r.validate(&sample{}))
}

func TestDataPtr(t *testing.T) {
	v := 1
	assert.Equal(t, &v, (*int)(dataPtr(&v)))
	assert.Nil(t, dataPtr((*int)(nil)))
	assert.Nil(t, dataPtr(nil))
}
