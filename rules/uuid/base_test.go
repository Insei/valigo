package uuid_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo"
	vuuid "github.com/insei/valigo/rules/uuid"
	"github.com/insei/valigo/shared"
)

type ids struct {
	ID      uuid.UUID
	IDPtr   *uuid.UUID
	IDs     []uuid.UUID
	IDPtrs  []*uuid.UUID
	IDsPtr  *[]uuid.UUID
	Parents *[]*uuid.UUID
}

var (
	id1 = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	id2 = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

func validate(t *testing.T, cfg func(c valigo.Configurator[ids], m *ids), obj *ids) []shared.Error {
	t.Helper()
	v := valigo.New()
	valigo.Configure[ids](v, cfg)
	return v.ValidateTyped(context.Background(), obj)
}

func TestRules(t *testing.T) {
	tests := []struct {
		name string
		cfg  func(c valigo.Configurator[ids], m *ids)
		obj  *ids
		want []shared.Error
	}{
		{
			name: "required nil uuid",
			cfg:  func(c valigo.Configurator[ids], m *ids) { c.Field(&m.ID).As(vuuid.Of).Required() },
			obj:  &ids{},
			want: []shared.Error{{Location: "ID", Message: "Should be fulfilled", Value: uuid.Nil, Code: "uuid.required"}},
		},
		{
			name: "required ok",
			cfg:  func(c valigo.Configurator[ids], m *ids) { c.Field(&m.ID).As(vuuid.Of).Required() },
			obj:  &ids{ID: id1},
		},
		{
			name: "required nil pointer",
			cfg:  func(c valigo.Configurator[ids], m *ids) { c.Field(&m.IDPtr).As(vuuid.Of).Required() },
			obj:  &ids{},
			want: []shared.Error{{Location: "IDPtr", Message: "Should be fulfilled", Code: "uuid.required"}},
		},
		{
			name: "any of",
			cfg:  func(c valigo.Configurator[ids], m *ids) { c.Field(&m.ID).As(vuuid.Of).AnyOf(id1) },
			obj:  &ids{ID: id2},
			want: []shared.Error{{Location: "ID", Message: "Only [" + id1.String() + "] values is allowed", Value: id2, Code: "uuid.any_of", Params: map[string]any{"values": []uuid.UUID{id1}}}},
		},
		{
			name: "any of skips nil pointer",
			cfg:  func(c valigo.Configurator[ids], m *ids) { c.Field(&m.IDPtr).As(vuuid.Of).AnyOf(id1) },
			obj:  &ids{},
		},
		{
			name: "custom",
			cfg: func(c valigo.Configurator[ids], m *ids) {
				c.Field(&m.ID).As(vuuid.Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v uuid.UUID) []shared.Error {
					if v.Version() != 7 {
						return []shared.Error{h.ErrorT(ctx, v, "Should be v7")}
					}
					return nil
				})
			},
			obj:  &ids{ID: id1},
			want: []shared.Error{{Location: "ID", Message: "Should be v7", Value: id1, Code: "Should be v7"}},
		},
		{
			name: "when",
			cfg: func(c valigo.Configurator[ids], m *ids) {
				c.Field(&m.ID).As(vuuid.Of).When(func(_ context.Context, v uuid.UUID) bool { return v != id1 }).AnyOf(id2)
			},
			obj: &ids{ID: id1},
		},
		{
			name: "slice length",
			cfg:  func(c valigo.Configurator[ids], m *ids) { c.Slice(&m.IDs).As(vuuid.Slice).MinLen(1).MaxLen(1) },
			obj:  &ids{IDs: []uuid.UUID{id1, id2}},
			want: []shared.Error{{Location: "IDs", Message: "Cannot contain more than 1 element", Value: 2, Code: "slice.max_len", Params: map[string]any{"max": 1}}},
		},
		{
			name: "slice required",
			cfg: func(c valigo.Configurator[ids], m *ids) {
				c.Field(&m.IDsPtr).As(shared.SlicePtr).As(vuuid.Slice).Required()
			},
			obj:  &ids{},
			want: []shared.Error{{Location: "IDsPtr", Message: "Should be fulfilled", Code: "slice.required"}},
		},
		{
			name: "slice any of",
			cfg:  func(c valigo.Configurator[ids], m *ids) { c.Slice(&m.IDPtrs).As(vuuid.Slice).AnyOf(id1) },
			obj:  &ids{IDPtrs: []*uuid.UUID{&id1, nil, &id2}},
			want: []shared.Error{{Location: "IDPtrs[2]", Message: "Only [" + id1.String() + "] values is allowed", Value: id2, Code: "uuid.any_of", Params: map[string]any{"values": []uuid.UUID{id1}}}},
		},
		{
			name: "slice each",
			cfg: func(c valigo.Configurator[ids], m *ids) {
				c.Field(&m.Parents).As(shared.SlicePtr).As(vuuid.Slice).Each(func(e *vuuid.BaseConfigurator) { e.Required() })
			},
			obj:  &ids{Parents: &[]*uuid.UUID{nil, &id1}},
			want: []shared.Error{{Location: "Parents[0]", Message: "Should be fulfilled", Code: "uuid.required"}},
		},
		{
			name: "slice custom",
			cfg: func(c valigo.Configurator[ids], m *ids) {
				c.Slice(&m.IDs).As(vuuid.Slice).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v shared.Elems[uuid.UUID]) []shared.Error {
					seen := map[uuid.UUID]bool{}
					for _, id := range v.All() {
						if seen[*id] {
							return []shared.Error{h.ErrorT(ctx, *id, "Should be unique")}
						}
						seen[*id] = true
					}
					return nil
				})
			},
			obj:  &ids{IDs: []uuid.UUID{id1, id2, id1}},
			want: []shared.Error{{Location: "IDs", Message: "Should be unique", Value: id1, Code: "Should be unique"}},
		},
		{
			name: "slice when",
			cfg: func(c valigo.Configurator[ids], m *ids) {
				c.Slice(&m.IDs).As(vuuid.Slice).When(func(_ context.Context, v shared.Elems[uuid.UUID]) bool { return v.Len() > 2 }).AnyOf(id1)
			},
			obj: &ids{IDs: []uuid.UUID{id2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, validate(t, tt.cfg, tt.obj))
		})
	}
}

func TestZeroAllocationsOnValid(t *testing.T) {
	if raceEnabled {
		t.Skip("allocations are not stable under the race detector")
	}
	v := valigo.New()
	valigo.Configure[ids](v, func(c valigo.Configurator[ids], m *ids) {
		c.Field(&m.ID).As(vuuid.Of).Required().AnyOf(id1, id2)
		c.Field(&m.IDPtr).As(vuuid.Of).Required()
		c.Slice(&m.IDPtrs).As(vuuid.Slice).Required().AnyOf(id1)
	})
	obj := &ids{ID: id1, IDPtr: &id2, IDPtrs: []*uuid.UUID{&id1}}
	ctx := context.Background()
	assert.Nil(t, v.Validate(ctx, obj))
	assert.Zero(t, testing.AllocsPerRun(100, func() { v.Validate(ctx, obj) }))
}
