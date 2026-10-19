// Package goi18n translates the messages of valigo by a bundle of go-i18n (github.com/nicksnyder/go-i18n/v2).
//
// Source is a source of translations of the translator of valigo, the built-in YAML locales translate the rest:
//
//	tr := translator.New(translator.WithSources(goi18n.Source(bundle), translator.Embedded()))
//	v := valigo.New(valigo.WithTranslator(tr))
//
// The message IDs are the locale keys of valigo, validation:string:Cannot be longer than %d characters.
// A message can use the named parameters of the rule, {{.max}}, or the format verbs of the key, %d;
// the plural forms are chosen by the first integer argument. Messages exports the built-in messages
// to translate them with your messages.
package goi18n

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"

	"github.com/insei/valigo/translator"
)

// source is the translator.Source of a bundle.
type source struct {
	bundle *i18n.Bundle
	// localizers are the localizers of the languages, created once
	localizers sync.Map
}

// Source returns the translator.Source of the messages of the bundle. A language without the message in the bundle
// has no translation, the fallback of go-i18n to the default language isn't used: the next sources translate it.
func Source(bundle *i18n.Bundle) translator.Source {
	return &source{bundle: bundle}
}

type localizer struct {
	l    *i18n.Localizer
	base language.Base
}

func (s *source) localizer(lang string) localizer {
	if l, ok := s.localizers.Load(lang); ok {
		return l.(localizer)
	}
	base, _ := language.Make(lang).Base()
	l, _ := s.localizers.LoadOrStore(lang, localizer{l: i18n.NewLocalizer(s.bundle, lang), base: base})
	return l.(localizer)
}

// Translate returns the message of the key in the language, see translator.Source.
func (s *source) Translate(lang, key string, args []any, params map[string]any) (string, bool) {
	l := s.localizer(lang)
	cfg := &i18n.LocalizeConfig{MessageID: key, TemplateData: params}
	if n, ok := count(args); ok {
		cfg.PluralCount = n
	}
	msg, tag, err := l.l.LocalizeWithTag(cfg)
	if err != nil {
		return "", false
	}
	if base, _ := tag.Base(); base != l.base {
		// the fallback of go-i18n to another language
		return "", false
	}
	if len(args) > 0 && strings.IndexByte(msg, '%') >= 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return msg, true
}

// count returns the first integer argument, the count of the plural forms.
func count(args []any) (int64, bool) {
	for _, arg := range args {
		switch n := arg.(type) {
		case int:
			return int64(n), true
		case int8:
			return int64(n), true
		case int16:
			return int64(n), true
		case int32:
			return int64(n), true
		case int64:
			return n, true
		case uint:
			return int64(n), true
		case uint8:
			return int64(n), true
		case uint16:
			return int64(n), true
		case uint32:
			return int64(n), true
		case uint64:
			return int64(n), true
		}
	}
	return 0, false
}

// Messages returns the built-in messages of valigo in the language as the messages of go-i18n, with their plural
// forms, e.g. to add them to a bundle or to write them to the files of the translators.
// The messages keep the format verbs of the keys, %d.
func Messages(lang string) ([]*i18n.Message, error) {
	// the embedded locales are valid, the tests of the translator check them
	locales, _ := translator.LocalesFromFS(translator.EmbedFSLocalesYAML)
	translations, ok := locales[lang]
	if !ok {
		return nil, errors.New("goi18n: no built-in messages of " + lang)
	}
	byID := map[string]*i18n.Message{}
	var ids []string
	for key, value := range translations {
		id, category, plural := strings.Cut(key, "#")
		m, ok := byID[id]
		if !ok {
			m = &i18n.Message{ID: id}
			byID[id] = m
			ids = append(ids, id)
		}
		if !plural {
			category = "other"
		}
		setForm(m, category, value)
	}
	slices.Sort(ids)
	messages := make([]*i18n.Message, len(ids))
	for i, id := range ids {
		messages[i] = byID[id]
	}
	return messages, nil
}

// setForm sets the plural form of the CLDR category of the message, other for an unknown category.
func setForm(m *i18n.Message, category, value string) {
	switch category {
	case "zero":
		m.Zero = value
	case "one":
		m.One = value
	case "two":
		m.Two = value
	case "few":
		m.Few = value
	case "many":
		m.Many = value
	default:
		m.Other = value
	}
}
