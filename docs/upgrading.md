# Upgrading from v0.0.x

The configuration style is the same, `b.String(&u.Name).Required()`, but the configurator is typed now. Go 1.27 is required.

| v0.0.x | Now |
|---|---|
| `Configurator[T]` is an interface | a struct with generic methods, the signature of the `Configure` function is the same |
| field configurators are interfaces: `str.BaseConfigurator`, `num.BaseConfigurator` | pointers to concrete types: `*str.BaseConfigurator`, `*num.BaseConfigurator[N]`, `*uuid.BaseConfigurator` |
| `b.Number(&u.Score)` for `*int` | `b.NumberPtr(&u.Score)` |
| `Max(any)`, `Min(any)`, `AnyOf(...any)`, `AnyOfInterval(any, any)` | the arguments have the type of the field, `Max(5)` works for `int64` |
| `Custom(func(ctx, h, value any) ...)`, `When(func(ctx, value any) bool)` of fields | the value is typed: `string`, the number type, `uuid.UUID` |
| `Slice(any)` with `Custom(func(ctx, h, value []*any) ...)` | `Slice[E](*[]E)` with `Custom(func(ctx, h, value shared.Elems[E]) ...)` |
| `StringSlice(...).MaxLen(n)` | the same: the number of elements; rules on elements are set by `Each` |
| `shared.NewSliceFieldConfigurator(params)`, `shared.FieldConfigurator` | `shared.NewSliceFieldConfigurator(field, helper, add)`, `shared.Field[T]` |

## Behavior changes

- `Validate` returns `nil` for a valid object, it returned an empty slice. Check `len(errs) > 0` or `errs != nil`, both work now.
- A struct value is validated by a copy, it was silently skipped.
- Nil pointers and the fields of nil nested structs are skipped by rules, they panicked; `Required` fails on them.
- `Number(...).Required()` fails on zero, it did nothing. `NumberPtr(...).Required()` fails on nil.
- `AnyOfInterval(begin, end)` includes both ends.
- The errors of slice elements have the index in the location: `Tags[1]`.
- A field `When` isn't called for an absent value, it got the zero value.
- A repeated `Configure` of the same type replaces its rules, it appended them.
