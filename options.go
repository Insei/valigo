package valigo

import (
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/translator"
)

// Option is an interface that defines a single method, apply,
// which is used to apply an option to a Validator.
type Option interface {
	apply(v *Validator)
}

// optionFunc is a type that implements the Option interface.
// It is a function that takes a Validator as an argument.
type optionFunc func(*Validator)

// apply implements the Option interface for optionFunc.
// It calls the underlying function with the given Validator.
func (f optionFunc) apply(v *Validator) {
	f(v)
}

// WithTranslator returns an Option that sets the translator for the Validator.
func WithTranslator(t translator.Translator) Option {
	return optionFunc(func(v *Validator) {
		if t != nil {
			v.helper.t = t
		}
	})
}

// WithFieldLocationNamingFn returns an Option that sets the field location
// naming function for the Validator.
func WithFieldLocationNamingFn(fn func(field fmap.Field) string) Option {
	return optionFunc(func(v *Validator) {
		if fn != nil {
			v.helper.getFieldLocation = fn
		}
	})
}
