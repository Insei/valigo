# valigo

**valigo** is a validation library for Go structs. Rules are configured once, in code, by pointers to the fields of a model: `b.Field(&u.Name).As(str.Of).Required().MaxLen(64)`. The compiler checks that a rule fits the field type and that its arguments have the field type, the validation itself doesn't allocate on valid objects, and error messages are translated.

!!! tip "One-line ideology"
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
	"github.com/insei/valigo/rules/num"
	"github.com/insei/valigo/rules/str"
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
		b.Field(&u.Name).As(str.Of).Trim().Required().MaxLen(64)
		b.Field(&u.Email).As(str.Of).Required().Email()
		b.Field(&u.Age).As(num.Of).AnyOfInterval(18, 120)
		b.Field(&u.Phone).As(str.Of).Regexp(regexp.MustCompile(`^\+\d{10,15}$`)) // a nil phone is skipped
		b.Slice(&u.Tags).As(str.Slice).MaxLen(5).Each(func(tag *str.BaseConfigurator) {
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

1. `valigo.New()` creates a validator with the messages in English, German, Spanish, French and Russian.
2. `valigo.Configure[User]` calls the function once with a *model*: a `*User` whose nested pointer structs are allocated. Rules are bound to the fields of the model by their pointers, `&u.Name`.
3. Every configurator method returns the configurator of the field, so rules are chained. They run in the order of the chain: `Trim` changes the value before `Required` checks it.
4. `v.Validate` runs the rules of the `*User` type and returns `nil` for a valid object. Every error has the message, the location of the field (`Tags[1]` for the element of a slice) and the value.

## Where to go next

- [Why valigo?](why-valigo.md) — what it gives compared to tag based validators.
- [Features](features/index.md) — the rules of every type, conditions, custom rules, errors, translations.
- [How it works](architecture/index.md) — the model, the field pointers and why valid objects don't allocate.
- [Performance](performance.md) — benchmarks with go-playground/validator, ozzo-validation, zog, gookit/validate and hand-written code.
- [HTTP problem details](features/problem.md) and [huma](features/huma.md) — valigo in HTTP APIs.
- [Upgrading from v0.0.x](upgrading.md) — the configurator API has changed.
- [For AI agents](ai-agents.md) — the skill and the documentation in one file.
