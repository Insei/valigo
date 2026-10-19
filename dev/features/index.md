# Features

The map of valigo, in the order of the menu.

## Configuration

How the rules are configured and composed.

| Page                                                                                      | What's inside                                                                                       |
| ----------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| [Validator and Configure](https://insei.github.io/valigo/dev/features/configure/index.md) | `New` and its options, `Configure`, the model, field pointers, nested and recursive structs         |
| [Types and conversions](https://insei.github.io/valigo/dev/features/conversions/index.md) | `Field`, `Slice`, `Map` and `As`, the built-in conversions, conversions of your types               |
| [Absent values](https://insei.github.io/valigo/dev/features/absent/index.md)              | nil pointers, nil slices and nil nested structs: which rules skip them                              |
| [Conditions](https://insei.github.io/valigo/dev/features/conditions/index.md)             | `When` on the struct and on a field, conditional requiredness                                       |
| [Custom rules](https://insei.github.io/valigo/dev/features/custom/index.md)               | `Custom` on a field, on a slice and on the struct, rules with limits from data                      |
| [Rules on other fields](https://insei.github.io/valigo/dev/features/fields/index.md)      | `EqualField`, `MinField`, `MaxField`, `BeforeField`, `AfterField` by the pointer of the other field |
| [Rules from tags](https://insei.github.io/valigo/dev/features/tags/index.md)              | `Tags`: the tags of go-playground/validator, `validate:"required,max=64"`, for a migration          |

## Rules

The built-in rules of every type.

| Page                                                                             | What's inside                                                                                                                                                                              |
| -------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [Strings](https://insei.github.io/valigo/dev/features/strings/index.md)          | `str.Of`: lengths, `AnyOf`, substrings, `Regexp`, `Email`, `URL`, `IP`, `CIDR`, `MAC`, `Hostname`, `UUID`, `E164`, `Base64`, `JSON`, `Semver`, `Datetime`, ISO codes, `CreditCard`, `Trim` |
| [Numbers](https://insei.github.io/valigo/dev/features/numbers/index.md)          | `num.Of`, `num.Ptr`: `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval`, `Latitude`, `Longitude`                                                                                           |
| [Time](https://insei.github.io/valigo/dev/features/time/index.md)                | `timex.Of`: `Required`, `Before`, `After`, `Between`, `InPast`, `InFuture`                                                                                                                 |
| [Comparable values](https://insei.github.io/valigo/dev/features/values/index.md) | `eq.Of`, `eq.Ptr`, `eq.Slice`: enums, bools — `OneOf`, `NotOneOf`, `Equal`, `NotEqual`, `Unique`                                                                                           |
| [UUID](https://insei.github.io/valigo/dev/features/uuid/index.md)                | `uuid.Of` of the separate module: `Required`, `AnyOf`                                                                                                                                      |
| [Slices](https://insei.github.io/valigo/dev/features/slices/index.md)            | `Slice` and the slice conversions: length rules, `Each` for the elements                                                                                                                   |
| [Maps](https://insei.github.io/valigo/dev/features/maps/index.md)                | `Map`: length rules, `EachKey` and `EachValue` with conversions                                                                                                                            |

## Errors

What a validation returns.

| Page                                                                                 | What's inside                                                                                          |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------ |
| [Errors](https://insei.github.io/valigo/dev/features/errors/index.md)                | `Validate` and `ValidateTyped`, `shared.Error`, locations, codes for clients, errors transformer       |
| [Translations](https://insei.github.io/valigo/dev/features/translations/index.md)    | built-in locales, plural forms, own messages, languages from the context, `Accept-Language` middleware |
| [HTTP problem details](https://insei.github.io/valigo/dev/features/problem/index.md) | `problem`: the errors as `application/problem+json` of RFC 9457 with JSON pointers and codes           |

## Integrations

The separate modules for the frameworks.

| Page                                                                       | What's inside                                                                                                                                 |
| -------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| [JSON Schema](https://insei.github.io/valigo/dev/features/schema/index.md) | the rules as the keywords of JSON Schema for the generators of API schemas, `Constraints`                                                     |
| [huma](https://insei.github.io/valigo/dev/features/huma/index.md)          | `github.com/insei/valigo/integrations/huma`: the rules in the `Resolve` of huma inputs with the locations of huma and in its OpenAPI document |
| [gin](https://insei.github.io/valigo/dev/features/gin/index.md)            | `github.com/insei/valigo/integrations/gin`: the `binding.StructValidator` of gin and validation with the language of the request              |
| [go-i18n](https://insei.github.io/valigo/dev/features/goi18n/index.md)     | `github.com/insei/valigo/integrations/goi18n`: the messages from a bundle of go-i18n, the built-in ones for the rest                          |
| [echo](https://insei.github.io/valigo/dev/features/echo/index.md)          | `github.com/insei/valigo/integrations/echo`: the `echo.Validator` and validation with the language of the request                             |
