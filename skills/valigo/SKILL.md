---
name: valigo
description: Write Go code that validates structs with github.com/insei/valigo — typed rules configured by pointers to the fields of a model and converted to the configurators of their types (b.Field(&m.Name).As(str.Of), num.Of, num.Ptr, b.Slice(&m.Tags).As(str.Slice), uuid.Of of github.com/insei/valigo/uuid), Each for slice elements, When conditions, Custom rules of fields and structs, errors with field locations, translated messages and the Accept-Language middleware. Use when the code imports github.com/insei/valigo, or when a Go service validates request or config structs.
---

# valigo

Rules are configured once by `valigo.Configure[T]` with pointers to the fields of a model and are checked by the compiler. `Validate` runs them on objects of type `*T` and doesn't allocate on valid objects. Full documentation: https://insei.github.io/valigo/latest/, as one file for the context: https://insei.github.io/valigo/latest/llms-full.txt.

## Rules

- Go 1.27 or newer: the configurator methods are generic.
- Create one validator with `valigo.New(...)` on start and call `valigo.Configure[T]` once per type, before `Validate`. `Configure` panics on invalid configuration, it isn't safe concurrently with `Validate`, a repeated `Configure` of the same type replaces its rules.
- Field pointers must point to the fields of the model passed to the function: `func(b valigo.Configurator[T], m *T) { b.Field(&m.Name).As(str.Of) }`. A pointer to anything else panics. Nested structs stored by pointer are allocated in the model, `&m.Address.City` is valid.
- The configurator has only `Field`, `Slice`, `Map`, `When` and `Custom`. Convert a field to the configurator of its type by `As`: `b.Field(&m.Name).As(str.Of)` for `string`, `*string` and named strings; `.As(num.Of)` for numbers; `.As(num.Ptr)` for pointers to numbers; `.As(vuuid.Of)` for `uuid.UUID` and `*uuid.UUID` with `vuuid "github.com/insei/valigo/uuid"`, a separate module. Slices: `b.Slice(&m.Tags).As(str.Slice)`, `.As(vuuid.Slice)`; a pointer to a slice: `b.Field(&m.Tags).As(shared.SlicePtr).As(str.Slice)`. There are no `b.String`, `b.Number`, `b.UUID`, `b.StringSlice` methods.
- Conditional requiredness (`required_if`, `required_with`, `required_without` of tags) is `b.When(func(ctx, m *T) bool { return m.Kind == "company" }).Field(&m.TaxID).As(str.Of).Required()`.
- Rules on other fields take a pointer to the other field of the same model and type: `b.Field(&m.EndsAt).As(timex.Of).AfterField(&m.StartsAt)`, `num` `MinField`/`MaxField`, `str`/`eq` `EqualField`; the message has the location of the other field.
- A conversion that doesn't fit the field type doesn't compile: `b.Field(&m.Age).As(str.Of)` for an `int`, `.As(num.Of)` for a `*int` (use `num.Ptr`).
- Rule arguments have the type of the field: `b.Field(&m.ID).As(num.Of).Max(5)` for an `int64` field, never `Max(int64(5))` with a value of another type or `any`.
- `Custom` and `When` of a field get the typed value: `func(ctx context.Context, h *shared.FieldCustomHelper, value string) []shared.Error`, `func(ctx context.Context, value int) bool`. Of slices: `shared.Elems[E]` with `Len()`, `At(i)`, `All()`.
- `MaxLen`, `MinLen` and `Required` of slice configurators check the slice itself, rules on elements are set by `Each(func(e *str.BaseConfigurator) { ... })`.
- Absent values (nil pointers, nil `*[]T`, fields of nil nested structs, nil elements of `[]*T`) are skipped by every rule except `Required`.
- Rules run in the order of the chain, `Trim()` changes the field of the validated object.
- `Validate` returns `[]error`, `ValidateTyped` returns `[]shared.Error{Message, Location, Value, Code, Params}`, both are `nil` for valid objects. `Message` is translated on the backend; `Code` (`string.max_len`, `number.interval`, `time.after_field`...) and `Params` (`{"max": 64}`) are stable, for the clients that translate themselves. Custom rules: `h.ErrorT(ctx, v, "message")` has the message as the code, `h.ErrorMsg(ctx, v, shared.M("my.code", "message %v", shared.P("name", x)))` sets a code and params.

