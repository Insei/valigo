package str

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo"
)

type year struct {
	Month []string
}

func TestStringSliceBuilderRegexp(t *testing.T) {
	v := valigo.New()
	re := regexp.MustCompile(`April`)
	valigo.Configure[year](v, func(builder valigo.Builder[year], obj *year) {
		builder.StringSlice(&obj.Month).
			Regexp(re)
	})

	testCases := []struct {
		name     string
		month    []string
		expected error
	}{
		{
			name:     "Single month matching the regexp",
			month:    []string{"April"},
			expected: nil,
		},
		{
			name:     "Multiple months with one matching the regexp",
			month:    []string{"March", "April", "May"},
			expected: errors.New("Doesn't match required regexp pattern (Month: [March April May])"),
		},
		{
			name:     "Multiple months with none matching the regexp",
			month:    []string{"January", "February", "March"},
			expected: errors.New("Doesn't match required regexp pattern (Month: [January February March])"),
		},
		{
			name:     "Empty month slice",
			month:    []string{},
			expected: errors.New("Doesn't match required regexp pattern (Month: [])"),
		},
		{
			name:     "Month slice with all empty strings",
			month:    []string{"", "", ""},
			expected: errors.New("Doesn't match required regexp pattern (Month: [])"),
		},
		{
			name:     "Month slice with all whitespace strings",
			month:    []string{"   ", "  ", " "},
			expected: errors.New("Doesn't match required regexp pattern (Month: [])"),
		},
		{
			name:     "Month slice with multiple matching months",
			month:    []string{"April", "April", "April"},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &year{
				Month: tc.month,
			}
			errs := v.Validate(context.Background(), y)
			if len(errs) == 0 {
				assert.Empty(t, errs, "Should be no validation errors")
			} else {
				assert.Error(t, tc.expected, errs[0])
			}
		})
	}
}

func TestStringSliceBuilderTrim(t *testing.T) {
	v := valigo.New()
	valigo.Configure[year](v, func(builder valigo.Builder[year], obj *year) {
		builder.StringSlice(&obj.Month).
			Trim()
	})

	testCases := []struct {
		name     string
		month    []string
		expected []error
	}{
		{
			name:     "Single month with trailing space",
			month:    []string{"April  "},
			expected: []error{},
		},
		{
			name:     "Multiple months with leading and trailing spaces",
			month:    []string{"  March", "April   ", "   May  "},
			expected: []error{},
		},
		{
			name:     "Month with only leading spaces",
			month:    []string{"  June"},
			expected: []error{},
		},
		{
			name:     "Month with only trailing spaces",
			month:    []string{"July  "},
			expected: []error{},
		},
		{
			name:     "Month with no spaces",
			month:    []string{"August"},
			expected: []error{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &year{
				Month: tc.month,
			}
			errs := v.Validate(context.Background(), y)
			assert.Empty(t, errs, "Should be no validation errors")
		})
	}
}

type user struct {
	PhoneNumbers []string
}

func TestStringSliceBuilderMaxLen(t *testing.T) {
	v := valigo.New()
	valigo.Configure[user](v, func(builder valigo.Builder[user], obj *user) {
		builder.StringSlice(&obj.PhoneNumbers).
			MaxLen(12)
	})

	testCases := []struct {
		name         string
		phoneNumbers []string
		expected     []error
	}{
		{
			name:         "Valid phone numbers",
			phoneNumbers: []string{"8910928772", "8982168233", "03", "911"},
			expected:     []error{},
		},
		{
			name:         "Invalid phone numbers (too long)",
			phoneNumbers: []string{"897213684987276313", "03", "911"},
			expected:     []error{fmt.Errorf("Cannot be longer than 12 characters (PhoneNumbers: 897213684987276313)")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &user{
				PhoneNumbers: tc.phoneNumbers,
			}
			errs := v.Validate(context.Background(), y)
			if errs != nil && len(errs) > 0 {
				assert.Equal(t, tc.expected[0].Error(), errs[0].Error(), "Should have the expected validation errors")
			} else {
				assert.Equal(t, tc.expected, errs, "Should have the expected validation errors")
			}
		})
	}
}

func TestStringSliceBuilderMinLen(t *testing.T) {
	v := valigo.New()
	valigo.Configure[user](v, func(builder valigo.Builder[user], obj *user) {
		builder.StringSlice(&obj.PhoneNumbers).
			MinLen(9)
	})

	testCases := []struct {
		name         string
		phoneNumbers []string
		expected     []error
	}{
		{
			name:         "Valid phone numbers",
			phoneNumbers: []string{"8910928772", "8982168233", "89313672813"},
			expected:     []error{},
		},
		{
			name:         "Invalid phone numbers (too long)",
			phoneNumbers: []string{"897213684987276313", "8982168233", "911"},
			expected:     []error{fmt.Errorf("Cannot be longer than 9 characters (PhoneNumbers: 911)")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &user{
				PhoneNumbers: tc.phoneNumbers,
			}
			errs := v.Validate(context.Background(), y)
			if errs != nil && len(errs) > 0 {
				assert.Equal(t, tc.expected[0].Error(), errs[0].Error(), "Should have the expected validation errors")
			} else {
				assert.Equal(t, tc.expected, errs, "Should have the expected validation errors")
			}
		})
	}
}
