# huma

The module `github.com/insei/valigo/huma` validates the inputs of [huma](https://huma.rocks) operations by valigo. The errors are the error details of huma with its locations — `body.email`, `query.limit`, `path.id` — in the 422 response of huma, and the messages are translated to the languages of the `Accept-Language` header.

```bash
go get github.com/insei/valigo/huma@latest
```

```go
import vhuma "github.com/insei/valigo/huma"

var validator = valigo.New(valigo.WithFieldLocationNamingFn(vhuma.Location))

func (i *CreateUserInput) Resolve(ctx huma.Context) []error {
	return vhuma.Resolve(ctx, validator, i)
}
```

- `vhuma.Location` is the location function of the input structs of huma: the fields of `Body` are `body` and their JSON path, the parameters are their kind and name from the `path`, `query`, `header`, `cookie` and `form` tags.
- `vhuma.Resolve` validates the input with the languages of the `Accept-Language` header and returns the error details of huma, huma calls `Resolve` of the input after it has parsed the request.
- `vhuma.Errors` converts the errors of `ValidateTyped` to the error details of huma. The error details of huma have no codes, use [`problem`](problem.md) when the clients need them.

huma validates the input by its JSON schema first, so the rules of valigo extend the schema: formats, rules on other fields, rules from data, translated messages.

## Example

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/insei/valigo"
	vhuma "github.com/insei/valigo/huma"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/str"
)

type CreateUserInput struct {
	Limit int `query:"limit"`
	Body  struct {
		Email string   `json:"email"`
		Tags  []string `json:"tags"`
	}
}

var validator = valigo.New(valigo.WithFieldLocationNamingFn(vhuma.Location))

func init() {
	valigo.Configure[CreateUserInput](validator, func(b valigo.Configurator[CreateUserInput], i *CreateUserInput) {
		b.Field(&i.Limit).As(num.Of).Max(100)
		b.Field(&i.Body.Email).As(str.Of).Email()
		b.Slice(&i.Body.Tags).As(str.Slice).Each(func(tag *str.BaseConfigurator) { tag.Alphanumeric() })
	})
}

func (i *CreateUserInput) Resolve(ctx huma.Context) []error {
	return vhuma.Resolve(ctx, validator, i)
}

func main() {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Users", "1.0.0"))
	huma.Post(api, "/users", func(context.Context, *CreateUserInput) (*struct{}, error) { return nil, nil })

	r := httptest.NewRequest(http.MethodPost, "/users?limit=500", strings.NewReader(`{"email": "ann", "tags": ["go", "a b"]}`))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept-Language", "ru")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	var resp struct {
		Status int `json:"status"`
		Errors []struct {
			Location string `json:"location"`
			Message  string `json:"message"`
		} `json:"errors"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	fmt.Println(resp.Status)
	for _, err := range resp.Errors {
		fmt.Println(err.Location, "-", err.Message)
	}
}
```

Output:

```text
422
query.limit - Не может быть больше 100
body.email - Должно быть электронным адресом
body.tags[1] - Должно содержать только латинские буквы и цифры
```
