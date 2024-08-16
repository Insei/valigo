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

type car struct {
	Model string
}

func TestStringBuilderTrim(t *testing.T) {
	v := valigo.New()
	valigo.Configure[car](v, func(builder valigo.Builder[car], obj *car) {
		builder.String(&obj.Model).
			Trim()
	})

	testCases := []struct {
		name  string
		model string
	}{
		{
			name:  "Model with trailing space",
			model: "Ford  ",
		},
		{
			name:  "Multiple months with leading and trailing spaces",
			model: "  Ford Mustang  ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &car{
				Model: tc.model,
			}
			errs := v.Validate(context.Background(), y)
			assert.Empty(t, errs, "Should be no validation errors")
		})
	}
}

func TestStringBuilderRegexp(t *testing.T) {
	v := valigo.New()
	re := regexp.MustCompile(`^[a-zA-Z0-9]{3,20}$`)
	valigo.Configure[car](v, func(builder valigo.Builder[car], obj *car) {
		builder.String(&obj.Model).
			Regexp(re)
	})

	testCases := []struct {
		name     string
		model    string
		expected error
	}{
		{
			name:     "Valid model name",
			model:    "ToyotaCamry",
			expected: nil,
		},
		{
			name:     "Too short model name",
			model:    "Toy",
			expected: errors.New("regexp: regex did not match"),
		},
		{
			name:     "Too long model name",
			model:    "ThisIsAReallyLongModelNameThatShouldNotBeValid",
			expected: errors.New("regexp: regex did not match"),
		},
		{
			name:     "Model name with special characters",
			model:    "Toyota!Camry",
			expected: errors.New("regexp: regex did not match"),
		},
		{
			name:     "Model name with whitespace",
			model:    "Toyota Camry",
			expected: errors.New("regexp: regex did not match"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &car{
				Model: tc.model,
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

type admin struct {
	PhoneNumber string
}

func TestStringBuilderMaxLen(t *testing.T) {
	v := valigo.New()
	valigo.Configure[admin](v, func(builder valigo.Builder[admin], obj *admin) {
		builder.String(&obj.PhoneNumber).
			MaxLen(12)
	})

	testCases := []struct {
		name        string
		phoneNumber string
		expected    []error
	}{
		{
			name:        "Valid phone number",
			phoneNumber: "12345678901",
			expected:    []error{},
		},
		{
			name:        "Phone number too long",
			phoneNumber: "1234567890123",
			expected:    []error{fmt.Errorf("Cannot be longer than 12 characters (PhoneNumber: 1234567890123)")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &admin{
				PhoneNumber: tc.phoneNumber,
			}
			errs := v.Validate(context.Background(), y)
			if len(errs) == 0 {
				assert.Empty(t, errs, "Should have no validation errors")
			} else {
				assert.Equal(t, tc.expected[0].Error(), errs[0].Error(), "Should have the expected validation errors")
			}
		})
	}
}
