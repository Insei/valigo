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

type baseConfigurator[T strPtr] struct {
	c     *shared.FieldConfigurator[T]
	field fmap.Field
	h     shared.Helper
}

// Trim removes leading and trailing whitespace from the string value.
func (i *baseConfigurator[T]) Trim() BaseConfigurator {
	i.c.Append(func(v T) bool {
		if v != nil {
			*v = strings.TrimSpace(*v)
		}
		return true
	}, "")
	return i
}

// MaxLen checks if the string length exceeds the maximum allowed length.
func (i *baseConfigurator[T]) MaxLen(maxLen int) BaseConfigurator {
	i.c.Append(func(v T) bool {
		return len(*v) <= maxLen
	}, maxLengthLocaleKey, maxLen)

	return i
}

// MinLen checks if the string length is not less than the given minimum length.
func (i *baseConfigurator[T]) MinLen(minLen int) BaseConfigurator {
	i.c.Append(func(v T) bool {
		return len(*v) >= minLen
	}, minLengthLocaleKey, minLen)

	return i
}

// Required checks if the string is not empty.
func (i *baseConfigurator[T]) Required() BaseConfigurator {
	i.c.Append(func(v T) bool {
		return len(*v) > 0
	}, requiredLocaleKey)

	return i
}

// Custom allows for custom validation logic to be applied to the string value.
func (i *baseConfigurator[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error) BaseConfigurator {
	customHelper := shared.NewFieldCustomHelper(i.field, i.h)
	i.c.CustomAppend(func(ctx context.Context, h shared.Helper, value any) []shared.Error {
		return f(ctx, customHelper, value)
	})
	return i
}

// Regexp checks if the string value matches the given regular expression.
func (i *baseConfigurator[T]) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) BaseConfigurator {
	options := regexpOptions{
		localeKey: regexpLocaleKey,
	}
	for _, opt := range opts {
		opt.apply(&options)
	}
	i.c.Append(func(v T) bool {
		return regexp.MatchString(*v)
	}, options.localeKey)
	return i
}

// AnyOf checks if the string value is one of the allowed values.
func (i *baseConfigurator[T]) AnyOf(allowed ...string) BaseConfigurator {
	i.c.Append(func(v T) bool {
		return slices.Contains(allowed, *v)
	}, anyOfLocaleKey)
	return i
}

// When allows for conditional validation logic to be applied to the string value.
func (i *baseConfigurator[T]) When(whenFn func(ctx context.Context, value any) bool) BaseConfigurator {
	if whenFn == nil {
		return i
	}
	base := i.c.NewWithWhen(func(ctx context.Context, value any) bool {
		v, ok := value.(**T)
		if !ok {
			return false
		}
		return whenFn(ctx, v)
	})
	return &baseConfigurator[T]{
		c:     base,
		field: i.field,
		h:     i.h,
	}
}
