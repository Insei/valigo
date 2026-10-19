package shared

import (
	"reflect"

	"github.com/insei/fmap/v3"
)

// The keywords of JSON Schema in Constraint.Path from the schema of a field to the schema of its parts.
const (
	// PathItems is the schema of the elements of a slice.
	PathItems = "items"
	// PathAdditionalProperties is the schema of the values of a map.
	PathAdditionalProperties = "additionalProperties"
	// PathPropertyNames is the schema of the keys of a map.
	PathPropertyNames = "propertyNames"
)

// KeywordRequired is the Keyword of the constraint of Required: the field is in the required properties of its object.
const KeywordRequired = "required"

// Constraint is a rule of a field that JSON Schema can express, e.g. maxLength 64 of MaxRunes(64),
// for the generators of the schemas of the API. Only the unconditional rules are described:
// the rules under When are not; after a Modify rule, e.g. Trim, only Required is, the schema describes the raw value.
// The constraints describe the rules approximately: valigo stays the validation of the values.
type Constraint struct {
	// Struct is the configured struct type.
	Struct reflect.Type
	// Field is the field of the rule, a field of Struct or of its nested structs.
	Field fmap.Field
	// Path is the keywords from the schema of the field to the schema of the constrained value:
	// [items] for the elements of a slice, [additionalProperties] for the values of a map, [propertyNames] for its keys.
	Path []string
	// Keyword is the keyword of JSON Schema: maxLength, minimum, enum, pattern, format, required...
	Keyword string
	// Value is the value of the keyword.
	Value any
}

// Describe registers the constraints of a configured struct, it is nil for conditional rules.
type Describe func(c Constraint)

// constraints describes the rules of one value, it is the zero value for the values without constraints.
type constraints struct {
	describe Describe
	field    fmap.Field
	path     []string
	// modified is true after a Modify rule, only Required is described then
	modified *bool
}

func (c constraints) add(keyword string, value any) {
	if c.describe == nil || (c.modified != nil && *c.modified && keyword != KeywordRequired) {
		return
	}
	c.describe(Constraint{Field: c.field, Path: c.path, Keyword: keyword, Value: value})
}

// child returns the constraints of the parts of the value by the keyword of the path, e.g. items.
func (c constraints) child(keyword string) constraints {
	if c.describe == nil {
		return constraints{}
	}
	return constraints{describe: c.describe, field: c.field, path: append(c.path[:len(c.path):len(c.path)], keyword), modified: new(bool)}
}

// Values converts values to []any, e.g. for the enum keyword.
func Values[T any](values []T) []any {
	result := make([]any, len(values))
	for i, v := range values {
		result[i] = v
	}
	return result
}
