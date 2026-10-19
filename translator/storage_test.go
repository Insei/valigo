package translator

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemStorageMerge(t *testing.T) {
	storage := NewInMemStorage(WithInMemData(map[string]map[string]string{
		"en": {"hello": "Hello"},
	}))
	storage.Merge(map[string]map[string]string{
		"en": {"bye": "Bye"},
		"de": {"hello": "Hallo"},
	})
	assert.Equal(t, "Hello", storage.Get([]string{"en"}, "hello"))
	assert.Equal(t, "Bye", storage.Get([]string{"en"}, "bye"))
	assert.Equal(t, "Hallo", storage.Get([]string{"de", "en"}, "hello"))
}

func TestInMemStorageMergeWithoutTranslations(t *testing.T) {
	storage := &inMemTranslatorStorage{}
	storage.Merge(map[string]map[string]string{"en": {"hello": "Hello"}})
	assert.Nil(t, storage.translations)
}

func TestInMemStorageAddNil(t *testing.T) {
	storage := NewInMemStorage(WithInMemData(map[string]map[string]string{}))
	storage.Add("en", nil)
	assert.Equal(t, "hello", storage.Get([]string{"en"}, "hello"))
}

func TestWithInMemDataNil(t *testing.T) {
	storage := NewInMemStorage(WithInMemData(nil))
	assert.Equal(t, "Should be fulfilled", storage.Get([]string{"en"}, "validation:string:Should be fulfilled"),
		"nil data keeps the embedded locales")
}

func TestLocalesFromFSNested(t *testing.T) {
	locales, err := LocalesFromFS(fstest.MapFS{
		"locales/en/data.yaml": {Data: []byte("a:\n  b: B\n  c:\n    d: D\n")},
		"locales/de/README":    {Data: []byte("no data.yaml in the locale")},
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]map[string]string{"en": {"a:b": "B", "a:c:d": "D"}}, locales)
}

func TestLocalesFromFSErrors(t *testing.T) {
	_, err := LocalesFromFS(fstest.MapFS{})
	assert.ErrorContains(t, err, "can't read locales")

	_, err = LocalesFromFS(fstest.MapFS{"locales/en/data.yaml": {Data: []byte("a: [")}})
	assert.ErrorContains(t, err, "can't parse locales/en/data.yaml")

	_, err = LocalesFromFS(fstest.MapFS{"locales/en/data.yaml/x": {Data: []byte("")}})
	assert.ErrorContains(t, err, "locales/en/data.yaml")
}
