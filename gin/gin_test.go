package gin_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dgin "github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insei/valigo"
	vgin "github.com/insei/valigo/gin"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/shared"
)

type SignUp struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func newValidator() *valigo.Validator {
	v := valigo.New()
	valigo.Configure[SignUp](v, func(b valigo.Configurator[SignUp], s *SignUp) {
		b.Field(&s.Email).As(str.Of).Required().Email()
		b.Field(&s.Name).As(str.Of).MaxLen(3)
	})
	return v
}

func TestValidateStruct(t *testing.T) {
	g := vgin.New(newValidator())
	assert.Nil(t, g.ValidateStruct(&SignUp{Email: "a@b.cc"}))
	assert.Nil(t, g.ValidateStruct(42), "not a struct")
	assert.Nil(t, g.ValidateStruct(nil))

	err := g.ValidateStruct(&SignUp{Email: "a", Name: "Anna"})
	var errs shared.Errors
	require.ErrorAs(t, err, &errs)
	assert.Equal(t, "Should be email address (Email: a); Cannot be longer than 3 characters (Name: Anna)", err.Error())

	err = g.ValidateStruct([]SignUp{{Email: "a@b.cc"}, {Email: "b"}})
	require.ErrorAs(t, err, &errs)
	assert.Equal(t, "[1].Email", errs[0].Location)

	err = g.ValidateStruct(&[]*SignUp{{Email: "c"}})
	require.ErrorAs(t, err, &errs)
	assert.Equal(t, "[0].Email", errs[0].Location)

	err = g.ValidateStruct([1]SignUp{{Email: "d"}})
	require.ErrorAs(t, err, &errs)
	assert.Equal(t, "[0].Email", errs[0].Location)

	assert.Same(t, g.Engine(), g.Engine())
	_, ok := g.Engine().(*valigo.Validator)
	assert.True(t, ok)
}

func request(t *testing.T, handler dgin.HandlerFunc, body string, headers ...string) *httptest.ResponseRecorder {
	t.Helper()
	dgin.SetMode(dgin.TestMode)
	r := dgin.New()
	r.POST("/sign-up", handler)
	req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestShouldBindJSON(t *testing.T) {
	previous := binding.Validator
	binding.Validator = vgin.New(newValidator())
	t.Cleanup(func() { binding.Validator = previous })

	w := request(t, func(c *dgin.Context) {
		var req SignUp
		if err := c.ShouldBindJSON(&req); err != nil {
			var errs shared.Errors
			if errors.As(err, &errs) {
				c.JSON(http.StatusUnprocessableEntity, errs)
				return
			}
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusCreated)
	}, `{"email": "ann", "name": "Ann"}`)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.JSONEq(t, `[{"Message": "Should be email address", "Location": "Email", "Value": "ann", "Code": "string.email"}]`, w.Body.String())
}

func TestValidateWithAcceptLanguage(t *testing.T) {
	v := newValidator()
	w := request(t, func(c *dgin.Context) {
		var req SignUp
		_ = c.ShouldBindBodyWith(&req, binding.JSON)
		if err := vgin.Validate(c, v, &req); err != nil {
			c.String(http.StatusUnprocessableEntity, err.Error())
			return
		}
		c.Status(http.StatusCreated)
	}, `{"email": "ann"}`, "Accept-Language", "ru")
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "Должно быть электронным адресом (Email: ann)", w.Body.String())

	w = request(t, func(c *dgin.Context) {
		if err := vgin.Validate(c, v, &SignUp{Email: "a@b.cc"}); err != nil {
			c.Status(http.StatusUnprocessableEntity)
			return
		}
		c.Status(http.StatusCreated)
	}, `{}`)
	assert.Equal(t, http.StatusCreated, w.Code)
}
