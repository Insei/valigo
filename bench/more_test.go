package bench

import (
	"errors"
	"regexp"
	"testing"
	"unicode/utf8"

	z "github.com/Oudwins/zog"
	"github.com/gookit/validate"
)

// GookitUser has the rules of User in the tags of gookit/validate, its tag key is validate too.
type GookitUser struct {
	Name    string   `validate:"required|maxLen:64"`
	Email   string   `validate:"required|email"`
	Age     int      `validate:"min:18|max:120"`
	Tags    []string `validate:"maxLen:5"`
	Address struct {
		City string `validate:"required"`
	}
}

func gookitValidate(u *GookitUser) validate.Errors {
	v := validate.Struct(u)
	// all errors, like the other libraries
	v.StopOnError = false
	v.StringRule("Tags.*", "required|maxLen:16")
	v.Validate()
	return v.Errors
}

func toGookit(u *User) *GookitUser {
	g := &GookitUser{Name: u.Name, Email: u.Email, Age: u.Age, Tags: u.Tags}
	g.Address.City = u.Address.City
	return g
}

var zogSchema = z.Struct(z.Shape{
	"Name":  z.String().Required().Max(64),
	"Email": z.String().Required().Email(),
	"Age":   z.Int().GTE(18).LTE(120),
	"Tags":  z.Slice(z.String().Required().Max(16)).Max(5),
	"Address": z.Struct(z.Shape{
		"City": z.String().Required(),
	}),
})

// manualEmail is the email expression of valigo, so the baseline does the same work.
var manualEmail = regexp.MustCompile(`^[^\s,@#$%^&*!()]+@([a-zA-Z0-9]+[.])+[a-zA-Z]{2,8}$`)

// validateManually is the hand-written baseline: the same rules as plain Go code.
func validateManually(u *User) []error {
	var errs []error
	if u.Name == "" || utf8.RuneCountInString(u.Name) > 64 {
		errs = append(errs, errors.New("name"))
	}
	if u.Email == "" || !manualEmail.MatchString(u.Email) {
		errs = append(errs, errors.New("email"))
	}
	if u.Age < 18 || u.Age > 120 {
		errs = append(errs, errors.New("age"))
	}
	if len(u.Tags) > 5 {
		errs = append(errs, errors.New("tags"))
	}
	for _, tag := range u.Tags {
		if tag == "" || utf8.RuneCountInString(tag) > 16 {
			errs = append(errs, errors.New("tag"))
		}
	}
	if u.Address.City == "" {
		errs = append(errs, errors.New("city"))
	}
	return errs
}

func TestTheSameRulesOfMoreLibraries(t *testing.T) {
	if errs := gookitValidate(toGookit(valid)); len(errs) != 0 {
		t.Errorf("gookit/validate: %v", errs)
	}
	if issues := zogSchema.Validate(valid); len(issues) != 0 {
		t.Errorf("zog: %v", issues)
	}
	if errs := validateManually(valid); errs != nil {
		t.Errorf("manual: %v", errs)
	}
	// Name, Email, Age, Tags.1, Tags.2, Address.City; gookit/validate reports the elements of Tags under one key Tags.*
	if errs := gookitValidate(toGookit(invalid)); len(errs) != 5 {
		t.Errorf("gookit/validate: %d fields: %v", len(errs), errs)
	}
	if issues := zogSchema.Validate(invalid); len(issues) != 6 {
		t.Errorf("zog: %d issues: %v", len(issues), issues)
	}
	if errs := validateManually(invalid); len(errs) != 6 {
		t.Errorf("manual: %d errors", len(errs))
	}
}

func BenchmarkValidMore(b *testing.B) {
	g := toGookit(valid)
	b.Run("gookit", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = gookitValidate(g)
		}
	})
	b.Run("zog", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = zogSchema.Validate(valid)
		}
	})
	b.Run("manual", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = validateManually(valid)
		}
	})
}

func BenchmarkInvalidMore(b *testing.B) {
	g := toGookit(invalid)
	b.Run("gookit", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = gookitValidate(g)
		}
	})
	b.Run("zog", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = zogSchema.Validate(invalid)
		}
	})
	b.Run("manual", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = validateManually(invalid)
		}
	})
}