| Configurator | Rules |
|---|---|
| `b.Field(&m.S).As(str.Of)` | `Required`, `MinLen`/`MaxLen` (bytes), `MinRunes`/`MaxRunes` (characters), `AnyOf`, `EqualField(&m.Other)`, `Prefix`, `Suffix`, `Contains`, `Regexp(re, str.WithRegexpLocaleKey(key))`, `Email`, `URL`, `IP`, `IPv4`, `IPv6`, `CIDR`, `MAC`, `Hostname`, `UUID`, `Alphanumeric`, `ASCII`, `Lowercase`, `Uppercase`, `Base64`, `Base64URL`, `JSON`, `Semver`, `E164`, `HexColor`, `Datetime(layout)`, `CountryCode`, `CurrencyCode`, `Luhn`, `CreditCard`, `Trim`, `Custom`, `When`; formats fail on empty strings |
| `b.Field(&m.N).As(num.Of)`, `b.Field(&m.P).As(num.Ptr)` | `Required` (zero for `num.Of`, nil for `num.Ptr`), `Min`, `Max`, `MinField(&m.Other)`, `MaxField(&m.Other)`, `AnyOf`, `AnyOfInterval` (inclusive), `Latitude`, `Longitude`, `Custom`, `When`; slices: `b.Slice(&m.Sizes).As(num.Slice)` |
| `b.Field(&m.At).As(timex.Of)` | `time.Time` and `*time.Time`: `Required` (not zero), `Before`, `After`, `Between` (inclusive), `BeforeField(&m.Other)`, `AfterField(&m.Other)`, `InPast`, `InFuture`, `Custom`, `When`; slices: `b.Slice(&m.Dates).As(timex.Slice)` |
| `b.Field(&m.Status).As(eq.Of)`, `.As(eq.Ptr)` | any `comparable` (enums of named types, `bool`): `Required` (not zero / not nil), `OneOf`, `NotOneOf`, `Equal`, `NotEqual`, `EqualField(&m.Other)`, `Custom`, `When`; `b.Slice(&m.Roles).As(eq.Slice)`: `Unique`, `OneOf`, `Each` |
| `b.Field(&m.ID).As(vuuid.Of)` | `Required` (not `uuid.Nil`), `AnyOf`, `Custom`, `When` |
| `b.Slice(&m.Items)` | `Required`, `MinLen`, `MaxLen`, `Custom`, `When`, `Each(func(*shared.Field[E]))` with `Check(ok, message)`, `Required(ok, message)`, `As` to convert the element: `e.As(num.Of).Min(1)` |
| `b.Slice(&m.Tags).As(str.Slice)` | the rules of the slice and `Each(func(*str.BaseConfigurator))`, `Trim`, `Regexp`, `Email` |
| `b.Slice(&m.IDs).As(vuuid.Slice)` | the rules of the slice and `Each(func(*vuuid.BaseConfigurator))`, `AnyOf` |
| `b.Map(&m.Labels)` | `Required`, `MinLen`, `MaxLen`, `Custom`, `When`, `EachKey(func(*shared.Field[K]))`, `EachValue(func(*shared.Field[V]))` with `As`: `k.As(str.Of)`; locations `Labels[key]`, errors sorted by location |
| `b.Field(&m.X)` | the base rules without a conversion: `Check(ok, message)`, `Required(ok, message)`, `Modify`, `Custom`, `When` |
| `b.When(func(ctx, *T) bool)` | a configurator whose rules apply only when the function returns true |
| `b.Custom(func(ctx, h shared.StructCustomHelper, obj *T) []shared.Error)` | a rule on several fields, `h.ErrorT(ctx, &obj.Field, value, message)` |

