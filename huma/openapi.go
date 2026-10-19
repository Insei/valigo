package huma

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"

	dhuma "github.com/danielgtaylor/huma/v2"
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
)

// schemasPrefix is the prefix of the references of the schemas of the registry of huma by default.
const schemasPrefix = "#/components/schemas/"

// OpenAPI returns the OpenAPI document of the api as JSON values with the rules of the validator in the schemas
// of the request bodies and their nested structs: maxLength of MaxRunes, minimum of Min, the required properties...,
// see Validator.Constraints. Serve it instead of the document of huma, so the clients see the rules of valigo.
//
// The schemas of huma, which huma validates the requests by, are not changed, the validation stays by valigo
// with translated messages. The parameters keep the schemas of huma.
func OpenAPI(api dhuma.API, v *valigo.Validator) (map[string]any, error) {
	oapi := api.OpenAPI()
	data, err := json.Marshal(oapi)
	if err != nil {
		return nil, err
	}
	var doc map[string]any
	_ = json.Unmarshal(data, &doc) // the JSON of json.Marshal
	byStruct := map[reflect.Type][]shared.Constraint{}
	for _, c := range v.Constraints() {
		if owner, ok := ownerOf(c); ok {
			byStruct[owner] = append(byStruct[owner], c)
		}
	}
	components, _ := doc["components"].(map[string]any)
	schemas, _ := components["schemas"].(map[string]any)
	for name, s := range schemas {
		constraints := byStruct[oapi.Components.Schemas.TypeFromRef(schemasPrefix+name)]
		for _, c := range constraints {
			apply(s.(map[string]any), c)
		}
	}
	return doc, nil
}

// ownerOf returns the struct type of the schema the field of the constraint is a property of:
// the configured struct or the nested struct of the field, the embedded structs without a JSON name are flattened.
func ownerOf(c shared.Constraint) (reflect.Type, bool) {
	if _, ok := propertyName(c.Field); !ok {
		return nil, false
	}
	owner := c.Field.GetParent()
	for owner != nil && owner.GetAnonymous() && jsonName(owner) == "" {
		owner = owner.GetParent()
	}
	if owner == nil {
		return c.Struct, true
	}
	return owner.GetDereferencedType(), true
}

func jsonName(f fmap.Field) string {
	name, _, _ := strings.Cut(f.GetTag().Get("json"), ",")
	return name
}

// propertyName returns the name of the property of the field in the schema of its struct, false for skipped fields.
func propertyName(f fmap.Field) (string, bool) {
	name := jsonName(f)
	if name == "-" {
		return "", false
	}
	if name == "" {
		name = f.GetName()
	}
	return name, true
}

// apply adds the constraint to the schema of the struct of its field.
func apply(schema map[string]any, c shared.Constraint) {
	name, _ := propertyName(c.Field)
	properties, _ := schema["properties"].(map[string]any)
	target, ok := properties[name].(map[string]any)
	if !ok {
		return
	}
	if c.Keyword == shared.KeywordRequired && len(c.Path) == 0 {
		required, _ := schema["required"].([]any)
		if !slices.Contains(required, any(name)) {
			schema["required"] = append(required, name)
		}
		return
	}
	if c.Keyword == shared.KeywordRequired {
		// the elements of a slice and the values of a map can't be required
		return
	}
	for _, keyword := range c.Path {
		next, ok := target[keyword].(map[string]any)
		if !ok {
			// additionalProperties: true, or no propertyNames yet
			next = map[string]any{}
			target[keyword] = next
		}
		target = next
	}
	value := normalize(c.Value)
	existing, ok := target[c.Keyword]
	switch {
	case !ok:
		target[c.Keyword] = value
	case !reflect.DeepEqual(existing, value):
		// e.g. two patterns or maxLength of huma and of valigo: the value must satisfy both
		allOf, _ := target["allOf"].([]any)
		target["allOf"] = append(allOf, map[string]any{c.Keyword: value})
	}
}

// normalize converts the value to the JSON values of the document, e.g. ints to float64 and named strings to strings.
func normalize(v any) any {
	data, _ := json.Marshal(v) // the values of the constraints are JSON values
	var result any
	_ = json.Unmarshal(data, &result)
	return result
}
