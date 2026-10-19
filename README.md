[![codecov](https://codecov.io/gh/Insei/valigo/graph/badge.svg?token=TK4XN3I2Q6)](https://codecov.io/gh/Insei/valigo)
[![build](https://github.com/Insei/valigo/actions/workflows/go.yml/badge.svg)](https://github.com/Insei/valigo/actions/workflows/go.yml)
[![Goreport](https://goreportcard.com/badge/github.com/insei/valigo)](https://goreportcard.com/report/github.com/insei/valigo)
[![GoDoc](https://godoc.org/github.com/insei/valigo?status.svg)](https://godoc.org/github.com/insei/valigo)
# valigo
valigo validates Go structs by rules configured once, in code, with pointers to the fields of a model. The compiler checks
that a rule fits the field and that its arguments have the field type, valid objects are validated without allocations,
and the messages are translated.

Documentation: [insei.github.io/valigo](https://insei.github.io/valigo/) — get started, [why valigo](https://insei.github.io/valigo/latest/why-valigo/), features and how it works.

For AI agents: the [skill](skills/valigo/SKILL.md) with the current API, and the documentation as one file, [llms-full.txt](https://insei.github.io/valigo/latest/llms-full.txt).

# Installation
Install via go get. Note that Go 1.27 or newer is required, the configurator uses generic methods.
```sh
go get github.com/insei/valigo@latest
```
Upgrading from v0.0.x: see [Upgrading](https://insei.github.io/valigo/latest/upgrading/).

# Example
```go
type User struct {
	Name  string
	Email string
	Age   int
	Phone *string
	Tags  []string
}

v := valigo.New()
valigo.Configure[User](v, func(b valigo.Configurator[User], u *User) {
	b.String(&u.Name).Trim().Required().MaxLen(64)
	b.String(&u.Email).Required().Email()
	b.Number(&u.Age).AnyOfInterval(18, 120)
	b.String(&u.Phone).Regexp(phoneRe) // a nil phone is skipped
	b.StringSlice(&u.Tags).MaxLen(5).Each(func(tag *str.BaseConfigurator) {
		tag.Trim().Required().MaxLen(16)
	})
})

errs := v.Validate(ctx, &User{Name: "  Ann  ", Email: "ann.example.com", Age: 16, Tags: []string{"go", " "}})
// Should be email address (Email: ann.example.com)
// Only interval[18 - 120] is allowed (Age: 16)
// Should be fulfilled (Tags[1]: )
```

# Rules

| Configurator | Fields | Rules |
|---|---|---|
| `String` | `string`, `*string`, named strings | `Required`, `MinLen`, `MaxLen`, `AnyOf`, `Regexp`, `Email`, `Trim` |
| `Number`, `NumberPtr` | integers, floats, named numbers, pointers to them | `Required`, `Min`, `Max`, `AnyOf`, `AnyOfInterval` |
| `UUID` | `uuid.UUID`, `*uuid.UUID` | `Required`, `AnyOf` |
| `StringSlice`, `UUIDSlice`, `Slice` | `[]T`, `*[]T`, `[]*T`, `*[]*T` | `Required`, `MinLen`, `MaxLen`, `Each` for the rules of the elements |

Every configurator has `Custom` for your rules and `When` for conditions, the struct configurator has `When` and `Custom`
on the whole object. Absent values — nil pointers and the fields of nil nested structs — are skipped by every rule
except `Required`.

# Features
* Rules are checked by the compiler: `Number(&u.ID).Max("5")` and `String(&u.Age)` don't compile
* Zero allocations on valid objects
* Conditions on the struct and on fields, custom rules of fields and structs
* Errors with the locations of fields and slice elements, `Tags[1]`, or JSON paths
* Translated messages: English and Russian out of the box, own messages, `Accept-Language` middleware

# Roadmap
* [ ] Other default types validations: time, bool, maps
* [ ] Rules from tags
