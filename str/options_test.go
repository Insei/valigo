package str

import (
	"testing"
)

func TestWithRegexpLocaleKey(t *testing.T) {
	// Test case 1: empty locale key
	options := WithRegexpLocaleKey("")
	regexpOpts := &regexpOptions{}
	options.apply(regexpOpts)
	if regexpOpts.localeKey != "" {
		t.Errorf("expected empty locale key, got %v", regexpOpts.localeKey)
	}

	// Test case 2: non-empty locale key
	options = WithRegexpLocaleKey("my-locale-key")
	regexpOpts = &regexpOptions{}
	options.apply(regexpOpts)
	if regexpOpts.localeKey != "my-locale-key" {
		t.Errorf("expected locale key 'my-locale-key', got %v", regexpOpts.localeKey)
	}
}
