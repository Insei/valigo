package problem_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/insei/fmap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insei/valigo"
	"github.com/insei/valigo/problem"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
)

type address struct {
	City string `json:"city"`
}

type signUp struct {
	Email   string   `json:"email"`
	Address address  `json:"address"`
	Tags    []string `json:"tags"`
}

func newValidator(opts ...problem.Option) *valigo.Validator {
	v := valigo.New(
		valigo.WithFieldLocationNamingFn(func(f fmap.Field) string { return f.GetTagPath("json", false) }),
		valigo.WithErrorsTransformer(problem.Transformer(opts...)),
	)
	valigo.Configure[signUp](v, func(b valigo.Configurator[signUp], s *signUp) {
		b.Field(&s.Email).As(str.Of).Required().Email()
		b.Field(&s.Address.City).As(str.Of).MaxLen(3)
		b.Slice(&s.Tags).As(str.Slice).Each(func(e *str.BaseConfigurator) { e.Alphanumeric() })
	})
	return v
}

func TestTransformerAndWrite(t *testing.T) {
	errs := newValidator().Validate(context.Background(), &signUp{Email: "a", Address: address{City: "Berlin"}, Tags: []string{"go", "a/b"}})
	require.Len(t, errs, 1)
	d := problem.From(errs)
	require.NotNil(t, d)

	rec := httptest.NewRecorder()
	require.NoError(t, d.Write(rec))
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Equal(t, "application/problem+json", rec.Header().Get("Content-Type"))
	assert.JSONEq(t, `{
		"title": "Validation failed",
		"status": 422,
		"errors": [
			{"detail": "Should be email address", "pointer": "#/email", "location": "email", "code": "string.email"},
			{"detail": "Cannot be longer than 3 characters", "pointer": "#/address/city", "location": "address.city", "code": "string.max_len", "params": {"max": 3}},
			{"detail": "Should contain only letters and digits", "pointer": "#/tags/1", "location": "tags[1]", "code": "string.alphanumeric"}
		]
	}`, rec.Body.String())
}

func TestOptions(t *testing.T) {
	errs := newValidator(
		problem.WithType("https://example.com/problems/validation"),
		problem.WithTitle("Invalid request"),
		problem.WithStatus(http.StatusBadRequest),
		problem.WithDetail("The request has invalid fields"),
	).Validate(context.Background(), &signUp{Email: "a@b.cc"})
	assert.Nil(t, errs, "valid objects give no problem")

	d := problem.New([]shared.Error{{Location: "email", Message: "bad"}},
		problem.WithType("https://example.com/problems/validation"), problem.WithTitle("Invalid request"),
		problem.WithStatus(http.StatusBadRequest), problem.WithDetail("The request has invalid fields"))
	body, err := json.Marshal(d)
	require.NoError(t, err)
	assert.JSONEq(t, `{"type": "https://example.com/problems/validation", "title": "Invalid request", "status": 400,
		"detail": "The request has invalid fields", "errors": [{"detail": "bad", "pointer": "#/email", "location": "email"}]}`, string(body))
	d.Instance = "/sign-up"
	body, _ = json.Marshal(d)
	assert.Contains(t, string(body), `"instance":"/sign-up"`)
}

func TestError(t *testing.T) {
	assert.Equal(t, "Validation failed", problem.New(nil).Error())
	assert.Equal(t, "Validation failed: email: bad", problem.New([]shared.Error{{Location: "email", Message: "bad"}}).Error())
	assert.Equal(t, "Validation failed: email: bad and other errors",
		problem.New([]shared.Error{{Location: "email", Message: "bad"}, {Location: "name", Message: "worse"}}).Error())
}

func TestFrom(t *testing.T) {
	assert.Nil(t, problem.From(nil))
	assert.Nil(t, problem.From([]error{errors.New("other")}))
	d := problem.New(nil)
	assert.Same(t, d, problem.From([]error{errors.New("other"), d}))
}

func TestPointer(t *testing.T) {
	for location, want := range map[string]string{
		"":              "#",
		"email":         "#/email",
		"address.city":  "#/address/city",
		"tags[1]":       "#/tags/1",
		"items[2].name": "#/items/2/name",
		"Labels[a/b~c]": "#/Labels/a~1b~0c",
		"Matrix[1][2]":  "#/Matrix/1/2",
	} {
		assert.Equal(t, want, problem.Pointer(location), location)
	}
}
