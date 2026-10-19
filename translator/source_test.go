package translator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mapSource translates the keys of its languages and keeps the parameters of the last translation.
type mapSource struct {
	messages map[string]map[string]string
	params   map[string]any
}

func (s *mapSource) Translate(lang, key string, args []any, params map[string]any) (string, bool) {
	msg, ok := s.messages[lang][key]
	if ok {
		s.params = params
	}
	return format(msg, args), ok
}

func TestWithSources(t *testing.T) {
	own := &mapSource{messages: map[string]map[string]string{
		"ru": {"validation:string:Should be email address": "Нужен e-mail"},
		"en": {"validation:string:Should be URL": "A link, please", "greeting %s": "Hello, %s"},
	}}
	tr := New(WithSources(own, Embedded()))
	ru := ContextWithLanguages(context.Background(), "ru")
	assert.Equal(t, "Нужен e-mail", tr.T(ru, "validation:string:Should be email address"), "the first source")
	assert.Equal(t, "Должно быть URL", tr.T(ru, "validation:string:Should be URL"),
		"the built-in source of the language before the first source of the default language")
	assert.Equal(t, "Hello, Ann", tr.T(context.Background(), "greeting %s", "Ann"))
	assert.Equal(t, "no translation 1", tr.T(ru, "no translation %d", 1), "the key itself")
	assert.Equal(t, "plain", tr.T(ru, "plain"))
	assert.EqualError(t, tr.ErrorT(ru, "validation:string:Should be email address"), "Нужен e-mail")

	pt := tr.(*translator)
	params := map[string]any{"max": 3}
	assert.Equal(t, "Нужен e-mail", pt.TParams(ru, "validation:string:Should be email address", nil, params))
	assert.Equal(t, params, own.params, "the parameters reach the source")
}

// plainStorage is a storage without the parameters.
type plainStorage struct{}

func (plainStorage) Get(_ []string, key string, _ ...any) string { return "plain " + key }

func TestTParamsOfStorages(t *testing.T) {
	assert.Equal(t, "plain k", New(WithStorage(plainStorage{})).(*translator).TParams(context.Background(), "k", nil, nil))
	builtin := New().(*translator)
	assert.Equal(t, "Should be URL", builtin.TParams(context.Background(), "validation:string:Should be URL", nil, map[string]any{"x": 1}))
	assert.Equal(t, "Should be URL", NewInMemStorage().Get([]string{"en"}, "validation:string:Should be URL"))
	_, ok := Embedded().Translate("xx", "validation:string:Should be URL", nil, nil)
	assert.False(t, ok)
}
