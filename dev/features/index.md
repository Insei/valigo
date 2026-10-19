# Features

Reference of valigo capabilities grouped by area.

## Configuration

| Page                                                                                         | What's inside                                                                                                    |
| -------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| [Validator and Configure](https://insei.github.io/valigo/latest/features/configure/index.md) | `New` and its options, `Configure`, the model, field pointers, nested and recursive structs                      |
| [Types and conversions](https://insei.github.io/valigo/latest/features/conversions/index.md) | `Field`, `Slice` and `As`, the built-in conversions, conversions of your types                                   |
| [Absent values](https://insei.github.io/valigo/latest/features/absent/index.md)              | nil pointers, nil slices and nil nested structs: which rules skip them                                           |
| [Conditions](https://insei.github.io/valigo/latest/features/conditions/index.md)             | `When` on the struct and on a field                                                                              |
| [Rules on other fields](https://insei.github.io/valigo/latest/features/fields/index.md)      | `EqualField`, `MinField`, `MaxField`, `BeforeField`, `AfterField`: comparisons with another field by its pointer |
| [Custom rules](https://insei.github.io/valigo/latest/features/custom/index.md)               | `Custom` on a field, on a slice and on the struct                                                                |

## Rules

| Page                                                                                | What's inside                                                                                                                                      |
| ----------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Strings](https://insei.github.io/valigo/latest/features/strings/index.md)          | `str.Of`: `Required`, lengths, `AnyOf`, `Prefix`, `Suffix`, `Contains`, `Regexp`, `Email`, `URL`, `IP`, `Hostname`, `UUID`, `Alphanumeric`, `Trim` |
| [Numbers](https://insei.github.io/valigo/latest/features/numbers/index.md)          | `num.Of` and `num.Ptr`: `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval`                                                                         |
| [Time](https://insei.github.io/valigo/latest/features/time/index.md)                | `timex.Of`: `Required`, `Before`, `After`, `Between`, `InPast`, `InFuture`                                                                         |
| [Comparable values](https://insei.github.io/valigo/latest/features/values/index.md) | `eq.Of`, `eq.Ptr`, `eq.Slice`: enums, bools — `Required`, `OneOf`, `NotOneOf`, `Equal`, `NotEqual`, `Unique`                                       |
| [UUID](https://insei.github.io/valigo/latest/features/uuid/index.md)                | `uuid.Of` of the separate module: `Required`, `AnyOf`                                                                                              |
| [Maps](https://insei.github.io/valigo/latest/features/maps/index.md)                | `Map`: length rules, `EachKey` and `EachValue` with conversions                                                                                    |
| [Slices](https://insei.github.io/valigo/latest/features/slices/index.md)            | `Slice`, `str.Slice`, `uuid.Slice`: length rules, `Each` for elements                                                                              |

## Results

| Page                                                                                    | What's inside                                                                                        |
| --------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| [Errors](https://insei.github.io/valigo/latest/features/errors/index.md)                | `Validate` and `ValidateTyped`, `shared.Error`, locations of fields and elements, errors transformer |
| [HTTP problem details](https://insei.github.io/valigo/latest/features/problem/index.md) | `problem`: the errors as `application/problem+json` of RFC 9457 with JSON pointers and codes         |
| [huma](https://insei.github.io/valigo/latest/features/huma/index.md)                    | the module `github.com/insei/valigo/huma`: valigo rules in the `Resolve` of huma inputs              |
| [Translations](https://insei.github.io/valigo/latest/features/translations/index.md)    | Built-in locales, own messages, languages from the context, `Accept-Language` middleware             |
