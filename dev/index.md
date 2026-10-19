# valigo

**valigo** is a validation library for Go structs. Rules are configured once, in code, by pointers to the fields of a model: `b.String(&u.Name).Required().MaxLen(64)`. The compiler checks that a rule fits the field type and that its arguments have the field type, the validation itself doesn't allocate on valid objects, and error messages are translated.

One-line ideology

*Rules are Go code next to the struct, not strings in tags: a typo or a wrong type doesn't compile.*

## Install

```bash
go get github.com/insei/valigo@latest
```

Minimum Go version is **1.27**: the configurator uses generic methods.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "regexp"

    "github.com/insei/valigo"
    "github.com/insei/valigo/str"
)

type User struct {
    Name  string
    Email string
    Age   int
    Phone *string
    Tags  []string
}

func main() {
    v := valigo.New()
    valigo.Configure[User](v, func(b valigo.Configurator[User], u *User) {
        b.String(&u.Name).Trim().Required().MaxLen(64)
        b.String(&u.Email).Required().Email()
        b.Number(&u.Age).AnyOfInterval(18, 120)
        b.String(&u.Phone).Regexp(regexp.MustCompile(`^\+\d{10,15}$`)) // a nil phone is skipped
        b.StringSlice(&u.Tags).MaxLen(5).Each(func(tag *str.BaseConfigurator) {
            tag.Trim().Required().MaxLen(16)
        })
    })

    user := &User{Name: "  Ann  ", Email: "ann.example.com", Age: 16, Tags: []string{"go", " "}}
    for _, err := range v.Validate(context.Background(), user) {
        fmt.Println(err)
    }
    fmt.Printf("%q\n", user.Name)
}
```

Output:

```text
Should be email address (Email: ann.example.com)
Only interval[18 - 120] is allowed (Age: 16)
Should be fulfilled (Tags[1]: )
"Ann"
```

What happens here:

1. `valigo.New()` creates a validator with the English and Russian messages.
1. `valigo.Configure[User]` calls the function once with a *model*: a `*User` whose nested pointer structs are allocated. Rules are bound to the fields of the model by their pointers, `&u.Name`.
1. Every configurator method returns the configurator of the field, so rules are chained. They run in the order of the chain: `Trim` changes the value before `Required` checks it.
1. `v.Validate` runs the rules of the `*User` type and returns `nil` for a valid object. Every error has the message, the location of the field (`Tags[1]` for the element of a slice) and the value.

## Where to go next

- [Why valigo?](https://insei.github.io/valigo/latest/why-valigo/index.md) — what it gives compared to tag based validators.
- [Features](https://insei.github.io/valigo/latest/features/index.md) — the rules of every type, conditions, custom rules, errors, translations.
- [How it works](https://insei.github.io/valigo/latest/architecture/index.md) — the model, the field pointers and why valid objects don't allocate.
- [Upgrading from v0.0.x](https://insei.github.io/valigo/latest/upgrading/index.md) — the configurator API has changed.
- [For AI agents](https://insei.github.io/valigo/latest/ai-agents/index.md) — the skill and the documentation in one file.
