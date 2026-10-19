# go-i18n

The module `github.com/insei/valigo/integrations/goi18n` translates the messages of valigo by a bundle of [go-i18n](https://github.com/nicksnyder/go-i18n), so the messages of the validation live in the catalogs of your service with the rest of its texts.

```bash
go get github.com/insei/valigo/integrations/goi18n@latest
```

- `goi18n.Source(bundle)` is a [source of translations](translations.md#sources-of-translations): `translator.WithSources(goi18n.Source(bundle), translator.Embedded())` takes your messages from the bundle and the built-in YAML locales translate the rest.
- The message IDs are the locale keys of valigo, `validation:string:Cannot be longer than %d characters`, see [`translator/locales`](https://github.com/Insei/valigo/tree/main/translator/locales).
- A message uses the named parameters of the rule, `{{.max}}`, `{{.min}}`, `{{.field}}` — the params of the [codes](errors.md#codes-for-clients) — or the format verbs of the key, `%d`. The plural forms are chosen by the first integer argument, `{{.PluralCount}}` is it.
- A language without the message in the bundle has no translation: the fallback of go-i18n to its default language isn't used, the next sources translate the language of the request.
- `goi18n.Messages(lang)` exports the built-in messages of a language with their plural forms, to add them to a bundle or to write them to the files of the translators.

## Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/insei/valigo"
	"github.com/insei/valigo/integrations/goi18n"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/translator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type SignUp struct {
	Name  string
	Email string
}

func main() {
	bundle := i18n.NewBundle(language.English)
	_ = bundle.AddMessages(language.Russian, &i18n.Message{
		ID:    "validation:string:Cannot be longer than %d characters",
		One:   "Имя — не больше {{.max}} буквы",
		Few:   "Имя — не больше {{.max}} букв",
		Many:  "Имя — не больше {{.max}} букв",
		Other: "Имя — не больше {{.max}} буквы",
	})

	tr := translator.New(translator.WithSources(goi18n.Source(bundle), translator.Embedded()))
	v := valigo.New(valigo.WithTranslator(tr))
	valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
		b.Field(&s.Name).As(str.Of).MaxRunes(3)
		b.Field(&s.Email).As(str.Of).Email()
	})

	ctx := translator.ContextWithLanguages(context.Background(), "ru")
	for _, err := range v.ValidateTyped(ctx, &SignUp{Name: "Анастасия", Email: "ann"}) {
		fmt.Println(err.Location, "-", err.Message)
	}
}
```

Output:

```text
Name - Имя — не больше 3 букв
Email - Должно быть электронным адресом
```
