# Features

Reference of valigo capabilities grouped by area.

## Configuration

| Page | What's inside |
|---|---|
| [Validator and Configure](configure.md) | `New` and its options, `Configure`, the model, field pointers, nested and recursive structs |
| [Absent values](absent.md) | nil pointers, nil slices and nil nested structs: which rules skip them |
| [Conditions](conditions.md) | `When` on the struct and on a field |
| [Custom rules](custom.md) | `Custom` on a field, on a slice and on the struct |

## Rules

| Page | What's inside |
|---|---|
| [Strings](strings.md) | `String`: `Required`, `MinLen`, `MaxLen`, `AnyOf`, `Regexp`, `Email`, `Trim` |
| [Numbers](numbers.md) | `Number` and `NumberPtr`: `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval` |
| [UUID](uuid.md) | `UUID`: `Required`, `AnyOf` |
| [Slices](slices.md) | `Slice`, `StringSlice`, `UUIDSlice`: length rules, `Each` for elements |

## Results

| Page | What's inside |
|---|---|
| [Errors](errors.md) | `Validate` and `ValidateTyped`, `shared.Error`, locations of fields and elements, errors transformer |
| [Translations](translations.md) | Built-in locales, own messages, languages from the context, `Accept-Language` middleware |
