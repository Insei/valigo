package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/insei/valigo"
	"github.com/insei/valigo/str"
	"github.com/insei/valigo/translator"
	"regexp"
)

type Sender struct {
	Type          string
	SMTPHost      string
	SMTPPort      string
	HTTPAddress   string
	HTTPDestParam string
	Description   *string
	Int           int
	Id            uuid.UUID
	Templates     []string
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
	valigo.Configure[Sender](v, func(builder valigo.Configurator[Sender], obj *Sender) {
		builder.Number(&obj.Int).AnyOf(2)
		builder.String(&obj.Type).Trim()
		builder.String(&obj.SMTPHost).Trim().
			Regexp(regexp.MustCompile("^[a-zA-Z0-9.]+$"), str.WithRegexpLocaleKey(customRegexpLocaleKey))
		builder.UUID(&obj.Id).Required()
		builder.String(&obj.Description).Trim()
		builder.StringSlice(&obj.Templates).Trim().
			Regexp(regexp.MustCompile("^[a-zA-Z0-9.]+$"), str.WithRegexpLocaleKey(customRegexpLocaleKey))
	})
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	sender := &Sender{
		Type:          "123@123       ",
		Templates:     []string{"  correct  ", "incorrect&"},
		SMTPHost:      uuid.New().String() + "   ",
		SMTPPort:      uuid.New().String() + " ",
		HTTPAddress:   uuid.New().String() + " ",
		HTTPDestParam: uuid.New().String() + "  ",
		Id:            id,
		Int:           2,
	}
	errs := v.Validate(context.Background(), sender)
	errsJson, _ := json.Marshal(errs)
	fmt.Println(string(errsJson))
	fmt.Println(sender)
}
