# Features

Reference of valigo capabilities grouped by area.

## Configuration

| Page | What's inside |
|---|---|
| [Validator and Configure](configure.md) | `New` and its options, `Configure`, the model, field pointers, nested and recursive structs |
| [Types and conversions](conversions.md) | `Field`, `Slice` and `As`, the built-in conversions, conversions of your types |
| [Absent values](absent.md) | nil pointers, nil slices and nil nested structs: which rules skip them |
| [Conditions](conditions.md) | `When` on the struct and on a field |
| [Rules on other fields](fields.md) | `EqualField`, `MinField`, `MaxField`, `BeforeField`, `AfterField`: comparisons with another field by its pointer |
| [Custom rules](custom.md) | `Custom` on a field, on a slice and on the struct |

## Rules

| Page | What's inside |
|---|---|
| [Strings](strings.md) | `str.Of`: `Required`, lengths, `AnyOf`, `Prefix`, `Suffix`, `Contains`, `Regexp`, `Email`, `URL`, `IP`, `Hostname`, `UUID`, `Alphanumeric`, `Trim` |
| [Numbers](numbers.md) | `num.Of` and `num.Ptr`: `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval` |
| [Time](time.md) | `timex.Of`: `Required`, `Before`, `After`, `Between`, `InPast`, `InFuture` |
| [Comparable values](values.md) | `eq.Of`, `eq.Ptr`, `eq.Slice`: enums, bools — `Required`, `OneOf`, `NotOneOf`, `Equal`, `NotEqual`, `Unique` |
| [UUID](uuid.md) | `uuid.Of` of the separate module: `Required`, `AnyOf` |
| [Maps](maps.md) | `Map`: length rules, `EachKey` and `EachValue` with conversions |
| [Slices](slices.md) | `Slice`, `str.Slice`, `uuid.Slice`: length rules, `Each` for elements |

## Results

| Page | What's inside |
|---|---|
| [Errors](errors.md) | `Validate` and `ValidateTyped`, `shared.Error`, locations of fields and elements, errors transformer |
| [HTTP problem details](problem.md) | `problem`: the errors as `application/problem+json` of RFC 9457 with JSON pointers and codes |
| [huma](huma.md) | the module `github.com/insei/valigo/huma`: valigo rules in the `Resolve` of huma inputs |
| [Translations](translations.md) | Built-in locales, own messages, languages from the context, `Accept-Language` middleware |
