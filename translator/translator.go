package translator

import (
	"context"
	"fmt"
	"slices"
)

type translator struct {
	storage               TranslationStorageRO
	getPreferredLanguages func(ctx context.Context) []string
	defaultLang           string
}

func (t *translator) getLanguages(ctx context.Context) []string {
	preferred := t.getPreferredLanguages(ctx)
	if !slices.Contains(preferred, t.defaultLang) {
		preferred = append(preferred, t.defaultLang)
	}
	return preferred
}

func (t *translator) ErrorT(ctx context.Context, format string, args ...any) error {
	return fmt.Errorf(t.storage.GetTranslated(t.getLanguages(ctx), format, args...))
}

func (t *translator) T(ctx context.Context, format string, args ...any) string {
	return t.storage.GetTranslated(t.getLanguages(ctx), format, args...)
}

func New(opts ...Option) Translator {
	t := &translator{
		storage:               NewInMemStorage(),
		getPreferredLanguages: GetPreferredLanguagesFromContext,
		defaultLang:           "en",
	}
	for _, opt := range opts {
		opt.apply(t)
	}
	return t
}
