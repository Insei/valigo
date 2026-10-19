package translator

import (
	"context"
	"errors"
	"slices"
)

// Translator is a struct that handles translations.
type translator struct {
	storage               StorageRO
	getPreferredLanguages func(ctx context.Context) []string
	defaultLang           string
	// defaultLangs is the default language alone, for the contexts without languages
	defaultLangs []string
}

// getLanguages returns the list of languages to use for translation,
// including the default language if necessary.
func (t *translator) getLanguages(ctx context.Context) []string {
	preferred := t.getPreferredLanguages(ctx)
	if len(preferred) == 0 && t.defaultLangs != nil {
		return t.defaultLangs
	}
	if !slices.Contains(preferred, t.defaultLang) {
		// a full slice expression: the languages of the context are never changed by append
		preferred = append(preferred[:len(preferred):len(preferred)], t.defaultLang)
	}
	return preferred
}

// ErrorT returns a translated error message.
func (t *translator) ErrorT(ctx context.Context, format string, args ...any) error {
	return errors.New(t.storage.Get(t.getLanguages(ctx), format, args...))
}

// T returns a translated string.
func (t *translator) T(ctx context.Context, format string, args ...any) string {
	return t.storage.Get(t.getLanguages(ctx), format, args...)
}

// TParams returns the translated message, the named parameters are passed to the sources, see Source.
func (t *translator) TParams(ctx context.Context, key string, args []any, params map[string]any) string {
	if s, ok := t.storage.(paramsStorage); ok {
		return s.getParams(t.getLanguages(ctx), key, args, params)
	}
	return t.storage.Get(t.getLanguages(ctx), key, args...)
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
	t.defaultLangs = []string{t.defaultLang}
	return t
}
