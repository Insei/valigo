# Translations

Messages are translated by a `translator.Translator`. The default one of `valigo.New` has the English, German, Spanish, French and Russian messages of the built-in rules, English is the default language.

## Locale keys

The rules build messages by *locale keys*: `validation:string:Cannot be longer than %d characters`. The translator looks the key up in the languages of the context, then in the default language, and formats the found message with the arguments of the rule. A key without a translation is used as is: so the messages of [custom rules](https://insei.github.io/valigo/0.2/features/custom/index.md) work without translations.

The keys of the built-in rules are in [`translator/locales`](https://github.com/Insei/valigo/tree/main/translator/locales):

```yaml
validation:
  string:
    "Cannot be longer than %d characters": Не может быть длиннее %d символов
```

## Plural forms

A message with a count has plural forms: the value of its key is a map of the plural categories of [CLDR](https://cldr.unicode.org/index/cldr-spec/plural-rules) — `one`, `few`, `many`, `other`. The form is selected by the first integer argument of the message, the categories of the language are built in for English and the default, French, Russian, Ukrainian and Belarusian; a missing category falls back to `other`.

```yaml
validation:
  string:
    "Cannot be longer than %d characters":
      one: Не может быть длиннее %d символа
      few: Не может быть длиннее %d символов
      many: Не может быть длиннее %d символов
      other: Не может быть длиннее %d символа
```

The forms are stored as the keys `key#one`, `key#few`, `key#many` and `key#other`, so they are added by `storage.Add` too:

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/shared"
    "github.com/insei/valigo/translator"
)

type Cart struct {
    Items []string
    Name  string
}

func main() {
    storage := translator.NewInMemStorage()
    storage.Add("ru", map[string]string{
        "Up to %d items in a cart#one":  "В корзине может быть не больше %d товара",
        "Up to %d items in a cart#many": "В корзине может быть не больше %d товаров",
    })
    t := translator.New(translator.WithStorage(storage), translator.WithPreferredLanguagesFn(func(context.Context) []string {
        return []string{"ru"}
    }))

    v := valigo.New(valigo.WithTranslator(t))
    valigo.Configure[Cart](v, func(b valigo.Configurator[Cart], c *Cart) {
        b.Slice(&c.Items).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, items shared.Elems[string]) []shared.Error {
            if items.Len() > 5 {
                return []shared.Error{h.ErrorT(ctx, items.Len(), "Up to %d items in a cart", 5)}
            }
            return nil
        })
        b.Slice(&c.Items).MaxLen(1)
    })

    for _, err := range v.Validate(context.Background(), &Cart{Items: []string{"a", "b", "c", "d", "e", "f"}}) {
        fmt.Println(err)
    }
}
```

Output:

```text
В корзине может быть не больше 5 товаров (Items: 6)
Не может содержать больше 1 элемента (Items: 6)
```

## Languages from the context

The languages of the request are taken from the context, `translator.NewAcceptLanguageMiddleware` puts the languages of the `Accept-Language` header there, ordered by their `q`:

```go
handler := translator.NewAcceptLanguageMiddleware()(mux)
```

Without the middleware of `net/http`, e.g. in the `Resolve` of a huma input, put the languages of the header to the context by `translator.ContextWithAcceptLanguage(ctx, header)`, or any languages by `translator.ContextWithLanguages(ctx, "ru")`, e.g. the language of the user profile. Another source of the languages is set by `translator.WithPreferredLanguagesFn`.

## Own messages

Add the translations of your messages to a storage and pass the translator to the validator:

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/str"
    "github.com/insei/valigo/shared"
    "github.com/insei/valigo/translator"
)

type language struct{}

type Account struct {
    Login string
}

const reservedKey = "validation:login:Login %q is reserved"

func main() {
    storage := translator.NewInMemStorage() // the built-in locales are loaded
    storage.Add("en", map[string]string{reservedKey: "Login %q is reserved"})
    storage.Add("ru", map[string]string{reservedKey: "Логин %q зарезервирован"})

    t := translator.New(
        translator.WithStorage(storage),
        translator.WithDefaultLang("en"),
        translator.WithPreferredLanguagesFn(func(ctx context.Context) []string {
            lang, _ := ctx.Value(language{}).(string)
            return []string{lang}
        }),
    )

    v := valigo.New(valigo.WithTranslator(t))
    valigo.Configure[Account](v, func(b valigo.Configurator[Account], a *Account) {
        b.Field(&a.Login).As(str.Of).MinLen(3).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, login string) []shared.Error {
            if login == "root" {
                return []shared.Error{h.ErrorT(ctx, login, reservedKey, login)}
            }
            return nil
        })
    })

    ctx := context.WithValue(context.Background(), language{}, "ru")
    fmt.Println(v.Validate(ctx, &Account{Login: "ro"}))
    fmt.Println(v.Validate(ctx, &Account{Login: "root"}))
    fmt.Println(v.Validate(context.Background(), &Account{Login: "root"}))
}
```

Output:

```text
[Не может быть короче 3 символов (Login: ro)]
[Логин "root" зарезервирован (Login: root)]
[Login "root" is reserved (Login: root)]
```

- `storage.Merge(locales)` adds the translations of several languages at once, e.g. from your own YAML files read by `translator.LocalesFromFS(fsys)`: `locales/<lang>/data.yaml` with the same nested format as the built-in locales.
- `translator.NewInMemStorage(translator.WithInMemData(data))` replaces the built-in locales with yours.
- [`str.WithRegexpLocaleKey`](https://insei.github.io/valigo/0.2/features/strings/#regexp-message) sets the key of a `Regexp` rule.

## Sources of translations

The translator looks the messages up in its sources: the built-in YAML locales by default, your catalogs of an i18n library, your own source. `translator.WithSources` chains them — for every language of a request in order, the first source with a translation of the key wins:

```go
tr := translator.New(translator.WithSources(goi18n.Source(bundle), translator.Embedded()))
v := valigo.New(valigo.WithTranslator(tr))
```

Your messages in go-i18n override the built-in ones, the built-in YAML translates the rest; a German request gets the built-in German before your English. `translator.Embedded(opts...)` is the built-in YAML locales as a source, with your additions by `WithInMemData`.

A source implements `translator.Source`:

```go
type Source interface {
    Translate(lang, key string, args []any, params map[string]any) (string, bool)
}
```

`args` are the format arguments of the key, `params` are the same values by their names, `{"max": 64}`, for the libraries with named placeholders; `false` passes the language to the next source. The sources:

| Source                                                                 | Module                                        |
| ---------------------------------------------------------------------- | --------------------------------------------- |
| the built-in YAML locales                                              | `translator.Embedded()`                       |
| [go-i18n](https://insei.github.io/valigo/0.2/features/goi18n/index.md) | `github.com/insei/valigo/integrations/goi18n` |
