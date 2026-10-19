# Why valigo?

**valigo is for Go services that validate their input structs and want the rules checked by the compiler, the validation free on valid input, and the messages in the language of the user.**

## What you get

### Rules are code, checked by the compiler

Rules are configured by pointers to the fields of the model, so the field and its type are known at compile time. `As` converts the field to the configurator of its type, and the conversion has a type constraint:

```go
b.Field(&u.Name).As(str.Of).Required().MaxLen(64) // str.Of accepts string, *string and named strings
b.Field(&u.ID).As(num.Of).Min(1)                  // Min takes an int64 for an int64 field
b.Field(&u.ID).As(num.Of).Min("1")                // doesn't compile
b.Field(&u.Age).As(str.Of)                        // doesn't compile: int does not satisfy ~string | *string
```

A renamed field is renamed in the rules by the IDE, a removed field breaks the build instead of a silently skipped rule. Tag based validators keep the rules in strings, `validate:"required,max=64"`, and report a typo or a rule of a wrong type only at runtime.

### Any type, without changes of the library

The configurator knows nothing about strings, numbers or uuids: they are conversions in their packages. A type of your domain, `time.Time` or a uuid of another package gets its rules by a conversion function of a few lines, see [Types and conversions](features/conversions.md).

### Validation of valid objects doesn't allocate

The rules are built once by `Configure`. `Validate` reads the fields through precomputed accessors and calls the rules of every field, allocations happen only to build errors. A struct with string, number, uuid and slice rules, conditions and a struct level custom rule is validated with 0 allocations — the tests check it. On the same rules a valid object is validated about 4 times faster than by go-playground/validator and 1.6 times slower than by the same rules written by hand, an invalid one with translated messages — about 2 times faster than go-playground/validator builds its messages, see [Performance](performance.md).

### Absent values are handled, not panicked on

A nil pointer field, a nil pointer to a slice and the fields of a nil nested struct are *absent*: every rule skips them, only `Required` fails. Optional fields need no special code:

```go
b.Field(&u.Phone).As(str.Of).Regexp(phoneRe)    // nil Phone is valid, a set one must match
b.Field(&u.Discount).As(num.Ptr).Max(50)         // nil Discount is valid
b.Field(&u.Address.City).As(str.Of).Required()  // fails for a nil Address too
```

### Conditions and custom rules at every level

`When` on the struct enables rules by the state of the object, `When` on a field — by the value of the field. `Custom` on a field gets its typed value, `Custom` on the struct gets the whole object and reports errors on any of its fields:

```go
smtp := b.When(func(_ context.Context, s *Sender) bool { return s.Type == "smtp" })
smtp.Field(&s.Host).As(str.Of).Required()
smtp.Field(&s.Port).As(num.Of).AnyOfInterval(1, 65535)
```

See [Conditions](features/conditions.md) and [Custom rules](features/custom.md).

### Messages in the language of the user, or codes for the client

One validator serves all languages: the language of every request comes from the context, `NewAcceptLanguageMiddleware` puts the languages of the `Accept-Language` header there. The messages of the built-in rules are translated to English, German, Spanish, French and Russian with plural forms, your own messages are added to the same storage and work untranslated until you add translations. A client with its own locales gets a stable `Code` and named `Params` of every error instead. See [Translations](features/translations.md) and [Codes for clients](features/errors.md#codes-for-clients).

How the libraries compare, by their documentation:

| | valigo | go-playground/validator | ozzo-validation |
|---|---|---|---|
| Translations out of the box | en, de, es, fr, ru | many languages, the packages of `validator/translations` | English messages |
| Language of a request | from the context, one validator for all languages | a translator of `universal-translator` is picked and passed to `Translate` for every request | the messages are global |
| Setup | none, `valigo.New()` | a translator per language and the registration of the translations of every tag | replacing the message templates of the errors |
| Plural forms | CLDR categories in the locale files | cardinal rules of `universal-translator` | — |
| Codes for clients | `Code` and named `Params` of every error | the tag and one string parameter of a field error | the code and the parameters of an error |
| Messages of own rules | the message is its locale key, translated when a translation is added | a translation is registered for every tag | the message of the error |

### Errors that point to the field

Every error has the message, the location and the value: `Tags[1]` for an element of a slice, `Address.City` for a nested field, or the JSON path of the field with a custom location function. See [Errors](features/errors.md).

## When valigo is not the right tool

- **The fields aren't known at compile time**, e.g. a form builder with fields defined by users: valigo validates the fields of Go structs. The limits and the conditions of the rules can come from the data — a database, a cache, the settings of a tenant — see [Rules from data](features/custom.md#rules-from-data).
- **The input is not a struct**: valigo validates structs, decode maps and JSON into a struct first.
- **Go older than 1.27**: the configurator uses generic methods.
