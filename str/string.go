package str

import (
	"context"
	"slices"
	"strings"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/helper"
)

type stringBuilder[T string | *string] struct {
	field    fmap.Field
	appendFn func(field fmap.Field, fn func(ctx context.Context, h *helper.Helper, v any) []error)
	enabler  func(ctx context.Context, value *T) bool
}

func (s *stringBuilder[T]) Trim() StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *helper.Helper, value any) []error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
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

func (s *stringBuilder[T]) MaxLen(maxLen uint) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *helper.Helper, value any) []error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
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

func (s *stringBuilder[T]) MinLen(minLen int) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *helper.Helper, value any) []error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			if len(*strVal) < minLen {

			}
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
	tagRequiredFormatKey := "validation:Field should be fulfilled"
	s.appendFn(s.field, func(ctx context.Context, h *helper.Helper, value any) []error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			if len(*strVal) < 1 {
				return []error{h.ErrorT(ctx, tagRequiredFormatKey)}
			}
		case **string:
			if *strVal == nil || len(**strVal) < 1 {
				return []error{h.ErrorT(ctx, tagRequiredFormatKey)}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) AnyOf(vals ...string) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *helper.Helper, value any) []error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			contains := slices.Contains(vals, *strVal)
			if !contains {
				return []error{h.ErrorT(ctx, "")}
			}
		case **string:
			if *strVal == nil {
				return []error{h.ErrorT(ctx, "")}
			}
			contains := slices.Contains(vals, **strVal)
			if !contains {
				return []error{h.ErrorT(ctx, "")}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) Custom(f func(ctx context.Context, h *helper.Helper, value *T) []error) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *helper.Helper, value any) []error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		return f(ctx, h, value.(*T))
	})
	return s
}

func (s *stringBuilder[T]) When(f func(ctx context.Context, value *T) bool) StringBuilder[T] {
	fn := f
	if s.enabler != nil {
		fn = func(ctx context.Context, value *T) bool {
			if s.enabler(ctx, value) {
				return f(ctx, value)
			}
			return false
		}
	}
	s.enabler = fn
	return s
}
