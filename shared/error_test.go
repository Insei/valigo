package shared

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name          string
		error         Error
		expectedError bool
	}{
		{
			name: "Simple error",
			error: Error{
				Message: "Test error",
			},
			expectedError: true,
		},
		{
			name: "Error with value",
			error: Error{
				Message: "Test error",
				Value:   "some value",
			},
			expectedError: true,
		},
		{
			name: "Error with location and value",
			error: Error{
				Message:  "Test error",
				Location: "file.go",
				Value:    "some value",
			},
			expectedError: true,
		},
		{
			name: "Error with empty message",
			error: Error{
				Location: "file.go",
				Value:    "some value",
			},
			expectedError: true,
		},
		{
			name:          "No error",
			error:         Error{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.error.Error()
			if (result != "") != tt.expectedError {
				t.Errorf("Expected error: %v, but got: %s", tt.expectedError, result)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	errs := Errors{{Message: "a", Location: "x", Value: 1}, {Message: "b"}}
	if got := errs.Error(); got != "a (x: 1); b" {
		t.Errorf("Error() = %q", got)
	}
	var target Errors
	if !errors.As(fmt.Errorf("wrapped: %w", errs), &target) || len(target) != 2 {
		t.Errorf("errors.As failed: %v", target)
	}
	if got := Errors(nil).Error(); got != "" {
		t.Errorf("Error() of no errors = %q", got)
	}
}
