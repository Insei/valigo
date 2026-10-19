package eq_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo"
	"github.com/insei/valigo/internal/testenv"
	"github.com/insei/valigo/rules/eq"
	"github.com/insei/valigo/shared"
)

type status string

const (
	active  status = "active"
	blocked status = "blocked"
	deleted status = "deleted"
)

type account struct {
	Previous status
	Status   status
	Agreed   bool
	Level    int
	Role     *status
	Roles    []status
	Scores   []int
	Statuses []*status
}

func ptr[T any](v T) *T {
	return &v
}

func validate(t *testing.T, cfg func(c valigo.Configurator[account], m *account), obj *account) []string {
	t.Helper()
	v := valigo.New()
	valigo.Configure[account](v, cfg)
	var msgs []string
	for _, err := range v.ValidateTyped(context.Background(), obj) {
		msgs = append(msgs, err.Location+": "+err.Message)
	}
	return msgs
}

func TestRules(t *testing.T) {
	many := make([]int, 40)
	for i := range many {
		many[i] = i
	}
	manyDup := append(append([]int{}, many...), 7)

	tests := []struct {
		name string
		cfg  func(c valigo.Configurator[account], m *account)
		obj  *account
		want []string
	}{
		{"required zero", func(c valigo.Configurator[account], m *account) { c.Field(&m.Status).As(eq.Of).Required() }, &account{}, []string{"Status: Should be fulfilled"}},
		{"required ok", func(c valigo.Configurator[account], m *account) { c.Field(&m.Level).As(eq.Of).Required() }, &account{Level: 1}, nil},
		{"ptr required nil", func(c valigo.Configurator[account], m *account) { c.Field(&m.Role).As(eq.Ptr).Required() }, &account{}, []string{"Role: Should be fulfilled"}},
		{"ptr required zero allowed", func(c valigo.Configurator[account], m *account) { c.Field(&m.Role).As(eq.Ptr).Required() }, &account{Role: ptr(status(""))}, nil},
		{"ptr nil skipped", func(c valigo.Configurator[account], m *account) {
			c.Field(&m.Role).As(eq.Ptr).OneOf(active).Equal(blocked)
		}, &account{}, nil},
		{"ptr checked", func(c valigo.Configurator[account], m *account) { c.Field(&m.Role).As(eq.Ptr).OneOf(active) }, &account{Role: ptr(deleted)}, []string{"Role: Only [active] values is allowed"}},
		{"one of", func(c valigo.Configurator[account], m *account) { c.Field(&m.Status).As(eq.Of).OneOf(active, blocked) }, &account{Status: deleted}, []string{"Status: Only [active blocked] values is allowed"}},
		{"one of ok", func(c valigo.Configurator[account], m *account) { c.Field(&m.Status).As(eq.Of).OneOf(active, blocked) }, &account{Status: blocked}, nil},
		{"not one of", func(c valigo.Configurator[account], m *account) { c.Field(&m.Status).As(eq.Of).NotOneOf(deleted) }, &account{Status: deleted}, []string{"Status: Values [deleted] are not allowed"}},
		{"equal", func(c valigo.Configurator[account], m *account) { c.Field(&m.Agreed).As(eq.Of).Equal(true) }, &account{}, []string{"Agreed: Should be true"}},
		{"equal ok", func(c valigo.Configurator[account], m *account) { c.Field(&m.Agreed).As(eq.Of).Equal(true) }, &account{Agreed: true}, nil},
		{"not equal", func(c valigo.Configurator[account], m *account) { c.Field(&m.Level).As(eq.Of).NotEqual(13) }, &account{Level: 13}, []string{"Level: Cannot be 13"}},
		{"equal field", func(c valigo.Configurator[account], m *account) { c.Field(&m.Status).As(eq.Of).EqualField(&m.Previous) }, &account{Status: active, Previous: blocked}, []string{"Status: Should be equal to Previous"}},
		{"equal field ok", func(c valigo.Configurator[account], m *account) { c.Field(&m.Status).As(eq.Of).EqualField(&m.Previous) }, &account{Status: active, Previous: active}, nil},
		{"when", func(c valigo.Configurator[account], m *account) {
			c.Field(&m.Status).As(eq.Of).When(func(_ context.Context, v status) bool { return v != "" }).OneOf(active)
		}, &account{}, nil},
		{"ptr when keeps optional", func(c valigo.Configurator[account], m *account) {
			c.Field(&m.Role).As(eq.Ptr).When(func(context.Context, status) bool { return true }).Required()
		}, &account{Role: ptr(status(""))}, nil},
		{"custom", func(c valigo.Configurator[account], m *account) {
			c.Field(&m.Level).As(eq.Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v int) []shared.Error {
				return []shared.Error{h.ErrorT(ctx, v, "custom")}
			})
		}, &account{}, []string{"Level: custom"}},
		{"slice one of", func(c valigo.Configurator[account], m *account) {
			c.Slice(&m.Roles).As(eq.Slice).OneOf(active, blocked)
		}, &account{Roles: []status{active, deleted}}, []string{"Roles[1]: Only [active blocked] values is allowed"}},
		{"slice each", func(c valigo.Configurator[account], m *account) {
			c.Slice(&m.Scores).As(eq.Slice).Required().Each(func(e *eq.BaseConfigurator[int]) { e.NotEqual(0) })
		}, &account{Scores: []int{1, 0}}, []string{"Scores[1]: Cannot be 0"}},
		{"unique", func(c valigo.Configurator[account], m *account) { c.Slice(&m.Roles).As(eq.Slice).Unique() }, &account{Roles: []status{active, blocked, active}}, []string{"Roles: Should contain unique values"}},
		{"unique ok", func(c valigo.Configurator[account], m *account) { c.Slice(&m.Roles).As(eq.Slice).Unique() }, &account{Roles: []status{active, blocked}}, nil},
		{"unique big ok", func(c valigo.Configurator[account], m *account) { c.Slice(&m.Scores).As(eq.Slice).Unique() }, &account{Scores: many}, nil},
		{"unique big", func(c valigo.Configurator[account], m *account) { c.Slice(&m.Scores).As(eq.Slice).Unique() }, &account{Scores: manyDup}, []string{"Scores: Should contain unique values"}},
		{"unique of absent", func(c valigo.Configurator[account], m *account) { c.Slice(&m.Roles).As(eq.Slice).Unique() }, &account{}, nil},
		{"slice of pointers each", func(c valigo.Configurator[account], m *account) {
			c.Slice(&m.Statuses).Each(func(e *shared.Field[*status]) { e.As(eq.Ptr).Required().OneOf(active) })
		}, &account{Statuses: []*status{ptr(active), nil, ptr(blocked)}}, []string{
			"Statuses[1]: Should be fulfilled",
			"Statuses[2]: Only [active] values is allowed",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, validate(t, tt.cfg, tt.obj))
		})
	}
}

