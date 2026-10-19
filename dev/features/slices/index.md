# Slices

| Configurator  | Fields                                                         | Rules on elements                                                                                                                 |
| ------------- | -------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| `Slice`       | `[]E` of any `E`                                               | `Each` with `*shared.Field[E]`: `Check`, `Required`, `Modify`, `Custom`, `When`                                                   |
| `StringSlice` | `[]string`, `*[]string`, `[]*string`, `*[]*string`             | `Each` with the [string configurator](https://insei.github.io/valigo/latest/features/strings/index.md), `Trim`, `Regexp`, `Email` |
| `UUIDSlice`   | `[]uuid.UUID`, `*[]uuid.UUID`, `[]*uuid.UUID`, `*[]*uuid.UUID` | `Each` with the [uuid configurator](https://insei.github.io/valigo/latest/features/uuid/index.md), `AnyOf`                        |

## The slice and its elements

The rules of the configurator check **the slice itself**, the rules inside `Each` check **every element**:

| Rule         | Fails when                           | Message                                   |
| ------------ | ------------------------------------ | ----------------------------------------- |
| `Required()` | the slice is absent or empty         | Should be fulfilled                       |
| `MinLen(n)`  | the slice has less than `n` elements | Cannot contain less than %d elements      |
| `MaxLen(n)`  | the slice has more than `n` elements | Cannot contain more than %d elements      |
| `Each(fn)`   | —                                    | `fn` configures the rules of the elements |

`Trim`, `Regexp` and `Email` of `StringSlice` and `AnyOf` of `UUIDSlice` are shortcuts for the same rules in `Each`.

```go
package main

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/insei/valigo"
    "github.com/insei/valigo/shared"
    "github.com/insei/valigo/str"
    vuuid "github.com/insei/valigo/uuid"
)

type Post struct {
    Tags       []string
    Editors    []*uuid.UUID
    Quantities []int
}

func main() {
    v := valigo.New()
    valigo.Configure[Post](v, func(b valigo.Configurator[Post], p *Post) {
        b.StringSlice(&p.Tags).Required().MaxLen(3).Each(func(tag *str.BaseConfigurator) {
            tag.Trim().Required().MaxLen(10)
        })
        b.UUIDSlice(&p.Editors).Each(func(id *vuuid.BaseConfigurator) {
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

The location of an element error has the index of the element, the errors of the slice itself have the location of the field. Nil elements of slices of pointers are [absent](https://insei.github.io/valigo/latest/features/absent/index.md): only `Required` fails on them.

## Custom and When

`Custom` and `When` of the slice configurators get the elements as `shared.Elems[E]`, a view over `[]E` and `[]*E` without copying:

```go
b.StringSlice(&p.Tags).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, tags shared.Elems[string]) []shared.Error {
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

## Slices of pointers in Slice

`Slice` works with the declared element type: for `[]*int` the element is `*int`, and nil elements are regular values for its rules:

```go
b.Slice(&p.Limits).Each(func(limit *shared.Field[*int]) {
    limit.Check(func(v *int) bool { return v != nil }, "Should be set")
})
```
