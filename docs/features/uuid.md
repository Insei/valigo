# UUID

`UUID` configures the rules of a [`uuid.UUID`](https://pkg.go.dev/github.com/google/uuid) field or a `*uuid.UUID` field:

```go
type Comment struct {
	ID       uuid.UUID
	ParentID *uuid.UUID
}

b.UUID(&c.ID).Required()
b.UUID(&c.ParentID).Custom(notDeleted) // a nil ParentID is skipped
```

## Rules

| Rule | Fails when | Message |
|---|---|---|
| `Required()` | the value is absent or `uuid.Nil` | Should be fulfilled |
| `AnyOf(values...)` | the value isn't one of `values` | Only %s values is allowed |
| `Custom(fn)` | `fn` returns errors | your messages, see [Custom rules](custom.md) |
| `When(fn)` | — | the next rules are applied only if `fn` returns true, see [Conditions](conditions.md) |

Slices of uuids are configured by `UUIDSlice`, see [Slices](slices.md).

## Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
)

type Comment struct {
	ID       uuid.UUID
	ParentID *uuid.UUID
}

func main() {
	v := valigo.New()
	valigo.Configure[Comment](v, func(b valigo.Configurator[Comment], c *Comment) {
		b.UUID(&c.ID).Required().Custom(func(ctx context.Context, h *shared.FieldCustomHelper, id uuid.UUID) []shared.Error {
			if id.Version() != 7 {
				return []shared.Error{h.ErrorT(ctx, id, "Should be UUID v7")}
			}
			return nil
		})
	})

	fmt.Println(v.Validate(context.Background(), &Comment{}))
	fmt.Println(v.Validate(context.Background(), &Comment{ID: uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")}))
}
```

Output:

```text
[Should be fulfilled (ID: 00000000-0000-0000-0000-000000000000) Should be UUID v7 (ID: 00000000-0000-0000-0000-000000000000)]
[Should be UUID v7 (ID: 6ba7b810-9dad-11d1-80b4-00c04fd430c8)]
```
