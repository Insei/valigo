package valigo

import (
	"log"
	"testing"
	"time"
)

type User struct {
	Name string
}

type Sender struct {
	Type          string
	SMTPHost      string
	SMTPPort      string
	HTTPAddress   string
	HTTPDestParam string
	PtrString     *string
}

func BenchmarkName(b *testing.B) {
	sender := &Sender{
		Type:        "SMTP",
		SMTPHost:    "123     ",
		SMTPPort:    "123    ",
		HTTPAddress: "123   ",
	}
	AddValidation[Sender](func(builder Builder[Sender], temp *Sender) {
		//smtpValidator := builder.When(func(obj *Sender) bool {
		//	return obj.Type == "SMTP"
		//})
		//smtpValidator.String(&temp.SMTPHost).Trim().Required()
		//smtpValidator.String(&temp.SMTPPort).Trim().Required()
		builder.StringPtr(&temp.PtrString).Trim().Required()

		//httpValidator := builder.When(func(temp *Sender) bool {
		//	return temp.Type == "HTTP"
		//})
		//httpValidator.String(&temp.HTTPAddress).Trim().Required()
		//httpValidator.String(&temp.HTTPDestParam).Trim().Required()
	})
	start := time.Now()
	for i := 0; i < b.N; i++ {
		_ = Validate(sender)
	}
	elapsed := time.Since(start)
	log.Printf("op/ns %d", elapsed.Nanoseconds()/int64(b.N))
}
