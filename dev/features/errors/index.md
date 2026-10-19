# Errors

## Validate and ValidateTyped

```go
func (v *Validator) Validate(ctx context.Context, obj any) []error
func (v *Validator) ValidateTyped(ctx context.Context, obj any) []shared.Error
```

- `ValidateTyped` returns the errors as `shared.Error` values, `Validate` passes them through the [errors transformer](#errors-transformer), by default every `shared.Error` becomes an `error`.
- Both return `nil` when the object is valid, when the type of the object has no rules, and for a nil pointer.
- `obj` is a pointer to the configured struct. A struct value is validated too, by a shallow copy: [`Trim`](https://insei.github.io/valigo/latest/features/strings/#trim) of its fields isn't visible to the caller, but the values behind its pointers and the elements of its slices are shared with the caller and can be changed.
- The errors go in the order of the configuration: the rules of a field in the order of the chain, the fields in the order of their first rule.

## shared.Error

```go
type Error struct {
    Message  string         // translated message
    Location string         // location of the field, Address.City or Tags[1] by default
    Value    any            // the checked value, nil for absent values
    Code     string         // stable code of the rule, e.g. string.max_len
    Params   map[string]any // named parameters of the message, e.g. {"max": 64}
}
```

`Error()` returns `Message (Location: Value)`, or `Message` alone when both the location and the value are empty. `Error` has JSON field names of its own: `{"Message":"...","Location":"...","Value":...,"Code":"...","Params":{...}}`, `Code` and `Params` are omitted when empty. Convert it to your API format by the [transformer](#errors-transformer).

## Codes for clients

The message is translated on the backend, see [Translations](https://insei.github.io/valigo/latest/features/translations/index.md). A client that translates the messages itself — a web or a mobile application with its own locales — uses `Code` and `Params` instead: the code is stable, the parameters are named.

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/insei/valigo"
    "github.com/insei/valigo/num"
    "github.com/insei/valigo/str"
    "github.com/insei/valigo/timex"
)

type Event struct {
    Title    string
    Seats    int
    StartsAt time.Time
    EndsAt   time.Time
}

func main() {
    v := valigo.New()
    valigo.Configure[Event](v, func(b valigo.Configurator[Event], e *Event) {
        b.Field(&e.Title).As(str.Of).MaxLen(5)
        b.Field(&e.Seats).As(num.Of).AnyOfInterval(1, 100)
        b.Field(&e.EndsAt).As(timex.Of).AfterField(&e.StartsAt)
    })

    start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
    for _, err := range v.ValidateTyped(context.Background(), &Event{Title: "Conference", StartsAt: start, EndsAt: start}) {
        body, _ := json.Marshal(map[string]any{"field": err.Location, "code": err.Code, "params": err.Params})
        fmt.Println(string(body))
    }
}
```

Output:

```text
{"code":"string.max_len","field":"Title","params":{"max":5}}
{"code":"number.interval","field":"Seats","params":{"max":100,"min":1}}
{"code":"time.after_field","field":"EndsAt","params":{"field":"StartsAt"}}
```

| Code                                                                                                                                                                                                                                                                            | Params                          |
| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| `string.required`, `string.email`, `string.url`, `string.ip`, `string.ipv4`, `string.ipv6`, `string.hostname`, `string.uuid`, `string.alphanumeric`, `string.regexp`                                                                                                            | —                               |
| `string.min_len`, `string.min_runes`                                                                                                                                                                                                                                            | `min`                           |
| `string.max_len`, `string.max_runes`                                                                                                                                                                                                                                            | `max`                           |
| `string.any_of`                                                                                                                                                                                                                                                                 | `values`                        |
| `string.prefix`, `string.suffix`, `string.contains`                                                                                                                                                                                                                             | `prefix`, `suffix`, `substring` |
| `string.equal_field`                                                                                                                                                                                                                                                            | `field`                         |
| `string.base64`, `string.base64url`, `string.json`, `string.semver`, `string.cidr`, `string.mac`, `string.e164`, `string.hex_color`, `string.ascii`, `string.lowercase`, `string.uppercase`, `string.country_code`, `string.currency_code`, `string.luhn`, `string.credit_card` | —                               |
| `string.datetime`                                                                                                                                                                                                                                                               | `layout`                        |
| `number.required`                                                                                                                                                                                                                                                               | —                               |
| `number.min`, `number.max`                                                                                                                                                                                                                                                      | `min`, `max`                    |
| `number.min_field`, `number.max_field`                                                                                                                                                                                                                                          | `field`                         |
| `number.any_of`                                                                                                                                                                                                                                                                 | `values`                        |
| `number.interval`                                                                                                                                                                                                                                                               | `min`, `max`                    |
| `number.latitude`, `number.longitude`                                                                                                                                                                                                                                           | —                               |
| `time.required`, `time.in_past`, `time.in_future`                                                                                                                                                                                                                               | —                               |
| `time.before`, `time.after`                                                                                                                                                                                                                                                     | `time` (RFC 3339)               |
| `time.between`                                                                                                                                                                                                                                                                  | `from`, `to` (RFC 3339)         |
| `time.before_field`, `time.after_field`                                                                                                                                                                                                                                         | `field`                         |
| `value.required`, `value.unique`                                                                                                                                                                                                                                                | —                               |
| `value.one_of`, `value.not_one_of`                                                                                                                                                                                                                                              | `values`                        |
| `value.equal`, `value.not_equal`                                                                                                                                                                                                                                                | `value`                         |
| `value.equal_field`                                                                                                                                                                                                                                                             | `field`                         |
| `uuid.required`                                                                                                                                                                                                                                                                 | —                               |
| `uuid.any_of`                                                                                                                                                                                                                                                                   | `values`                        |
| `slice.required`, `map.required`                                                                                                                                                                                                                                                | —                               |
| `slice.min_len`, `map.min_len`                                                                                                                                                                                                                                                  | `min`                           |
| `slice.max_len`, `map.max_len`                                                                                                                                                                                                                                                  | `max`                           |

The `field` parameter is the location of the other field by the location function of the validator. The errors of [custom rules](https://insei.github.io/valigo/latest/features/custom/index.md) have the locale key as the code, or the code of `shared.M` with `ErrorMsg`.

## Location of fields

The location is the path of the field in the struct by default: `Address.City`. `WithFieldLocationNamingFn` replaces it, e.g. with the path of the JSON tags, so the client sees the names of its request:

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/fmap/v3"
    "github.com/insei/valigo"
    "github.com/insei/valigo/str"
)

type Address struct {
    City string `json:"city"`
}

type SignUp struct {
    Email   string   `json:"email"`
    Address Address  `json:"address"`
    Tags    []string `json:"tags"`
}

func main() {
    v := valigo.New(valigo.WithFieldLocationNamingFn(func(field fmap.Field) string {
        return field.GetTagPath("json", false)
    }))
    valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
        b.Field(&s.Email).As(str.Of).Email()
        b.Field(&s.Address.City).As(str.Of).Required()
        b.Slice(&s.Tags).As(str.Slice).Email()
    })

    for _, err := range v.ValidateTyped(context.Background(), &SignUp{Email: "a", Tags: []string{"b"}}) {
        fmt.Println(err.Location, err.Message)
    }
}
```

Output:

```text
email Should be email address
address.city Should be fulfilled
tags[0] Should be email address
```

The index of a slice element is appended to the location returned by the function.

## Errors transformer

`WithErrorsTransformer` sets how `Validate` converts the errors, e.g. to one error of the API:

```go
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type ValidationError struct {
    Fields []FieldError `json:"fields"`
}

func (e *ValidationError) Error() string { return "validation failed" }

v := valigo.New(valigo.WithErrorsTransformer(func(errs []shared.Error) []error {
    fields := make([]FieldError, len(errs))
    for i, err := range errs {
        fields[i] = FieldError{Field: err.Location, Message: err.Message}
    }
    return []error{&ValidationError{Fields: fields}}
}))
```

The transformer is called only when there are errors. For an HTTP API the [`problem`](https://insei.github.io/valigo/latest/features/problem/index.md) package has the transformer of the problem details of RFC 9457.
