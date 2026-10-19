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

The configurator knows nothing about strings, numbers or uuids: they are conversions in their packages. A type of your domain, `time.Time` or a uuid of another package gets its rules by a conversion function of a few lines, see [Types and conversions](https://insei.github.io/valigo/latest/features/conversions/index.md).

### Validation of valid objects doesn't allocate

The rules are built once by `Configure`. `Validate` reads the fields through precomputed accessors and calls the rules of every field, allocations happen only to build errors. A struct with string, number, uuid and slice rules, conditions and a struct level custom rule is validated with 0 allocations — the tests check it.

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

See [Conditions](https://insei.github.io/valigo/latest/features/conditions/index.md) and [Custom rules](https://insei.github.io/valigo/latest/features/custom/index.md).

### Translated messages

Messages of the built-in rules are translated to English, German, Spanish, French and Russian, your own messages are added to the same storage. The language is taken from the context, `NewAcceptLanguageMiddleware` puts the languages of the `Accept-Language` header there. See [Translations](https://insei.github.io/valigo/latest/features/translations/index.md).

### Errors that point to the field

Every error has the message, the location and the value: `Tags[1]` for an element of a slice, `Address.City` for a nested field, or the JSON path of the field with a custom location function. See [Errors](https://insei.github.io/valigo/latest/features/errors/index.md).

## When valigo is not the right tool

- **Rules come from the data**, e.g. a form builder or rules stored in a database: valigo rules are Go code compiled into the service.
- **The input is not a struct**: valigo validates structs, decode maps and JSON into a struct first.
- **Go older than 1.27**: the configurator uses generic methods.
