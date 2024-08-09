package valigo

import "github.com/insei/valigo/translator"

type Options struct {
	storage storage
}

type Option interface {
	apply(v *Validator)
}

type translatorOption struct {
	t translator.Translator
}

func (t translatorOption) apply(v *Validator) {
	if t.t != nil {
		v.helper.Translator = t.t
	}
}

func WithTranslator(t translator.Translator) Option {
	return &translatorOption{
		t: t,
	}
}
