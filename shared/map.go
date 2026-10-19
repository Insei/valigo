package shared

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/insei/fmap/v3"
)

const (
	mapRequiredLocaleKey = "validation:map:Should be fulfilled"
	mapMaxLenLocaleKey   = "validation:map:Cannot contain more than %d elements"
	mapMinLenLocaleKey   = "validation:map:Cannot contain less than %d elements"
)

// MapFieldConfigurator configures rules for a map field: the map itself, its keys and its values.
type MapFieldConfigurator[K comparable, V any] struct {
	field  fmap.Field
	h      Helper
	fields Fields
	get    func(obj any) *map[K]V
	add    func(fn ValidationFn)
}

// NewMapFieldConfigurator creates a MapFieldConfigurator for the field of type map[K]V,
// fields resolves the pointers to other fields of the model, add registers a validation function of the configured struct.
func NewMapFieldConfigurator[K comparable, V any](field fmap.Field, h Helper, fields Fields, add func(fn ValidationFn)) *MapFieldConfigurator[K, V] {
	return &MapFieldConfigurator[K, V]{field: field, h: h, fields: fields, get: valueGetter[map[K]V](field), add: add}
}

// As converts the map configurator to a typed one by conv.
func (m *MapFieldConfigurator[K, V]) As[R any](conv func(m *MapFieldConfigurator[K, V]) R) R {
	return conv(m)
}

func (m *MapFieldConfigurator[K, V]) check(ok func(n int) bool, msg Msg) {
	field, get := m.field, m.get
	m.add(func(ctx context.Context, h Helper, obj any) []Error {
		v := get(obj)
		if v == nil || ok(len(*v)) {
			return nil
		}
		return []Error{msg.Error(ctx, h, field, len(*v))}
	})
}

// MaxLen checks that the map contains no more than maxLen elements. Absent map is skipped.
func (m *MapFieldConfigurator[K, V]) MaxLen(maxLen int) *MapFieldConfigurator[K, V] {
	m.check(func(n int) bool { return n <= maxLen }, M("map.max_len", mapMaxLenLocaleKey, P("max", maxLen)))
	return m
}

// MinLen checks that the map contains at least minLen elements. Absent map is skipped.
func (m *MapFieldConfigurator[K, V]) MinLen(minLen int) *MapFieldConfigurator[K, V] {
	m.check(func(n int) bool { return n >= minLen }, M("map.min_len", mapMinLenLocaleKey, P("min", minLen)))
	return m
}

// Required checks that the map is present and not empty.
func (m *MapFieldConfigurator[K, V]) Required() *MapFieldConfigurator[K, V] {
	field, get, msg := m.field, m.get, M("map.required", mapRequiredLocaleKey)
	m.add(func(ctx context.Context, h Helper, obj any) []Error {
		if v := get(obj); v != nil && len(*v) > 0 {
			return nil
		}
		return []Error{msg.Error(ctx, h, field, nil)}
	})
	return m
}

// Custom allows for custom validation logic on the whole map. Absent map is skipped.
func (m *MapFieldConfigurator[K, V]) Custom(f func(ctx context.Context, h *FieldCustomHelper, value map[K]V) []Error) *MapFieldConfigurator[K, V] {
	customHelper := NewFieldCustomHelper(m.field, m.h)
	get := m.get
	m.add(func(ctx context.Context, _ Helper, obj any) []Error {
		v := get(obj)
		if v == nil {
			return nil
		}
		return f(ctx, customHelper, *v)
	})
	return m
}

// When returns a MapFieldConfigurator whose rules are applied only if whenFn returns true.
// whenFn receives a nil map when the map is absent.
func (m *MapFieldConfigurator[K, V]) When(whenFn func(ctx context.Context, value map[K]V) bool) *MapFieldConfigurator[K, V] {
	if whenFn == nil {
		return m
	}
	add, get := m.add, m.get
	return &MapFieldConfigurator[K, V]{field: m.field, h: m.h, fields: m.fields, get: get, add: func(fn ValidationFn) {
		add(func(ctx context.Context, h Helper, obj any) []Error {
			var v map[K]V
			if p := get(obj); p != nil {
				v = *p
			}
			if !whenFn(ctx, v) {
				return nil
			}
			return fn(ctx, h, obj)
		})
	}}
}

// EachKey configures rules applied to every key of the map, the location of an error is Labels[key].
func (m *MapFieldConfigurator[K, V]) EachKey(fn func(key *Field[K])) *MapFieldConfigurator[K, V] {
	fn(&Field[K]{field: m.field, h: m.h, fields: m.fields, add: mapSource(m.add, m.get, func(k K, _ V) K { return k })})
	return m
}

// EachValue configures rules applied to every value of the map, the location of an error is Labels[key].
// The values of a map aren't addressable: Modify rules, e.g. Trim, change a copy that is checked by the next rules,
// the map keeps the original value.
func (m *MapFieldConfigurator[K, V]) EachValue(fn func(value *Field[V])) *MapFieldConfigurator[K, V] {
	fn(&Field[V]{field: m.field, h: m.h, fields: m.fields, add: mapSource(m.add, m.get, func(_ K, v V) V { return v })})
	return m
}

// mapSource applies checks to the keys or the values of the map read by get, pick selects the checked one.
// All checks are registered as a single validation function. The checked value is copied to a pooled slot,
// so valid maps are validated without allocations. Errors are sorted by location, the order of a map is random.
func mapSource[K comparable, V, T any](add func(fn ValidationFn), get func(obj any) *map[K]V, pick func(k K, v V) T) func(fn checkFn[T]) {
	var checks []checkFn[T]
	slots := &sync.Pool{New: func() any { return new(T) }}
	return func(fn checkFn[T]) {
		checks = append(checks, fn)
		if len(checks) > 1 {
			return
		}
		add(func(ctx context.Context, h Helper, obj any) []Error {
			m := get(obj)
			if m == nil {
				return nil
			}
			slot := slots.Get().(*T)
			var errs []Error
			for k, v := range *m {
				*slot = pick(k, v)
				n := len(errs)
				for _, check := range checks {
					errs = append(errs, check(ctx, h, obj, slot)...)
				}
				for j := n; j < len(errs); j++ {
					errs[j].Location += "[" + fmt.Sprint(k) + "]"
				}
			}
			var zero T
			*slot = zero
			slots.Put(slot)
			slices.SortStableFunc(errs, func(a, b Error) int { return cmp.Compare(a.Location, b.Location) })
			return errs
		})
	}
}
