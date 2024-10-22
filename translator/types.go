package translator

import (
	"context"
)

// Translator interface defines methods for translating text.
type Translator interface {
	// ErrorT returns an error message translated into the desired language.
	ErrorT(ctx context.Context, format string, args ...any) error

	// T returns a translated string
	T(ctx context.Context, format string, args ...any) string
}

// StorageRO interface defines a method for retrieving translated text.
type StorageRO interface {
	// Get returns a translated string based on the provided preferences, format, and arguments.
	Get(prefer []string, format string, args ...any) string
}

// Storage interface extends StorageRO and adds a method for adding new translations.
type Storage interface {
	StorageRO

	// Add adds translations for a specific language.
	Add(lang string, data map[string]string)
	// Merge already existing translations with new
	Merge(locales map[string]map[string]string)
}

// Option interface defines a method for applying options to a translator.
type Option interface {
	// apply applies the option to a translator object.
	apply(t *translator)
}

// InMemStorageOption interface defines a method for applying options to an in-memory translator storage.
type InMemStorageOption interface {
	// apply applies the option to an inMemTranslatorStorage object.
	apply(t *inMemTranslatorStorage)
}
