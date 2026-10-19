# How it works

This page is for the contributors and for those who want to know the cost of a rule.

## Configure: from pointers to fields

`Configure[T]` creates the *model*, a `*T` with all exported nested structs stored by pointer allocated (pointers to the types being allocated are left nil, so recursive types stop). The model gives every field a unique address, so a pointer `&u.Address.City` identifies the field: [fmap](https://github.com/Insei/fmap) finds the field descriptor by the address relative to the model.

The descriptor knows the offset of the field and the chain of the nested pointer structs, so it reads the field of any object of the type: `Address.City` of a `*User` follows the `Address` pointer and returns a typed nil pointer when it is nil. That's how [absent values](../features/absent.md) appear.

## Typed rules

The configurator methods are generic: `Field[F any](p *F) *shared.Field[F]`, `Slice[E any](p *[]E)`. The type of the field becomes the type parameter of the base configurator. `As` is a generic method too, it takes a conversion with a type constraint, `str.Of[S ~string | *string](f *shared.Field[S])`: the compiler checks the constraint against the type of the field, and the conversion fixes the types of the rule arguments.

Every typed configurator is a thin wrapper over one core type, `shared.Field[T]`. It has four primitives, all built-in rules are made of them:

| Primitive | Behavior |
|---|---|
| `Check(ok, key, args...)` | an error with the key when `ok(value)` is false, absent values are skipped |
| `Required(ok, key)` | an error when the value is absent or `ok(value)` is false |
| `Modify(fn)` | changes the value in place, e.g. `Trim` |
| `Custom(fn)` | your function with the typed value and the field helper |

`When(cond)` wraps the next rules of the field with the condition. `MaxLen` of a string is `Check(func(v string) bool { return len(v) <= n }, ...)`.

`shared.Field[T]` doesn't know where the value comes from. For a struct field the source reads the field, for `Each` the source iterates over the elements of the slice and adds the index to the locations of the errors, for the keys and the values of a map it copies them to a pooled slot. All rules of a source share one read of the value. That's why the same string configurator validates a `Name string` field and every element of `Tags []*string`.

## Validate: the hot path

`Configure` registers one validation function per field source in the validator, keyed by the type `*T`. `Validate`:

1. finds the functions of the type of the object, a map lookup;
2. calls every function: it reads the field once through the descriptor, a pointer into the object, and runs the rules of the field;
3. a rule returns `nil` for a valid value, and builds a `shared.Error` with the translated message otherwise.

Nothing in this path allocates for valid objects: the values are read by pointers, the rules are closures created by `Configure`, the regexps are compiled once, the messages and their parameters are built by `shared.M` on the configuration, and the helper of struct `Custom` rules is pooled. A test of every package checks it with `testing.AllocsPerRun`.

The rules append their errors to one pooled buffer of the validation, and `Validate` copies it once: an invalid object allocates the copy, the translated messages and the boxed values of its errors.

## Packages

| Package | Contents |
|---|---|
| `valigo` | `Validator`, `Configure`, `Configurator` with `Field`, `Slice`, `When` and `Custom`, options |
| `shared` | the rules core `Field[T]`, slice configurator and `Elems[E]`, `As` and the conversions `Deref`, `Retype`, `RetypeSlice`, `SlicePtr`, `SliceRules`, `Error`, helpers |
| `str`, `num`, `timex`, `eq` | configurators of strings, numbers, times and comparable values built on `shared.Field[T]`, their conversions `str.Of`, `num.Of`, `num.Ptr`, `timex.Of`, `eq.Of`... and slice conversions |
| `problem` | the problem details of RFC 9457 |
| `github.com/insei/valigo/uuid` | the separate module with the rules of `github.com/google/uuid`, the core doesn't depend on a uuid package |
| `github.com/insei/valigo/huma` | the separate module of the inputs of huma operations |
| `translator` | translator, in-memory storage, built-in locales, `Accept-Language` middleware |

## Contributing

- Commits follow [Conventional Commits](https://www.conventionalcommits.org/), the changelog is generated from them by git-cliff.
- `go test -race ./...` and the examples modules run in CI. Programs in the documentation and in the skill for AI agents are built and run by `scripts/check-examples.sh`.
- A new rule of an existing type is a method of its configurator over `Check` or `Required`, with a locale key in `translator/locales/*/data.yaml`.
