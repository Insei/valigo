package num_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
)

type score uint16

type numbers struct {
	Min     int
	Max     int
	Int     int
	Int8    int8
	Int64   int64
	Uint    uint
	Float32 float32
	Float64 float64
	Score   score
	IntPtr  *int
	Sizes   []int
	Float   *float64
}

func ptr[T any](v T) *T {
	return &v
}

func validate(t *testing.T, cfg func(c valigo.Configurator[numbers], m *numbers), obj *numbers) []shared.Error {
	t.Helper()
	v := valigo.New()
	valigo.Configure[numbers](v, cfg)
	return v.ValidateTyped(context.Background(), obj)
}

func TestRules(t *testing.T) {
	tests := []struct {
		name string
		cfg  func(c valigo.Configurator[numbers], m *numbers)
		obj  *numbers
		want []string
	}{
		{"max ok", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).Max(5) }, &numbers{Int: 5}, nil},
		{"max fail", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).Max(5) }, &numbers{Int: 6}, []string{"Cannot be greater than 5"}},
		{"min ok", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int8).As(num.Of).Min(-5) }, &numbers{Int8: -5}, nil},
		{"min fail", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int8).As(num.Of).Min(-5) }, &numbers{Int8: -6}, []string{"Cannot be less than -5"}},
		{"int64", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int64).As(num.Of).Max(1 << 40) }, &numbers{Int64: 1<<40 + 1}, []string{"Cannot be greater than 1099511627776"}},
		{"uint", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Uint).As(num.Of).Min(1) }, &numbers{}, []string{"Cannot be less than 1"}},
		{"float32", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Float32).As(num.Of).Max(0.5) }, &numbers{Float32: 0.75}, []string{"Cannot be greater than 0.5"}},
		{"float64", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Float64).As(num.Of).Min(0.5) }, &numbers{Float64: 0.25}, []string{"Cannot be less than 0.5"}},
		{"named type", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Score).As(num.Of).Max(100) }, &numbers{Score: 101}, []string{"Cannot be greater than 100"}},
		{"any of ok", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).AnyOf(1, 2) }, &numbers{Int: 2}, nil},
		{"any of fail", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).AnyOf(1, 2) }, &numbers{Int: 3}, []string{"Only [1 2] values is allowed"}},
		{"interval lower bound", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).AnyOfInterval(1, 3) }, &numbers{Int: 1}, nil},
		{"interval upper bound", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).AnyOfInterval(1, 3) }, &numbers{Int: 3}, nil},
		{"interval below", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).AnyOfInterval(1, 3) }, &numbers{Int: 0}, []string{"Only interval[1 - 3] is allowed"}},
		{"interval above", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).AnyOfInterval(1, 3) }, &numbers{Int: 4}, []string{"Only interval[1 - 3] is allowed"}},
		{"required zero", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).Required() }, &numbers{}, []string{"Should be fulfilled"}},
		{"required non zero", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).Required() }, &numbers{Int: -1}, nil},
		{"ptr required nil", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.IntPtr).As(num.Ptr).Required() }, &numbers{}, []string{"Should be fulfilled"}},
		{"ptr required zero", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.IntPtr).As(num.Ptr).Required() }, &numbers{IntPtr: ptr(0)}, nil},
		{"ptr nil skips rules", func(c valigo.Configurator[numbers], m *numbers) {
			c.Field(&m.Float).As(num.Ptr).Min(1).Max(0).AnyOf(5).AnyOfInterval(2, 3)
		}, &numbers{}, nil},
		{"ptr checked", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Float).As(num.Ptr).Max(1.5) }, &numbers{Float: ptr(2.5)}, []string{"Cannot be greater than 1.5"}},
		{"min field", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).MinField(&m.Min) }, &numbers{Int: 1, Min: 2}, []string{"Cannot be less than Min"}},
		{"min field ok", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).MinField(&m.Min) }, &numbers{Int: 2, Min: 2}, nil},
		{"max field", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).MaxField(&m.Max) }, &numbers{Int: 3, Max: 2}, []string{"Cannot be greater than Max"}},
		{"max field ok", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.Int).As(num.Of).MaxField(&m.Max) }, &numbers{Int: 2, Max: 2}, nil},
		{"ptr field with min field", func(c valigo.Configurator[numbers], m *numbers) { c.Field(&m.IntPtr).As(num.Ptr).MinField(&m.Min) }, &numbers{Min: 5}, nil},
		{"slice", func(c valigo.Configurator[numbers], m *numbers) {
			c.Slice(&m.Sizes).As(num.Slice).MaxLen(2).Each(func(e *num.BaseConfigurator[int]) { e.Min(1).MaxField(&m.Max) })
		}, &numbers{Sizes: []int{1, 0, 9}, Max: 5}, []string{"Cannot contain more than 2 elements", "Cannot be less than 1", "Cannot be greater than Max"}},
		{"when true", func(c valigo.Configurator[numbers], m *numbers) {
			c.Field(&m.Int).As(num.Of).When(func(_ context.Context, v int) bool { return v > 10 }).Max(20)
		}, &numbers{Int: 21}, []string{"Cannot be greater than 20"}},
		{"when false", func(c valigo.Configurator[numbers], m *numbers) {
			c.Field(&m.Int).As(num.Of).When(func(_ context.Context, v int) bool { return v > 100 }).Max(20)
		}, &numbers{Int: 21}, nil},
		{"ptr when keeps optional", func(c valigo.Configurator[numbers], m *numbers) {
			c.Field(&m.IntPtr).As(num.Ptr).When(func(context.Context, int) bool { return true }).Required()
		}, &numbers{IntPtr: ptr(0)}, nil},
		{"custom", func(c valigo.Configurator[numbers], m *numbers) {
			c.Field(&m.Int).As(num.Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v int) []shared.Error {
				if v%2 != 0 {
					return []shared.Error{h.ErrorT(ctx, v, "Should be even")}
				}
				return nil
			})
		}, &numbers{Int: 3}, []string{"Should be even"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, err := range validate(t, tt.cfg, tt.obj) {
				got = append(got, err.Message)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestErrorValueAndLocation(t *testing.T) {
	errs := validate(t, func(c valigo.Configurator[numbers], m *numbers) {
		c.Field(&m.Score).As(num.Of).Max(10)
		c.Field(&m.IntPtr).As(num.Ptr).Required()
	}, &numbers{Score: 11})
	assert.Equal(t, []shared.Error{
		{Location: "Score", Message: "Cannot be greater than 10", Value: score(11), Code: "number.max", Params: map[string]any{"max": score(10)}},
		{Location: "IntPtr", Message: "Should be fulfilled", Code: "number.required"},
	}, errs)
}

func TestCodes(t *testing.T) {
	v := valigo.New()
	valigo.Configure[numbers](v, func(c valigo.Configurator[numbers], m *numbers) {
		c.Field(&m.Int).As(num.Of).Required().Min(5).Max(-1).MinField(&m.Min).MaxField(&m.Max).AnyOf(7).AnyOfInterval(1, 2)
	})
	var codes []string
	for _, err := range v.ValidateTyped(context.Background(), &numbers{Min: 1, Max: -1}) {
		codes = append(codes, err.Code)
	}
	assert.Equal(t, []string{"number.required", "number.min", "number.max", "number.min_field", "number.max_field", "number.any_of", "number.interval"}, codes)
}
