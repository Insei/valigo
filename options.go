package valigo

import (
	"github.com/insei/fmap/v3"
	"github.com/insei/valigo/translator"
)

type Option interface {
	apply(v *Validator)
}

type optionFunc func(*Validator)

func (f optionFunc) apply(v *Validator) {
	f(v)
}

func WithTranslator(t translator.Translator) Option {
	return optionFunc(func(v *Validator) {
		if t != nil {
			v.helper.Translator = t
		}
	})
}

func WithFieldLocationNamingFn(fn func(field fmap.Field) string) Option {
	return optionFunc(func(v *Validator) {
		if fn != nil {
			v.helper.GetFieldLocation = fn
		}
	})
}
