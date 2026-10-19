# Translations

Messages are translated by a `translator.Translator`. The default one of `valigo.New` has the English and Russian messages of the built-in rules, English is the default language.

## Locale keys

The rules build messages by *locale keys*: `validation:string:Cannot be longer than %d characters`. The translator looks the key up in the languages of the context, then in the default language, and formats the found message with the arguments of the rule. A key without a translation is used as is: so the messages of [custom rules](custom.md) work without translations.

The keys of the built-in rules are in [`translator/locales`](https://github.com/Insei/valigo/tree/main/translator/locales):

```yaml
validation:
  string:
    "Cannot be longer than %d characters": Не может быть длиннее %d символов
```

## Languages from the context

The languages of the request are taken from the context, `translator.NewAcceptLanguageMiddleware` puts the languages of the `Accept-Language` header there, ordered by their `q`:

```go
handler := translator.NewAcceptLanguageMiddleware()(mux)
```

Another source of the languages is set by `translator.WithPreferredLanguagesFn`, e.g. the language of the user profile.

## Own messages

Add the translations of your messages to a storage and pass the translator to the validator:

```go
package main

import (
	"context"
	"fmt"

	"github.com/insei/valigo"
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
		b.String(&a.Login).MinLen(3).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, login string) []shared.Error {
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
- [`str.WithRegexpLocaleKey`](strings.md#regexp-message) sets the key of a `Regexp` rule.
