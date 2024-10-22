package translator

import (
	"testing"
)

func TestInMemTranslatorStorage(t *testing.T) {
	storage := NewInMemStorage()
	lang := "en"
	data := map[string]string{
		"hello": "Hello, %s!",
	}
	storage.Add(lang, data)
	if _, ok := storage.(*inMemTranslatorStorage).translations[lang]; !ok {
		t.Errorf("expected translations to be added for lang %s", lang)
	}
	if _, ok := storage.(*inMemTranslatorStorage).translations[lang]["hello"]; !ok {
		t.Errorf("expected translation to be added for key 'hello' in lang %s", lang)
	}

	prefer := []string{"en", "fr"}
	format := "hello"
	args := []any{"John"}
	translated := storage.Get(prefer, format, args...)
	if translated != "Hello, John!" {
		t.Errorf("expected translated string to be 'Hello, John!', got '%s'", translated)
	}

	prefer = []string{"es", "fr"}
	format = "hello"
	args = []any{"John"}
	translated = storage.Get(prefer, format, args...)
	if translated != "hello%!(EXTRA string=John)" {
		t.Errorf("expected translated string to be empty, got '%s'", translated)
	}
}

func TestInMemTranslatorStorageAddTranslations_ExistingLang(t *testing.T) {
	storage := NewInMemStorage()

	lang := "en"
	data := map[string]string{
		"hello": "Hello, %s!",
	}
	storage.Add(lang, data)

	newData := map[string]string{
		"goodbye": "Goodbye, %s!",
	}
	storage.Add(lang, newData)

	if _, ok := storage.(*inMemTranslatorStorage).translations[lang]["goodbye"]; !ok {
		t.Errorf("expected new translation to be added for key 'goodbye' in lang %s", lang)
	}
}

func TestInMemTranslatorStorage_GetTranslated_NoPrefer(t *testing.T) {
	storage := NewInMemStorage()

	format := "hello"
	args := []any{"John"}
	translated := storage.Get(nil, format, args...)
	if translated != "hello%!(EXTRA string=John)" {
		t.Errorf("expected translated string to be empty, got '%s'", translated)
	}
}
