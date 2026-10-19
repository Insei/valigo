# How it works

This page is for the contributors and for those who want to know the cost of a rule.

## Configure: from pointers to fields

`Configure[T]` creates the *model*, a `*T` with all exported nested structs stored by pointer allocated (pointers to the types being allocated are left nil, so recursive types stop). The model gives every field a unique address, so a pointer `&u.Address.City` identifies the field: [fmap](https://github.com/Insei/fmap) finds the field descriptor by the address relative to the model.

The descriptor knows the offset of the field and the chain of the nested pointer structs, so it reads the field of any object of the type: `Address.City` of a `*User` follows the `Address` pointer and returns a typed nil pointer when it is nil. That's how [absent values](../features/absent.md) appear.

## Typed rules

Configurator methods are generic: `String[F ~string | *string](p *F)`, `Number[N num.Number](p *N)`. The type of the pointer selects the method and the types of the rule arguments, so a wrong field or argument type doesn't compile.

Every field configurator is a thin wrapper over one core type, `shared.Field[T]`. It has four primitives, all built-in rules are made of them:

| Primitive | Behavior |
|---|---|
| `Check(ok, key, args...)` | an error with the key when `ok(value)` is false, absent values are skipped |
| `Required(ok, key)` | an error when the value is absent or `ok(value)` is false |
| `Modify(fn)` | changes the value in place, e.g. `Trim` |
| `Custom(fn)` | your function with the typed value and the field helper |

`When(cond)` wraps the next rules of the field with the condition. `MaxLen` of a string is `Check(func(v string) bool { return len(v) <= n }, ...)`.

`shared.Field[T]` doesn't know where the value comes from. For a struct field the source reads the field, for `Each` the source iterates over the elements of the slice and adds the index to the locations of the errors. That's why the same string configurator validates a `Name string` field and every element of `Tags []*string`.

## Validate: the hot path

`Configure` registers one validation function per field source in the validator, keyed by the type `*T`. `Validate`:

1. finds the functions of the type of the object, a map lookup;
2. calls every function: it reads the field once through the descriptor, a pointer into the object, and runs the rules of the field;
3. a rule returns `nil` for a valid value, and builds a `shared.Error` with the translated message otherwise.

Nothing in this path allocates for valid objects: the values are read by pointers, the rules are closures created by `Configure`, the regexps are compiled once, and the helper of struct `Custom` rules is pooled. A test checks it with `testing.AllocsPerRun`.

## Packages

| Package | Contents |
|---|---|
| `valigo` | `Validator`, `Configure`, `Configurator` with its generic methods, options |
| `shared` | the rules core `Field[T]`, slice configurator and `Elems[E]`, `Error`, helpers |
| `str`, `num`, `uuid` | configurators of strings, numbers and uuids built on `shared.Field[T]` |
| `translator` | translator, in-memory storage, built-in locales, `Accept-Language` middleware |

## Contributing

- Commits follow [Conventional Commits](https://www.conventionalcommits.org/), the changelog is generated from them by git-cliff.
- `go test -race ./...` and the examples modules run in CI. Programs in the documentation and in the skill for AI agents are built and run by `scripts/check-examples.sh`.
- A new rule of an existing type is a method of its configurator over `Check` or `Required`, with a locale key in `translator/locales/*/data.yaml`.
