// package main created only for benchmarking without go test bench.
// Seems that go bench works not fully correctly in our case, we have unreal op/s increases in each iteration of bench test.
// pprof is not shows any CPU or MEM bad things, from bench memdump and cpudump.
// In all operations we have zero allocations.
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

func addValidation() {
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

func validate(data []Sender) {
	start := time.Now()
	for i, _ := range data {
		_ = valigo.Validate(&data[i])
	}
	elapsed := time.Since(start)
	log.Printf("ops %d op/ns %d", len(data), elapsed.Nanoseconds()/int64(len(data)))
}

func createData(iterations int64) []Sender {
	senders := make([]Sender, 0, iterations)
	for i := 0; i < int(iterations); i++ {
		senders = append(senders, Sender{
			Type:          "SMTP",
			SMTPHost:      uuid.New().String() + "   ",
			SMTPPort:      uuid.New().String() + " ",
			HTTPAddress:   uuid.New().String() + " ",
			HTTPDestParam: uuid.New().String() + "  ",
		})
	}
	return senders
}

func main() {
	iterations := int64(99999)
	addValidation()
	for i := 0; i < 200000; i++ {
		senders := createData(iterations)
		validate(senders)
	}
}
