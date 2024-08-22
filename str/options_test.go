package str

import (
	"testing"
)

func TestWithRegexpLocaleKey(t *testing.T) {
	testCases := []struct {
		name          string
		localeKey     string
		expectedValue string
		expectedError bool
	}{
		{
			name:          "Empty locale key",
			localeKey:     "",
			expectedValue: "",
			expectedError: false,
		},
		{
			name:          "Non-empty locale key",
			localeKey:     "my-locale-key",
			expectedValue: "my-locale-key",
			expectedError: false,
		},
		{
			name:          "Locale key with special characters",
			localeKey:     "my-locale-key-with-special-chars!@#$%^&*()",
			expectedValue: "my-locale-key-with-special-chars!@#$%^&*()",
			expectedError: false,
		},
		{
			name:          "Locale key with whitespace",
			localeKey:     "my locale key with whitespace",
			expectedValue: "my locale key with whitespace",
			expectedError: false,
		},
		{
			name:          "Locale key with non-ASCII characters",
			localeKey:     "my-løcale-key",
			expectedValue: "my-løcale-key",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := WithRegexpLocaleKey(tc.localeKey)
			regexpOpts := &regexpOptions{}
			options.apply(regexpOpts)
			if regexpOpts.localeKey != tc.expectedValue {
				t.Errorf("expected locale key '%v', got '%v'", tc.expectedValue, regexpOpts.localeKey)
			}
		})
	}
}
