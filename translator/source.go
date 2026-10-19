package translator

import (
	"fmt"
	"strings"
)

// Source is a source of the translations of the messages: the YAML locales of valigo, a bundle of go-i18n...
// A translator looks the languages of a request up in its sources in order, see WithSources.
type Source interface {
	// Translate returns the message of the key in the language formatted with the arguments, false if the source
	// has no translation of the key in the language. args are the format arguments of the key, params are the same
	// values by their names, e.g. {"max": 64}, for the sources with named placeholders; params can be nil.
	Translate(lang, key string, args []any, params map[string]any) (string, bool)
}

// WithSources returns an option that looks the translations up in the sources: for every language of a request
// in order, the first source with a translation of the key wins. Your messages go before the built-in ones,
// the built-in ones translate the rest:
//
//	translator.New(translator.WithSources(goi18n.Source(bundle), translator.Embedded()))
//
// A key without a translation in any language is the message itself.
func WithSources(sources ...Source) Option {
	return &storageOption{storage: &sourcesStorage{sources: sources}}
}

// Embedded returns the source of the built-in YAML locales of valigo, the options add your translations to them.
func Embedded(opts ...InMemStorageOption) Source {
	return newInMemStorage(EmbedFSLocalesYAML, opts...).(*inMemTranslatorStorage)
}

// paramsStorage is a storage that passes the named parameters of the messages to its sources.
type paramsStorage interface {
	getParams(prefer []string, key string, args []any, params map[string]any) string
}

// sourcesStorage is the storage of the translations of the sources, see WithSources.
type sourcesStorage struct {
	sources []Source
}

// Get returns the translation of the key in the first of the preferred languages that has one.
func (s *sourcesStorage) Get(prefer []string, key string, args ...any) string {
	return s.getParams(prefer, key, args, nil)
}

func (s *sourcesStorage) getParams(prefer []string, key string, args []any, params map[string]any) string {
	for _, lang := range prefer {
		for _, source := range s.sources {
			if msg, ok := source.Translate(lang, key, args, params); ok {
				return msg
			}
		}
	}
	return format(key, args)
}

// format formats the message with the arguments, a message without arguments and verbs is returned as is.
func format(msg string, args []any) string {
	if len(args) == 0 && strings.IndexByte(msg, '%') < 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}
