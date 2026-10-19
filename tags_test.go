package valigo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo/internal/testenv"
	"github.com/insei/valigo/str"
)

type tagsAddress struct {
	City string `validate:"required,max=8"`
}

// tagsPair compares the fields of a nested struct.
type tagsPair struct {
	A string
	B string `validate:"eqfield=A"`
}

type tagsStatus string

type tagsModel struct {
	Name     string            `validate:"required,min=2,max=4"`
	Code     tagsStatus        `validate:"oneof=new done"`
	Email    *string           `validate:"omitempty,email"`
	Nick     string            `validate:"omitempty,alphanum,len=3"`
	Password string            `validate:"required"`
	Confirm  string            `validate:"eqfield=Password"`
	Site     string            `validate:"startswith=https://,endswith=.com,contains=ex,url"`
	Day      string            `validate:"datetime=2006-01-02"`
	Kind     string            `validate:"eq=person"`
	Age      int               `validate:"gte=18,lte=120"`
	Min      uint8             `validate:"min=1,max=9"`
	Score    *float64          `validate:"required,oneof=1.5 2.5"`
	Level    int64             `validate:"eq=3"`
	Lat      float32           `validate:"latitude"`
	Lng      float64           `validate:"longitude"`
	Low      int               `validate:"ltefield=Age"`
	High     int               `validate:"gtefield=Age"`
	Agree    bool              `validate:"required"`
	StartsAt time.Time         `validate:"required"`
	EndsAt   time.Time         `validate:"gtfield=StartsAt"`
	BeforeAt *time.Time        `validate:"ltfield=StartsAt"`
	Tags     []string          `validate:"required,max=2,dive,required,max=3"`
	Roles    []string          `validate:"unique"`
	Nums     []*int            `validate:"min=1,dive,min=1"`
	Times    []time.Time       `validate:"dive,required"`
	Labels   map[string]string `validate:"omitempty,max=1,dive,lowercase"`
	Counts   *map[string]int   `validate:"required"`
	Items    []tagsAddress     `validate:"len=1"`
	Address  tagsAddress
	Home     *tagsAddress `validate:"required"`
	Ignored  string       `validate:"-"`
	Untagged string
	Skipped  *[]string    `validate:"omitempty,min=2"`
	Optional *tagsAddress `validate:"omitempty"`
	Pair     tagsPair
}

func validTagsModel() *tagsModel {
	email, score, n := "a@b.cc", 1.5, 1
	starts := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	before := starts.Add(-time.Hour)
	counts := map[string]int{"a": 1}
	return &tagsModel{
		Name: "Ann", Code: "new", Email: &email, Password: "p", Confirm: "p", Site: "https://ex.com", Day: "2026-01-02",
		Kind: "person", Age: 30, Min: 1, Score: &score, Level: 3, Low: 30, High: 30, Agree: true,
		StartsAt: starts, EndsAt: starts.Add(time.Hour), BeforeAt: &before,
		Tags: []string{"go"}, Roles: []string{"a", "b"}, Nums: []*int{&n}, Times: []time.Time{starts},
		Labels: map[string]string{"a": "x"}, Counts: &counts, Items: []tagsAddress{{City: "x"}},
		Address: tagsAddress{City: "Paris"}, Home: &tagsAddress{City: "Rome"},
	}
}

func newTagsValidator() *Validator {
	v := New()
	Configure[tagsModel](v, func(b Configurator[tagsModel], m *tagsModel) {
		b.Tags("validate")
		b.Field(&m.Untagged).As(str.Of).MaxLen(1) // the rules in code with the tags
	})
	return v
}

