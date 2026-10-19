# Upgrading from v0.0.x

Rules are still configured by pointers to the fields of the model, but the configurator is typed now and knows nothing about the types of the values: a field is converted to the configurator of its type by `As`. Go 1.27 is required.

| v0.0.x | Now |
|---|---|
| `github.com/insei/valigo/str`, `github.com/insei/valigo/num` | `github.com/insei/valigo/rules/str`, `github.com/insei/valigo/rules/num`, also `rules/eq` and `rules/timex` |
| `github.com/insei/valigo/uuid` of the core module | the separate module `github.com/insei/valigo/rules/uuid` |
| `b.String(&u.Name)` | `b.Field(&u.Name).As(str.Of)` |
| `b.Number(&u.Age)` | `b.Field(&u.Age).As(num.Of)` |
| `b.Number(&u.Score)` for `*int` | `b.Field(&u.Score).As(num.Ptr)` |
| `b.UUID(&u.ID)` | `b.Field(&u.ID).As(vuuid.Of)` of the separate module `github.com/insei/valigo/rules/uuid` |
| `b.StringSlice(&u.Tags)` | `b.Slice(&u.Tags).As(str.Slice)`, for `*[]string`: `b.Field(&u.Tags).As(shared.SlicePtr).As(str.Slice)` |
| `b.UUIDSlice(&u.IDs)` | `b.Slice(&u.IDs).As(vuuid.Slice)` |
| `b.Slice(any)` with `Custom(func(ctx, h, value []*any) ...)` | `b.Slice(&u.Items)` with `Custom(func(ctx, h, value shared.Elems[E]) ...)` and `Each` |
| `Configurator[T]` is an interface | a struct with generic methods, the signature of the `Configure` function is the same |
| field configurators are interfaces: `str.BaseConfigurator`, `num.BaseConfigurator` | pointers to concrete types: `*str.BaseConfigurator`, `*num.BaseConfigurator[N]`, `*vuuid.BaseConfigurator` |
| `Max(any)`, `Min(any)`, `AnyOf(...any)`, `AnyOfInterval(any, any)` | the arguments have the type of the field, `Max(5)` works for `int64` |
| `Custom(func(ctx, h, value any) ...)`, `When(func(ctx, value any) bool)` of fields | the value is typed: `string`, the number type, `uuid.UUID` |
| `StringSlice(...).MaxLen(n)` | the same: the number of elements; rules on elements are set by `Each` |
| `shared.NewSliceFieldConfigurator(params)`, `shared.FieldConfigurator` | `shared.NewSliceFieldConfigurator(field, helper, fields, add, describe)`, `shared.Field[T]` |

## Behavior changes

- `Validate` returns `nil` for a valid object, it returned an empty slice. Check `len(errs) > 0` or `errs != nil`, both work now.
- A struct value is validated by a copy, it was silently skipped.
- Nil pointers and the fields of nil nested structs are skipped by rules, they panicked; `Required` fails on them.
- `Required` of `num.Of` fails on zero, it did nothing. `Required` of `num.Ptr` fails on nil.
- `AnyOfInterval(begin, end)` includes both ends.
- The errors of slice elements have the index in the location: `Tags[1]`.
- A field `When` isn't called for an absent value, it got the zero value.
- A repeated `Configure` of the same type replaces its rules, it appended them.
- `shared.Error` has `Code` and `Params`, its JSON has them when they aren't empty; the messages of counts have plural forms: "longer than 1 character".
