package echo_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	decho "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insei/valigo"
	vecho "github.com/insei/valigo/echo"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/shared"
)

type SignUp struct {
	Email string `json:"email"`
}

func newValidator() *valigo.Validator {
	v := valigo.New()
	valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
		b.Field(&s.Email).As(str.Of).Required().Email()
	})
	return v
}

func serve(t *testing.T, e *decho.Echo, handler decho.HandlerFunc, body string, headers ...string) *httptest.ResponseRecorder {
	t.Helper()
	e.POST("/sign-up", handler)
	req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(body))
	req.Header.Set(decho.HeaderContentType, decho.MIMEApplicationJSON)
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	return w
}

func TestValidator(t *testing.T) {
	e := decho.New()
	e.Validator = vecho.New(newValidator())
	w := serve(t, e, func(c decho.Context) error {
		var req SignUp
		if err := c.Bind(&req); err != nil {
			return err
		}
		if err := c.Validate(&req); err != nil {
			var errs shared.Errors
			require.True(t, errors.As(err, &errs))
			return c.JSON(http.StatusUnprocessableEntity, errs)
		}
		return c.NoContent(http.StatusCreated)
	}, `{"email": "ann"}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `[{"Message": "Should be email address", "Location": "Email", "Value": "ann", "Code": "string.email"}]`, w.Body.String())

	assert.Nil(t, vecho.New(newValidator()).Validate(&SignUp{Email: "a@b.cc"}))
}

func TestValidateWithAcceptLanguage(t *testing.T) {
	v := newValidator()
	handler := func(c decho.Context) error {
		var req SignUp
		if err := c.Bind(&req); err != nil {
			return err
		}
		if err := vecho.Validate(c, v, &req); err != nil {
			return c.String(http.StatusUnprocessableEntity, err.Error())
		}
		return c.NoContent(http.StatusCreated)
	}
	w := serve(t, decho.New(), handler, `{"email": "ann"}`, "Accept-Language", "es")
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "Debe ser una dirección de correo electrónico (Email: ann)", w.Body.String())

	w = serve(t, decho.New(), handler, `{"email": "a@b.cc"}`)
	assert.Equal(t, http.StatusCreated, w.Code)
}