A conversion of your type is a function `func Of[T MyType | *MyType](f *shared.Field[T]) *MyConfigurator` over `shared.Deref` and `shared.Retype`, see https://insei.github.io/valigo/latest/features/conversions/.

## A request with typed rules

```go
package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/insei/valigo"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
	vuuid "github.com/insei/valigo/uuid"
)

type Address struct {
	City string
	Zip  string
}

type CreateOrder struct {
	CustomerID uuid.UUID
	Email      string
	Comment    *string
	Quantity   int
	Discount   *uint8
	Address    *Address
	Tags       []string
	Items      []int64
}

var zipRe = regexp.MustCompile(`^\d{5}$`)

func main() {
	v := valigo.New()
	valigo.Configure[CreateOrder](v, func(b valigo.Configurator[CreateOrder], o *CreateOrder) {
		b.Field(&o.CustomerID).As(vuuid.Of).Required()
		b.Field(&o.Email).As(str.Of).Trim().Required().Email()
		b.Field(&o.Comment).As(str.Of).MaxLen(500) // nil Comment is valid
		b.Field(&o.Quantity).As(num.Of).AnyOfInterval(1, 100)
		b.Field(&o.Discount).As(num.Ptr).Max(50) // nil Discount is valid
		b.Field(&o.Address.City).As(str.Of).Required()
		b.Field(&o.Address.Zip).As(str.Of).Regexp(zipRe, str.WithRegexpLocaleKey("Should be 5 digits"))
		b.Slice(&o.Tags).As(str.Slice).MaxLen(5).Each(func(tag *str.BaseConfigurator) {
			tag.Trim().Required().MaxLen(20)
		})
		b.Slice(&o.Items).Required().Each(func(item *shared.Field[int64]) {
			item.Check(func(id int64) bool { return id > 0 }, "Should be positive")
		})
	})

	order := &CreateOrder{
		Email:    " buyer@example.com ",
		Quantity: 0,
		Address:  &Address{City: "Berlin", Zip: "1"},
		Tags:     []string{"gift", " "},
		Items:    []int64{10, -1},
	}
	for _, err := range v.ValidateTyped(context.Background(), order) {
		fmt.Printf("%s: %s\n", err.Location, err.Message)
	}
	fmt.Printf("%q\n", order.Email) // trimmed in place
}
```

## Conditions, custom rules and errors for an API

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
)

type Notification struct {
	Channel string    `json:"channel"`
	Email   string    `json:"email"`
	Phone   string    `json:"phone"`
	SendAt  time.Time `json:"send_at"`
	Retries int       `json:"retries"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func main() {
	v := valigo.New(
		// locations are JSON paths: "email", "address.city", "tags[1]"
		valigo.WithFieldLocationNamingFn(func(field fmap.Field) string { return field.GetTagPath("json", false) }),
	)
	valigo.Configure[Notification](v, func(b valigo.Configurator[Notification], n *Notification) {
		b.Field(&n.Channel).As(str.Of).Required().AnyOf("email", "sms")

		email := b.When(func(_ context.Context, n *Notification) bool { return n.Channel == "email" })
		email.Field(&n.Email).As(str.Of).Required().Email()

		sms := b.When(func(_ context.Context, n *Notification) bool { return n.Channel == "sms" })
		sms.Field(&n.Phone).As(str.Of).Required()

		b.Field(&n.Retries).As(num.Of).When(func(_ context.Context, retries int) bool { return retries != 0 }).AnyOfInterval(1, 5)

		b.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *Notification) []shared.Error {
			if !obj.SendAt.IsZero() && obj.SendAt.Before(time.Now()) {
				return []shared.Error{h.ErrorT(ctx, &obj.SendAt, obj.SendAt, "Should be in the future")}
			}
			return nil
		})
	})

	n := &Notification{Channel: "email", Email: "nope", Retries: 9, SendAt: time.Now().Add(-time.Hour)}
	errs := v.ValidateTyped(context.Background(), n)
	fields := make([]FieldError, len(errs))
	for i, err := range errs {
		fields[i] = FieldError{Field: err.Location, Message: err.Message}
	}
	body, _ := json.Marshal(fields)
	fmt.Println(string(body))
}
```

## HTTP problem details

`valigo.New(valigo.WithErrorsTransformer(problem.Transformer()))` makes `Validate` return one `*problem.Details` of RFC 9457, write it by `problem.From(errs).Write(w)`: 422, `application/problem+json`, `errors` with `detail`, `pointer` (`#/tags/1`), `location`, `code`, `params`. Options: `problem.WithType`, `WithTitle`, `WithStatus`, `WithDetail`.

## JSON Schema

`v.Constraints()` returns the unconditional rules as JSON Schema keywords (`shared.Constraint`: Struct, Field, Path `items`/`additionalProperties`/`propertyNames`, Keyword, Value); rules under `When` aren't described, after `Trim`/`Modify` only `required`. `vhuma.OpenAPI(api, v)` returns the huma document with them in the body schemas (huma's own schemas stay unchanged). Own conversions describe rules by `f.Schema("minimum", 0)`.

