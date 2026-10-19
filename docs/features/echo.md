# echo

The module `github.com/insei/valigo/echo` validates the requests of [echo](https://echo.labstack.com) by valigo.

```bash
go get github.com/insei/valigo/echo@latest
```

- `vecho.New(v)` is the `echo.Validator`: set `e.Validator = vecho.New(v)`, and `c.Validate(&req)` validates by the rules of valigo.
- echo doesn't pass the request to its validator, so these messages are in the default language. `vecho.Validate(c, v, &req)` validates with the languages of the `Accept-Language` header of the request.
- The error is [`shared.Errors`](errors.md#errors-transformer), get the errors by `errors.As` and write them, e.g. as [problem details](problem.md).

## Example

```go
package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/insei/valigo"
	vecho "github.com/insei/valigo/echo"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
	"github.com/labstack/echo/v4"
)

type SignUp struct {
	Email string `json:"email"`
}

func main() {
	v := valigo.New()
	valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
		b.Field(&s.Email).As(str.Of).Required().Email()
	})

	e := echo.New()
	e.Validator = vecho.New(v)
	e.POST("/sign-up", func(c echo.Context) error {
		var req SignUp
		if err := c.Bind(&req); err != nil {
			return err
		}
		if err := c.Validate(&req); err != nil {
			var errs shared.Errors
			if errors.As(err, &errs) {
				return c.JSON(http.StatusUnprocessableEntity, errs)
			}
			return err
		}
		return c.NoContent(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(`{"email": "ann"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	fmt.Println(w.Code)
	fmt.Print(w.Body.String())
}
```

Output:

```text
422
[{"Message":"Should be email address","Location":"Email","Value":"ann","Code":"string.email"}]
```
