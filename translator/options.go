package translator

import "context"

type storageOption struct {
	storage TranslationStorageRO
}

func (o *storageOption) apply(t *translator) {
	t.storage = o.storage
}

func WithStorage(storage TranslationStorageRO) Option {
	return &storageOption{storage}
}

type defaultLangOption struct {
	lang string
}

func (o *defaultLangOption) apply(s *translator) {
	if len(o.lang) > 1 {
		s.defaultLang = o.lang
	}
}

func WithDefaultLang(lang string) Option {
	return &defaultLangOption{lang: lang}
}

type inMemDataOption struct {
	data map[string]map[string]string
}

func (o *inMemDataOption) apply(s *inMemTranslatorStorage) {
	if o.data != nil {
		s.translations = o.data
	}
}

func WithInMemData(data map[string]map[string]string) InMemStorageOption {
	return &inMemDataOption{data: data}
}

type preferredLanguagesFnOption struct {
	fn func(ctx context.Context) []string
}

func (o *preferredLanguagesFnOption) apply(s *translator) {
	if o.fn != nil {
		s.getPreferredLanguages = o.fn
	}
}

func WithPreferredLanguagesFn(fn func(ctx context.Context) []string) Option {
	return &preferredLanguagesFnOption{fn: fn}
}
