package main

import (
	"context"
	"regexp"

	"github.com/google/uuid"
	"github.com/insei/valigo"
	"github.com/insei/valigo/str"
	"github.com/insei/valigo/translator"
)

type Sender struct {
	Type          string
	SMTPHost      string
	SMTPPort      string
	HTTPAddress   string
	HTTPDestParam string
}

const (
	customRegexpLocaleMsg = "Only numbers and words is allowed"
	customRegexpLocaleKey = "validation:string:" + customRegexpLocaleMsg
)

func manualValidatorSettings() *valigo.Validator {
	tStorage := translator.NewInMemStorage()
	// if you want to add new translation for you custom validators
	tStorage.AddTranslations("en", map[string]string{
		customRegexpLocaleKey: customRegexpLocaleMsg,
	})
	t := translator.New(translator.WithStorage(tStorage), translator.WithDefaultLang("en"))
	v := valigo.New(valigo.WithTranslator(t))
	return v
}

func main() {
	v := manualValidatorSettings() //v := valigo.New()
	valigo.Configure[Sender](v, func(builder valigo.Builder[Sender], obj *Sender) {
		builder.String(&obj.Type).Required()
		builder.String(&obj.SMTPHost).
			Regexp(regexp.MustCompile("^[a-zA-Z0-9.]+$"), str.WithRegexpLocaleKey(customRegexpLocaleKey)).
			AnyOf("TEST", "TEST1")
	})
	sender := &Sender{
		Type:          "123@123",
		SMTPHost:      uuid.New().String() + "   ",
		SMTPPort:      uuid.New().String() + " ",
		HTTPAddress:   uuid.New().String() + " ",
		HTTPDestParam: uuid.New().String() + "  ",
	}
	errs := v.Validate(context.Background(), sender)
	_ = errs
}
