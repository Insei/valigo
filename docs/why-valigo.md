# Why valigo?

**valigo is for Go services that validate their input structs and want the rules checked by the compiler, the validation free on valid input, and the messages in the language of the user.**

## What you get

### Rules are code, checked by the compiler

Rules are configured by pointers to the fields of the model, so the field and its type are known at compile time:

```go
b.String(&u.Name).Required().MaxLen(64)   // String accepts only string and *string fields
b.Number(&u.ID).Min(1)                    // Min takes an int64 for an int64 field
b.Number(&u.ID).Min("1")                  // doesn't compile
b.String(&u.Age)                          // doesn't compile: Age is an int
```

A renamed field is renamed in the rules by the IDE, a removed field breaks the build instead of a silently skipped rule. Tag based validators keep the rules in strings, `validate:"required,max=64"`, and report a typo or a rule of a wrong type only at runtime.

### Validation of valid objects doesn't allocate

The rules are built once by `Configure`. `Validate` reads the fields through precomputed accessors and calls the rules of every field, allocations happen only to build errors. A struct with string, number, uuid and slice rules, conditions and a struct level custom rule is validated with 0 allocations — the tests check it.

### Absent values are handled, not panicked on

A nil pointer field, a nil pointer to a slice and the fields of a nil nested struct are *absent*: every rule skips them, only `Required` fails. Optional fields need no special code:

```go
b.String(&u.Phone).Regexp(phoneRe)  // nil Phone is valid, a set one must match
b.NumberPtr(&u.Discount).Max(50)    // nil Discount is valid
b.String(&u.Address.City).Required() // fails for a nil Address too
```

### Conditions and custom rules at every level

`When` on the struct enables rules by the state of the object, `When` on a field — by the value of the field. `Custom` on a field gets its typed value, `Custom` on the struct gets the whole object and reports errors on any of its fields:

```go
smtp := b.When(func(_ context.Context, s *Sender) bool { return s.Type == "smtp" })
smtp.String(&s.Host).Required()
smtp.Number(&s.Port).AnyOfInterval(1, 65535)
```

See [Conditions](features/conditions.md) and [Custom rules](features/custom.md).

### Translated messages

Messages of the built-in rules are translated to English and Russian, your own messages are added to the same storage. The language is taken from the context, `NewAcceptLanguageMiddleware` puts the languages of the `Accept-Language` header there. See [Translations](features/translations.md).

### Errors that point to the field

Every error has the message, the location and the value: `Tags[1]` for an element of a slice, `Address.City` for a nested field, or the JSON path of the field with a custom location function. See [Errors](features/errors.md).

## When valigo is not the right tool

- **Rules come from the data**, e.g. a form builder or rules stored in a database: valigo rules are Go code compiled into the service.
- **The input is not a struct**: valigo validates structs, decode maps and JSON into a struct first.
- **Go older than 1.27**: the configurator uses generic methods.
