package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/insei/valigo"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/str"
	"github.com/insei/valigo/translator"
	vuuid "github.com/insei/valigo/uuid"
)

type Sender struct {
	Type          string
	SMTPHost      string
	SMTPPort      string
	HTTPAddress   string
	HTTPDestParam string
	Email         string
	Emails        []string
	Description   *string
	Int           int
	Id            uuid.UUID
	Templates     []string
	ClientIDs     []*uuid.UUID
}

const (
	customRegexpLocaleMsg = "Only numbers and words is allowed"
	customRegexpLocaleKey = "validation:string:" + customRegexpLocaleMsg
)

func manualValidatorSettings() *valigo.Validator {
	tStorage := translator.NewInMemStorage()
	// if you want to add new translation for you custom validators
	tStorage.Add("en", map[string]string{
		customRegexpLocaleKey: customRegexpLocaleMsg,
	})
	t := translator.New(translator.WithStorage(tStorage), translator.WithDefaultLang("en"))
	v := valigo.New(valigo.WithTranslator(t))
	return v
}

func main() {
	v := manualValidatorSettings() //v := valigo.New()
	allowedClientID := uuid.New()
	valigo.Configure[Sender](v, func(builder valigo.Configurator[Sender], obj *Sender) {
		builder.Field(&obj.Int).As(num.Of).AnyOf(2)
		builder.Field(&obj.Type).As(str.Of).When(func(ctx context.Context, value string) bool {
			return true
		}).Trim()
		builder.Field(&obj.Description).As(str.Of).When(func(ctx context.Context, value string) bool {
			return true
		}).Trim()
		builder.Field(&obj.SMTPHost).As(str.Of).Trim().
			Regexp(regexp.MustCompile("^[a-zA-Z0-9.]+$"), str.WithRegexpLocaleKey(customRegexpLocaleKey))
		builder.Field(&obj.Id).As(vuuid.Of).When(func(ctx context.Context, value uuid.UUID) bool {
			return true
		}).AnyOf(uuid.Nil)
		builder.Field(&obj.Email).As(str.Of).Email()
		builder.Slice(&obj.Emails).As(str.Slice).Email()
		builder.Slice(&obj.Templates).As(str.Slice).Trim().
			Regexp(regexp.MustCompile("^[a-zA-Z0-9.]+$"), str.WithRegexpLocaleKey(customRegexpLocaleKey))
		builder.Slice(&obj.ClientIDs).As(vuuid.Slice).AnyOf(allowedClientID)
	})
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	desc := "desc        "
	sender := &Sender{
		Type:          "type        ",
		Description:   &desc,
		Templates:     []string{"  correct  ", "incorrect&"},
		ClientIDs:     []*uuid.UUID{&id},
		SMTPHost:      uuid.New().String() + "   ",
		SMTPPort:      uuid.New().String() + " ",
		HTTPAddress:   uuid.New().String() + " ",
		HTTPDestParam: uuid.New().String() + "  ",
		Email:         "correct@email.com",
		Emails:        []string{"correct@email.com", "incorrect.email.com"},
		Id:            id,
		Int:           5,
	}

	errs := v.Validate(context.Background(), sender)
	errsJson, _ := json.Marshal(errs)
	fmt.Println(string(errsJson))
}
