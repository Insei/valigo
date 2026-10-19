# Validator and Configure

## Validator

`valigo.New` creates a validator. It holds the rules of all configured types and the translator of the messages, create one per application and share it.

```go
v := valigo.New(
    valigo.WithTranslator(t),                     // translator of the messages, see Translations
    valigo.WithFieldLocationNamingFn(jsonPath),   // location of the field in errors, see Errors
    valigo.WithErrorsTransformer(toProblemFields), // []shared.Error to []error for Validate, see Errors
)
```

All options are optional, a `nil` argument keeps the default.

## Configure

`valigo.Configure[T]` configures the rules of the struct type `T`:

```go
valigo.Configure[User](v, func(b valigo.Configurator[User], u *User) {
    b.Field(&u.Name).As(str.Of).Required()
})
```

The function is called once, with:

- `b` — the configurator of `T`: `Field`, `Slice` and `Map` return the base configurators of fields, converted to the configurators of their types by `As`, see [Types and conversions](https://insei.github.io/valigo/0.2/features/conversions/index.md); [`When`](https://insei.github.io/valigo/0.2/features/conditions/index.md) and [`Custom`](https://insei.github.io/valigo/0.2/features/custom/index.md) are the rules of the struct;
- `u` — the *model*, a `*User` created by `Configure`. Rules are bound to the fields of the model by their pointers.

The rules are applied to objects of type `*T`: `v.Validate(ctx, &user)`. A struct value is validated too, by a copy, see [Errors](https://insei.github.io/valigo/0.2/features/errors/#validate-and-validatetyped).

`Configure` is called on the application start. It isn't safe to call concurrently with `Validate`, and it **panics** on invalid configuration: a configuration error is a bug of the code, not a state of the running service. A repeated `Configure` of the same type replaces its rules.

## Field pointers

A field pointer must point to a field of the model passed to the function:

```go
valigo.Configure[User](v, func(b valigo.Configurator[User], u *User) {
    b.Field(&u.Name)         // a field of the model
    b.Field(&u.Address.City) // a field of a nested struct of the model

    var name string
    b.Field(&name)           // panic: *string is not a pointer to a field of *main.User: field not found
})
```

Types of the fields are checked by the compiler: `b.Field(&u.Age).As(str.Of)` doesn't compile for an `int` field, the conversions accept only the types they support.

## Nested structs

Fields of nested structs are configured by the same pointers. The model allocates all nested structs stored by pointer, so `&u.Address.City` is valid even if `Address` is `*Address`:

```go
package main

import (
    "context"
    "fmt"

    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/str"
)

type Address struct {
    City string
}

type User struct {
    Name    string
    Home    Address
    Work    *Address
    Manager *User // recursive types are supported, the model doesn't allocate Manager
}

func main() {
    v := valigo.New()
    valigo.Configure[User](v, func(b valigo.Configurator[User], u *User) {
        b.Field(&u.Home.City).As(str.Of).Required()
        b.Field(&u.Work.City).As(str.Of).Required()
    })

    fmt.Println(v.Validate(context.Background(), &User{Work: &Address{}}))
    fmt.Println(v.Validate(context.Background(), &User{Home: Address{City: "Berlin"}}))
}
```

Output:

```text
[Should be fulfilled (Home.City: ) Should be fulfilled (Work.City: )]
[Should be fulfilled (Work.City: <nil>)]
```

When the validated object has a nil `Work`, the field `Work.City` is absent: rules skip it and `Required` fails, see [Absent values](https://insei.github.io/valigo/0.2/features/absent/index.md).

Embedded structs work the same way, `&u.ID` for a field of an embedded `Base`. Unexported fields are allowed in the model, and unexported fields of the model itself can be configured too.

## Reusing rules

Rules of nested types aren't inherited: configuring `Address` doesn't add its rules to `User`. Share the rules by a function:

```go
func addressRules[T any](b valigo.Configurator[T], a *Address) {
    b.Field(&a.City).As(str.Of).Required()
}

valigo.Configure[User](v, func(b valigo.Configurator[User], u *User) {
    addressRules(b, &u.Home)
    addressRules(b, u.Work)
})
```
