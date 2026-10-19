package goi18n_test

import (
	"context"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/insei/valigo"
	"github.com/insei/valigo/integrations/goi18n"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/translator"
)

const maxRunesKey = "validation:string:Cannot be longer than %d characters"

type signUp struct {
	Name  string
	Email string
	Nick  string
}

func newBundle(t *testing.T) *i18n.Bundle {
	t.Helper()
	bundle := i18n.NewBundle(language.English)
	require.NoError(t, bundle.AddMessages(language.Russian,
		&i18n.Message{ID: maxRunesKey, One: "Не длиннее {{.max}} буквы", Few: "Не длиннее {{.max}} букв", Many: "Не длиннее {{.max}} букв", Other: "Не длиннее {{.max}} буквы"},
		&i18n.Message{ID: "nick", Other: "Ник не %s"},
	))
	require.NoError(t, bundle.AddMessages(language.English, &i18n.Message{ID: "validation:string:Should be email address", Other: "An e-mail, please"}))
	return bundle
}

func TestSource(t *testing.T) {
	tr := translator.New(translator.WithSources(goi18n.Source(newBundle(t)), translator.Embedded()))
	v := valigo.New(valigo.WithTranslator(tr))
	valigo.Configure[signUp](v, func(b valigo.Configurator[signUp], s *signUp) {
		b.Field(&s.Name).As(str.Of).MaxRunes(3).MaxRunes(5)
		b.Field(&s.Email).As(str.Of).Email()
	})
	ctx := translator.ContextWithLanguages(context.Background(), "ru")
	var msgs []string
	for _, err := range v.ValidateTyped(ctx, &signUp{Name: "Анастасия", Email: "x"}) {
		msgs = append(msgs, err.Message)
	}
	assert.Equal(t, []string{"Не длиннее 3 букв", "Не длиннее 5 букв", "Должно быть электронным адресом"}, msgs,
		"the named parameters and the plural forms of go-i18n, the built-in Russian for the rest")

	de := translator.ContextWithLanguages(context.Background(), "de")
	msgs = nil
	for _, err := range v.ValidateTyped(de, &signUp{Email: "x"}) {
		msgs = append(msgs, err.Message)
	}
	assert.Equal(t, []string{"Muss eine E-Mail-Adresse sein"}, msgs, "the English of the bundle isn't the fallback of German")

	msgs = nil
	for _, err := range v.ValidateTyped(context.Background(), &signUp{Email: "x"}) {
		msgs = append(msgs, err.Message)
	}
	assert.Equal(t, []string{"An e-mail, please"}, msgs, "the default language")

	assert.Equal(t, "Ник не x", tr.T(ctx, "nick", "x"), "the format verbs of a key")
	assert.Equal(t, "unknown", tr.T(ctx, "unknown"))
}

func TestMessages(t *testing.T) {
	messages, err := goi18n.Messages("ru")
	require.NoError(t, err)
	var maxRunes *i18n.Message
	for _, m := range messages {
		if m.ID == maxRunesKey {
			maxRunes = m
		}
	}
	require.NotNil(t, maxRunes)
	assert.Equal(t, "Не может быть длиннее %d символа", maxRunes.One)
	assert.Equal(t, "Не может быть длиннее %d символов", maxRunes.Many)

	// the exported messages in a bundle translate like the built-in ones
	bundle := i18n.NewBundle(language.English)
	require.NoError(t, bundle.AddMessages(language.Russian, messages...))
	tr := translator.New(translator.WithSources(goi18n.Source(bundle)))
	ru := translator.ContextWithLanguages(context.Background(), "ru")
	assert.Equal(t, "Не может быть длиннее 21 символа", tr.T(ru, maxRunesKey, 21))
	assert.Equal(t, "Должно быть электронным адресом", tr.T(ru, "validation:string:Should be email address"))

	all, err := goi18n.Messages("en")
	require.NoError(t, err)
	assert.Len(t, all, len(messages), "every language has every message")

	_, err = goi18n.Messages("xx")
	assert.EqualError(t, err, "goi18n: no built-in messages of xx")
}
