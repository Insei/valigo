package str

import (
	"context"
	"regexp"
	"slices"
	"strings"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/shared"
)

const (
	minLengthLocaleKey = "validation:string:Cannot be longer than %d characters"
	maxLengthLocaleKey = "validation:string:Cannot be longer than %d characters"
	requiredLocaleKey  = "validation:string:Should be fulfilled"
	regexpLocaleKey    = "validation:string:Doesn't match required regexp pattern"
	anyOfLocaleKey     = "validation:string:Only %s values is allowed"
)

type stringBuilder[T string | *string] struct {
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	enabler  func(ctx context.Context, value *T) bool
}

func (s *stringBuilder[T]) Trim() StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *shared.Helper, value any) []shared.Error {
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

func (s *stringBuilder[T]) MaxLen(maxLen int) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *shared.Helper, value any) []shared.Error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			if len(*strVal) > maxLen {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, maxLengthLocaleKey, maxLen)}
			}
		case **string:
			if *strVal == nil || len(**strVal) > maxLen {
				return []shared.Error{h.ErrorT(ctx, s.field, **strVal, maxLengthLocaleKey, maxLen)}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) MinLen(minLen int) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *shared.Helper, value any) []shared.Error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			if len(*strVal) < minLen {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, minLengthLocaleKey, minLen)}
			}
		case **string:
			if *strVal == nil || len(**strVal) < minLen {
				return []shared.Error{h.ErrorT(ctx, s.field, **strVal, minLengthLocaleKey, minLen)}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) Required() StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *shared.Helper, value any) []shared.Error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			if len(*strVal) < 1 {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, requiredLocaleKey)}
			}
		case **string:
			if *strVal == nil || len(**strVal) < 1 {
				return []shared.Error{h.ErrorT(ctx, s.field, **strVal, requiredLocaleKey)}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *shared.Helper, value any) []shared.Error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		options := regexpOptions{
			localeKey: regexpLocaleKey,
		}
		for _, opt := range opts {
			opt.apply(&options)
		}
		switch strVal := value.(type) {
		case *string:
			if regexp.FindString(*strVal) == "" {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, options.localeKey)}
			}
		case **string:
			if *strVal == nil || regexp.FindString(**strVal) == "" {
				return []shared.Error{h.ErrorT(ctx, s.field, **strVal, options.localeKey)}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) AnyOf(allowed ...string) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *shared.Helper, value any) []shared.Error {
		if s.enabler != nil && !s.enabler(ctx, value.(*T)) {
			return nil
		}
		switch strVal := value.(type) {
		case *string:
			if !slices.Contains(allowed, *strVal) {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, anyOfLocaleKey, "\""+strings.Join(allowed, "\",\"")+"\"")}
			}
		case **string:
			if *strVal == nil || !slices.Contains(allowed, **strVal) {
				return []shared.Error{h.ErrorT(ctx, s.field, **strVal, anyOfLocaleKey, "\""+strings.Join(allowed, "\",\"")+"\"")}
			}
		}
		return nil
	})
	return s
}

func (s *stringBuilder[T]) Custom(f func(ctx context.Context, h *shared.Helper, value *T) []shared.Error) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h *shared.Helper, value any) []shared.Error {
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
