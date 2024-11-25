package valigo

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/google/uuid"

	"github.com/insei/valigo/shared"
)

type Sender struct {
	Type          string
	SMTPHost      string
	SMTPPort      string
	SMTPInt       int
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
		SMTPInt:       10,
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
	Configure[Sender](validator, func(builder Configurator[Sender], temp *Sender) {
		builder.String(&temp.HTTPDestParam).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error {
			v := *value.(*string)
			if len(v) < 20 {
				return []shared.Error{h.ErrorT(ctx, v, "Should be longer than 20 characters")}
			}
			return nil
		})
		// Custom validation on struct gives 1 alloc per operation
		//builder.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *Sender) []shared.Error {
		//	if obj.Type == "" {
		//		return []shared.Error{h.ErrorT(ctx, &obj.Type, obj.Type, "Should be fulfilled")}
		//	}
		//	return nil
		//})
		smtpValidator := builder.When(func(_ context.Context, obj *Sender) bool {
			return obj.Type == "SMTP"
		})
		smtpValidator.String(&temp.SMTPHost).Trim().Required()
		smtpValidator.String(&temp.SMTPPort).Trim().Required()
		smtpValidator.Number(&temp.SMTPInt).Max(20)
		builder.String(&temp.PtrString).Trim().Required()
		//builder.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *Sender) []shared.Error {
		//	h.ErrorT(ctx)
		//})
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

func TestValidatorValidateTyped(t *testing.T) {
	tests := []struct {
		name         string
		obj          any
		validator    func(ctx context.Context, h shared.Helper, obj any) []shared.Error
		expectedErrs int
	}{
		{
			name: "single error",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return []shared.Error{{Message: "test error"}}
			},
			expectedErrs: 1,
		},
		{
			name: "multiple errors",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return []shared.Error{{Message: "test error 1"}, {Message: "test error 2"}}
			},
			expectedErrs: 2,
		},
		{
			name: "no errors",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return nil
			},
			expectedErrs: 0,
		},
		{
			name: "nil validator",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return nil
			},
			expectedErrs: 0,
		},
		{
			name: "validator returns empty slice",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return []shared.Error{}
			},
			expectedErrs: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := New()
			v.storage.validators[reflect.TypeOf(test.obj)] = []structValidationFn{test.validator}
			errs := v.ValidateTyped(context.Background(), test.obj)
			if len(errs) != test.expectedErrs {
				t.Errorf("expected %d errors, got %d", test.expectedErrs, len(errs))
			}
		})
	}
}

func TestValidatorValidate(t *testing.T) {
	tests := []struct {
		name           string
		obj            any
		validator      func(ctx context.Context, h shared.Helper, obj any) []shared.Error
		transformError func(errs []shared.Error) []error
		expectedErrs   int
	}{
		{
			name: "single error",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return []shared.Error{{Message: "test error"}}
			},
			expectedErrs: 1,
		},
		{
			name: "multiple errors",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return []shared.Error{{Message: "test error 1"}, {Message: "test error 2"}}
			},
			expectedErrs: 2,
		},
		{
			name: "no errors",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return nil
			},
			expectedErrs: 0,
		},
		{
			name: "nil validator",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return nil
			},
			expectedErrs: 0,
		},
		{
			name: "validator returns empty slice",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return []shared.Error{}
			},
			expectedErrs: 0,
		},
		{
			name: "transform error",
			obj:  struct{}{},
			validator: func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
				return []shared.Error{{Message: "test error"}}
			},
			transformError: func(errs []shared.Error) []error {
				return []error{fmt.Errorf("transformed error")}
			},
			expectedErrs: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := New()
			v.storage.validators[reflect.TypeOf(test.obj)] = []structValidationFn{test.validator}
			v.helper.transformError = test.transformError
			errs := v.Validate(context.Background(), test.obj)
			if len(errs) != test.expectedErrs {
				t.Errorf("expected %d errors, got %d", test.expectedErrs, len(errs))
			}
		})
	}
}

func TestValidatorGetHelper(t *testing.T) {
	v := New()
	h := v.GetHelper()
	if h == nil {
		t.Errorf("expected helper to not be nil")
	}
}

func TestNewValidator(t *testing.T) {
	v := New()
	if v.storage == nil {
		t.Errorf("expected storage to not be nil")
	}
	if v.helper == nil {
		t.Errorf("expected helper to not be nil")
	}
}

func TestConfigure(t *testing.T) {
	str := "ptr-string"
	ptrString := &str
	re := regexp.MustCompile(`SMTP`)
	tests := []struct {
		name          string
		sender        *Sender
		expectedError bool
	}{
		{
			name: " sender invalid",
			sender: &Sender{
				Type:          "SMTP",
				SMTPHost:      "example.com",
				SMTPPort:      "25",
				PtrString:     ptrString,
				HTTPDestParam: "test",
			},
			expectedError: true,
		},
		{
			name: "sender valid",
			sender: &Sender{
				Type:          "SMTP",
				SMTPHost:      "example.com",
				SMTPPort:      "25",
				PtrString:     ptrString,
				HTTPDestParam: uuid.New().String(),
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Configure[Sender](validator, func(builder Configurator[Sender], temp *Sender) {
				builder.String(&temp.HTTPDestParam).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error {
					v := *value.(*string)
					if len(v) < 20 {
						return []shared.Error{h.ErrorT(ctx, v, "Should be longer than 20 characters")}
					}
					return nil
				})

				smtpValidator := builder.When(func(_ context.Context, obj *Sender) bool {
					return obj.Type == "SMTP"
				})
				smtpValidator.String(&temp.SMTPHost).Trim().Required()
				smtpValidator.String(&temp.SMTPPort).Trim().Required()

				builder.String(&temp.PtrString).Trim().Required()

				httpValidator := builder.When(func(_ context.Context, temp *Sender) bool {
					return temp.Type == "HTTP"
				})
				httpValidator.String(&temp.HTTPAddress).Trim().Required()
				httpValidator.String(&temp.HTTPDestParam).Trim().Required()
				httpValidator.String(&temp.Type).Trim().Regexp(re).MaxLen(5).MinLen(3)
			})
			errs := validator.Validate(context.Background(), tt.sender)
			if tt.expectedError && !(errs != nil && len(errs) > 0) {
				t.Errorf("expected error, got nil")
			} else if !tt.expectedError && (errs != nil && len(errs) > 0) {
				t.Errorf("expected no error, got %v", errs)
			}
		})
	}
}
