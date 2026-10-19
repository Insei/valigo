package str_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo"
	"github.com/insei/valigo/internal/testenv"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
)

type status string

type texts struct {
	Name    string
	NamePtr *string
	Status  status
	Tags    []string
	TagPtrs []*string
	TagsPtr *[]string
}

func ptr[T any](v T) *T {
	return &v
}

func validate(t *testing.T, cfg func(c valigo.Configurator[texts], m *texts), obj *texts) []string {
	t.Helper()
	v := valigo.New()
	valigo.Configure[texts](v, cfg)
	var msgs []string
	for _, err := range v.ValidateTyped(context.Background(), obj) {
		msgs = append(msgs, err.Location+": "+err.Message)
	}
	return msgs
}

func TestRules(t *testing.T) {
	digits := regexp.MustCompile(`^\d+$`)
	tests := []struct {
		name string
		cfg  func(c valigo.Configurator[texts], m *texts)
		obj  *texts
		want []string
	}{
		{"max len", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).MaxLen(2) }, &texts{Name: "abc"}, []string{"Name: Cannot be longer than 2 characters"}},
		{"min len", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).MinLen(4) }, &texts{Name: "abc"}, []string{"Name: Cannot be shorter than 4 characters"}},
		{"len in bytes", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).MaxLen(2) }, &texts{Name: "яя"}, []string{"Name: Cannot be longer than 2 characters"}},
		{"required empty", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Required() }, &texts{}, []string{"Name: Should be fulfilled"}},
		{"required nil", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.NamePtr).As(str.Of).Required() }, &texts{}, []string{"NamePtr: Should be fulfilled"}},
		{"required empty pointer", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.NamePtr).As(str.Of).Required() }, &texts{NamePtr: ptr("")}, []string{"NamePtr: Should be fulfilled"}},
		{"any of", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Status).As(str.Of).AnyOf("on", "off") }, &texts{Status: "idle"}, []string{"Status: Only [on off] values is allowed"}},
		{"regexp", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Regexp(digits) }, &texts{Name: "12a"}, []string{"Name: Doesn't match required regexp pattern"}},
		{"regexp custom key", func(c valigo.Configurator[texts], m *texts) {
			c.Field(&m.Name).As(str.Of).Regexp(digits, str.WithRegexpLocaleKey("Only digits"))
		}, &texts{Name: "12a"}, []string{"Name: Only digits"}},
		{"regexp empty key keeps default", func(c valigo.Configurator[texts], m *texts) {
			c.Field(&m.Name).As(str.Of).Regexp(digits, str.WithRegexpLocaleKey(""))
		}, &texts{Name: "a"}, []string{"Name: Doesn't match required regexp pattern"}},
		{"email ok", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Email() }, &texts{Name: "user@mail.example.com"}, nil},
		{"email fail", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Email() }, &texts{Name: "user@"}, []string{"Name: Should be email address"}},
		{"url", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).URL() }, &texts{Name: "example.com"}, []string{"Name: Should be URL"}},
		{"url ok", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).URL() }, &texts{Name: "https://example.com"}, nil},
		{"ip", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).IP() }, &texts{Name: "1.2.3"}, []string{"Name: Should be IP address"}},
		{"ipv4", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).IPv4() }, &texts{Name: "::1"}, []string{"Name: Should be IPv4 address"}},
		{"ipv6", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).IPv6() }, &texts{Name: "127.0.0.1"}, []string{"Name: Should be IPv6 address"}},
		{"hostname", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Hostname() }, &texts{Name: "a_b"}, []string{"Name: Should be hostname"}},
		{"uuid", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).UUID() }, &texts{Name: "nope"}, []string{"Name: Should be UUID"}},
		{"alphanumeric", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Alphanumeric() }, &texts{Name: "a-b"}, []string{"Name: Should contain only letters and digits"}},
		{"prefix", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Prefix("ab") }, &texts{Name: "ba"}, []string{"Name: Should start with ab"}},
		{"suffix", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Suffix("ab") }, &texts{Name: "ba"}, []string{"Name: Should end with ab"}},
		{"contains", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Contains("x") }, &texts{Name: "ba"}, []string{"Name: Should contain x"}},
		{"substring rules ok", func(c valigo.Configurator[texts], m *texts) {
			c.Field(&m.Name).As(str.Of).Prefix("a").Suffix("c").Contains("b")
		}, &texts{Name: "abc"}, nil},
		{"max runes", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).MaxRunes(2) }, &texts{Name: "яяя"}, []string{"Name: Cannot be longer than 2 characters"}},
		{"max runes counts runes", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).MaxRunes(2) }, &texts{Name: "яя"}, nil},
		{"min runes", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).MinRunes(3) }, &texts{Name: "яя"}, []string{"Name: Cannot be shorter than 3 characters"}},
		{"nil pointer skipped", func(c valigo.Configurator[texts], m *texts) {
			c.Field(&m.NamePtr).As(str.Of).MinLen(1).MaxLen(0).AnyOf("x").Regexp(digits).Email()
		}, &texts{}, nil},
		{"trim before checks", func(c valigo.Configurator[texts], m *texts) { c.Field(&m.Name).As(str.Of).Trim().Required() }, &texts{Name: "  "}, []string{"Name: Should be fulfilled"}},
		{"when", func(c valigo.Configurator[texts], m *texts) {
			c.Field(&m.Name).As(str.Of).When(func(_ context.Context, v string) bool { return v != "" }).Email()
		}, &texts{}, nil},
		{"custom", func(c valigo.Configurator[texts], m *texts) {
			c.Field(&m.Name).As(str.Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v string) []shared.Error {
				if v == "admin" {
					return []shared.Error{h.ErrorT(ctx, v, "Is reserved")}
				}
				return nil
			})
		}, &texts{Name: "admin"}, []string{"Name: Is reserved"}},

		{"slice length", func(c valigo.Configurator[texts], m *texts) { c.Slice(&m.Tags).As(str.Slice).MinLen(3).MaxLen(1) }, &texts{Tags: []string{"a", "b"}}, []string{
			"Tags: Cannot contain less than 3 elements", "Tags: Cannot contain more than 1 elements",
		}},
		{"slice required", func(c valigo.Configurator[texts], m *texts) {
			c.Field(&m.TagsPtr).As(shared.SlicePtr).As(str.Slice).Required()
		}, &texts{}, []string{"TagsPtr: Should be fulfilled"}},
		{"slice each", func(c valigo.Configurator[texts], m *texts) {
			c.Slice(&m.TagPtrs).As(str.Slice).Each(func(e *str.BaseConfigurator) { e.Required().MaxLen(1) })
		}, &texts{TagPtrs: []*string{ptr("ab"), nil, ptr("c")}}, []string{
			"TagPtrs[0]: Cannot be longer than 1 characters", "TagPtrs[1]: Should be fulfilled",
		}},
		{"slice regexp", func(c valigo.Configurator[texts], m *texts) { c.Slice(&m.Tags).As(str.Slice).Regexp(digits) }, &texts{Tags: []string{"1", "a"}}, []string{
			"Tags[1]: Doesn't match required regexp pattern",
		}},
		{"slice email", func(c valigo.Configurator[texts], m *texts) { c.Slice(&m.Tags).As(str.Slice).Email() }, &texts{Tags: []string{"a@b.cc", "a"}}, []string{
			"Tags[1]: Should be email address",
		}},
		{"slice custom", func(c valigo.Configurator[texts], m *texts) {
			c.Slice(&m.Tags).As(str.Slice).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v shared.Elems[string]) []shared.Error {
				if v.Len()%2 != 0 {
					return []shared.Error{h.ErrorT(ctx, v.Len(), "Should contain pairs")}
				}
				return nil
			})
		}, &texts{Tags: []string{"a"}}, []string{"Tags: Should contain pairs"}},
		{"slice when", func(c valigo.Configurator[texts], m *texts) {
			c.Slice(&m.Tags).As(str.Slice).When(func(_ context.Context, v shared.Elems[string]) bool { return v.Len() > 0 }).Email()
		}, &texts{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, validate(t, tt.cfg, tt.obj))
		})
	}
}

