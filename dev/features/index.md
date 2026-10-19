# Features

The map of valigo, in the order of the menu.

## Configuration

How the rules are configured and composed.

| Page                                                                                         | What's inside                                                                                       |
| -------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| [Validator and Configure](https://insei.github.io/valigo/latest/features/configure/index.md) | `New` and its options, `Configure`, the model, field pointers, nested and recursive structs         |
| [Types and conversions](https://insei.github.io/valigo/latest/features/conversions/index.md) | `Field`, `Slice`, `Map` and `As`, the built-in conversions, conversions of your types               |
| [Absent values](https://insei.github.io/valigo/latest/features/absent/index.md)              | nil pointers, nil slices and nil nested structs: which rules skip them                              |
| [Conditions](https://insei.github.io/valigo/latest/features/conditions/index.md)             | `When` on the struct and on a field, conditional requiredness                                       |
| [Custom rules](https://insei.github.io/valigo/latest/features/custom/index.md)               | `Custom` on a field, on a slice and on the struct, rules with limits from data                      |
| [Rules on other fields](https://insei.github.io/valigo/latest/features/fields/index.md)      | `EqualField`, `MinField`, `MaxField`, `BeforeField`, `AfterField` by the pointer of the other field |

## Rules

The built-in rules of every type.

| Page                                                                                | What's inside                                                                                                                                                                              |
| ----------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [Strings](https://insei.github.io/valigo/latest/features/strings/index.md)          | `str.Of`: lengths, `AnyOf`, substrings, `Regexp`, `Email`, `URL`, `IP`, `CIDR`, `MAC`, `Hostname`, `UUID`, `E164`, `Base64`, `JSON`, `Semver`, `Datetime`, ISO codes, `CreditCard`, `Trim` |
| [Numbers](https://insei.github.io/valigo/latest/features/numbers/index.md)          | `num.Of`, `num.Ptr`: `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval`, `Latitude`, `Longitude`                                                                                           |
| [Time](https://insei.github.io/valigo/latest/features/time/index.md)                | `timex.Of`: `Required`, `Before`, `After`, `Between`, `InPast`, `InFuture`                                                                                                                 |
| [Comparable values](https://insei.github.io/valigo/latest/features/values/index.md) | `eq.Of`, `eq.Ptr`, `eq.Slice`: enums, bools — `OneOf`, `NotOneOf`, `Equal`, `NotEqual`, `Unique`                                                                                           |
| [UUID](https://insei.github.io/valigo/latest/features/uuid/index.md)                | `uuid.Of` of the separate module: `Required`, `AnyOf`                                                                                                                                      |
| [Slices](https://insei.github.io/valigo/latest/features/slices/index.md)            | `Slice` and the slice conversions: length rules, `Each` for the elements                                                                                                                   |
| [Maps](https://insei.github.io/valigo/latest/features/maps/index.md)                | `Map`: length rules, `EachKey` and `EachValue` with conversions                                                                                                                            |

## Errors

What a validation returns.

| Page                                                                                    | What's inside                                                                                          |
| --------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| [Errors](https://insei.github.io/valigo/latest/features/errors/index.md)                | `Validate` and `ValidateTyped`, `shared.Error`, locations, codes for clients, errors transformer       |
| [Translations](https://insei.github.io/valigo/latest/features/translations/index.md)    | built-in locales, plural forms, own messages, languages from the context, `Accept-Language` middleware |
| [HTTP problem details](https://insei.github.io/valigo/latest/features/problem/index.md) | `problem`: the errors as `application/problem+json` of RFC 9457 with JSON pointers and codes           |

## Integrations

The separate modules for the frameworks.

| Page                                                                 | What's inside                                                                                                       |
| -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- |
| [huma](https://insei.github.io/valigo/latest/features/huma/index.md) | `github.com/insei/valigo/huma`: the rules in the `Resolve` of huma inputs with the locations of huma                |
| [gin](https://insei.github.io/valigo/latest/features/gin/index.md)   | `github.com/insei/valigo/gin`: the `binding.StructValidator` of gin and validation with the language of the request |
| [echo](https://insei.github.io/valigo/latest/features/echo/index.md) | `github.com/insei/valigo/echo`: the `echo.Validator` and validation with the language of the request                |
