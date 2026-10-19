package translator

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// messageKeys returns the keys of the messages of a locale, the plural forms of a key are one message.
func messageKeys(translations map[string]string) []string {
	seen := map[string]bool{}
	for key := range translations {
		base, _, _ := strings.Cut(key, "#")
		seen[base] = true
	}
	keys := make([]string, 0, len(seen))
	for key := range seen {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func TestEmbeddedLocalesAreComplete(t *testing.T) {
	locales, err := LocalesFromFS(EmbedFSLocalesYAML)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"de", "en", "es", "fr", "ru"}, func() []string {
		langs := make([]string, 0, len(locales))
		for lang := range locales {
			langs = append(langs, lang)
		}
		return langs
	}())
	want := messageKeys(locales["en"])
	for lang, translations := range locales {
		assert.Equal(t, want, messageKeys(translations), "the messages of %s", lang)
		for key, value := range translations {
			assert.Equal(t, strings.Count(key, "%"), strings.Count(value, "%"), "the arguments of %s in %s", key, lang)
			if base, category, ok := strings.Cut(key, "#"); ok {
				assert.Contains(t, translations, base+"#other", "the other form of %s in %s, the form of %s", base, lang, category)
			}
		}
	}
}

func TestEmbeddedLocalesTranslate(t *testing.T) {
	storage := NewInMemStorage()
	key := "validation:slice:Cannot contain more than %d elements"
	assert.Equal(t, "Darf nicht mehr als 1 Element enthalten", storage.Get([]string{"de"}, key, 1))
	assert.Equal(t, "No puede contener más de 3 elementos", storage.Get([]string{"es"}, key, 3))
	assert.Equal(t, "Ne peut pas contenir plus de 0 élément", storage.Get([]string{"fr"}, key, 0), "0 is one in French")
	assert.Equal(t, "Doit être une adresse e-mail", storage.Get([]string{"fr"}, "validation:string:Should be email address"))
}
