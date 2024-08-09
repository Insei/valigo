package translator

import (
	"context"
)

type Translator interface {
	ErrorT(ctx context.Context, format string, args ...any) error
	T(ctx context.Context, format string, args ...any) string
}

type TranslationStorageRO interface {
	GetTranslated(prefer []string, format string, args ...any) string
}

type TranslationStorage interface {
	TranslationStorageRO
	AddTranslations(lang string, data map[string]string) error
}

type Option interface {
	apply(t *translator)
}

type InMemStorageOption interface {
	apply(t *inMemTranslatorStorage)
}
