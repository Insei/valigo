# Features

Reference of valigo capabilities grouped by area.

## Configuration

| Page                                                                                         | What's inside                                                                               |
| -------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| [Validator and Configure](https://insei.github.io/valigo/latest/features/configure/index.md) | `New` and its options, `Configure`, the model, field pointers, nested and recursive structs |
| [Types and conversions](https://insei.github.io/valigo/latest/features/conversions/index.md) | `Field`, `Slice` and `As`, the built-in conversions, conversions of your types              |
| [Absent values](https://insei.github.io/valigo/latest/features/absent/index.md)              | nil pointers, nil slices and nil nested structs: which rules skip them                      |
| [Conditions](https://insei.github.io/valigo/latest/features/conditions/index.md)             | `When` on the struct and on a field                                                         |
| [Custom rules](https://insei.github.io/valigo/latest/features/custom/index.md)               | `Custom` on a field, on a slice and on the struct                                           |

## Rules

| Page                                                                       | What's inside                                                                                                                                      |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Strings](https://insei.github.io/valigo/latest/features/strings/index.md) | `str.Of`: `Required`, lengths, `AnyOf`, `Prefix`, `Suffix`, `Contains`, `Regexp`, `Email`, `URL`, `IP`, `Hostname`, `UUID`, `Alphanumeric`, `Trim` |
| [Numbers](https://insei.github.io/valigo/latest/features/numbers/index.md) | `num.Of` and `num.Ptr`: `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval`                                                                         |
| [Time](https://insei.github.io/valigo/latest/features/time/index.md)       | `timex.Of`: `Required`, `Before`, `After`, `Between`, `InPast`, `InFuture`                                                                         |
| [UUID](https://insei.github.io/valigo/latest/features/uuid/index.md)       | `uuid.Of` of the separate module: `Required`, `AnyOf`                                                                                              |
| [Slices](https://insei.github.io/valigo/latest/features/slices/index.md)   | `Slice`, `str.Slice`, `uuid.Slice`: length rules, `Each` for elements                                                                              |

## Results

| Page                                                                                 | What's inside                                                                                        |
| ------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------- |
| [Errors](https://insei.github.io/valigo/latest/features/errors/index.md)             | `Validate` and `ValidateTyped`, `shared.Error`, locations of fields and elements, errors transformer |
| [Translations](https://insei.github.io/valigo/latest/features/translations/index.md) | Built-in locales, own messages, languages from the context, `Accept-Language` middleware             |
