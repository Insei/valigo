package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/insei/valigo"
	"github.com/insei/valigo/translator"
)

type Sender struct {
	Type          string
	SMTPHost      string
	SMTPPort      string
	HTTPAddress   string
	HTTPDestParam string
}

func autoValidatorSettings() {
	validator := valigo.New()
	valigo.Configure[Sender](validator, func(builder valigo.Builder[Sender], obj *Sender) {
		builder.String(&obj.Type).Required()
	})
	sender := &Sender{
		Type:          "",
		SMTPHost:      uuid.New().String() + "   ",
		SMTPPort:      uuid.New().String() + " ",
		HTTPAddress:   uuid.New().String() + " ",
		HTTPDestParam: uuid.New().String() + "  ",
	}
	ctx := context.Background()
	errs := validator.Validate(ctx, sender)
	_ = errs
}

func manualValidatorSettings() {
	tStorage := translator.NewInMemStorage()
	// if you want to add new translation for you custom validators
	//translationStorage.AddTranslations()
	t := translator.New(translator.WithStorage(tStorage), translator.WithDefaultLang("en"))
	v := valigo.New(valigo.WithTranslator(t))
	valigo.Configure[Sender](v, func(builder valigo.Builder[Sender], obj *Sender) {
		builder.String(&obj.Type).Required()
	})
	sender := &Sender{
		Type:          "",
		SMTPHost:      uuid.New().String() + "   ",
		SMTPPort:      uuid.New().String() + " ",
		HTTPAddress:   uuid.New().String() + " ",
		HTTPDestParam: uuid.New().String() + "  ",
	}
	errs := v.Validate(context.Background(), sender)
	_ = errs
}

func main() {
	manualValidatorSettings()
}
