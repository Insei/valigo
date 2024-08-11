package valigo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/insei/valigo/shared"
)

type Sender struct {
	Type          string
	SMTPHost      string
	SMTPPort      string
	HTTPAddress   string
	HTTPDestParam string
	PtrString     *string
}

func createSender() *Sender {
	ptr := uuid.New().String()
	return &Sender{
		Type:          "SMTP",
		SMTPHost:      uuid.New().String() + "   ",
		SMTPPort:      uuid.New().String() + " ",
		HTTPAddress:   uuid.New().String() + " ",
		HTTPDestParam: uuid.New().String() + "  ",
		PtrString:     &ptr,
	}
}

var (
	validator   = New()
	initialized = false
	sender      = createSender()
	ctx         = context.Background()
)

func benchValidateInit() {
	if initialized {
		return
	}
	initialized = true
	Configure[Sender](validator, func(builder Builder[Sender], temp *Sender) {
		builder.Custom(func(ctx context.Context, h shared.Helper, obj *Sender) []shared.Error {
			return nil
		})
		smtpValidator := builder.When(func(_ context.Context, obj *Sender) bool {
			return obj.Type == "SMTP"
		})
		smtpValidator.String(&temp.SMTPHost).Trim().Required()
		smtpValidator.String(&temp.SMTPPort).Trim().Required()

		builder.StringPtr(&temp.PtrString).Trim().Required()

		httpValidator := builder.When(func(_ context.Context, temp *Sender) bool {
			return temp.Type == "HTTP"
		})
		httpValidator.String(&temp.HTTPAddress).Trim().Required()
		httpValidator.String(&temp.HTTPDestParam).Trim().Required()
	})
}

func BenchmarkValidate(b *testing.B) {
	// initialize validation rules only one at bench start
	benchValidateInit()
	for i := 0; i < b.N; i++ {
		_ = validator.Validate(ctx, sender)
	}
}
