package translator

import (
	"context"
	"fmt"
	"slices"
)

// Translator is a struct that handles translations.
type translator struct {
	storage               StorageRO
	getPreferredLanguages func(ctx context.Context) []string
	defaultLang           string
}

// getLanguages returns the list of languages to use for translation,
// including the default language if necessary.
func (t *translator) getLanguages(ctx context.Context) []string {
	preferred := t.getPreferredLanguages(ctx)
	if !slices.Contains(preferred, t.defaultLang) {
		preferred = append(preferred, t.defaultLang)
	}
	return preferred
}

// ErrorT returns a translated error message.
func (t *translator) ErrorT(ctx context.Context, format string, args ...any) error {
	return fmt.Errorf(t.storage.Get(t.getLanguages(ctx), format, args...))
}

// T returns a translated string.
func (t *translator) T(ctx context.Context, format string, args ...any) string {
	return t.storage.Get(t.getLanguages(ctx), format, args...)
}

// New returns a new Translator instance with default settings.
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
