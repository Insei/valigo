package translator

import "fmt"

// inMemTranslatorStorage represents an in-memory storage for translations.
// It maps languages to maps of translation keys to translated values.
type inMemTranslatorStorage struct {
	translations map[string]map[string]string
}

// NewInMemStorage creates a new in-memory translation storage.
// It loads initial data from the embedded locales YAML file.
// If an error occurs during loading, it panics.
func NewInMemStorage(opts ...InMemStorageOption) Storage {
	data, err := LocalesFromFS(EmbedFSLocalesYAML)
	if err != nil {
		panic(err)
	}
	storage := &inMemTranslatorStorage{translations: data}
	for _, opt := range opts {
		opt.apply(storage)
	}
	return storage
}

// Add adds new translations for a given language to the storage.
func (t *inMemTranslatorStorage) Add(lang string, data map[string]string) {
	if data == nil {
		return
	}
	if _, ok := t.translations[lang]; !ok {
		t.translations[lang] = data
	}
	for key, value := range data {
		t.translations[lang][key] = value
	}
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
func (t *inMemTranslatorStorage) Get(prefer []string, format string, args ...any) string {
	translatedFormat := format
	for _, preferLang := range prefer {
		langFormat, ok := t.translations[preferLang][format]
		if ok {
			translatedFormat = langFormat
			break
		}
	}
	return fmt.Sprintf(translatedFormat, args...)
}
