// Package gin validates the requests of gin (github.com/gin-gonic/gin) by valigo.
//
// New is the binding.StructValidator of gin, the binding methods of the context, c.ShouldBindJSON and others,
// validate the requests by valigo:
//
//	binding.Validator = vgin.New(validator)
//
// gin doesn't pass the request to the validator, so its messages are in the default language;
// Validate validates with the languages of the Accept-Language header of the request.
package gin

import (
	"context"
	"reflect"
	"strconv"

	dgin "github.com/gin-gonic/gin"

	"github.com/insei/valigo"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/translator"
)

// Validator is the binding.StructValidator of gin by valigo.
type Validator struct {
	v *valigo.Validator
}

// New returns the binding.StructValidator of gin by the validator.
func New(v *valigo.Validator) *Validator {
	return &Validator{v: v}
}

// ValidateStruct validates a struct, a pointer to a struct or every element of a slice or an array,
// the errors of an element have its index in the location: [1].Email. It returns shared.Errors or nil.
func (g *Validator) ValidateStruct(obj any) error {
	return validate(obj, func(obj any) []shared.Error { return g.v.ValidateTyped(context.Background(), obj) })
}

// Engine returns the validator of valigo.
func (g *Validator) Engine() any {
	return g.v
}

// Validate validates obj like ValidateStruct, with the languages of the Accept-Language header of the request.
func Validate(c *dgin.Context, v *valigo.Validator, obj any) error {
	ctx := translator.ContextWithAcceptLanguage(c.Request.Context(), c.GetHeader("Accept-Language"))
	return validate(obj, func(obj any) []shared.Error { return v.ValidateTyped(ctx, obj) })
}

func validate(obj any, validateTyped func(obj any) []shared.Error) error {
	val := reflect.ValueOf(obj)
	for val.Kind() == reflect.Pointer && !val.IsNil() && (val.Elem().Kind() == reflect.Slice || val.Elem().Kind() == reflect.Array) {
		val = val.Elem()
	}
	var errs []shared.Error
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.Kind() != reflect.Pointer && elem.CanAddr() {
				elem = elem.Addr()
			}
			prefix := "[" + strconv.Itoa(i) + "]"
			for _, err := range validateTyped(elem.Interface()) {
				err.Location = prefix + "." + err.Location
				errs = append(errs, err)
			}
		}
	} else {
		errs = validateTyped(obj)
	}
	if len(errs) == 0 {
		return nil
	}
	return shared.Errors(errs)
}
