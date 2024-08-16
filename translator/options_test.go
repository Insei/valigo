package translator

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const mockTranslatedString = "mock translated string"

type mockStorage struct{}

func (m *mockStorage) GetTranslated(prefer []string, format string, args ...any) string {
	return mockTranslatedString
}

func TestWithStorage(t *testing.T) {
	storage := &mockStorage{}
	opt := WithStorage(storage)
	assert.NotNil(t, opt)
	assert.IsType(t, &storageOption{}, opt)

	so, ok := opt.(*storageOption)
	assert.True(t, ok)
	assert.Equal(t, storage, so.storage)
}

func TestStorageOptionApply(t *testing.T) {
	storage := &mockStorage{}
	opt := &storageOption{storage}
	tr := &translator{}
	opt.apply(tr)
	assert.Equal(t, storage, tr.storage)
}

func TestDefaultLangOption(t *testing.T) {
	lang := "en-US"
	option := WithDefaultLang(lang)
	tr := &translator{}
	option.apply(tr)
	if tr.defaultLang != lang {
		t.Errorf("expected translator.defaultLang to be set to %q, but got %q", lang, tr.defaultLang)
	}

	option = WithDefaultLang("")
	tr.defaultLang = ""
	option.apply(tr)
	if tr.defaultLang != "" {
		t.Errorf("expected translator.defaultLang to be empty, but got %q", tr.defaultLang)
	}
}

func TestInMemDataOption(t *testing.T) {
	data := map[string]map[string]string{
		"en-US": {
			"hello": "Hello",
		},
	}
	option := WithInMemData(data)
	storage := &inMemTranslatorStorage{}
	option.apply(storage)
	if !reflect.DeepEqual(storage.translations, data) {
		t.Errorf("expected storage.translations to be set to %v, but got %v", data, storage.translations)
	}

	option = WithInMemData(nil)
	option.apply(storage)
	if storage.translations == nil {
		t.Errorf("expected storage.translations to be %v, but got nil", data)
	}
}

func TestPreferredLanguagesFnOption(t *testing.T) {
	expectedLanguages := []string{"en-US", "fr-FR"}
	ctx := context.Background()

	fn := func(ctx context.Context) []string {
		return expectedLanguages
	}
	option := WithPreferredLanguagesFn(fn)
	tr := &translator{}
	option.apply(tr)
	actualLanguages := tr.getPreferredLanguages(ctx)
	if !reflect.DeepEqual(actualLanguages, expectedLanguages) {
		t.Errorf("expected translator.getPreferredLanguages to return %v, but got %v",
			expectedLanguages, actualLanguages)
	}

	option = WithPreferredLanguagesFn(nil)
	tr.getPreferredLanguages = nil
	option.apply(tr)
	if tr.getPreferredLanguages != nil {
		t.Errorf("expected translator.getPreferredLanguages to be nil, but got %v",
			tr.getPreferredLanguages(ctx))
	}
}
