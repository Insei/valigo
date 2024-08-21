package str

import (
	"context"
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
			Required().
			Regexp(re)
	})

	testCases := []struct {
		name          string
		month         []string
		expectedError bool
	}{
		{
			name:          "Single month matching the regexp",
			month:         []string{"April"},
			expectedError: false,
		},
		{
			name:          "Multiple months with one matching the regexp",
			month:         []string{"March", "April", "May"},
			expectedError: true,
		},
		{
			name:          "Multiple months with none matching the regexp",
			month:         []string{"January", "February", "March"},
			expectedError: true,
		},
		{
			name:          "Empty month slice",
			month:         []string{},
			expectedError: true,
		},
		{
			name:          "Month slice with all empty strings",
			month:         []string{"", "", ""},
			expectedError: true,
		},
		{
			name:          "Month slice with all whitespace strings",
			month:         []string{"   ", "  ", " "},
			expectedError: true,
		},
		{
			name:          "Month slice with multiple matching months",
			month:         []string{"April", "April", "April"},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &year{
				Month: tc.month,
			}
			errs := v.Validate(context.Background(), y)
			if tc.expectedError {
				assert.NotEmpty(t, errs, "Should be validation errors")
			} else {
				assert.Empty(t, errs, "Should be no validation errors")
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
		name          string
		month         []string
		expectedError bool
	}{
		{
			name:          "Single month with trailing space",
			month:         []string{"April  "},
			expectedError: false,
		},
		{
			name:          "Multiple months with leading and trailing spaces",
			month:         []string{"  March", "April   ", "   May  "},
			expectedError: false,
		},
		{
			name:          "Month with only leading spaces",
			month:         []string{"  June"},
			expectedError: false,
		},
		{
			name:          "Month with only trailing spaces",
			month:         []string{"July  "},
			expectedError: false,
		},
		{
			name:          "Month with no spaces",
			month:         []string{"August"},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &year{
				Month: tc.month,
			}
			errs := v.Validate(context.Background(), y)
			if tc.expectedError {
				assert.NotEmpty(t, errs, "Should be validation errors")
			} else {
				assert.Empty(t, errs, "Should be no validation errors")
			}
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
		name          string
		phoneNumbers  []string
		expectedError bool
	}{
		{
			name:          "Valid phone numbers",
			phoneNumbers:  []string{"8910928772", "8982168233", "03", "911"},
			expectedError: false,
		},
		{
			name:          "Invalid phone numbers (too long)",
			phoneNumbers:  []string{"897213684987276313", "03", "911"},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &user{
				PhoneNumbers: tc.phoneNumbers,
			}
			errs := v.Validate(context.Background(), y)
			if tc.expectedError {
				assert.NotEmpty(t, errs, "Should be validation errors")
			} else {
				assert.Empty(t, errs, "Should be no validation errors")
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
		name          string
		phoneNumbers  []string
		expectedError bool
	}{
		{
			name:          "Valid phone numbers",
			phoneNumbers:  []string{"8910928772", "8982168233", "89313672813"},
			expectedError: false,
		},
		{
			name:          "Invalid phone numbers (too long)",
			phoneNumbers:  []string{"897213684987276313", "8982168233", "911"},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &user{
				PhoneNumbers: tc.phoneNumbers,
			}
			errs := v.Validate(context.Background(), y)
			if tc.expectedError {
				assert.NotNil(t, errs, "Expected validation error")
			} else {
				assert.Empty(t, errs, "Expected no validation error")
			}
		})
	}
}