func TestTrim(t *testing.T) {
	obj := &texts{Name: " a ", NamePtr: ptr(" b "), Tags: []string{" c "}, TagPtrs: []*string{ptr(" d "), nil}}
	assert.Empty(t, validate(t, func(c valigo.Configurator[texts], m *texts) {
		c.Field(&m.Name).As(str.Of).Trim()
		c.Field(&m.NamePtr).As(str.Of).Trim()
		c.Slice(&m.Tags).As(str.Slice).Trim()
		c.Slice(&m.TagPtrs).As(str.Slice).Trim()
	}, obj))
	assert.Equal(t, "a", obj.Name)
	assert.Equal(t, "b", *obj.NamePtr)
	assert.Equal(t, []string{"c"}, obj.Tags)
	assert.Equal(t, "d", *obj.TagPtrs[0])
}

func TestFormatRulesDontAllocate(t *testing.T) {
	if testenv.RaceEnabled {
		t.Skip("allocations are not stable under the race detector")
	}
	v := valigo.New()
	valigo.Configure[texts](v, func(c valigo.Configurator[texts], m *texts) {
		c.Field(&m.Name).As(str.Of).URL().Prefix("https://").Contains("example").MaxRunes(100)
		c.Field(&m.NamePtr).As(str.Of).Hostname()
		c.Field(&m.Status).As(str.Of).UUID()
		c.Slice(&m.Tags).As(str.Slice).Each(func(e *str.BaseConfigurator) { e.IP().Alphanumeric() })
	})
	obj := &texts{Name: "https://example.com", NamePtr: ptr("example.com"), Status: "6ba7b810-9dad-11d1-80b4-00c04fd430c8"}
	ctx := context.Background()
	assert.Nil(t, v.Validate(ctx, obj))
	assert.Zero(t, testing.AllocsPerRun(100, func() { v.Validate(ctx, obj) }))
}
