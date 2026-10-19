# Rules from tags

`b.Tags("validate")` configures the rules of the tags of the fields of the struct and of its nested structs in the syntax of [go-playground/validator](https://github.com/go-playground/validator): `validate:"required,max=64"`. It is for a migration: the structs keep their tags, and the rules move to code one by one — the tags and the rules in code are mixed in one `Configure`.

The rules of the tags are the rules of valigo: the same messages in the languages of the users, the [codes](https://insei.github.io/valigo/dev/features/errors/#codes-for-clients), the [JSON Schema](https://insei.github.io/valigo/dev/features/schema/index.md) and the speed — a valid object is validated without allocations, see [Performance](https://insei.github.io/valigo/dev/performance/index.md).

## Example

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/str"
)

type Address struct {
    City string `validate:"required,max=32"`
}

type SignUp struct {
    Email    string   `validate:"required,email"`
    Password string   `validate:"required,min=8"`
    Confirm  string   `validate:"eqfield=Password"`
    Age      int      `validate:"omitempty,gte=18,lte=120"`
    Tags     []string `validate:"max=3,dive,required,max=16"`
    Address  *Address // the rules of City apply when the address is present
    Nick     string
}

func main() {
    v := valigo.New()
    valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
        b.Tags("validate")
        // the rules in code with the tags
        b.Field(&s.Nick).As(str.Of).Trim().Alphanumeric()
    })

    errs := v.ValidateTyped(context.Background(), &SignUp{
        Email: "ann", Password: "secret", Confirm: "secrets", Age: 16, Tags: []string{"go", ""}, Nick: " a b ",
    })
    for _, err := range errs {
        fmt.Println(err.Location, "-", err.Message, "-", err.Code)
    }
}
```

Output:

```text
Email - Should be email address - string.email
Password - Cannot be shorter than 8 characters - string.min_runes
Confirm - Should be equal to Password - string.equal_field
Age - Cannot be less than 18 - number.min
Tags[1] - Should be fulfilled - string.required
Nick - Should contain only letters and digits - string.alphanumeric
```

## Rules

| Fields                     | Tags of go-playground/validator                                                                                                           | Rules of valigo                                                                                                                                                                     |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| all                        | `required`                                                                                                                                | `Required`: a string is not empty, a number is not zero, a bool is `true`, a slice and a map have elements, a pointer is not nil                                                    |
| all                        | `omitempty`, `omitnil`                                                                                                                    | the next rules are skipped for the zero value; a nil pointer is absent anyway                                                                                                       |
| all                        | `required_if=F v`, `required_unless=F v`, `required_with=F G`, `required_with_all`, `required_without`, `required_without_all`            | `Required` under the condition of the fields `F`, `G` of the same struct: their values, or their values are not zero                                                                |
| all                        | `a\|b`                                                                                                                                    | `Or`: one of the alternatives passes, the errors of the first one when none does                                                                                                    |
| strings                    | `min`, `max`, `len`, `gt`, `gte`, `lt`, `lte`                                                                                             | `MinRunes`, `MaxRunes`: the lengths in characters                                                                                                                                   |
| strings                    | `eq`, `ne`, `oneof=a b`, `eq_ignore_case`                                                                                                 | `AnyOf`, `eq.NotEqual`, `AnyOf`, `EqualFold`                                                                                                                                        |
| strings                    | `contains`, `excludes`, `startswith`, `endswith`, `startsnotwith`, `endsnotwith`, `containsany`, `excludesall`                            | `Contains`, `Excludes`, `Prefix`, `Suffix`, `NotPrefix`, `NotSuffix`, `ContainsAny`, `ExcludesAny`                                                                                  |
| strings                    | `alpha`, `alphanum`, `number`, `numeric`, `hexadecimal`, `ascii`, `lowercase`, `uppercase`                                                | `Alpha`, `Alphanumeric`, `Digits`, `Numeric`, `Hexadecimal`, `ASCII`, `Lowercase`, `Uppercase`                                                                                      |
| strings                    | `email`, `url`, `http_url`, `uuid`, `ip`, `ipv4`, `ipv6`, `cidr`, `mac`, `hostname`, `hostname_rfc1123`, `e164`                           | `Email`, `URL`, `URL`, `UUID`, `IP`, `IPv4`, `IPv6`, `CIDR`, `MAC`, `Hostname`, `Hostname`, `E164`                                                                                  |
| strings                    | `base64`, `base64url`, `json`, `semver`, `hexcolor`, `datetime=2006-01-02`, `iso3166_1_alpha2`, `iso4217`, `credit_card`, `luhn_checksum` | `Base64`, `Base64URL`, `JSON`, `Semver`, `HexColor`, `Datetime`, `CountryCode`, `CurrencyCode`, `CreditCard`, `Luhn`                                                                |
| strings                    | `eqfield`, `nefield`                                                                                                                      | `EqualField`, `NotEqualField`                                                                                                                                                       |
| numbers                    | `min`, `gte`, `max`, `lte`, `gt`, `lt`, `eq`, `ne`, `oneof=1 2`                                                                           | `Min`, `Max`, `GreaterThan`, `LessThan`, `AnyOf`, `eq.NotEqual`, `AnyOf`                                                                                                            |
| numbers                    | `gtefield`, `ltefield`, `gtfield`, `ltfield`, `eqfield`, `nefield`                                                                        | `MinField`, `MaxField`, `GreaterField`, `LessField`, `eq.EqualField`, `eq.NotEqualField`                                                                                            |
| numbers                    | `latitude`, `longitude`                                                                                                                   | `Latitude`, `Longitude`                                                                                                                                                             |
| `bool`                     | `eq=true`, `ne=true`                                                                                                                      | `eq.Equal`, `eq.NotEqual`                                                                                                                                                           |
| `time.Time`                | `gt`, `gte`, `lt`, `lte` without a parameter, `gtfield`, `ltfield`                                                                        | `InFuture`, `InPast`, `AfterField`, `BeforeField`                                                                                                                                   |
| slices, maps               | `min`, `max`, `len`, `gt`, `gte`, `lt`, `lte`, `eq`                                                                                       | the rules of the length                                                                                                                                                             |
| slices                     | `unique`                                                                                                                                  | `eq.Slice(...).Unique()` of the elements of comparable types                                                                                                                        |
| slices                     | `dive`                                                                                                                                    | the next rules are the rules of the elements; the elements of struct types are validated by the tags of their type, `Items[1].Title`, `dive,required` checks nil elements of `[]*T` |
| maps                       | `dive`, `dive,keys,...,endkeys,...`                                                                                                       | the rules of the values and of the keys of `map[string]V`                                                                                                                           |
| pointers to nested structs | `required`                                                                                                                                | the fields of a nested struct have their own tags                                                                                                                                   |

Pointers to the values have the rules of the values, a nil pointer is absent. The rules on fields and the conditions take the names of the fields of the same struct, or their nested fields `Inner.Field`. `0x2C` in a parameter is a comma, `0x7C` is a pipe, like in go-playground/validator.

## Not supported

These tags panic in `Configure`: a typo or an unsupported rule doesn't turn into a missing check. Write their rules in code.

| Tags                                                                                                                                                                                                                                                                                                            | In code                                                                                                              |
| --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| `excluded_if`, `excluded_unless`, `excluded_with`, `excluded_with_all`, `excluded_without`, `excluded_without_all`                                                                                                                                                                                              | `b.When(cond).Field(&m.X).Custom(...)` that fails on a present value                                                 |
| `eqcsfield`, `necsfield`, `gtcsfield`, `gtecsfield`, `ltcsfield`, `ltecsfield`, `fieldcontains`, `fieldexcludes`                                                                                                                                                                                                | the rules on other fields with the pointers to any field of the model, `EqualField(&m.Other.Field)`, or `b.Custom`   |
| `alphaunicode`, `alphanumunicode`, `boolean`, `containsrune`, `excludesrune`, `multibyte`, `printascii`, `ne_ignore_case`, `oneofci`                                                                                                                                                                            | `Custom` or `Regexp`                                                                                                 |
| `uuid3`, `uuid4`, `uuid5`, `uuid_rfc4122`... , `ulid`, `cve`                                                                                                                                                                                                                                                    | `UUID` with `Custom` of the version, `Regexp`                                                                        |
| `uri`, `urn_rfc2141`, `url_encoded`, `datauri`, `fqdn`, `hostname_port`, `tcp_addr`, `udp_addr`, `unix_addr`, `ip_addr`... , `cidrv4`, `cidrv6`                                                                                                                                                                 | `URL`, `Hostname`, `IP`, `CIDR` with `Custom`, or `Custom` with the parsers of the standard library                  |
| `isbn`, `issn`, `bic`, `ssn`, `jwt`, `html`, `html_encoded`, `rgb`, `rgba`, `hsl`, `hsla`, `iscolor`, `md5`, `sha256`..., `btc_addr`, `eth_addr`, `mongodb`, `cron`, `timezone`, `bcp47_language_tag`, `iso3166_1_alpha3`, `iso3166_1_alpha_numeric`, `iso3166_2`, `iso4217_numeric`, `postcode_iso3166_alpha2` | `Custom` or `Regexp`                                                                                                 |
| `dir`, `dirpath`, `file`, `filepath`, `image`                                                                                                                                                                                                                                                                   | `Custom`: valigo doesn't touch the file system                                                                       |
| `isdefault`, `omitzero`, `structonly`, `nostructlevel`, `skip_unless`                                                                                                                                                                                                                                           | `When` and `Custom`                                                                                                  |
| `dive` into the elements of slices of slices, of maps with other keys than strings, of maps of pointers and structs; `unique` of the slices of pointers and structs, `unique=Field`                                                                                                                             | `Each`, `EachKey`, `EachValue` and `Custom` in code                                                                  |
| the custom tags of `RegisterValidation`, the aliases of `RegisterAlias`                                                                                                                                                                                                                                         | the rules of the type in code, see [Your types](https://insei.github.io/valigo/dev/features/conversions/#your-types) |

## Differences from go-playground/validator

- `required` is `Required` of valigo: a string is not empty, a number is not zero, a slice and a map have elements — go-playground/validator accepts an empty non-nil slice and map, and a pointer to an empty string.
- The fields of a nil nested struct are skipped, like by go-playground/validator. The rules in code check `Required` of these fields, see [Absent values](https://insei.github.io/valigo/dev/features/absent/index.md).
- The messages are the messages of valigo in the languages of the users, the errors have the [codes](https://insei.github.io/valigo/dev/features/errors/#codes-for-clients) of valigo; the alternatives `a|b` report the errors of the first one.
- The formats are the formats of valigo, see [Strings](https://insei.github.io/valigo/dev/features/strings/index.md): `email`, `url` and `hostname` can differ from go-playground/validator in details.
- `gt`, `lt` of `time.Time` compare with the current time without a parameter, like go-playground/validator, and fail on a parameter.
- The rules of the tags after `omitempty` and the conditional ones aren't in the [JSON Schema](https://insei.github.io/valigo/dev/features/schema/index.md): they don't apply to every value.
- A valid object is validated without allocations, also with `dive` into the elements of structs; the alternatives allocate the errors of the failed ones.
