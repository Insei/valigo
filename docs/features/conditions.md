# Conditions

## When on the struct

`When` of the configurator returns a configurator whose rules are applied only if the function returns true for the validated object:

```go
package main

import (
	"context"
	"fmt"

	"github.com/insei/valigo"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/str"
)

type Sender struct {
	Type     string
	SMTPHost string
	SMTPPort int
	HTTPURL  string
}

func main() {
	v := valigo.New()
	valigo.Configure[Sender](v, func(b valigo.Configurator[Sender], s *Sender) {
		b.Field(&s.Type).As(str.Of).AnyOf("smtp", "http")

		smtp := b.When(func(_ context.Context, s *Sender) bool { return s.Type == "smtp" })
		smtp.Field(&s.SMTPHost).As(str.Of).Required()
		smtp.Field(&s.SMTPPort).As(num.Of).AnyOfInterval(1, 65535)

		http := b.When(func(_ context.Context, s *Sender) bool { return s.Type == "http" })
		http.Field(&s.HTTPURL).As(str.Of).Required()
	})

	fmt.Println(v.Validate(context.Background(), &Sender{Type: "smtp"}))
	fmt.Println(v.Validate(context.Background(), &Sender{Type: "http"}))
}
```

Output:

```text
[Should be fulfilled (SMTPHost: ) Only interval[1 - 65535] is allowed (SMTPPort: 0)]
[Should be fulfilled (HTTPURL: )]
```

- The returned configurator has all methods of the original one, the conditions of chained `When` calls are combined with AND: `b.When(a).When(b)`.
- `When` doesn't change the configurator it is called on: `b` keeps adding unconditional rules.
- The context passed to `Validate` is passed to the condition, e.g. to enable rules by the role of the user.

## When on a field

`When` of a field configurator returns a configurator of the same field whose rules are applied only if the function returns true for the value of the field. The rules before `When` stay unconditional:

```go
b.Field(&u.Website).As(str.Of).
	Trim().                                                                  // always
	When(func(_ context.Context, v string) bool { return v != "" }).
	Regexp(urlRe)                                                            // only for a non-empty value
```

The value is typed: `string` for `str.Of`, the number type for `num.Of` and `num.Ptr`, `uuid.UUID` for `uuid.Of`, and `shared.Elems[E]` for slices, see [Slices](slices.md#custom-and-when).

For an [absent value](absent.md) the condition isn't called and the rules after it are skipped, including `Required`.

## Conditions and custom rules

A condition on several fields can also be written as a [struct custom rule](custom.md#struct-rules), which reports errors itself. Prefer `When`: the rules after it are the regular ones, with their messages and translations.
