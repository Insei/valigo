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
	minLengthLocaleKey = "validation:string:Cannot be shorter than %d characters"
	maxLengthLocaleKey = "validation:string:Cannot be longer than %d characters"
	requiredLocaleKey  = "validation:string:Should be fulfilled"
	regexpLocaleKey    = "validation:string:Doesn't match required regexp pattern"
	anyOfLocaleKey     = "validation:string:Only %s values is allowed"
)

type stringBuilder[T string | *string] struct {
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	h        shared.Helper
}

// Trim removes leading and trailing whitespace from the string value.
func (s *stringBuilder[T]) Trim() StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
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

// MaxLen checks if the string length exceeds the maximum allowed length.
func (s *stringBuilder[T]) MaxLen(maxLen int) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
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

// MinLen checks if the string length is less than the minimum allowed length.
func (s *stringBuilder[T]) MinLen(minLen int) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch strVal := value.(type) {
		case *string:
			if len(*strVal) < minLen {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, minLengthLocaleKey, minLen)}
			}
		case **string:
			if *strVal == nil || len(**strVal) < minLen {
				return []shared.Error{h.ErrorT(ctx, s.field, "", minLengthLocaleKey, minLen)}
			}
		}
		return nil
	})
	return s
}

// Required checks if the string value is not empty.
func (s *stringBuilder[T]) Required() StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch strVal := value.(type) {
		case *string:
			if len(*strVal) < 1 {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, requiredLocaleKey)}
			}
		case **string:
			if *strVal == nil || len(**strVal) < 1 {
				return []shared.Error{h.ErrorT(ctx, s.field, "", requiredLocaleKey)}
			}
		}
		return nil
	})
	return s
}

// Regexp checks if the string value matches the given regular expression.
func (s *stringBuilder[T]) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		options := regexpOptions{
			localeKey: regexpLocaleKey,
		}
		for _, opt := range opts {
			opt.apply(&options)
		}
		switch strVal := value.(type) {
		case *string:
			if !regexp.MatchString(*strVal) {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, options.localeKey)}
			}
		case **string:
			if *strVal == nil || !regexp.MatchString(**strVal) {
				return []shared.Error{h.ErrorT(ctx, s.field, "", options.localeKey)}
			}
		}
		return nil
	})
	return s
}

// AnyOf checks if the string value is one of the allowed values.
func (s *stringBuilder[T]) AnyOf(allowed ...string) StringBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch strVal := value.(type) {
		case *string:
			if !slices.Contains(allowed, *strVal) {
				return []shared.Error{h.ErrorT(ctx, s.field, *strVal, anyOfLocaleKey, "\""+strings.Join(allowed, "\",\"")+"\"")}
			}
		case **string:
			if *strVal == nil || !slices.Contains(allowed, **strVal) {
				return []shared.Error{h.ErrorT(ctx, s.field, "", anyOfLocaleKey, "\""+strings.Join(allowed, "\",\"")+"\"")}
			}
		}
		return nil
	})
	return s
}

// Custom allows for custom validation logic to be applied to the string value.
func (s *stringBuilder[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) StringBuilder[T] {
	customHelper := shared.NewFieldCustomHelper(s.field, s.h)
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(*T))
	})
	return s
}

// When allows for conditional validation logic to be applied to the string value.
func (s *stringBuilder[T]) When(whenFn func(ctx context.Context, value *T) bool) StringBuilder[T] {
	if whenFn == nil {
		return s
	}
	s.appendFn = func(field fmap.Field, fn shared.FieldValidationFn) {
		fnWithEnabler := func(ctx context.Context, h shared.Helper, v any) []shared.Error {
			if !whenFn(ctx, v.(*T)) {
				return nil
			}
			return fn(ctx, h, v)
		}
		s.appendFn(field, fnWithEnabler)
	}
	return s
}