func TestTags(t *testing.T) {
	v := newTagsValidator()
	ctx := context.Background()
	assert.Empty(t, v.ValidateTyped(ctx, validTagsModel()))

	zero := 0
	bad := &tagsModel{
		Name: "A", Code: "old", Nick: "a b", Confirm: "x", Site: "ftp://a.org", Day: "02.01.2026", Kind: "robot",
		Age: 17, Min: 10, Level: 4, Lat: 91, Lng: 181, Low: 18, High: 16,
		EndsAt: time.Time{}.Add(-time.Hour), Tags: []string{"", "long", "x"}, Roles: []string{"a", "a"},
		Nums: []*int{&zero, nil}, Times: []time.Time{{}}, Labels: map[string]string{"A": "B", "c": "d"},
		Items: []tagsAddress{}, Address: tagsAddress{City: "Amsterdam"}, Untagged: "ab", Skipped: &[]string{"a"},
		Pair: tagsPair{A: "a", B: "b"},
	}
	var got []string
	for _, err := range v.ValidateTyped(ctx, bad) {
		got = append(got, err.Location+" "+err.Code)
	}
	assert.Equal(t, []string{
		"Name string.min_runes",
		"Code string.any_of",
		"Nick string.alphanumeric",
		"Password string.required",
		"Confirm string.equal_field",
		"Site string.prefix", "Site string.suffix", "Site string.contains",
		"Day string.datetime",
		"Kind string.any_of",
		"Age number.min",
		"Min number.max",
		"Score number.required",
		"Level number.any_of",
		"Lat number.latitude",
		"Lng number.longitude",
		"Low number.max_field",
		"High number.min_field",
		"Agree value.required",
		"StartsAt time.required",
		"EndsAt time.after_field",
		"Tags slice.max_len", "Tags[0] string.required", "Tags[1] string.max_runes",
		"Roles value.unique",
		"Nums[0] number.min",
		"Times[0] time.required",
		"Labels map.max_len", "Labels[A] string.lowercase",
		"Counts map.required",
		"Items slice.min_len",
		"Address.City string.max_runes",
		"Home value.required",
		"Skipped slice.min_len",
		"Pair.B string.equal_field",
		"Untagged string.max_len",
	}, got)

	empty := map[string]int{}
	model := validTagsModel()
	model.Counts, model.Tags = &empty, nil
	var errs []string
	for _, err := range v.ValidateTyped(ctx, model) {
		errs = append(errs, err.Location+" "+err.Code)
		assert.Equal(t, 0, err.Value, "the errors of the lengths have the lengths")
	}
	assert.Equal(t, []string{"Tags slice.required", "Counts map.required"}, errs)

	model = validTagsModel()
	model.Home.City = ""
	errsOfHome := v.ValidateTyped(ctx, model)
	if assert.Len(t, errsOfHome, 1) {
		assert.Equal(t, "Home.City", errsOfHome[0].Location, "the fields of a present nested struct are validated")
	}

	if !testenv.RaceEnabled {
		valid := validTagsModel()
		assert.Zero(t, testing.AllocsPerRun(10, func() { v.ValidateTyped(ctx, valid) }))
	}
}

func TestTagsConstraints(t *testing.T) {
	var got []string
	for _, c := range newTagsValidator().Constraints() {
		if c.Field.GetStructPath() == "Tags" || c.Field.GetStructPath() == "Labels" || c.Field.GetStructPath() == "Home" {
			got = append(got, c.Field.GetStructPath()+" "+c.Keyword)
		}
	}
	assert.Equal(t, []string{
		"Tags required", "Tags minItems", "Tags maxItems", "Tags required", "Tags minLength", "Tags maxLength",
		"Home required",
	}, got, "omitempty rules aren't described")
}

