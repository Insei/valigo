package translator

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insei/valigo/internal/testenv"
)

func TestPluralCategory(t *testing.T) {
	tests := []struct {
		lang string
		n    int64
		want string
	}{
		{"en", 1, "one"}, {"en", 0, "other"}, {"en", 2, "other"}, {"en", 21, "other"}, {"en-US", 1, "one"},
		{"fr", 0, "one"}, {"fr", 1, "one"}, {"fr", 2, "other"},
		{"ru", 1, "one"}, {"ru", 21, "one"}, {"ru", 101, "one"}, {"ru", 11, "many"},
		{"ru", 2, "few"}, {"ru", 4, "few"}, {"ru", 22, "few"}, {"ru", 12, "many"}, {"ru", 14, "many"},
		{"ru", 5, "many"}, {"ru", 0, "many"}, {"ru", 111, "many"}, {"ru", -1, "one"}, {"ru-RU", 3, "few"},
		{"uk", 2, "few"}, {"be", 5, "many"},
		{"de", 1, "one"}, {"de", 3, "other"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, pluralCategory(tt.lang, tt.n), "%s %d", tt.lang, tt.n)
	}
}

func TestIsPluralForms(t *testing.T) {
	assert.True(t, isPluralForms(map[string]any{"one": "a", "few": "b", "many": "c", "other": "d", "zero": "e", "two": "f"}))
	assert.False(t, isPluralForms(map[string]any{}))
	assert.False(t, isPluralForms(map[string]any{"one": "a", "string": map[string]any{}}))
}

func TestPluralCount(t *testing.T) {
	for _, arg := range []any{int(3), int8(3), int16(3), int32(3), int64(3), uint(3), uint8(3), uint16(3), uint32(3), uint64(3)} {
		n, ok := pluralCount([]any{"x", arg})
		assert.True(t, ok, "%T", arg)
		assert.Equal(t, int64(3), n, "%T", arg)
	}
	_, ok := pluralCount([]any{"x", 1.5})
	assert.False(t, ok)
}

func TestGetPluralForms(t *testing.T) {
	storage := NewInMemStorage(WithInMemData(map[string]map[string]string{}))
	storage.Add("ru", map[string]string{
		"%d items#one":   "%d элемент",
		"%d items#few":   "%d элемента",
		"%d items#other": "%d элементов",
		"plain %d":       "просто %d",
		"items#other":    "элементы",
	})
	assert.Equal(t, "1 элемент", storage.Get([]string{"ru"}, "%d items", 1))
	assert.Equal(t, "3 элемента", storage.Get([]string{"ru"}, "%d items", 3))
	assert.Equal(t, "5 элементов", storage.Get([]string{"ru"}, "%d items", 5), "many falls back to other")
	assert.Equal(t, "просто 1", storage.Get([]string{"ru"}, "plain %d", 1), "a plain translation wins")
	assert.Equal(t, "items", storage.Get([]string{"ru"}, "items"), "no count, no plural form")
	assert.Equal(t, "2 items", storage.Get([]string{"de", "en"}, "%d items", 2), "no translation in any language")
}

func TestLocalesFromFSPluralForms(t *testing.T) {
	locales, err := LocalesFromFS(fstest.MapFS{
		"locales/ru/data.yaml": {Data: []byte("validation:\n  slice:\n    \"%d items\":\n      one: \"%d элемент\"\n      other: \"%d элементов\"\n")},
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]map[string]string{"ru": {
		"validation:slice:%d items#one":   "%d элемент",
		"validation:slice:%d items#other": "%d элементов",
	}}, locales)
}

func TestEmbeddedPluralForms(t *testing.T) {
	storage := NewInMemStorage()
	key := "validation:string:Cannot be longer than %d characters"
	assert.Equal(t, "Не может быть длиннее 1 символа", storage.Get([]string{"ru"}, key, 1))
	assert.Equal(t, "Не может быть длиннее 3 символов", storage.Get([]string{"ru"}, key, 3))
	assert.Equal(t, "Не может быть длиннее 21 символа", storage.Get([]string{"ru"}, key, 21))
	assert.Equal(t, "Cannot be longer than 1 character", storage.Get([]string{"en"}, key, 1))
	assert.Equal(t, "Cannot be longer than 5 characters", storage.Get([]string{"en"}, key, 5))
	assert.Equal(t, "Не может содержать больше 2 элементов", storage.Get([]string{"ru"}, "validation:map:Cannot contain more than %d elements", 2))
}

func TestPluralIndex(t *testing.T) {
	storage := NewInMemStorage(WithInMemData(map[string]map[string]string{"en": {"%d files#one": "%d file", "%d files#other": "%d files"}}))
	assert.Equal(t, "1 file", storage.Get([]string{"en"}, "%d files", 1), "the forms of WithInMemData are indexed")
	assert.Equal(t, "2 files", storage.Get([]string{"en"}, "%d files", 2))

	bare := &inMemTranslatorStorage{translations: map[string]map[string]string{}}
	bare.Add("de", map[string]string{"%d Dateien#one": "%d Datei"})
	assert.Equal(t, "1 Datei", bare.Get([]string{"de"}, "%d Dateien", 1))
}

func TestGetWithoutArgs(t *testing.T) {
	storage := NewInMemStorage(WithInMemData(map[string]map[string]string{"en": {"done": "Done", "100%% sure": "100%% sure"}}))
	assert.Equal(t, "Done", storage.Get([]string{"en"}, "done"))
	assert.Equal(t, "100% sure", storage.Get([]string{"en"}, "100%% sure"), "a format with %% is formatted")
	en := []string{"en"}
	if !testenv.RaceEnabled {
		assert.Zero(t, testing.AllocsPerRun(10, func() { storage.Get(en, "done") }))
	}
}
