package translator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFromMap(t *testing.T) {
	tests := []struct {
		name     string
		source   map[string]any
		prefix   string
		expected map[string]string
	}{
		{
			name: "simple",
			source: map[string]any{
				"key": "value",
			},
			prefix: "",
			expected: map[string]string{
				"key": "value",
			},
		},
		{
			name: "nested",
			source: map[string]any{
				"key": map[string]any{
					"nestedKey": "nestedValue",
				},
			},
			prefix: "",
			expected: map[string]string{
				"key:nestedKey": "nestedValue",
			},
		},
		{
			name: "multiple levels",
			source: map[string]any{
				"key": map[string]any{
					"nestedKey": map[string]any{
						"deepKey": "deepValue",
					},
				},
			},
			prefix: "",
			expected: map[string]string{
				"key:nestedKey:deepKey": "deepValue",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dest := make(map[string]string)
			fromMap(test.source, dest, test.prefix)
			if !mapEqual(dest, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, dest)
			}
		})
	}
}

func TestLocalesFromFS(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "locales")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, file := range []struct {
		name string
		data string
	}{
		{
			name: "locales/en/data.yaml",
			data: `key: value`,
		},
		{
			name: "locales/fr/data.yaml",
			data: `key: value2`,
		},
	} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(tmpDir, file.name)), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, file.name), []byte(file.data), 0644); err != nil {
			t.Fatal(err)
		}
	}

	fsys := os.DirFS(tmpDir)

	locales, err := LocalesFromFS(fsys)
	if err != nil {
		t.Fatal(err)
	}

	if len(locales) != 2 {
		t.Errorf("expected 2 locales, got %d", len(locales))
	}
	if _, ok := locales["en"]; !ok {
		t.Errorf("expected locale 'en', not found")
	}
	if _, ok := locales["fr"]; !ok {
		t.Errorf("expected locale 'fr', not found")
	}
	if !mapEqual(locales["en"], map[string]string{"key": "value"}) {
		t.Errorf("expected locale 'en' data %v, got %v", map[string]string{"key": "value"}, locales["en"])
	}
	if !mapEqual(locales["fr"], map[string]string{"key": "value2"}) {
		t.Errorf("expected locale 'fr' data %v, got %v", map[string]string{"key": "value2"}, locales["fr"])
	}
}

func mapEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
