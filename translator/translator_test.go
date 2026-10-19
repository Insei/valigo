package translator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo/internal/testenv"
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

func TestTranslatorErrorTKeepsPercent(t *testing.T) {
	storage := NewInMemStorage()
	storage.Add("en", map[string]string{"discount": "Discount %d%% is too big"})
	tr := New(WithStorage(storage))
	err := tr.ErrorT(context.Background(), "discount", 50)
	assert.EqualError(t, err, "Discount 50% is too big")
}

func TestTranslatorGetLanguagesWithoutPreferred(t *testing.T) {
	tr := New(WithPreferredLanguagesFn(func(context.Context) []string { return nil })).(*translator)
	assert.Equal(t, []string{"en"}, tr.getLanguages(context.Background()))
	if !testenv.RaceEnabled {
		assert.Zero(t, testing.AllocsPerRun(10, func() { tr.getLanguages(context.Background()) }))
	}

	bare := &translator{getPreferredLanguages: func(context.Context) []string { return nil }, defaultLang: "fr"}
	assert.Equal(t, []string{"fr"}, bare.getLanguages(context.Background()))
}

func TestTranslatorGetLanguagesDoesntChangeTheContextSlice(t *testing.T) {
	langs := make([]string, 1, 4)
	langs[0] = "ru"
	tr := New(WithPreferredLanguagesFn(func(context.Context) []string { return langs })).(*translator)
	assert.Equal(t, []string{"ru", "en"}, tr.getLanguages(context.Background()))
	assert.Equal(t, []string{"ru", ""}, langs[:2], "append must not write to the slice of the context")
}
