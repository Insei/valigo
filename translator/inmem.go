package translator

import "fmt"

type inMemTranslatorStorage struct {
	translations map[string]map[string]string
}

func NewInMemStorage(opts ...InMemStorageOption) TranslationStorage {
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

func (t *inMemTranslatorStorage) AddTranslations(lang string, data map[string]string) error {
	if _, ok := t.translations[lang]; !ok {
		t.translations[lang] = data
		return nil
	}
	for key, value := range data {
		t.translations[lang][key] = value
	}
	return nil
}

func (t *inMemTranslatorStorage) GetTranslated(prefer []string, format string, args ...any) string {
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
