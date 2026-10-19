---
name: valigo
description: Write Go code that validates structs with github.com/insei/valigo — typed rules configured by pointers to the fields of a model (String, Number, NumberPtr, UUID, Slice, StringSlice, UUIDSlice), Each for slice elements, When conditions, Custom rules of fields and structs, errors with field locations, translated messages and the Accept-Language middleware. Use when the code imports github.com/insei/valigo, or when a Go service validates request or config structs.
---

# valigo

Rules are configured once by `valigo.Configure[T]` with pointers to the fields of a model and are checked by the compiler. `Validate` runs them on objects of type `*T` and doesn't allocate on valid objects. Full documentation: https://insei.github.io/valigo/latest/, as one file for the context: https://insei.github.io/valigo/latest/llms-full.txt.

## Rules

- Go 1.27 or newer: the configurator methods are generic.
- Create one validator with `valigo.New(...)` on start and call `valigo.Configure[T]` once per type, before `Validate`. `Configure` panics on invalid configuration, it isn't safe concurrently with `Validate`, a repeated `Configure` of the same type replaces its rules.
- Field pointers must point to the fields of the model passed to the function: `func(b valigo.Configurator[T], m *T) { b.String(&m.Name) }`. A pointer to anything else panics. Nested structs stored by pointer are allocated in the model, `&m.Address.City` is valid.
- Pick the configurator by the field type: `String` for `string`, `*string` and named strings; `Number` for numbers and named numbers; `NumberPtr` for pointers to numbers; `UUID` for `uuid.UUID` and `*uuid.UUID`; `StringSlice`, `UUIDSlice` and `Slice` for slices. `Number(&m.P)` doesn't compile for `P *int`, use `NumberPtr`.
- Rule arguments have the type of the field: `b.Number(&m.ID).Max(5)` for an `int64` field, never `Max(int64(5))` with a value of another type or `any`.
- `Custom` and `When` of a field get the typed value: `func(ctx context.Context, h *shared.FieldCustomHelper, value string) []shared.Error`, `func(ctx context.Context, value int) bool`. Of slices: `shared.Elems[E]` with `Len()`, `At(i)`, `All()`.
- `MaxLen`, `MinLen` and `Required` of slice configurators check the slice itself, rules on elements are set by `Each(func(e *str.BaseConfigurator) { ... })`.
- Absent values (nil pointers, nil `*[]T`, fields of nil nested structs, nil elements of `[]*T`) are skipped by every rule except `Required`.
- Rules run in the order of the chain, `Trim()` changes the field of the validated object.
- `Validate` returns `[]error`, `ValidateTyped` returns `[]shared.Error{Message, Location, Value}`, both are `nil` for valid objects.

| Configurator | Rules |
|---|---|
| `String(&m.S)` | `Required`, `MinLen`, `MaxLen`, `AnyOf`, `Regexp(re, str.WithRegexpLocaleKey(key))`, `Email`, `Trim`, `Custom`, `When` |
| `Number(&m.N)`, `NumberPtr(&m.P)` | `Required` (zero for `Number`, nil for `NumberPtr`), `Min`, `Max`, `AnyOf`, `AnyOfInterval` (inclusive), `Custom`, `When` |
| `UUID(&m.ID)` | `Required` (not `uuid.Nil`), `AnyOf`, `Custom`, `When` |
| `StringSlice(&m.Tags)` | `Required`, `MinLen`, `MaxLen`, `Each(func(*str.BaseConfigurator))`, `Trim`, `Regexp`, `Email`, `Custom`, `When` |
| `UUIDSlice(&m.IDs)` | `Required`, `MinLen`, `MaxLen`, `Each(func(*uuid.BaseConfigurator))`, `AnyOf`, `Custom`, `When` |
| `Slice(&m.Items)` | `Required`, `MinLen`, `MaxLen`, `Each(func(*shared.Field[E]))` with `Check(ok, message)`, `Required(ok, message)`, `Custom`, `When` |
| `b.When(func(ctx, *T) bool)` | a configurator whose rules apply only when the function returns true |
| `b.Custom(func(ctx, h shared.StructCustomHelper, obj *T) []shared.Error)` | a rule on several fields, `h.ErrorT(ctx, &obj.Field, value, message)` |

## A request with typed rules

```go
package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
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
		b.UUID(&o.CustomerID).Required()
		b.String(&o.Email).Trim().Required().Email()
		b.String(&o.Comment).MaxLen(500) // nil Comment is valid
		b.Number(&o.Quantity).AnyOfInterval(1, 100)
		b.NumberPtr(&o.Discount).Max(50) // nil Discount is valid
		b.String(&o.Address.City).Required()
		b.String(&o.Address.Zip).Regexp(zipRe, str.WithRegexpLocaleKey("Should be 5 digits"))
		b.StringSlice(&o.Tags).MaxLen(5).Each(func(tag *str.BaseConfigurator) {
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
	"github.com/insei/valigo/shared"
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
		b.String(&n.Channel).Required().AnyOf("email", "sms")

		email := b.When(func(_ context.Context, n *Notification) bool { return n.Channel == "email" })
		email.String(&n.Email).Required().Email()

		sms := b.When(func(_ context.Context, n *Notification) bool { return n.Channel == "sms" })
		sms.String(&n.Phone).Required()

		b.Number(&n.Retries).When(func(_ context.Context, retries int) bool { return retries != 0 }).AnyOfInterval(1, 5)

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

## Translations

The built-in messages are in English and Russian. The language comes from the context: wrap the HTTP handler with `translator.NewAcceptLanguageMiddleware()`, or pass `translator.WithPreferredLanguagesFn(fn)` to `translator.New`. Messages of your rules are locale keys: a key without a translation is printed as is, translations are added to a storage:

```go
storage := translator.NewInMemStorage() // with the built-in locales
storage.Add("ru", map[string]string{"Should be in the future": "Должно быть в будущем"})
v := valigo.New(valigo.WithTranslator(translator.New(translator.WithStorage(storage))))
```

## Common mistakes

- `b.Number(&m.Score)` for `Score *int` doesn't compile: use `b.NumberPtr(&m.Score)`.
- `b.String(&local)` with a variable outside the model panics in `Configure`: use the model passed to the function.
- `Configure` in a request handler: configure once on start, `Configure` replaces the rules and isn't safe concurrently with `Validate`.
- `b.StringSlice(&m.Tags).MaxLen(10)` limits the number of tags, not their length: use `Each(func(t *str.BaseConfigurator) { t.MaxLen(10) })`.
- `Required` after a field `When`: the condition isn't called for an absent value and the rules after it, including `Required`, are skipped. Put `Required` before `When`, or use `b.When` of the struct.
- Validating `T` instead of `*T` works on a copy: `Trim` doesn't change the caller's struct.
