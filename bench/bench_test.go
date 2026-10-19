// Package bench compares valigo with go-playground/validator and ozzo-validation on the same rules.
// It is a separate module, so the compared libraries aren't dependencies of valigo.
package bench

import (
	"context"
	"testing"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	playground "github.com/go-playground/validator/v10"

	"github.com/insei/valigo"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/str"
)

type Address struct {
	City string `validate:"required"`
}

// User has the same rules in all libraries: the tags for go-playground/validator,
// Validate for ozzo-validation and the configuration for valigo.
type User struct {
	Name    string   `validate:"required,max=64"`
	Email   string   `validate:"required,email"`
	Age     int      `validate:"gte=18,lte=120"`
	Tags    []string `validate:"max=5,dive,required,max=16"`
	Address Address
}

func (u *User) Validate() error {
	return ozzo.ValidateStruct(u,
		ozzo.Field(&u.Name, ozzo.Required, ozzo.RuneLength(0, 64)),
		ozzo.Field(&u.Email, ozzo.Required, is.EmailFormat),
		ozzo.Field(&u.Age, ozzo.Min(18), ozzo.Max(120)),
		ozzo.Field(&u.Tags, ozzo.Length(0, 5), ozzo.Each(ozzo.Required, ozzo.RuneLength(0, 16))),
		ozzo.Field(&u.Address),
	)
}

func (a Address) Validate() error {
	return ozzo.ValidateStruct(&a, ozzo.Field(&a.City, ozzo.Required))
}

var (
	valid   = &User{Name: "Ann", Email: "ann@example.com", Age: 30, Tags: []string{"go", "api"}, Address: Address{City: "Berlin"}}
	invalid = &User{Name: "", Email: "ann.example.com", Age: 16, Tags: []string{"go", "", "a-very-long-tag-name"}}
)

func newValigo() *valigo.Validator {
	v := valigo.New()
	valigo.Configure[User](v, func(b valigo.Configurator[User], u *User) {
		b.Field(&u.Name).As(str.Of).Required().MaxRunes(64)
		b.Field(&u.Email).As(str.Of).Required().Email()
		b.Field(&u.Age).As(num.Of).AnyOfInterval(18, 120)
		b.Slice(&u.Tags).As(str.Slice).MaxLen(5).Each(func(tag *str.BaseConfigurator) { tag.Required().MaxRunes(16) })
		b.Field(&u.Address.City).As(str.Of).Required()
	})
	return v
}

// newValigoTags configures valigo by the tags of go-playground/validator.
func newValigoTags() *valigo.Validator {
	v := valigo.New()
	valigo.Configure[User](v, func(b valigo.Configurator[User], _ *User) { b.Tags("validate") })
	return v
}

func TestTheSameRules(t *testing.T) {
	ctx := context.Background()
	v, pg := newValigo(), playground.New()
	if errs := v.Validate(ctx, valid); errs != nil {
		t.Errorf("valigo: %v", errs)
	}
	if err := pg.Struct(valid); err != nil {
		t.Errorf("go-playground/validator: %v", err)
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("ozzo-validation: %v", err)
	}
	// Name, Email, Age, Tags[1], Tags[2], Address.City
	if errs := v.Validate(ctx, invalid); len(errs) != 6 {
		t.Errorf("valigo: %d errors: %v", len(errs), errs)
	}
	tags := newValigoTags()
	if errs := tags.Validate(ctx, valid); errs != nil {
		t.Errorf("valigo tags: %v", errs)
	}
	if errs := tags.Validate(ctx, invalid); len(errs) != 6 {
		t.Errorf("valigo tags: %d errors: %v", len(errs), errs)
	}
	if errs, _ := pg.Struct(invalid).(playground.ValidationErrors); len(errs) != 6 {
		t.Errorf("go-playground/validator: %d errors: %v", len(errs), errs)
	}
	if err := invalid.Validate(); err == nil {
		t.Errorf("ozzo-validation: no errors")
	}
}

func BenchmarkValid(b *testing.B) {
	ctx := context.Background()
	b.Run("valigo", func(b *testing.B) {
		v := newValigo()
		b.ReportAllocs()
		for b.Loop() {
			_ = v.Validate(ctx, valid)
		}
	})
	b.Run("valigo-tags", func(b *testing.B) {
		v := newValigoTags()
		b.ReportAllocs()
		for b.Loop() {
			_ = v.Validate(ctx, valid)
		}
	})
	b.Run("go-playground", func(b *testing.B) {
		v := playground.New()
		b.ReportAllocs()
		for b.Loop() {
			_ = v.Struct(valid)
		}
	})
	b.Run("ozzo", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = valid.Validate()
		}
	})
}

func BenchmarkInvalid(b *testing.B) {
	ctx := context.Background()
	b.Run("valigo", func(b *testing.B) {
		v := newValigo()
		b.ReportAllocs()
		for b.Loop() {
			_ = v.Validate(ctx, invalid)
		}
	})
	b.Run("valigo-tags", func(b *testing.B) {
		v := newValigoTags()
		b.ReportAllocs()
		for b.Loop() {
			_ = v.Validate(ctx, invalid)
		}
	})
	b.Run("go-playground", func(b *testing.B) {
		v := playground.New()
		b.ReportAllocs()
		for b.Loop() {
			_ = v.Struct(invalid)
		}
	})
	// valigo and ozzo build the messages of the errors, go-playground/validator builds them by Error
	b.Run("go-playground-messages", func(b *testing.B) {
		v := playground.New()
		b.ReportAllocs()
		for b.Loop() {
			for _, err := range v.Struct(invalid).(playground.ValidationErrors) {
				_ = err.Error()
			}
		}
	})
	b.Run("ozzo", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = invalid.Validate().Error()
		}
	})
}
