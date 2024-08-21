package translator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslatorGetLanguages(t *testing.T) {
	tr := &translator{
		storage:               &mockStorage{},
		getPreferredLanguages: func(ctx context.Context) []string { return []string{"fr", "es"} },
		defaultLang:           "en",
	}

	languages := tr.getLanguages(context.Background())
	assert.ElementsMatch(t, languages, []string{"fr", "es", "en"})

	tr.defaultLang = "fr"
	languages = tr.getLanguages(context.Background())
	assert.ElementsMatch(t, languages, []string{"fr", "es"})
}

func TestTranslatorErrorT(t *testing.T) {
	tr := &translator{
		storage:               &mockStorage{},
		getPreferredLanguages: func(ctx context.Context) []string { return []string{"fr", "es"} },
		defaultLang:           "en",
	}
	err := tr.ErrorT(context.Background(), "error format", "arg1", "arg2")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestTranslatorT(t *testing.T) {
	tr := &translator{
		storage:               &mockStorage{},
		getPreferredLanguages: func(ctx context.Context) []string { return []string{"fr", "es"} },
		defaultLang:           "en",
	}

	translated := tr.T(context.Background(), "format", "arg1", "arg2")
	assert.Equal(t, translated, mockTranslatedString)
}

func TestNewTranslator(t *testing.T) {
	tr := New()
	assert.NotNil(t, tr)

	tr = New(WithStorage(&mockStorage{}), WithDefaultLang("fr"))
	assert.NotNil(t, tr)
}
