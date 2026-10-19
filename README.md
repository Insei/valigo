[![codecov](https://codecov.io/gh/Insei/valigo/graph/badge.svg?token=TK4XN3I2Q6)](https://codecov.io/gh/Insei/valigo)
[![build](https://github.com/Insei/valigo/actions/workflows/go.yml/badge.svg)](https://github.com/Insei/valigo/actions/workflows/go.yml)
[![Goreport](https://goreportcard.com/badge/github.com/insei/valigo)](https://goreportcard.com/report/github.com/insei/valigo)
[![GoDoc](https://godoc.org/github.com/insei/valigo?status.svg)](https://godoc.org/github.com/insei/valigo)
# valigo
valigo validates Go structs by rules written in Go, next to the struct, by pointers to its fields. The compiler checks
every rule, a valid object is validated without allocations about 4 times faster than by go-playground/validator, and
the errors come translated to the language of the user, with stable codes for the clients that translate themselves.

```go
b.Field(&u.Email).As(str.Of).Required().Email()
b.Field(&u.Age).As(num.Of).AnyOfInterval(18, 120) // the arguments are int: Max("120") doesn't compile
b.Field(&u.Name).As(num.Of)                       // doesn't compile: string does not satisfy num.Number
```

Documentation: [insei.github.io/valigo](https://insei.github.io/valigo/) — get started, [why valigo](https://insei.github.io/valigo/latest/why-valigo/), the rules of every type, [performance](https://insei.github.io/valigo/latest/performance/).

For AI agents: the [skill](skills/valigo/SKILL.md) with the current API, and the documentation as one file, [llms-full.txt](https://insei.github.io/valigo/latest/llms-full.txt).

# Installation
Go 1.27 or newer is required, the configurator uses generic methods.
```sh
go get github.com/insei/valigo@latest
```
Upgrading from v0.0.x: see [Upgrading](https://insei.github.io/valigo/latest/upgrading/).

# Example
```go
type User struct {
	Name     string
	Email    string
	Age      int
	Phone    *string
	Tags     []string
	StartsAt time.Time
	EndsAt   *time.Time
}

v := valigo.New()
valigo.Configure[User](v, func(b valigo.Configurator[User], u *User) {
	b.Field(&u.Name).As(str.Of).Trim().Required().MaxRunes(64)
	b.Field(&u.Email).As(str.Of).Required().Email()
	b.Field(&u.Age).As(num.Of).AnyOfInterval(18, 120)
	b.Field(&u.Phone).As(str.Of).E164() // a nil phone is skipped
	b.Slice(&u.Tags).As(str.Slice).MaxLen(5).Each(func(tag *str.BaseConfigurator) {
		tag.Trim().Required().MaxLen(16)
	})
	b.Field(&u.EndsAt).As(timex.Of).AfterField(&u.StartsAt) // a rule on another field
})

errs := v.ValidateTyped(ctx, &User{Name: "  Ann  ", Email: "ann.example.com", Age: 16, Tags: []string{"go", " "}})
// Should be email address (Email: ann.example.com)   code string.email
// Only interval[18 - 120] is allowed (Age: 16)       code number.interval, params {"min": 18, "max": 120}
// Should be fulfilled (Tags[1]: )                    code string.required
```

# Rules

The configurator has `Field`, `Slice` and `Map` for the fields, `As` converts them to the configurators of their types,
and a conversion that doesn't fit the field doesn't compile:

| Conversion | Fields | Rules |
|---|---|---|
| `str.Of` | `string`, `*string`, named strings | `Required`, lengths in bytes and runes, `AnyOf`, `Prefix`, `Suffix`, `Contains`, `Regexp`, `Email`, `URL`, `IP`, `CIDR`, `MAC`, `Hostname`, `UUID`, `E164`, `Base64`, `JSON`, `Semver`, `HexColor`, `Datetime`, `CountryCode`, `CurrencyCode`, `CreditCard`, `Trim`... |
| `num.Of`, `num.Ptr` | integers, floats, named numbers, pointers to them | `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval`, `Latitude`, `Longitude` |
| `timex.Of` | `time.Time`, `*time.Time` | `Required`, `Before`, `After`, `Between`, `InPast`, `InFuture` |
| `eq.Of`, `eq.Ptr` | any `comparable`: enums, bools | `Required`, `OneOf`, `NotOneOf`, `Equal`, `NotEqual` |
| `vuuid.Of` | `uuid.UUID`, `*uuid.UUID`, the separate module `github.com/insei/valigo/uuid` | `Required`, `AnyOf` |
| `str.Slice`, `num.Slice`, `timex.Slice`, `eq.Slice`, `vuuid.Slice` | slices, also of pointers | `Required`, `MinLen`, `MaxLen`, `Each` for the elements, `Unique` |
| `b.Map(&m.X)` | maps | `Required`, `MinLen`, `MaxLen`, `EachKey`, `EachValue` |

- **Rules on other fields** by their pointers: `AfterField`, `MinField`, `MaxField`, `EqualField`.
- **Conditions and custom rules**: `When` and `Custom` on fields and on the struct, closures with the context of the request — the limits can come from a database or the settings of a tenant.
- **Absent values**: nil pointers and the fields of nil nested structs are skipped by every rule except `Required`.
- **Your types**: a conversion function of a few lines, the core has no dependencies on uuid or other packages of types.

# Errors and translations
* Every error has the translated message, the location (`Tags[1]`, `address.city` with JSON paths), the value, a stable `Code` and named `Params`
* Messages in English, German, Spanish, French and Russian with plural forms, the language of every request from the context, an `Accept-Language` middleware
* `problem`: the errors as `application/problem+json` of RFC 9457 with JSON pointers
* `github.com/insei/valigo/huma`: the rules of valigo in the inputs of [huma](https://huma.rocks) operations

# Performance
The same struct and rules, Intel Core i9-14900HX, Go 1.27.0, see [Performance](https://insei.github.io/valigo/latest/performance/):

| | valid object | invalid object, 6 errors |
|---|---|---|
| hand-written code | 286 ns, 0 allocs | 590 ns, without messages |
| **valigo** | **458 ns, 0 allocs** | **2 205 ns, 14 allocs, translated messages** |
| go-playground/validator | 1 785 ns, 16 allocs | 4 150 ns, 69 allocs, with `Error()` |
| ozzo-validation | 2 950 ns, 38 allocs | 12 665 ns, 140 allocs, with `Error()` |

# Roadmap
* [ ] Rules from tags
