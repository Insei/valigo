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

func iterate(iterations int64) {
	start := time.Now()
	for i := 0; i <= int(iterations); i++ {
		sender := &Sender{
			Type:          "SMTP",
			SMTPHost:      uuid.New().String() + "   ",
			SMTPPort:      uuid.New().String() + " ",
			HTTPAddress:   uuid.New().String() + " ",
			HTTPDestParam: uuid.New().String() + "  ",
		}
		//start := time.Now()
		//if i == 59999999 || i == 0 {
		//	_ = Validate(sender)
		//}
		_ = valigo.Validate(sender)
		//elapsed := time.Since(start)
		//fmt.Println(elapsed)
	}
	elapsed := time.Since(start)
	log.Printf("ops %d op/ns %d", iterations, elapsed.Nanoseconds()/iterations)
}

func main() {
	iterations := int64(99999)
	//diff := int64(10304)
	for i := 0; i < 200000; i++ {
		//if iterations < diff {
		//	break
		//}
		iterate(iterations)
		//iterations -= diff
	}
}
