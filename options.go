package valigo

type Options struct {
	storage storage
}

type Option interface {
	apply(v *Validator)
}

type translatorOption struct {
	t Translator
}

func (t translatorOption) apply(v *Validator) {
	if t.t != nil {
		v.helper.Translator = t.t
	}
}

func WithTranslator(t Translator) Option {
	return &translatorOption{
		t: t,
	}
}
