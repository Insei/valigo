# HTTP problem details

The `problem` package writes the validation errors as the problem details of [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) (it obsoletes RFC 7807): an `application/problem+json` response with the errors of the fields in the `errors` member.

```go
v := valigo.New(valigo.WithErrorsTransformer(problem.Transformer()))

if errs := v.Validate(r.Context(), &req); errs != nil {
	problem.From(errs).Write(w) // 422 Unprocessable Content
	return
}
```

- `problem.Transformer(opts...)` makes `Validate` return one error, the `*problem.Details` of all errors; `problem.From(errs)` gets it back, `problem.New(errs)` builds it from the errors of `ValidateTyped`.
- `Details.Write(w)` writes the status and the JSON with the `application/problem+json` content type.
- Every error of a field has the translated message `detail`, the JSON Pointer `pointer` of RFC 6901 (built from the location: `tags[1]` is `#/tags/1`), the `location`, and the [`code` and `params`](errors.md#codes-for-clients) for the clients that translate the messages themselves.
- Options: `WithType` (the URI of the problem type, `about:blank` when it is empty), `WithTitle` ("Validation failed" by default), `WithStatus` (422 by default), `WithDetail`.

Use JSON locations, so the pointers point to the fields of the request, see [Errors](errors.md#location-of-fields).

## Example

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/insei/fmap/v3"
	"github.com/insei/valigo"
	"github.com/insei/valigo/problem"
	"github.com/insei/valigo/str"
	"github.com/insei/valigo/translator"
)

type SignUp struct {
	Email string   `json:"email"`
	Tags  []string `json:"tags"`
}

func main() {
	v := valigo.New(
		valigo.WithFieldLocationNamingFn(func(f fmap.Field) string { return f.GetTagPath("json", false) }),
		valigo.WithErrorsTransformer(problem.Transformer(problem.WithType("https://example.com/problems/validation"))),
	)
	valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
		b.Field(&s.Email).As(str.Of).Required().Email()
		b.Slice(&s.Tags).As(str.Slice).MaxLen(2).Each(func(tag *str.BaseConfigurator) { tag.Alphanumeric() })
	})

	handler := translator.NewAcceptLanguageMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req SignUp
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errs := v.Validate(r.Context(), &req); errs != nil {
			_ = problem.From(errs).Write(w)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))

	r := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(`{"email": "ann", "tags": ["go", "a b", "api"]}`))
	r.Header.Set("Accept-Language", "ru")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	var body bytes.Buffer
	_ = json.Indent(&body, w.Body.Bytes(), "", "  ")
	fmt.Println(w.Code, w.Header().Get("Content-Type"))
	fmt.Print(body.String())
}
```

Output:

```text
422 application/problem+json
{
  "type": "https://example.com/problems/validation",
  "title": "Validation failed",
  "status": 422,
  "errors": [
    {
      "detail": "Должно быть электронным адресом",
      "pointer": "#/email",
      "location": "email",
      "code": "string.email"
    },
    {
      "detail": "Не может содержать больше 2 элементов",
      "pointer": "#/tags",
      "location": "tags",
      "code": "slice.max_len",
      "params": {
        "max": 2
      }
    },
    {
      "detail": "Должно содержать только латинские буквы и цифры",
      "pointer": "#/tags/1",
      "location": "tags[1]",
      "code": "string.alphanumeric"
    }
  ]
}
```

The title isn't translated: set it by `WithTitle`, e.g. with a translator of your messages.
