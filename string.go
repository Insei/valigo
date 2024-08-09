package valigo

import (
	"fmt"
	"slices"
	"strings"

	"github.com/insei/fmap/v3"
)

type stringBuilder[T string | *string] struct {
	h        Helper
	field    fmap.Field
	appendFn func(field fmap.Field, fn func(h Helper, v any) []error)
	enabler  func(value *T) bool
}

func (s *stringBuilder[T]) Trim() StringBuilder[T] {
	s.appendFn(s.field, func(h Helper, value any) []error {
		if s.enabler != nil && !s.enabler(value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			*strVal = strings.TrimSpace(*strVal)
		case **string:
			if *strVal != nil {
				**strVal = strings.TrimSpace(**strVal)
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) Required() StringBuilder[T] {
	s.appendFn(s.field, func(h Helper, value any) []error {
		if s.enabler != nil && !s.enabler(value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			if len(*strVal) < 1 {
				return []error{fmt.Errorf("")}
			}
		case **string:
			if *strVal == nil || len(**strVal) < 1 {
				return []error{fmt.Errorf("")}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) AnyOf(vals ...string) StringBuilder[T] {
	s.appendFn(s.field, func(h Helper, value any) []error {
		if s.enabler != nil && !s.enabler(value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			contains := slices.Contains(vals, *strVal)
			if !contains {
				return []error{fmt.Errorf("")}
			}
		case **string:
			if *strVal == nil {
				return []error{fmt.Errorf("")}
			}
			contains := slices.Contains(vals, **strVal)
			if !contains {
				return []error{fmt.Errorf("")}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) Custom(f func(h Helper, value *T) []error) StringBuilder[T] {
	s.appendFn(s.field, func(h Helper, value any) []error {
		if s.enabler != nil && !s.enabler(value.(*T)) {
			return nil
		}
		return f(h, value.(*T))
	})
	return s
}

func (s *stringBuilder[T]) When(f func(value *T) bool) StringBuilder[T] {
	fn := f
	if s.enabler != nil {
		fn = func(value *T) bool {
			if s.enabler(value) {
				return f(value)
			}
			return false
		}
	}
	s.enabler = fn
	return s
}
