package shared

import (
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		error    Error
		expected string
	}{
		{
			name: "Simple error",
			error: Error{
				Message: "Test error",
			},
			expected: "Test error",
		},
		{
			name: "Error with value",
			error: Error{
				Message: "Test error",
				Value:   "some value",
			},
			expected: "Test error (: some value)",
		},
		{
			name: "Error with location and value",
			error: Error{
				Message:  "Test error",
				Location: "file.go",
				Value:    "some value",
			},
			expected: "Test error (file.go: some value)",
		},
		{
			name: "Error with empty message",
			error: Error{
				Location: "file.go",
				Value:    "some value",
			},
			expected: " (file.go: some value)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.error.Error()
			if result != tt.expected {
				t.Errorf("Expected: %s, but got: %s", tt.expected, result)
			}
		})
	}
}
