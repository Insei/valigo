package shared

import (
	"context"
	"fmt"
	"iter"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/insei/fmap/v3"
)

const (
	sliceRequiredLocaleKey = "validation:slice:Should be fulfilled"
	sliceMaxLenLocaleKey   = "validation:slice:Cannot contain more than %d elements"
	sliceMinLenLocaleKey   = "validation:slice:Cannot contain less than %d elements"
)

// Elems is a view over the elements of a slice field, the slice can be []E or []*E.
type Elems[E any] struct {
	data unsafe.Pointer
	n    int
	ptrs bool
}

// Len returns the number of elements.
func (e Elems[E]) Len() int {
	return e.n
}

// At returns a pointer to the i-th element, nil for nil elements of []*E.
func (e Elems[E]) At(i int) *E {
	if e.ptrs {
		return *(**E)(unsafe.Add(e.data, uintptr(i)*unsafe.Sizeof(uintptr(0))))
	}
	return (*E)(unsafe.Add(e.data, uintptr(i)*unsafe.Sizeof(*new(E))))
}

// All iterates over indexes and pointers to the elements, see At.
func (e Elems[E]) All() iter.Seq2[int, *E] {
	return func(yield func(int, *E) bool) {
		for i := range e.n {
			if !yield(i, e.At(i)) {
				return
			}
		}
	}
}

// elemsGetter returns a function that reads elements of the slice field from obj.
// The field type can be []E, *[]E, []*E or *[]*E, named types with the same underlying layout are allowed.
// The returned function reports false when the slice is absent.
func elemsGetter[E any](field fmap.Field) func(obj any) (Elems[E], bool) {
	t := field.GetType()
	slicePtr := t.Kind() == reflect.Pointer
	if slicePtr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Slice {
		panic(fmt.Sprintf("valigo: field %s has type %s, slice expected", field.GetStructPath(), field.GetType()))
	}
	elem := t.Elem()
	// []*E is read as pointers to E, unless E is a pointer type itself, e.g. Slice(&m.X) for X []*int
	ptrs := elem.Kind() == reflect.Pointer && reflect.TypeFor[E]().Kind() != reflect.Pointer
	if ptrs {
		elem = elem.Elem()
	}
	mustLayout[E](field, elem)
	return func(obj any) (Elems[E], bool) {
		p := dataPtr(field.GetPtr(obj))
		if p != nil && slicePtr {
			p = *(*unsafe.Pointer)(p)
		}
		if p == nil {
			return Elems[E]{}, false
		}
		if ptrs {
			s := *(*[]*E)(p)
			return Elems[E]{data: unsafe.Pointer(unsafe.SliceData(s)), n: len(s), ptrs: true}, true
		}
		s := *(*[]E)(p)
		return Elems[E]{data: unsafe.Pointer(unsafe.SliceData(s)), n: len(s)}, true
	}
}

// eachSource applies checks to every element of the slice read by get.
// All checks are registered as a single validation function, so the slice is read once.
// Locations of the elements errors are suffixed with the element index, e.g. Tags[1].
func eachSource[E any](add func(fn ValidationFn), get func(obj any) (Elems[E], bool)) func(fn checkFn[E]) {
	var checks []checkFn[E]
	return func(fn checkFn[E]) {
		checks = append(checks, fn)
		if len(checks) > 1 {
			return
		}
		add(func(ctx context.Context, h Helper, obj any) []Error {
			elems, _ := get(obj)
			var errs []Error
			for i := range elems.n {
				v := elems.At(i)
				n := len(errs)
				for _, check := range checks {
					errs = append(errs, check(ctx, h, v)...)
				}
				for j := n; j < len(errs); j++ {
					errs[j].Location += "[" + strconv.Itoa(i) + "]"
				}
			}
			return errs
		})
	}
}

