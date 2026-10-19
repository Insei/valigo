# Slices

`b.Slice(&m.X)` returns the configurator of a `[]E` field with the rules of the slice itself. `As` converts it to a typed slice configurator, the rules of the slice stay available after the conversion:

| Configurator                        | Elements                           | Rules on elements                                                                                                                                                                                       |
| ----------------------------------- | ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `b.Slice(&m.X)`                     | any `E`                            | `Each` with `*shared.Field[E]`: `Check`, `Required`, `Modify`, `Custom`, `When`, and `As` to convert the element                                                                                        |
| `.As(str.Slice)`                    | `string`, `*string`, named strings | `Each` with the [string configurator](https://insei.github.io/valigo/0.1/features/strings/index.md), `Trim`, `Regexp`, `Email`                                                                          |
| `.As(num.Slice)`                    | numbers and named numbers          | `Each` with the [number configurator](https://insei.github.io/valigo/0.1/features/numbers/index.md)                                                                                                     |
| `.As(timex.Slice)`, `.As(eq.Slice)` | times, comparable values           | `Each` with the [time](https://insei.github.io/valigo/0.1/features/time/index.md) and the [equality](https://insei.github.io/valigo/0.1/features/values/index.md) configurators, `Unique` of `eq.Slice` |
| `.As(vuuid.Slice)`                  | `uuid.UUID`, `*uuid.UUID`          | `Each` with the [uuid configurator](https://insei.github.io/valigo/0.1/features/uuid/index.md), `AnyOf`                                                                                                 |

A pointer to a slice, `*[]E`, is converted from its field: `b.Field(&m.Tags).As(shared.SlicePtr).As(str.Slice)`, a nil pointer is [absent](https://insei.github.io/valigo/0.1/features/absent/index.md).

## The slice and its elements

The rules of the configurator check **the slice itself**, the rules inside `Each` check **every element**:

| Rule         | Fails when                           | Message                                   |
| ------------ | ------------------------------------ | ----------------------------------------- |
| `Required()` | the slice is absent or empty         | Should be fulfilled                       |
| `MinLen(n)`  | the slice has less than `n` elements | Cannot contain less than %d elements      |
| `MaxLen(n)`  | the slice has more than `n` elements | Cannot contain more than %d elements      |
| `Each(fn)`   | —                                    | `fn` configures the rules of the elements |

`Trim`, `Regexp` and `Email` of `str.Slice` and `AnyOf` of `uuid.Slice` are shortcuts for the same rules in `Each`. Any element of `b.Slice` is converted in `Each` too: `e.As(num.Ptr).Max(100)` for `[]*int`.

```go
package main

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/insei/valigo"
    "github.com/insei/valigo/rules/str"
    vuuid "github.com/insei/valigo/rules/uuid"
    "github.com/insei/valigo/shared"
)

type Post struct {
    Tags       []string
    Editors    []*uuid.UUID
    Quantities []int
}

func main() {
    v := valigo.New()
    valigo.Configure[Post](v, func(b valigo.Configurator[Post], p *Post) {
        b.Slice(&p.Tags).As(str.Slice).Required().MaxLen(3).Each(func(tag *str.BaseConfigurator) {
            tag.Trim().Required().MaxLen(10)
        })
        b.Slice(&p.Editors).As(vuuid.Slice).Each(func(id *vuuid.BaseConfigurator) {
            id.Required()
        })
        b.Slice(&p.Quantities).Each(func(q *shared.Field[int]) {
            q.Check(func(v int) bool { return v > 0 }, "Should be positive")
        })
    })

    editor := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
    post := &Post{
        Tags:       []string{"go", "  ", "validation-library"},
        Editors:    []*uuid.UUID{&editor, nil},
        Quantities: []int{1, 0},
    }
    for _, err := range v.Validate(context.Background(), post) {
        fmt.Println(err)
    }
}
```

Output:

```text
Should be fulfilled (Tags[1]: )
Cannot be longer than 10 characters (Tags[2]: validation-library)
Should be fulfilled (Editors[1]: <nil>)
Should be positive (Quantities[1]: 0)
```

The location of an element error has the index of the element, the errors of the slice itself have the location of the field. Nil elements of slices of pointers are [absent](https://insei.github.io/valigo/0.1/features/absent/index.md): only `Required` fails on them.

## Custom and When

`Custom` and `When` of the slice configurators get the elements as `shared.Elems[E]`, a view over `[]E` and `[]*E` without copying:

```go
b.Slice(&p.Tags).As(str.Slice).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, tags shared.Elems[string]) []shared.Error {
    seen := make(map[string]bool, tags.Len())
    for _, tag := range tags.All() { // tag is *string, nil for nil elements of []*string
        if tag != nil && seen[*tag] {
            return []shared.Error{h.ErrorT(ctx, *tag, "Tags should be unique")}
        }
        if tag != nil {
            seen[*tag] = true
        }
    }
    return nil
})
```

`Len()` returns the number of elements, `At(i)` the pointer to the element `i`, `All()` iterates over the indexes and the pointers.

## Slices of pointers

`Slice` works with the declared element type: for `[]*int` the element is `*int`. Convert the element to check the pointed values, nil elements are absent then:

```go
b.Slice(&p.Limits).Each(func(limit *shared.Field[*int]) {
    limit.As(num.Ptr).Required().Max(100)
})
```
