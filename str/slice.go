package str

import (
	"context"
	"github.com/insei/valigo/shared"
	"regexp"
	"strings"
)

const (
	typeLocaleKey = "validation:string:Only string values is allowed"
)

type StringSliceFieldConfigurator struct {
	*shared.SliceFieldConfigurator
}

func NewStringSliceFieldConfigurator(p shared.SliceFieldConfiguratorParams) *StringSliceFieldConfigurator {
	return &StringSliceFieldConfigurator{
		shared.NewSliceFieldConfigurator(p),
	}
}

func (s *StringSliceFieldConfigurator) Trim() *StringSliceFieldConfigurator {
	s.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v []*any) []shared.Error {
		values := shared.UnsafeValigoSliceCast[string](v)
		for _, val := range values {
			*val = strings.TrimSpace(*val)
		}
		return nil
	})
	return s
}

func (s *StringSliceFieldConfigurator) Regexp(regexp *regexp.Regexp, opts ...RegexpOption) *StringSliceFieldConfigurator {
	options := regexpOptions{
		localeKey: regexpLocaleKey,
	}
	for _, opt := range opts {
		opt.apply(&options)
	}
	s.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v []*any) []shared.Error {
		values := shared.UnsafeValigoSliceCast[string](v)
		var errs []shared.Error
		for _, val := range values {
			if !regexp.MatchString(*val) {
				errs = append(errs, h.ErrorT(ctx, *val, options.localeKey))
			}
		}

		return errs
	})
	return s
}

func (s *StringSliceFieldConfigurator) Email() *StringSliceFieldConfigurator {
	s.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, v []*any) []shared.Error {
		values := shared.UnsafeValigoSliceCast[string](v)
		var errs []shared.Error
		r := regexp.MustCompile(emailRegexp)
		for _, val := range values {
			if !r.MatchString(*val) {
				errs = append(errs, h.ErrorT(ctx, *val, emailLocaleKey))
			}
		}

		return errs
	})
	return s
}