## gin

The module `github.com/insei/valigo/gin` (`vgin`): `binding.Validator = vgin.New(v)` makes `c.ShouldBindJSON` validate by valigo (messages in the default language, gin passes no request), `vgin.Validate(c, v, &req)` validates with the languages of `Accept-Language`; the error is `shared.Errors`, get it by `errors.As`.

## echo

The module `github.com/insei/valigo/echo` (`vecho`): `e.Validator = vecho.New(v)` makes `c.Validate(&req)` validate by valigo (messages in the default language), `vecho.Validate(c, v, &req)` validates with the languages of `Accept-Language`; the error is `shared.Errors`.

## huma

The module `github.com/insei/valigo/huma` (`vhuma`): `valigo.New(valigo.WithFieldLocationNamingFn(vhuma.Location))` gives the locations of huma (`body.email`, `query.limit`, `path.id`), `func (i *Input) Resolve(ctx huma.Context) []error { return vhuma.Resolve(ctx, validator, i) }` validates the input with the languages of `Accept-Language`.

## Translations

The built-in messages are in English, German, Spanish, French and Russian, with plural forms. The language comes from the context: wrap the HTTP handler with `translator.NewAcceptLanguageMiddleware()`, or pass `translator.WithPreferredLanguagesFn(fn)` to `translator.New`. Messages of your rules are locale keys: a key without a translation is printed as is, translations are added to a storage:

```go
storage := translator.NewInMemStorage() // with the built-in locales
storage.Add("ru", map[string]string{"Should be in the future": "Должно быть в будущем"})
v := valigo.New(valigo.WithTranslator(translator.New(translator.WithStorage(storage))))
```

## Rules from data

Rules are configured once, but a `Custom` closure can read its limits from a database, a cache or the settings of a tenant on every validation, it gets the context of the request; a `When` condition too. Reuse a built-in message for the translation and the code: `h.ErrorMsg(ctx, v, shared.M("string.max_runes", "validation:string:Cannot be longer than %d characters", shared.P("max", maxLen)))`.

## Common mistakes

- `b.Field(&m.Score).As(num.Of)` for `Score *int` doesn't compile: use `b.Field(&m.Score).As(num.Ptr)`.
- `b.String(&m.Name)`, `b.Number(...)`, `b.UUID(...)`, `b.StringSlice(...)` are the old API, they don't exist: use `b.Field(...).As(...)` and `b.Slice(...).As(...)`.
- `b.Field(&local).As(str.Of)` with a variable outside the model panics in `Configure`: use the model passed to the function.
- `Configure` in a request handler: configure once on start, `Configure` replaces the rules and isn't safe concurrently with `Validate`.
- `b.Slice(&m.Tags).As(str.Slice).MaxLen(10)` limits the number of tags, not their length: use `Each(func(t *str.BaseConfigurator) { t.MaxLen(10) })`.
- `Required` after a field `When`: the condition isn't called for an absent value and the rules after it, including `Required`, are skipped. Put `Required` before `When`, or use `b.When` of the struct.
- Validating `T` instead of `*T` works on a copy: `Trim` doesn't change the caller's struct.
