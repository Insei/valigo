package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/insei/valigo"
)

type Sender struct {
	Type          string
	SMTPHost      string
	SMTPPort      string
	HTTPAddress   string
	HTTPDestParam string
}

func main() {
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
