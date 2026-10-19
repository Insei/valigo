# Features

The map of valigo, in the order of the menu.

## Configuration

How the rules are configured and composed.

| Page | What's inside |
|---|---|
| [Validator and Configure](configure.md) | `New` and its options, `Configure`, the model, field pointers, nested and recursive structs |
| [Types and conversions](conversions.md) | `Field`, `Slice`, `Map` and `As`, the built-in conversions, conversions of your types |
| [Absent values](absent.md) | nil pointers, nil slices and nil nested structs: which rules skip them |
| [Conditions](conditions.md) | `When` on the struct and on a field, conditional requiredness |
| [Custom rules](custom.md) | `Custom` on a field, on a slice and on the struct, rules with limits from data |
| [Rules on other fields](fields.md) | `EqualField`, `MinField`, `MaxField`, `BeforeField`, `AfterField` by the pointer of the other field |

## Rules

The built-in rules of every type.

| Page | What's inside |
|---|---|
| [Strings](strings.md) | `str.Of`: lengths, `AnyOf`, substrings, `Regexp`, `Email`, `URL`, `IP`, `CIDR`, `MAC`, `Hostname`, `UUID`, `E164`, `Base64`, `JSON`, `Semver`, `Datetime`, ISO codes, `CreditCard`, `Trim` |
| [Numbers](numbers.md) | `num.Of`, `num.Ptr`: `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval`, `Latitude`, `Longitude` |
| [Time](time.md) | `timex.Of`: `Required`, `Before`, `After`, `Between`, `InPast`, `InFuture` |
| [Comparable values](values.md) | `eq.Of`, `eq.Ptr`, `eq.Slice`: enums, bools — `OneOf`, `NotOneOf`, `Equal`, `NotEqual`, `Unique` |
| [UUID](uuid.md) | `uuid.Of` of the separate module: `Required`, `AnyOf` |
| [Slices](slices.md) | `Slice` and the slice conversions: length rules, `Each` for the elements |
| [Maps](maps.md) | `Map`: length rules, `EachKey` and `EachValue` with conversions |

## Errors

What a validation returns.

| Page | What's inside |
|---|---|
| [Errors](errors.md) | `Validate` and `ValidateTyped`, `shared.Error`, locations, codes for clients, errors transformer |
| [Translations](translations.md) | built-in locales, plural forms, own messages, languages from the context, `Accept-Language` middleware |
| [HTTP problem details](problem.md) | `problem`: the errors as `application/problem+json` of RFC 9457 with JSON pointers and codes |

## Integrations

The separate modules for the frameworks.

| Page | What's inside |
|---|---|
| [JSON Schema](schema.md) | the rules as the keywords of JSON Schema for the generators of API schemas, `Constraints` |
| [huma](huma.md) | `github.com/insei/valigo/huma`: the rules in the `Resolve` of huma inputs with the locations of huma and in its OpenAPI document |
| [gin](gin.md) | `github.com/insei/valigo/gin`: the `binding.StructValidator` of gin and validation with the language of the request |
| [echo](echo.md) | `github.com/insei/valigo/echo`: the `echo.Validator` and validation with the language of the request |