// SliceFieldConfigurator configures rules for the slice itself, e.g. its length.
type SliceFieldConfigurator[E any] struct {
	field fmap.Field
	h     Helper
	get   func(obj any) (Elems[E], bool)
	add   func(fn ValidationFn)
}

// NewSliceFieldConfigurator creates a SliceFieldConfigurator for the field of type []E, *[]E, []*E or *[]*E,
// named types with the same underlying layout are allowed.
// add registers a validation function of the configured struct.
func NewSliceFieldConfigurator[E any](field fmap.Field, h Helper, add func(fn ValidationFn)) *SliceFieldConfigurator[E] {
	return &SliceFieldConfigurator[E]{field: field, h: h, get: elemsGetter[E](field), add: add}
}

// Each configures rules applied to every element of the slice.
func (s *SliceFieldConfigurator[E]) Each(fn func(elem *Field[E])) *SliceFieldConfigurator[E] {
	fn(s.EachField())
	return s
}

// EachField returns a Field that applies rules to every element of the slice.
func (s *SliceFieldConfigurator[E]) EachField() *Field[E] {
	return &Field[E]{field: s.field, h: s.h, add: eachSource(s.add, s.get)}
}

func (s *SliceFieldConfigurator[E]) check(ok func(n int) bool, localeKey string, args ...any) {
	field, get := s.field, s.get
	s.add(func(ctx context.Context, h Helper, obj any) []Error {
		elems, present := get(obj)
		if !present || ok(elems.n) {
			return nil
		}
		return []Error{h.ErrorT(ctx, field, elems.n, localeKey, args...)}
	})
}

// MaxLen checks that the slice contains no more than maxLen elements. Absent slice is skipped.
func (s *SliceFieldConfigurator[E]) MaxLen(maxLen int) *SliceFieldConfigurator[E] {
	s.check(func(n int) bool { return n <= maxLen }, sliceMaxLenLocaleKey, maxLen)
	return s
}

// MinLen checks that the slice contains at least minLen elements. Absent slice is skipped.
func (s *SliceFieldConfigurator[E]) MinLen(minLen int) *SliceFieldConfigurator[E] {
	s.check(func(n int) bool { return n >= minLen }, sliceMinLenLocaleKey, minLen)
	return s
}

// Required checks that the slice is present and not empty.
func (s *SliceFieldConfigurator[E]) Required() *SliceFieldConfigurator[E] {
	field, get := s.field, s.get
	s.add(func(ctx context.Context, h Helper, obj any) []Error {
		if elems, _ := get(obj); elems.n > 0 {
			return nil
		}
		return []Error{h.ErrorT(ctx, field, nil, sliceRequiredLocaleKey)}
	})
	return s
}

// Custom allows for custom validation logic. Absent slice is skipped.
func (s *SliceFieldConfigurator[E]) Custom(f func(ctx context.Context, h *FieldCustomHelper, value Elems[E]) []Error) *SliceFieldConfigurator[E] {
	customHelper := NewFieldCustomHelper(s.field, s.h)
	get := s.get
	s.add(func(ctx context.Context, _ Helper, obj any) []Error {
		elems, present := get(obj)
		if !present {
			return nil
		}
		return f(ctx, customHelper, elems)
	})
	return s
}

// When returns a SliceFieldConfigurator whose rules are applied only if whenFn returns true.
// whenFn receives empty Elems when the slice is absent.
func (s *SliceFieldConfigurator[E]) When(whenFn func(ctx context.Context, value Elems[E]) bool) *SliceFieldConfigurator[E] {
	if whenFn == nil {
		return s
	}
	add, get := s.add, s.get
	return &SliceFieldConfigurator[E]{field: s.field, h: s.h, get: get, add: func(fn ValidationFn) {
		add(func(ctx context.Context, h Helper, obj any) []Error {
			elems, _ := get(obj)
			if !whenFn(ctx, elems) {
				return nil
			}
			return fn(ctx, h, obj)
		})
	}}
}
