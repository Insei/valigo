# Conditions

## When on the struct

`When` of the configurator returns a configurator whose rules are applied only if the function returns true for the validated object:

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/num"
    "github.com/insei/valigo/rules/str"
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

The value is typed: `string` for `str.Of`, the number type for `num.Of` and `num.Ptr`, `uuid.UUID` for `uuid.Of`, and `shared.Elems[E]` for slices, see [Slices](https://insei.github.io/valigo/0.2/features/slices/#custom-and-when).

For an [absent value](https://insei.github.io/valigo/0.2/features/absent/index.md) the condition isn't called and the rules after it are skipped, including `Required`.

## Conditional requiredness

Tag based validators have families of tags like `required_if`, `required_with` and `required_without`. In valigo they are `When` with a condition on the fields:

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/str"
    "github.com/insei/valigo/shared"
)

type Customer struct {
    Kind    string
    TaxID   string
    Phone   string
    Email   string
    Company string
    Title   string
}

func main() {
    v := valigo.New()
    valigo.Configure[Customer](v, func(b valigo.Configurator[Customer], c *Customer) {
        // required_if=Kind company
        b.When(func(_ context.Context, c *Customer) bool { return c.Kind == "company" }).Field(&c.TaxID).As(str.Of).Required()
        // required_without=Email: a phone or an email
        b.When(func(_ context.Context, c *Customer) bool { return c.Email == "" }).Field(&c.Phone).As(str.Of).Required()
        // required_with=Company: the title of the person in the company
        b.When(func(_ context.Context, c *Customer) bool { return c.Company != "" }).Field(&c.Title).As(str.Of).Required()
        // excluded_unless=Kind company
        b.When(func(_ context.Context, c *Customer) bool { return c.Kind != "company" }).Field(&c.Company).As(str.Of).
            Custom(func(ctx context.Context, h *shared.FieldCustomHelper, company string) []shared.Error {
                if company != "" {
                    return []shared.Error{h.ErrorT(ctx, company, "Only companies have a company")}
                }
                return nil
            })
    })

    fmt.Println(v.Validate(context.Background(), &Customer{Kind: "company"}))
    fmt.Println(v.Validate(context.Background(), &Customer{Kind: "person", Email: "a@b.cc", Company: "ACME"}))
}
```

Output:

```text
[Should be fulfilled (TaxID: ) Should be fulfilled (Phone: )]
[Should be fulfilled (Title: ) Only companies have a company (Company: ACME)]
```

## Conditions and custom rules

A condition on several fields can also be written as a [struct custom rule](https://insei.github.io/valigo/0.2/features/custom/#struct-rules), which reports errors itself. Prefer `When`: the rules after it are the regular ones, with their messages and translations.
