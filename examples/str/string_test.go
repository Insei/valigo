package str

import (
	"context"
	"errors"
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
	valigo.Configure[car](v, func(builder valigo.Configurator[car], obj *car) {
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
	valigo.Configure[car](v, func(builder valigo.Configurator[car], obj *car) {
		builder.String(&obj.Model).
			Regexp(re)
	})

	testCases := []struct {
		name          string
		model         string
		expectedError bool
	}{
		{
			name:          "Valid model name",
			model:         "ToyotaCamry",
			expectedError: false,
		},
		{
			name:          "Too short model name",
			model:         "Toy",
			expectedError: false,
		},
		{
			name:          "Too long model name",
			model:         "ThisIsAReallyLongModelNameThatShouldNotBeValid",
			expectedError: true,
		},
		{
			name:          "Model name with special characters",
			model:         "Toyota!Camry",
			expectedError: true,
		},
		{
			name:          "Model name with whitespace",
			model:         "Toyota Camry",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &car{
				Model: tc.model,
			}
			errs := v.Validate(context.Background(), y)
			if tc.expectedError {
				if len(errs) == 0 {
					t.Errorf("Expected error, but got none")
				} else {
					assert.Error(t, errors.New("regexp: regex did not match"), errs[0])
				}
			} else {
				assert.Empty(t, errs, "Should be no validation errors")
			}
		})
	}
}

type admin struct {
	PhoneNumber string
}

func TestStringBuilderMaxLen(t *testing.T) {
	v := valigo.New()
	valigo.Configure[admin](v, func(builder valigo.Configurator[admin], obj *admin) {
		builder.String(&obj.PhoneNumber).
			MaxLen(12)
	})

	testCases := []struct {
		name          string
		phoneNumber   string
		expectedError bool
	}{
		{
			name:          "Valid phone number",
			phoneNumber:   "12345678901",
			expectedError: false,
		},
		{
			name:          "Phone number too long",
			phoneNumber:   "1234567890123",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			y := &admin{
				PhoneNumber: tc.phoneNumber,
			}
			errs := v.Validate(context.Background(), y)
			if tc.expectedError {
				if len(errs) == 0 {
					t.Errorf("Expected error, but got none")
				}
			} else {
				assert.Empty(t, errs, "Should be no validation errors")
			}
		})
	}
}
