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
    Message  string // translated message
    Location string // location of the field, Address.City or Tags[1] by default
    Value    any    // the checked value, nil for absent values
}
```

`Error()` returns `Message (Location: Value)`, or `Message` alone when both the location and the value are empty. `Error` has JSON field names of its own: `{"Message":"...","Location":"...","Value":...}`, convert it to your API format by the transformer.

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

The transformer is called only when there are errors.
