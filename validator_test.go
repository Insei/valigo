package valigo

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
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
	ptr := "c6b7f3a2-0d4e-4f6a-9b8c-1d2e3f4a5b6c"
	return &Sender{
		Type:          "SMTP",
		SMTPHost:      "c6b7f3a2-0d4e-4f6a-9b8c-1d2e3f4a5b6c" + "   ",
		SMTPPort:      "c6b7f3a2-0d4e-4f6a-9b8c-1d2e3f4a5b6c" + " ",
		SMTPInt:       10,
		HTTPAddress:   "c6b7f3a2-0d4e-4f6a-9b8c-1d2e3f4a5b6c" + " ",
		HTTPDestParam: "c6b7f3a2-0d4e-4f6a-9b8c-1d2e3f4a5b6c" + "  ",
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
		builder.Field(&temp.HTTPDestParam).As(str.Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value string) []shared.Error {
			v := value
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
		smtpValidator.Field(&temp.SMTPHost).As(str.Of).Trim().Required()
		smtpValidator.Field(&temp.SMTPPort).As(str.Of).Trim().Required()
		smtpValidator.Field(&temp.SMTPInt).As(num.Of).Max(20)
		builder.Field(&temp.PtrString).As(str.Of).Trim().Required()
		//builder.Custom(func(ctx context.Context, h shared.StructCustomHelper, obj *Sender) []shared.Error {
		//	h.ErrorT(ctx)
		//})
		httpValidator := builder.When(func(_ context.Context, temp *Sender) bool {
			return temp.Type == "HTTP"
		})
		httpValidator.Field(&temp.HTTPAddress).As(str.Of).Trim().Required()
		httpValidator.Field(&temp.HTTPDestParam).As(str.Of).Trim().Required()
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
			v.validators[reflect.TypeOf(test.obj)] = []shared.ValidationFn{test.validator}
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
			v.validators[reflect.TypeOf(test.obj)] = []shared.ValidationFn{test.validator}
			if test.transformError != nil {
				v.transformError = test.transformError
			}
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
	if v.validators == nil {
		t.Errorf("expected validators to not be nil")
	}
	if v.helper == nil {
		t.Errorf("expected helper to not be nil")
	}
}

func TestConfigure(t *testing.T) {
	s := "ptr-string"
	ptrString := &s
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
				HTTPDestParam: "c6b7f3a2-0d4e-4f6a-9b8c-1d2e3f4a5b6c",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Configure[Sender](validator, func(builder Configurator[Sender], temp *Sender) {
				builder.Field(&temp.HTTPDestParam).As(str.Of).Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value string) []shared.Error {
					v := value
					if len(v) < 20 {
						return []shared.Error{h.ErrorT(ctx, v, "Should be longer than 20 characters")}
					}
					return nil
				})

				smtpValidator := builder.When(func(_ context.Context, obj *Sender) bool {
					return obj.Type == "SMTP"
				})
				smtpValidator.Field(&temp.SMTPHost).As(str.Of).Trim().Required()
				smtpValidator.Field(&temp.SMTPPort).As(str.Of).Trim().Required()

				builder.Field(&temp.PtrString).As(str.Of).Trim().Required()

				httpValidator := builder.When(func(_ context.Context, temp *Sender) bool {
					return temp.Type == "HTTP"
				})
				httpValidator.Field(&temp.HTTPAddress).As(str.Of).Trim().Required()
				httpValidator.Field(&temp.HTTPDestParam).As(str.Of).Trim().Required()
				httpValidator.Field(&temp.Type).As(str.Of).Trim().Regexp(re).MaxLen(5).MinLen(3)
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
