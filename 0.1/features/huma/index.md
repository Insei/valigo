# huma

The module `github.com/insei/valigo/integrations/huma` validates the inputs of [huma](https://huma.rocks) operations by valigo. The errors are the error details of huma with its locations — `body.email`, `query.limit`, `path.id` — in the 422 response of huma, and the messages are translated to the languages of the `Accept-Language` header.

```bash
go get github.com/insei/valigo/integrations/huma@latest
```

```go
import vhuma "github.com/insei/valigo/integrations/huma"

var validator = valigo.New(valigo.WithFieldLocationNamingFn(vhuma.Location))

func (i *CreateUserInput) Resolve(ctx huma.Context) []error {
    return vhuma.Resolve(ctx, validator, i)
}
```

- `vhuma.Location` is the location function of the input structs of huma: the fields of `Body` are `body` and their JSON path, the parameters are their kind and name from the `path`, `query`, `header`, `cookie` and `form` tags.
- `vhuma.Resolve` validates the input with the languages of the `Accept-Language` header and returns the error details of huma, huma calls `Resolve` of the input after it has parsed the request.
- `vhuma.Errors` converts the errors of `ValidateTyped` to the error details of huma. The error details of huma have no codes, use [`problem`](https://insei.github.io/valigo/0.1/features/problem/index.md) when the clients need them.

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
    vhuma "github.com/insei/valigo/integrations/huma"
    "github.com/insei/valigo/rules/num"
    "github.com/insei/valigo/rules/str"
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

## The rules in the OpenAPI document

`vhuma.OpenAPI(api, v)` returns the OpenAPI document of the api with the rules of valigo in the schemas of the request bodies and their nested structs, by [JSON Schema](https://insei.github.io/valigo/0.1/features/schema/index.md): `maxLength` of `MaxRunes`, `format: email` of `Email`, `items` of `Each`, the required properties of `Required`. Serve it, so the clients and the generators of SDKs see the rules:

```go
mux.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
    doc, err := vhuma.OpenAPI(api, validator)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    _ = json.NewEncoder(w).Encode(doc)
})
```

with `config.OpenAPIPath = ""` instead of the document of huma, or under another path.

- The schemas of huma are not changed: huma validates the requests by its tags, the rules of valigo are checked in `Resolve` with translated messages, without the duplicates of huma.
- A keyword with another value than the one of the tags of huma, e.g. two `maxLength`, is added to `allOf`: the value must satisfy both.
- The parameters keep the schemas of huma: an operation of huma doesn't know its input type.
- The schemas are found by the registry of huma with the prefix `#/components/schemas/`, the default one.

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/danielgtaylor/huma/v2"
    "github.com/danielgtaylor/huma/v2/adapters/humago"
    "github.com/insei/valigo"
    vhuma "github.com/insei/valigo/integrations/huma"
    "github.com/insei/valigo/rules/str"
)

type CreateUserInput struct {
    Body struct {
        Name  string   `json:"name"`
        Email string   `json:"email,omitempty"`
        Tags  []string `json:"tags,omitempty"`
    }
}

func main() {
    validator := valigo.New(valigo.WithFieldLocationNamingFn(vhuma.Location))
    valigo.Configure[CreateUserInput](validator, func(b valigo.Configurator[CreateUserInput], i *CreateUserInput) {
        b.Field(&i.Body.Name).As(str.Of).MaxRunes(64)
        b.Field(&i.Body.Email).As(str.Of).Required().Email()
        b.Slice(&i.Body.Tags).As(str.Slice).MaxLen(5).Each(func(tag *str.BaseConfigurator) { tag.Alphanumeric() })
    })

    api := humago.New(http.NewServeMux(), huma.DefaultConfig("Users", "1.0.0"))
    huma.Post(api, "/users", func(context.Context, *CreateUserInput) (*struct{}, error) { return nil, nil })

    doc, err := vhuma.OpenAPI(api, validator)
    if err != nil {
        panic(err)
    }
    body := doc["components"].(map[string]any)["schemas"].(map[string]any)["CreateUserInputBody"].(map[string]any)
    delete(body["properties"].(map[string]any), "$schema") // the property of huma for the link to the schema
    out, _ := json.MarshalIndent(body, "", "  ")
    fmt.Println(string(out))
}
```

Output:

```text
{
  "additionalProperties": false,
  "properties": {
    "email": {
      "format": "email",
      "minLength": 1,
      "type": "string"
    },
    "name": {
      "maxLength": 64,
      "type": "string"
    },
    "tags": {
      "items": {
        "pattern": "^[a-zA-Z0-9]+$",
        "type": "string"
      },
      "maxItems": 5,
      "type": [
        "array",
        "null"
      ]
    }
  },
  "required": [
    "name",
    "email"
  ],
  "type": "object"
}
```

`email` is `omitempty` for huma, `Required` of valigo makes it required in the document; `tags` has the rules of its elements in `items`.
