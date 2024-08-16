package translator

import "context"

// storageOption represents an option to set the storage for a translator.
type storageOption struct {
	storage TranslationStorageRO
}

// apply sets the storage for the given translator.
func (o *storageOption) apply(t *translator) {
	t.storage = o.storage
}

// WithStorage returns an option to set the storage for a translator.
func WithStorage(storage TranslationStorageRO) Option {
	return &storageOption{storage}
}

// defaultLangOption represents an option to set the default language for a translator.
type defaultLangOption struct {
	lang string
}

// apply sets the default language for the given translator.
func (o *defaultLangOption) apply(s *translator) {
	if len(o.lang) > 1 {
		s.defaultLang = o.lang
	}
}

// WithDefaultLang returns an option to set the default language for a translator.
func WithDefaultLang(lang string) Option {
	return &defaultLangOption{lang: lang}
}

// inMemDataOption represents an option to set the in-memory data for a translator storage.
type inMemDataOption struct {
	data map[string]map[string]string
}

// apply sets the in-memory data for the given translator storage.
func (o *inMemDataOption) apply(s *inMemTranslatorStorage) {
	if o.data != nil {
		s.translations = o.data
	}
}

// WithInMemData returns an option to set the in-memory data for a translator storage.
func WithInMemData(data map[string]map[string]string) InMemStorageOption {
	return &inMemDataOption{data: data}
}

// preferredLanguagesFnOption represents an option to set the function to get the preferred languages for a translator.
type preferredLanguagesFnOption struct {
	fn func(ctx context.Context) []string
}

// apply sets the function to get the preferred languages for the given translator.
func (o *preferredLanguagesFnOption) apply(s *translator) {
	if o.fn != nil {
		s.getPreferredLanguages = o.fn
	}
}

// WithPreferredLanguagesFn returns an option to set the function
// to get the preferred languages for a translator.
func WithPreferredLanguagesFn(fn func(ctx context.Context) []string) Option {
	return &preferredLanguagesFnOption{fn: fn}
}
