# gin

The module `github.com/insei/valigo/integrations/gin` validates the requests of [gin](https://gin-gonic.com) by valigo.

```bash
go get github.com/insei/valigo/integrations/gin@latest
```

- `vgin.New(v)` is the `binding.StructValidator` of gin: set `binding.Validator = vgin.New(v)`, and `c.ShouldBindJSON` and the other binding methods validate the requests by the rules of valigo. Slices are validated element by element, the location of an error has the index: `[1].email`.
- gin doesn't pass the request to its validator, so these messages are in the default language. `vgin.Validate(c, v, &req)` validates with the languages of the `Accept-Language` header of the request.
- The error is [`shared.Errors`](https://insei.github.io/valigo/0.1/features/errors/#errors-transformer), get the errors by `errors.As` and write them, e.g. as [problem details](https://insei.github.io/valigo/0.1/features/problem/index.md).

## Example

```go
package main

import (
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/insei/fmap/v3"
    "github.com/insei/valigo"
    vgin "github.com/insei/valigo/integrations/gin"
    "github.com/insei/valigo/problem"
    "github.com/insei/valigo/rules/str"
    "github.com/insei/valigo/shared"
)

type SignUp struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

func main() {
    v := valigo.New(valigo.WithFieldLocationNamingFn(func(f fmap.Field) string { return f.GetTagPath("json", false) }))
    valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
        b.Field(&s.Email).As(str.Of).Required().Email()
        b.Field(&s.Name).As(str.Of).Required().MaxRunes(16)
    })
    // the binding of gin validates by valigo, without the validation of its own tags
    binding.Validator = nil

    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.POST("/sign-up", func(c *gin.Context) {
        var req SignUp
        if err := c.ShouldBindJSON(&req); err != nil {
            c.Status(http.StatusBadRequest)
            return
        }
        if err := vgin.Validate(c, v, &req); err != nil {
            var errs shared.Errors
            if errors.As(err, &errs) {
                _ = problem.New(errs).Write(c.Writer)
                return
            }
        }
        c.Status(http.StatusCreated)
    })

    req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(`{"email": "ann", "name": ""}`))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept-Language", "de")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    fmt.Println(w.Code)
    fmt.Print(w.Body.String())
}
```

Output:

```text
422
{"title":"Validation failed","status":422,"errors":[{"detail":"Muss eine E-Mail-Adresse sein","pointer":"#/email","location":"email","code":"string.email"},{"detail":"Muss ausgefüllt sein","pointer":"#/name","location":"name","code":"string.required"}]}
```

To validate by the binding methods instead, set `binding.Validator = vgin.New(v)`: the messages are in the default language then.