func TestUniqueErrorValue(t *testing.T) {
	v := valigo.New()
	valigo.Configure[account](v, func(c valigo.Configurator[account], m *account) { c.Slice(&m.Roles).As(eq.Slice).Unique() })
	errs := v.ValidateTyped(context.Background(), &account{Roles: []status{active, blocked, blocked}})
	assert.Equal(t, []shared.Error{{Location: "Roles", Message: "Should contain unique values", Value: blocked, Code: "value.unique"}}, errs)
}

func TestZeroAllocationsOnValid(t *testing.T) {
	if testenv.RaceEnabled {
		t.Skip("allocations are not stable under the race detector")
	}
	v := valigo.New()
	valigo.Configure[account](v, func(c valigo.Configurator[account], m *account) {
		c.Field(&m.Status).As(eq.Of).Required().OneOf(active, blocked).NotOneOf(deleted)
		c.Field(&m.Agreed).As(eq.Of).Equal(true)
		c.Field(&m.Role).As(eq.Ptr).Required()
		c.Slice(&m.Roles).As(eq.Slice).Unique().OneOf(active, blocked)
	})
	obj := &account{Status: active, Agreed: true, Role: ptr(active), Roles: []status{active, blocked}}
	ctx := context.Background()
	assert.Nil(t, v.Validate(ctx, obj))
	assert.Zero(t, testing.AllocsPerRun(100, func() { v.Validate(ctx, obj) }))
}

func TestCodes(t *testing.T) {
	v := valigo.New()
	valigo.Configure[account](v, func(c valigo.Configurator[account], m *account) {
		c.Field(&m.Status).As(eq.Of).Required().OneOf(active).NotOneOf("").Equal(blocked).NotEqual("").EqualField(&m.Previous)
		c.Slice(&m.Roles).As(eq.Slice).Unique()
	})
	var codes []string
	for _, err := range v.ValidateTyped(context.Background(), &account{Previous: active, Roles: []status{active, active}}) {
		codes = append(codes, err.Code)
	}
	assert.Equal(t, []string{"value.required", "value.one_of", "value.not_one_of", "value.equal", "value.not_equal", "value.equal_field", "value.unique"}, codes)
}
