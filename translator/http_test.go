package translator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
)

func TestGetPriorityLanguages(t *testing.T) {
	tests := []struct {
		acceptLanguage string
		expected       []string
	}{
		{
			acceptLanguage: "en-US,en;q=0.9,es;q=0.8,fr;q=0.7",
			expected:       []string{"en", "es", "fr"},
		},
		{
			acceptLanguage: "en;q=0.9,es;q=0.8,fr;q=0.7,en-US",
			expected:       []string{"en", "es", "fr"},
		},
		{
			acceptLanguage: "en-US,en;q=0.9,es;q=0.8,fr;q=0.7,en-US",
			expected:       []string{"en", "es", "fr"},
		},
		{
			acceptLanguage: "en;q=0.9,es;q=0.8,fr;q=0.7,en-US",
			expected:       []string{"en", "es", "fr"},
		},
		{
			acceptLanguage: "",
			expected:       nil,
		},
	}

	for _, test := range tests {
		result := getPriorityLanguages(test.acceptLanguage)
		sort.Strings(result)
		sort.Strings(test.expected)
		if strings.Join(result, "") != strings.Join(test.expected, "") {
			t.Errorf("expected %v, got %v", test.expected, result)
		}
	}
}

func TestParseQuotient(t *testing.T) {
	tests := []struct {
		input    string
		expected float32
	}{
		{
			input:    "en;q=0.9",
			expected: 0.9,
		},
		{
			input:    "en;q=0.8",
			expected: 0.8,
		},
		{
			input:    "en;q=0.7",
			expected: 0.7,
		},
		{
			input:    "en",
			expected: 1,
		},
	}

	for _, test := range tests {
		result := parseQuotient(strings.Split(test.input, ";"))
		if result != test.expected {
			t.Errorf("expected %f, got %f", test.expected, result)
		}
	}
}

func TestNewAcceptLanguageMiddleware(t *testing.T) {
	middleware := NewAcceptLanguageMiddleware()
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		preferredLanguages := GetPreferredLanguagesFromContext(r.Context())
		if !slicesEqual(preferredLanguages, []string{"en", "ru"}) {
			t.Errorf("expected preferred languages to be [\"en\", \"ru\"], got %v", preferredLanguages)
		}
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
}

func TestGetPreferredLanguagesFromContext(t *testing.T) {
	ctx := context.Background()
	if GetPreferredLanguagesFromContext(ctx) != nil {
		t.Errorf("expected GetPreferredLanguagesFromContext(context.Background()) to be nil")
	}

	ctx = context.WithValue(ctx, languagesContextKeyVal, []string{"en-US", "en", "ru"})
	preferredLanguages := GetPreferredLanguagesFromContext(ctx)
	if !slicesEqual(preferredLanguages, []string{"en-US", "en", "ru"}) {
		t.Errorf("expected preferred languages to be [\"en-US\", \"en\", \"ru\"], got %v", preferredLanguages)
	}
}

func TestSortQuotient(t *testing.T) {
	q := sortQuotient{
		{quotient: 3.0},
		{quotient: 1.0},
		{quotient: 2.0},
	}

	// Test Len()
	if q.Len() != 3 {
		t.Errorf("Len() returned %d, want 3", q.Len())
	}

	// Test Swap()
	q.Swap(0, 1)
	if q[0].quotient != 1.0 || q[1].quotient != 3.0 {
		t.Errorf("Swap() did not swap elements correctly")
	}

	// Test Less()
	if q.Less(0, 1) {
		t.Errorf("Less() returned true, want false")
	}
	if !q.Less(1, 0) {
		t.Errorf("Less() returned false, want true")
	}

	// Test sorting
	sort.Sort(q)
	if q[0].quotient != 3.0 || q[1].quotient != 2.0 || q[2].quotient != 1.0 {
		t.Errorf("Sort() did not sort correctly")
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
