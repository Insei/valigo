// Package huma validates the inputs of huma (github.com/danielgtaylor/huma/v2) operations by valigo:
// the errors are the error details of huma with its locations, body.email or query.limit,
// and the messages are translated to the languages of the Accept-Language header.
//
//	var validator = valigo.New(valigo.WithFieldLocationNamingFn(vhuma.Location))
//
//	func (i *CreateUserInput) Resolve(ctx huma.Context) []error {
//		return vhuma.Resolve(ctx, validator, i)
//	}
package huma

import (
	"strings"

	dhuma "github.com/danielgtaylor/huma/v2"
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/translator"
)

// parameters are the tags of the parameters of huma inputs and the prefixes of their locations.
var parameters = []string{"path", "query", "header", "cookie", "form"}

// Location is the location function of valigo.WithFieldLocationNamingFn for the input structs of huma:
// the fields of Body are body and their JSON path, body.address.city, the parameters are their kind and name,
// query.limit, path.id, header.X-Request-Id. Other fields keep their struct path.
func Location(field fmap.Field) string {
	chain := []fmap.Field{field}
	for parent := field.GetParent(); parent != nil; parent = parent.GetParent() {
		chain = append(chain, parent)
	}
	root := chain[len(chain)-1]
	if len(chain) == 1 {
		for _, in := range parameters {
			if name, _, _ := strings.Cut(root.GetTag().Get(in), ","); name != "" {
				return in + "." + name
			}
		}
	}
	if root.GetName() != "Body" {
		return field.GetStructPath()
	}
	location := "body"
	// from the field of Body down to the field
	for i := len(chain) - 2; i >= 0; i-- {
		f := chain[i]
		name, _, _ := strings.Cut(f.GetTag().Get("json"), ",")
		if name == "" && f.GetAnonymous() {
			// embedded structs without a name are transparent, like in encoding/json
			continue
		}
		if name == "" {
			name = f.GetName()
		}
		location += "." + name
	}
	return location
}

// Errors converts the errors of valigo to the error details of huma, the response of huma has them in its errors.
func Errors(errs []shared.Error) []error {
	if len(errs) == 0 {
		return nil
	}
	details := make([]error, len(errs))
	for i, err := range errs {
		details[i] = &dhuma.ErrorDetail{Message: err.Message, Location: err.Location, Value: err.Value}
	}
	return details
}

// Resolve validates the input by the validator with the languages of the Accept-Language header of the request,
// call it from the Resolve method of the input of an operation.
func Resolve(ctx dhuma.Context, v *valigo.Validator, input any) []error {
	c := translator.ContextWithAcceptLanguage(ctx.Context(), ctx.Header("Accept-Language"))
	return Errors(v.ValidateTyped(c, input))
}
