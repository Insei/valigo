# UUID

The rules of [`github.com/google/uuid`](https://pkg.go.dev/github.com/google/uuid) are in the separate module `github.com/insei/valigo/uuid`, so the core doesn't depend on a uuid package:

```bash
go get github.com/insei/valigo/uuid@latest
```

`uuid.Of` converts a `uuid.UUID` field or a `*uuid.UUID` field to the uuid configurator, `uuid.Slice` — a slice of them. The module's package is also named `uuid`, import it with an alias next to `github.com/google/uuid`:

```go
import (
	"github.com/google/uuid"
	vuuid "github.com/insei/valigo/uuid"
)

type Comment struct {
	ID       uuid.UUID
	ParentID *uuid.UUID
}

b.Field(&c.ID).As(vuuid.Of).Required()
b.Field(&c.ParentID).As(vuuid.Of).Custom(notDeleted) // a nil ParentID is skipped
b.Slice(&c.Mentions).As(vuuid.Slice).MaxLen(10).AnyOf(allowed...)
```

## Rules

| Rule | Fails when | Message |
|---|---|---|
| `Required()` | the value is absent or `uuid.Nil` | Should be fulfilled |
| `AnyOf(values...)` | the value isn't one of `values` | Only %s values is allowed |
| `Custom(fn)` | `fn` returns errors | your messages, see [Custom rules](custom.md) |
| `When(fn)` | — | the next rules are applied only if `fn` returns true, see [Conditions](conditions.md) |

Another uuid package gets its rules by a conversion of the same kind, see [Types and conversions](conversions.md#your-types).

## Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
	vuuid "github.com/insei/valigo/uuid"
)

type Comment struct {
	ID       uuid.UUID
	ParentID *uuid.UUID
}

func main() {
	v := valigo.New()
	valigo.Configure[Comment](v, func(b valigo.Configurator[Comment], c *Comment) {
		b.Field(&c.ID).As(vuuid.Of).Required().Custom(func(ctx context.Context, h *shared.FieldCustomHelper, id uuid.UUID) []shared.Error {
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
