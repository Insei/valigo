package main

import (
	"log"
	"time"

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

func init() {
	valigo.AddValidation[Sender](func(builder valigo.Builder[Sender], obj *Sender) {
		builder.String(&obj.Type).AnyOf("SMTP", "HTTP")
		smtpValidator := builder.When(func(obj *Sender) bool {
			return obj.Type == "SMTP"
		})
		smtpValidator.String(&obj.SMTPHost).Trim().Required()
		smtpValidator.String(&obj.SMTPPort).Trim().Required()

		httpValidator := builder.When(func(obj *Sender) bool {
			return obj.Type == "HTTP"
		})
		httpValidator.String(&obj.HTTPAddress).Trim().Required()
		httpValidator.String(&obj.HTTPDestParam).Trim().Required()
	})
}

func iterate() {
	var iterations int64 = 9999999
	start := time.Now()
	for i := 0; i <= int(iterations); i++ {
		sender := &Sender{
			Type:          "SMTP",
			SMTPHost:      uuid.New().String() + "   ",
			SMTPPort:      uuid.New().String() + " ",
			HTTPAddress:   uuid.New().String() + " ",
			HTTPDestParam: uuid.New().String() + "  ",
		}
		_ = valigo.Validate(sender)
	}
	elapsed := time.Since(start)
	log.Printf("op/ns %d", elapsed.Nanoseconds()/iterations)
}

func main() {
	for i := 0; i < 10; i++ {
		iterate()
	}
}
