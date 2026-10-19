package translator

import (
	"fmt"
	"io/fs"
	"strings"
)

// inMemTranslatorStorage represents an in-memory storage for translations.
// It maps languages to maps of translation keys to translated values.
type inMemTranslatorStorage struct {
	translations map[string]map[string]string
	// plurals indexes the plural forms of translations, key#category, by language, key and category,
	// so a plural form is found without building its key
	plurals map[string]map[string]map[string]string
}

// NewInMemStorage creates a new in-memory translation storage.
// It loads initial data from the embedded locales YAML file.
// If an error occurs during loading, it panics.
func NewInMemStorage(opts ...InMemStorageOption) Storage {
	return newInMemStorage(EmbedFSLocalesYAML, opts...)
}

// newInMemStorage creates a new in-memory translation storage with the locales of fsys, it panics on invalid locales.
func newInMemStorage(fsys fs.FS, opts ...InMemStorageOption) Storage {
	data, err := LocalesFromFS(fsys)
	if err != nil {
		panic(err)
	}
	storage := &inMemTranslatorStorage{translations: data}
	for _, opt := range opts {
		opt.apply(storage)
	}
	storage.plurals = make(map[string]map[string]map[string]string)
	for lang, translations := range storage.translations {
		storage.indexPlurals(lang, translations)
	}
	return storage
}

// indexPlurals adds the plural forms of translations, key#category, to the index of plural forms of lang.
func (t *inMemTranslatorStorage) indexPlurals(lang string, translations map[string]string) {
	if t.plurals == nil {
		t.plurals = make(map[string]map[string]map[string]string)
	}
	for key, value := range translations {
		format, category, ok := strings.Cut(key, "#")
		if !ok {
			continue
		}
		if t.plurals[lang] == nil {
			t.plurals[lang] = make(map[string]map[string]string)
		}
		if t.plurals[lang][format] == nil {
			t.plurals[lang][format] = make(map[string]string)
		}
		t.plurals[lang][format][category] = value
	}
}

// Add adds new translations for a given language to the storage.
func (t *inMemTranslatorStorage) Add(lang string, data map[string]string) {
	if data == nil {
		return
	}
	if _, ok := t.translations[lang]; !ok {
		t.translations[lang] = make(map[string]string, len(data))
	}
	for key, value := range data {
		t.translations[lang][key] = value
	}
	t.indexPlurals(lang, data)
}

func (t *inMemTranslatorStorage) Merge(locales map[string]map[string]string) {
	if t.translations == nil {
		return
	}
	for lang, data := range locales {
		t.Add(lang, data)
	}
}

// Get returns the translated value for a given format and language preferences.
// A format with plural forms (format#one, format#few, format#many, format#other) is translated
// by the plural category of the first integer argument in the language.
func (t *inMemTranslatorStorage) Get(prefer []string, format string, args ...any) string {
	translatedFormat := format
	for _, preferLang := range prefer {
		if langFormat, ok := t.translate(preferLang, format, args); ok {
			translatedFormat = langFormat
			break
		}
	}
	if len(args) == 0 && strings.IndexByte(translatedFormat, '%') < 0 {
		return translatedFormat
	}
	return fmt.Sprintf(translatedFormat, args...)
}

// translate returns the translation of format in lang, the plural form if the format has plural forms.
func (t *inMemTranslatorStorage) translate(lang, format string, args []any) (string, bool) {
	translations := t.translations[lang]
	if langFormat, ok := translations[format]; ok {
		return langFormat, true
	}
	n, ok := pluralCount(args)
	if !ok {
		return "", false
	}
	forms := t.plurals[lang][format]
	if langFormat, ok := forms[pluralCategory(lang, n)]; ok {
		return langFormat, true
	}
	langFormat, ok := forms[pluralOther]
	return langFormat, ok
}
