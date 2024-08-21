package str

import (
	"context"
	"regexp"
	"strings"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

type stringSliceBuilder[T []string | *[]string] struct {
	field    fmap.Field
	appendFn func(field fmap.Field, fn shared.FieldValidationFn)
	h        shared.Helper
}

// Regexp applies a regular expression validation to the string slice.
// It checks if each string in the slice matches the given regular expression.
func (s *stringSliceBuilder[T]) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) StringSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		options := regexpOptions{
			localeKey: regexpLocaleKey,
		}
		for _, opt := range opts {
			opt.apply(&options)
		}
		switch strSliceVal := value.(type) {
		case *[]string:
			for _, strVal := range *strSliceVal {
				if !regexp.MatchString(strVal) {
					return []shared.Error{h.ErrorT(ctx, s.field, strVal, options.localeKey)}
				}
			}
		case **[]string:
			if *strSliceVal == nil {
				return []shared.Error{h.ErrorT(ctx, s.field, *strSliceVal, options.localeKey)}
			}
			for _, strVal := range **strSliceVal {
				if !regexp.MatchString(strVal) {
					return []shared.Error{h.ErrorT(ctx, s.field, strVal, options.localeKey)}
				}
			}
		}
		return nil
	})
	return s
}

// Trim trims each string in the slice.
func (s *stringSliceBuilder[T]) Trim() StringSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch strSliceVal := value.(type) {
		case *[]string:
			for i, _ := range *strSliceVal {
				(*strSliceVal)[i] = strings.TrimSpace((*strSliceVal)[i])
			}
		case **[]string:
			if *strSliceVal != nil {
				for i, _ := range **strSliceVal {
					(**strSliceVal)[i] = strings.TrimSpace((**strSliceVal)[i])
				}
			}
		}
		return nil
	})
	return s
}

// MaxLen checks if each string in the slice has a maximum length.
func (s *stringSliceBuilder[T]) MaxLen(maxLen int) StringSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch strSliceVal := value.(type) {
		case *[]string:
			for i, _ := range *strSliceVal {
				if len((*strSliceVal)[i]) > maxLen {
					return []shared.Error{h.ErrorT(ctx, s.field, (*strSliceVal)[i], maxLengthLocaleKey, maxLen)}
				}
			}
		case **[]string:
			if *strSliceVal != nil {
				for i, _ := range **strSliceVal {
					if len((**strSliceVal)[i]) > maxLen {
						return []shared.Error{h.ErrorT(ctx, s.field, (**strSliceVal)[i], maxLengthLocaleKey, maxLen)}
					}
				}
			}
		}
		return nil
	})
	return s
}

// MinLen checks if each string in the slice has a minimum length.
func (s *stringSliceBuilder[T]) MinLen(minLen int) StringSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch strSliceVal := value.(type) {
		case *[]string:
			for i, _ := range *strSliceVal {
				if len((*strSliceVal)[i]) < minLen {
					return []shared.Error{h.ErrorT(ctx, s.field, (*strSliceVal)[i], minLengthLocaleKey, minLen)}
				}
			}
		case **[]string:
			if *strSliceVal != nil {
				for i, _ := range **strSliceVal {
					if len((**strSliceVal)[i]) < minLen {
						return []shared.Error{h.ErrorT(ctx, s.field, (**strSliceVal)[i], minLengthLocaleKey, minLen)}
					}
				}
			}
		}
		return nil
	})
	return s
}

// Required checks if the string slice is not empty.
func (s *stringSliceBuilder[T]) Required() StringSliceBuilder[T] {
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		switch strSliceVal := value.(type) {
		case *[]string:
			if len(*strSliceVal) < 1 {
				return []shared.Error{h.ErrorT(ctx, s.field, strSliceVal, requiredLocaleKey)}
			}
		case **[]string:
			if *strSliceVal == nil {
				return []shared.Error{h.ErrorT(ctx, s.field, *strSliceVal, requiredLocaleKey)}
			}
		}
		return nil
	})
	return s
}

// Custom allows for custom validation logic.
func (s *stringSliceBuilder[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) StringSliceBuilder[T] {
	customHelper := shared.NewFieldCustomHelper(s.field, s.h)
	s.appendFn(s.field, func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value.(*T))
	})
	return s
}

// When allows for conditional validation based on a given condition.
func (s *stringSliceBuilder[T]) When(whenFn func(ctx context.Context, value *T) bool) StringSliceBuilder[T] {
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
