# JSON Schema

The rules that JSON Schema can express are described by its keywords, so the generators of the schemas of an API can show them to the clients: [huma](huma.md#the-rules-in-the-openapi-document) puts them into its OpenAPI document. `v.Constraints()` returns them for other generators.

```go
for _, c := range v.Constraints() {
	fmt.Println(c.Struct, c.Field.GetStructPath(), c.Path, c.Keyword, c.Value)
}
```

Every `shared.Constraint` has the configured struct, the field, the path from the schema of the field to the constrained value — `[items]` for the elements of a slice, `[additionalProperties]` for the values of a map, `[propertyNames]` for its keys — the keyword and its value. The keyword `required` means the field is in the required properties of its object.

| Rules | Keywords |
|---|---|
| `str`: `Required` | `required`, `minLength: 1` |
| `str`: `MinRunes`, `MaxRunes`, `MaxLen` | `minLength`, `maxLength` |
| `str`: `Regexp`, `Prefix`, `Suffix`, `Contains`, `Alphanumeric`, `Alpha`, `Digits`, `Numeric`, `Hexadecimal`, `Semver`, `E164`, `HexColor` | `pattern` |
| `str`: `Email`, `URL`, `IPv4`, `IPv6`, `Hostname`, `UUID` | `format`: `email`, `uri`, `ipv4`, `ipv6`, `hostname`, `uuid` |
| `str`: `Datetime(time.RFC3339)`, `Datetime(time.DateOnly)` | `format`: `date-time`, `date` |
| `str`: `Base64`, `Base64URL` | `contentEncoding` |
| `str`, `num`, `uuid`: `AnyOf`; `eq`: `OneOf` | `enum` |
| `eq`: `Equal` | `const` |
| `num`: `Min`, `Max`, `AnyOfInterval`, `Latitude`, `Longitude`, `GreaterThan`, `LessThan` | `minimum`, `maximum`, `exclusiveMinimum`, `exclusiveMaximum` |
| `num`, `eq`, `uuid`: `Required` | `required` |
| slices: `Required`, `MinLen`, `MaxLen`, `Unique`, `Each` | `required` and `minItems: 1`, `minItems`, `maxItems`, `uniqueItems`, `items` |
| maps: `Required`, `MinLen`, `MaxLen`, `EachKey`, `EachValue` | `required` and `minProperties: 1`, `minProperties`, `maxProperties`, `propertyNames`, `additionalProperties` |

The schema describes the rules, valigo validates the values: the descriptions are chosen not to be stricter than the rules, so a value valid for valigo is valid for the schema. The formats are annotations of JSON Schema, the validators that check them can differ in details.

- `MinLen` of strings counts bytes, `minLength` counts characters: it isn't described. `MaxLen` is `maxLength`, a string of `n` bytes has no more than `n` characters.
- The rules under `When` of the struct and of fields aren't described, they don't apply to every value.
- After `Trim` and other `Modify` rules only `Required` is described: the schema describes the value of the request, the rules check the modified one.
- Rules on other fields, custom rules and the formats without a keyword of JSON Schema aren't described.

## Your conversions

A [conversion](conversions.md#your-types) describes its rules by `Schema` of the field, `Schema` of a slice configurator describes the rules of the slice:

```go
// Percent checks that the number is in [0, 100].
func (c *RatioConfigurator) Percent() *RatioConfigurator {
	c.f.Schema("minimum", 0).Schema("maximum", 100)
	c.f.CheckMsg(func(v Ratio) bool { return v >= 0 && v <= 100 }, shared.M("ratio.percent", percentKey))
	return c
}
```