func TestTagsErrors(t *testing.T) {
	tests := []struct {
		name      string
		configure func()
		err       string
	}{
		{"unknown rule", tagsOf[struct {
			A string `validate:"required,isbn"`
		}], `valigo: tag validate:"required,isbn" of field A: the rule isbn isn't supported for the field, use the rules in code`},
		{"alternatives", tagsOf[struct {
			A string `validate:"email|url"`
		}], `valigo: tag validate:"email|url" of field A: the alternatives of | aren't supported, use When or Custom`},
		{"bad number", tagsOf[struct {
			A int `validate:"min=x"`
		}], `valigo: tag validate:"min=x" of field A: min=x: strconv.ParseInt: parsing "x": invalid syntax`},
		{"bad oneof", tagsOf[struct {
			A uint `validate:"oneof=1 -2"`
		}], `valigo: tag validate:"oneof=1 -2" of field A: oneof=1 -2: strconv.ParseUint: parsing "-2": invalid syntax`},
		{"bad float", tagsOf[struct {
			A float64 `validate:"max=y"`
		}], `valigo: tag validate:"max=y" of field A: max=y: strconv.ParseFloat: parsing "y": invalid syntax`},
		{"bad length", tagsOf[struct {
			A string `validate:"len=z"`
		}], `valigo: tag validate:"len=z" of field A: len=z: strconv.Atoi: parsing "z": invalid syntax`},
		{"bad slice length", tagsOf[struct {
			A []string `validate:"max=z"`
		}], `valigo: tag validate:"max=z" of field A: max=z: strconv.Atoi: parsing "z": invalid syntax`},
		{"slice rule", tagsOf[struct {
			A []string `validate:"email"`
		}], `valigo: tag validate:"email" of field A: the rule email isn't supported for the field, use the rules in code`},
		{"unique of a map", tagsOf[struct {
			A map[string]int `validate:"unique"`
		}], `valigo: tag validate:"unique" of field A: the rule unique isn't supported for the field, use the rules in code`},
		{"unique of pointers", tagsOf[struct {
			A []*string `validate:"unique"`
		}], `valigo: tag validate:"unique" of field A: unique of the slices of pointers isn't supported`},
		{"dive into structs", tagsOf[struct {
			A []tagsAddress `validate:"dive"`
		}], `valigo: tag validate:"dive" of field A: the elements of type valigo.tagsAddress aren't supported, configure the rules of their type`},
		{"dive into a map of ints", tagsOf[struct {
			A map[int]string `validate:"dive,required"`
		}], `valigo: tag validate:"dive,required" of field A: dive is supported for the maps of strings to values, map[string]V`},
		{"element rule", tagsOf[struct {
			A []string `validate:"dive,isbn"`
		}], `valigo: tag validate:"dive,isbn" of field A: the rule isbn isn't supported for the field, use the rules in code`},
		{"map value rule", tagsOf[struct {
			A map[string]string `validate:"dive,isbn"`
		}], `valigo: tag validate:"dive,isbn" of field A: the rule isbn isn't supported for the field, use the rules in code`},
		{"field of an element", tagsOf[struct {
			A []string `validate:"dive,eqfield=B"`
			B string
		}], `valigo: tag validate:"dive,eqfield=B" of field A: eqfield isn't supported for the elements`},
		{"no field", tagsOf[struct {
			A string `validate:"eqfield=B"`
		}], `valigo: tag validate:"eqfield=B" of field A: eqfield=B: no field B`},
		{"other type", tagsOf[struct {
			A int `validate:"gtefield=B"`
			B string
		}], `valigo: tag validate:"gtefield=B" of field A: gtefield=B: the field B has type string, int expected`},
		{"time field", tagsOf[struct {
			A time.Time `validate:"gtfield=B"`
		}], `valigo: tag validate:"gtfield=B" of field A: gtfield=B: no field B`},
		{"time rule", tagsOf[struct {
			A time.Time `validate:"min=1"`
		}], `valigo: tag validate:"min=1" of field A: the rule min isn't supported for the field, use the rules in code`},
		{"string field of a number", tagsOf[struct {
			A string `validate:"eqfield=B"`
			B int
		}], `valigo: tag validate:"eqfield=B" of field A: eqfield=B: the field B has type int, string expected`},
		{"number rule", tagsOf[struct {
			A int `validate:"email"`
		}], `valigo: tag validate:"email" of field A: the rule email isn't supported for the field, use the rules in code`},
		{"bool rule", tagsOf[struct {
			A bool `validate:"min=1"`
		}], `valigo: tag validate:"min=1" of field A: the rule min isn't supported for the field, use the rules in code`},
		{"struct rule", tagsOf[struct {
			A tagsAddress `validate:"required"`
		}], `valigo: tag validate:"required" of field A: the rule required isn't supported for the field, use the rules in code`},
		{"type", tagsOf[struct {
			A complex64 `validate:"required"`
		}], `valigo: tag validate:"required" of field A: the fields of type complex64 aren't supported`},
		{"number field", tagsOf[struct {
			A int `validate:"ltefield=B"`
			B string
		}], `valigo: tag validate:"ltefield=B" of field A: ltefield=B: the field B has type string, int expected`},
		{"elements of a type", tagsOf[struct {
			A []complex64 `validate:"dive,required"`
		}], `valigo: tag validate:"dive,required" of field A: the elements of type complex64 aren't supported, configure the rules of their type`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.PanicsWithValue(t, tt.err, tt.configure)
		})
	}
}

// tagsOf configures the tags of T.
func tagsOf[T any]() {
	Configure[T](New(), func(b Configurator[T], _ *T) { b.Tags("validate") })
}
