package timex

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo"
	"github.com/insei/valigo/internal/testenv"
	"github.com/insei/valigo/shared"
)

type event struct {
	Start  time.Time
	At     time.Time
	AtPtr  *time.Time
	Dates  []time.Time
	Breaks []*time.Time
}

var (
	day1 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	day2 = day1.AddDate(0, 0, 1)
	day3 = day1.AddDate(0, 0, 2)
)

func validate(t *testing.T, cfg func(c valigo.Configurator[event], m *event), obj *event) []string {
	t.Helper()
	v := valigo.New()
	valigo.Configure[event](v, cfg)
	var msgs []string
	for _, err := range v.ValidateTyped(context.Background(), obj) {
		msgs = append(msgs, err.Location+": "+err.Message)
	}
	return msgs
}

func TestRules(t *testing.T) {
	now = func() time.Time { return day2 }
	t.Cleanup(func() { now = time.Now })

	tests := []struct {
		name string
		cfg  func(c valigo.Configurator[event], m *event)
		obj  *event
		want []string
	}{
		{"required zero", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).Required() }, &event{}, []string{"At: Should be fulfilled"}},
		{"required nil", func(c valigo.Configurator[event], m *event) { c.Field(&m.AtPtr).As(Of).Required() }, &event{}, []string{"AtPtr: Should be fulfilled"}},
		{"required ok", func(c valigo.Configurator[event], m *event) { c.Field(&m.AtPtr).As(Of).Required() }, &event{AtPtr: &day1}, nil},
		{"before", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).Before(day2) }, &event{At: day2}, []string{"At: Should be before 2026-01-02T00:00:00Z"}},
		{"before ok", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).Before(day2) }, &event{At: day1}, nil},
		{"after", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).After(day2) }, &event{At: day2}, []string{"At: Should be after 2026-01-02T00:00:00Z"}},
		{"after ok", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).After(day2) }, &event{At: day3}, nil},
		{"between bounds", func(c valigo.Configurator[event], m *event) {
			c.Field(&m.At).As(Of).Between(day1, day2)
			c.Field(&m.AtPtr).As(Of).Between(day1, day2)
		}, &event{At: day1, AtPtr: &day2}, nil},
		{"between outside", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).Between(day1, day2) }, &event{At: day3},
			[]string{"At: Should be between 2026-01-01T00:00:00Z and 2026-01-02T00:00:00Z"}},
		{"in past", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).InPast() }, &event{At: day3}, []string{"At: Should be in the past"}},
		{"in past ok", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).InPast() }, &event{At: day1}, nil},
		{"in future", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).InFuture() }, &event{At: day1}, []string{"At: Should be in the future"}},
		{"in future ok", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).InFuture() }, &event{At: day3}, nil},
		{"before field", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).BeforeField(&m.Start) }, &event{At: day2, Start: day1}, []string{"At: Should be before Start"}},
		{"before field ok", func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).BeforeField(&m.Start) }, &event{At: day1, Start: day2}, nil},
		{"after field", func(c valigo.Configurator[event], m *event) { c.Field(&m.AtPtr).As(Of).AfterField(&m.Start) }, &event{AtPtr: &day1, Start: day2}, []string{"AtPtr: Should be after Start"}},
		{"after field ok", func(c valigo.Configurator[event], m *event) { c.Field(&m.AtPtr).As(Of).AfterField(&m.Start) }, &event{AtPtr: &day3, Start: day2}, nil},
		{"nil skipped", func(c valigo.Configurator[event], m *event) {
			c.Field(&m.AtPtr).As(Of).Before(day1).After(day3).InPast().InFuture()
		}, &event{}, nil},
		{"when", func(c valigo.Configurator[event], m *event) {
			c.Field(&m.At).As(Of).When(func(_ context.Context, v time.Time) bool { return !v.IsZero() }).After(day2)
		}, &event{}, nil},
		{"custom", func(c valigo.Configurator[event], m *event) {
			c.Field(&m.At).As(Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v time.Time) []shared.Error {
				if v.Weekday() == time.Thursday {
					return []shared.Error{h.ErrorT(ctx, v, "Not on Thursdays")}
				}
				return nil
			})
		}, &event{At: day1}, []string{"At: Not on Thursdays"}},
		{"slice", func(c valigo.Configurator[event], m *event) {
			c.Slice(&m.Dates).As(Slice).MaxLen(2).Each(func(e *BaseConfigurator) { e.After(day1) })
		}, &event{Dates: []time.Time{day2, day1, day3}}, []string{
			"Dates: Cannot contain more than 2 elements",
			"Dates[1]: Should be after 2026-01-01T00:00:00Z",
		}},
		{"slice of pointers", func(c valigo.Configurator[event], m *event) {
			c.Slice(&m.Breaks).As(Slice).Each(func(e *BaseConfigurator) { e.Required() })
		}, &event{Breaks: []*time.Time{&day1, nil}}, []string{"Breaks[1]: Should be fulfilled"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, validate(t, tt.cfg, tt.obj))
		})
	}
}

