package huma_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	dhuma "github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/insei/fmap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insei/valigo"
	vhuma "github.com/insei/valigo/integrations/huma"
	"github.com/insei/valigo/rules/num"
	"github.com/insei/valigo/rules/str"
)

type Address struct {
	City string `json:"city"`
}

type Meta struct {
	Source string `json:"source"`
}

type CreateUserInput struct {
	ID        string `path:"id"`
	Limit     int    `query:"limit"`
	RequestID string `header:"X-Request-Id"`
	Body      struct {
		Meta
		Email   string   `json:"email"`
		Address Address  `json:"address"`
		Tags    []string `json:"tags"`
		Note    string
	}
}

var validator = valigo.New(valigo.WithFieldLocationNamingFn(vhuma.Location))

func init() {
	valigo.Configure[CreateUserInput](validator, func(b valigo.Configurator[CreateUserInput], i *CreateUserInput) {
		b.Field(&i.ID).As(str.Of).UUID()
		b.Field(&i.Limit).As(num.Of).Max(100)
		b.Field(&i.RequestID).As(str.Of).Alphanumeric()
		b.Field(&i.Body.Email).As(str.Of).Email()
		b.Field(&i.Body.Address.City).As(str.Of).Required()
		b.Slice(&i.Body.Tags).As(str.Slice).Each(func(tag *str.BaseConfigurator) { tag.MaxLen(3) })
		b.Field(&i.Body.Source).As(str.Of).AnyOf("web", "api")
		b.Field(&i.Body.Note).As(str.Of).MaxLen(1)
	})
}

func (i *CreateUserInput) Resolve(ctx dhuma.Context) []error {
	return vhuma.Resolve(ctx, validator, i)
}

type errorModel struct {
	Status int `json:"status"`
	Errors []struct {
		Message  string `json:"message"`
		Location string `json:"location"`
	} `json:"errors"`
}

func call(t *testing.T, headers ...any) errorModel {
	t.Helper()
	_, api := humatest.New(t)
	dhuma.Post(api, "/users/{id}", func(context.Context, *CreateUserInput) (*struct{}, error) { return nil, nil })
	body := map[string]any{"email": "ann", "address": map[string]any{"city": ""}, "tags": []string{"go", "long"}, "source": "cli", "Note": "ab"}
	args := append([]any{"X-Request-Id: req-1"}, headers...)
	args = append(args, body)
	resp := api.Post("/users/abc?limit=500", args...)
	require.Equal(t, http.StatusUnprocessableEntity, resp.Code, resp.Body.String())
	var model errorModel
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &model))
	return model
}

func TestResolve(t *testing.T) {
	model := call(t)
	got := map[string]string{}
	for _, err := range model.Errors {
		got[err.Location] = err.Message
	}
	assert.Equal(t, map[string]string{
		"path.id":             "Should be UUID",
		"query.limit":         "Cannot be greater than 100",
		"header.X-Request-Id": "Should contain only letters and digits",
		"body.email":          "Should be email address",
		"body.address.city":   "Should be fulfilled",
		"body.tags[1]":        "Cannot be longer than 3 characters",
		"body.source":         "Only [web api] values is allowed",
		"body.Note":           "Cannot be longer than 1 character",
	}, got)
}

func TestResolveTranslatesByAcceptLanguage(t *testing.T) {
	model := call(t, "Accept-Language: ru")
	for _, err := range model.Errors {
		if err.Location == "body.email" {
			assert.Equal(t, "Должно быть электронным адресом", err.Message)
			return
		}
	}
	t.Fatal("no error of body.email")
}

func TestLocationOutsideOfInputs(t *testing.T) {
	type plain struct {
		Inner struct {
			Name string
		}
		Tagged string `query:",omitempty"`
	}
	fields, err := fmap.Get[plain]()
	require.NoError(t, err)
	assert.Equal(t, "Inner.Name", vhuma.Location(fields.MustFind("Inner.Name")))
	assert.Equal(t, "Tagged", vhuma.Location(fields.MustFind("Tagged")), "a parameter tag without a name")

	type withBody struct {
		Body string
	}
	bodyFields, err := fmap.Get[withBody]()
	require.NoError(t, err)
	assert.Equal(t, "body", vhuma.Location(bodyFields.MustFind("Body")))
}

func TestErrors(t *testing.T) {
	assert.Nil(t, vhuma.Errors(nil))
	errs := vhuma.Errors(validator.ValidateTyped(context.Background(), &CreateUserInput{ID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8"}))
	require.NotEmpty(t, errs)
	var detail dhuma.ErrorDetailer
	require.ErrorAs(t, errs[0], &detail)
	assert.True(t, strings.HasPrefix(detail.ErrorDetail().Location, "header.") || strings.HasPrefix(detail.ErrorDetail().Location, "body."))
}