func TestErrorValue(t *testing.T) {
	v := valigo.New()
	valigo.Configure[event](v, func(c valigo.Configurator[event], m *event) { c.Field(&m.At).As(Of).After(day2) })
	errs := v.ValidateTyped(context.Background(), &event{At: day1})
	assert.Equal(t, []shared.Error{{Location: "At", Message: "Should be after 2026-01-02T00:00:00Z", Value: day1, Code: "time.after", Params: map[string]any{"time": "2026-01-02T00:00:00Z"}}}, errs)
}

func TestZeroAllocationsOnValid(t *testing.T) {
	if testenv.RaceEnabled {
		t.Skip("allocations are not stable under the race detector")
	}
	v := valigo.New()
	valigo.Configure[event](v, func(c valigo.Configurator[event], m *event) {
		c.Field(&m.At).As(Of).Required().Between(day1, day3).InPast()
		c.Field(&m.AtPtr).As(Of).Before(day3)
		c.Slice(&m.Dates).As(Slice).Each(func(e *BaseConfigurator) { e.After(day1) })
	})
	obj := &event{At: day2, AtPtr: &day1, Dates: []time.Time{day2}}
	ctx := context.Background()
	assert.Nil(t, v.Validate(ctx, obj))
	assert.Zero(t, testing.AllocsPerRun(100, func() { v.Validate(ctx, obj) }))
}

func TestCodes(t *testing.T) {
	now = func() time.Time { return day2 }
	t.Cleanup(func() { now = time.Now })
	tests := []struct {
		code string
		cfg  func(c *BaseConfigurator, m *event)
		obj  *event
	}{
		{"time.required", func(c *BaseConfigurator, m *event) { c.Required() }, &event{}},
		{"time.before", func(c *BaseConfigurator, m *event) { c.Before(day1) }, &event{At: day2}},
		{"time.after", func(c *BaseConfigurator, m *event) { c.After(day3) }, &event{At: day2}},
		{"time.between", func(c *BaseConfigurator, m *event) { c.Between(day3, day3) }, &event{At: day2}},
		{"time.before_field", func(c *BaseConfigurator, m *event) { c.BeforeField(&m.Start) }, &event{At: day2, Start: day1}},
		{"time.after_field", func(c *BaseConfigurator, m *event) { c.AfterField(&m.Start) }, &event{At: day1, Start: day2}},
		{"time.in_past", func(c *BaseConfigurator, m *event) { c.InPast() }, &event{At: day3}},
		{"time.in_future", func(c *BaseConfigurator, m *event) { c.InFuture() }, &event{At: day1}},
	}
	for _, tt := range tests {
		v := valigo.New()
		valigo.Configure[event](v, func(c valigo.Configurator[event], m *event) { tt.cfg(c.Field(&m.At).As(Of), m) })
		errs := v.ValidateTyped(context.Background(), tt.obj)
		if assert.Len(t, errs, 1, tt.code) {
			assert.Equal(t, tt.code, errs[0].Code)
		}
	}
}
