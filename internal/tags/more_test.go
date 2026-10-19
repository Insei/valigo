package tags_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insei/valigo"
	"github.com/insei/valigo/internal/testenv"
)

// codes returns the locations and the codes of the errors of the tags of T.
func codes[T any](t *testing.T, obj *T) []string {
	t.Helper()
	v := valigo.New()
	valigo.Configure[T](v, func(b valigo.Configurator[T], _ *T) { b.Tags("validate") })
	var got []string
	for _, err := range v.ValidateTyped(context.Background(), obj) {
		got = append(got, err.Location+" "+err.Code)
	}
	return got
}

type moreStrings struct {
	Gt      string `validate:"gt=2"`
	Lt      string `validate:"lt=2"`
	Gte     string `validate:"gte=2,lte=3"`
	Ne      string `validate:"ne=x"`
	Old     string
	New     string `validate:"nefield=Old"`
	Alpha   string `validate:"alpha"`
	Number  string `validate:"number"`
	Numeric string `validate:"numeric"`
	Hex     string `validate:"hexadecimal"`
	Excl    string `validate:"excludes=x,containsany=0x2C;,excludesall=0x7C"`
	Starts  string `validate:"startsnotwith=a,endsnotwith=z"`
	Fold    string `validate:"eq_ignore_case=Yes"`
	Contact string `validate:"email|url"`
	Code    string `validate:"len=2|oneof=none"`
}

func TestMoreStrings(t *testing.T) {
	valid := &moreStrings{Gt: "abc", Gte: "ab", Old: "a", New: "b", Alpha: "ab", Number: "01", Numeric: "-1.5", Hex: "0xff",
		Excl: "a,b", Fold: "YES", Contact: "https://example.com", Code: "none"}
	assert.Empty(t, codes(t, valid))
	valid.Contact, valid.Code = "a@b.cc", "ab"
	assert.Empty(t, codes(t, valid), "the other alternatives")

	assert.Equal(t, []string{
		"Gt string.min_runes", "Lt string.max_runes", "Gte string.min_runes", "Ne value.not_equal", "New string.not_equal_field",
		"Alpha string.alpha", "Number string.digits", "Numeric string.numeric", "Hex string.hexadecimal",
		"Excl string.excludes", "Excl string.contains_any", "Excl string.excludes_any",
		"Starts string.not_prefix", "Starts string.not_suffix", "Fold string.equal_fold",
		"Contact string.email", "Code string.max_runes",
	}, codes(t, &moreStrings{Gt: "ab", Lt: "ab", Gte: "a", Ne: "x", Old: "a", New: "a", Alpha: "a1", Number: "-1", Numeric: "1.",
		Hex: "x", Excl: "x|", Starts: "abz", Fold: "no", Contact: "nope", Code: "abc"}))
}

type moreNumbers struct {
	Gt    int  `validate:"gt=1,lt=5"`
	Ne    uint `validate:"ne=3"`
	Min   float64
	Max   float64    `validate:"gtfield=Min"`
	Low   float64    `validate:"ltfield=Min"`
	Same  float64    `validate:"eqfield=Min"`
	Other float64    `validate:"nefield=Min"`
	On    bool       `validate:"eq=true"`
	Off   bool       `validate:"ne=true"`
	Later time.Time  `validate:"gt"`
	Ago   *time.Time `validate:"omitnil,lte"`
}

func TestMoreNumbers(t *testing.T) {
	past, future := time.Now().Add(-time.Hour), time.Now().Add(time.Hour)
	assert.Empty(t, codes(t, &moreNumbers{Gt: 2, Min: 1, Max: 2, Low: 0, Same: 1, Other: 2, On: true, Later: future, Ago: &past}))
	assert.Equal(t, []string{
		"Gt number.greater_than", "Ne value.not_equal", "Max number.greater_field", "Low number.less_field",
		"Same value.equal_field", "Other value.not_equal_field", "On value.equal", "Off value.not_equal",
		"Later time.in_future", "Ago time.in_past",
	}, codes(t, &moreNumbers{Gt: 1, Ne: 3, Min: 1, Max: 1, Low: 1, Same: 2, Other: 1, Off: true, Later: past, Ago: &future}))
	assert.Equal(t, []string{"Gt number.less_than"}, codes(t, &moreNumbers{Gt: 5, Min: 1, Max: 2, Same: 1, On: true, Later: future})[:1])
}

type moreLengths struct {
	Gt  []int          `validate:"gt=1"`
	Lt  []int          `validate:"lt=2"`
	Eq  map[string]int `validate:"eq=1"`
	Gte *[]int         `validate:"omitnil,gte=1,lte=1"`
}

func TestMoreLengths(t *testing.T) {
	one := []int{1}
	assert.Empty(t, codes(t, &moreLengths{Gt: []int{1, 2}, Eq: map[string]int{"a": 1}, Gte: &one}))
	assert.Equal(t, []string{"Gt slice.min_len", "Lt slice.max_len", "Eq map.min_len", "Gte slice.min_len"},
		codes(t, &moreLengths{Gt: []int{1}, Lt: []int{1, 2}, Gte: &[]int{}}))
}

type moreConditions struct {
	Kind    string
	Count   int
	Ok      bool
	Ptr     *string
	List    []string
	Labels  map[string]string
	At      time.Time
	If      string `validate:"required_if=Kind company Count 2"`
	Unless  string `validate:"required_unless=Kind person"`
	With    string `validate:"required_with=Ptr List"`
	WithAll string `validate:"required_with_all=Labels At"`
	Without string `validate:"required_without=Ok Ptr"`
	NoneOf  string `validate:"required_without_all=Count Ok"`
	Flag    *int   `validate:"required_if=Ok true"`
	Limit   *int
	Note    string `validate:"required_if=Limit 3"`
}

func TestConditions(t *testing.T) {
	s := "x"
	assert.Equal(t, []string{"Unless string.required", "Without string.required", "NoneOf string.required"},
		codes(t, &moreConditions{}), "nothing is present")
	assert.Equal(t, []string{"If string.required", "Unless string.required", "With string.required", "WithAll string.required", "Flag number.required"},
		codes(t, &moreConditions{Kind: "company", Count: 2, Ok: true, Ptr: &s, Labels: map[string]string{"a": "b"}, At: time.Now()}))
	three := 3
	assert.Equal(t, []string{"Unless string.required", "Without string.required", "NoneOf string.required", "Note string.required"},
		codes(t, &moreConditions{Limit: &three}), "a pointer field is compared with its value")
	one := 1
	assert.Empty(t, codes(t, &moreConditions{Kind: "person", Count: 1, Ok: true, Ptr: &s, List: []string{}, Labels: map[string]string{},
		With: "w", Flag: &one}))
}

type moreMaps struct {
	Labels map[string]string `validate:"dive,keys,alpha,endkeys,required"`
	Keys   map[string]int    `validate:"dive,keys,max=2,endkeys"`
	Values map[string]int    `validate:"dive,min=1"`
}

func TestKeys(t *testing.T) {
	assert.Empty(t, codes(t, &moreMaps{Labels: map[string]string{"ab": "x"}, Keys: map[string]int{"ab": 0}, Values: map[string]int{"a": 1}}))
	assert.Equal(t, []string{"Labels[a1] string.alpha", "Labels[b] string.required", "Keys[abc] string.max_runes", "Values[a] number.min"},
		codes(t, &moreMaps{Labels: map[string]string{"a1": "x", "b": ""}, Keys: map[string]int{"abc": 0}, Values: map[string]int{"a": 0}}))
}

type moreItem struct {
	Title string `validate:"required,max=3"`
	Min   int
	Max   int `validate:"gtefield=Min"`
}

type moreNode struct {
	Name     string     `validate:"required"`
	Children []moreNode `validate:"dive"`
}

type moreElems struct {
	Items    []moreItem   `validate:"min=1,dive"`
	Pointers []*moreItem  `validate:"dive,required"`
	Optional *[]*moreItem `validate:"omitempty,dive,omitempty"`
	Tree     moreNode
}

func TestStructElements(t *testing.T) {
	valid := &moreElems{Items: []moreItem{{Title: "a"}}, Pointers: []*moreItem{{Title: "b"}}, Optional: &[]*moreItem{nil},
		Tree: moreNode{Name: "root", Children: []moreNode{{Name: "a", Children: []moreNode{{Name: "b"}}}}}}
	assert.Empty(t, codes(t, valid))
	assert.Equal(t, []string{
		"Items[1].Title string.required", "Items[2].Title string.max_runes", "Items[2].Max number.min_field",
		"Pointers[0] value.required", "Pointers[1].Title string.required",
		"Optional[0].Title string.required",
		"Tree.Children[0].Children[0].Name string.required",
	}, codes(t, &moreElems{
		Items:    []moreItem{{Title: "a"}, {}, {Title: "long", Min: 2, Max: 1}},
		Pointers: []*moreItem{nil, {}},
		Optional: &[]*moreItem{{}},
		Tree:     moreNode{Name: "root", Children: []moreNode{{Name: "a", Children: []moreNode{{}}}}},
	}))

	if !testenv.RaceEnabled {
		v := valigo.New()
		valigo.Configure[moreElems](v, func(b valigo.Configurator[moreElems], _ *moreElems) { b.Tags("validate") })
		ctx := context.Background()
		assert.Zero(t, testing.AllocsPerRun(10, func() { v.ValidateTyped(ctx, valid) }), "the elements are validated without allocations")
		alts := valigo.New()
		valigo.Configure[moreStrings](alts, func(b valigo.Configurator[moreStrings], _ *moreStrings) { b.Tags("validate") })
		first := &moreStrings{Gt: "abc", Gte: "ab", Old: "a", New: "b", Alpha: "ab", Number: "01", Numeric: "1", Hex: "f",
			Excl: ",", Fold: "yes", Contact: "a@b.cc", Code: "ab"}
		assert.Zero(t, testing.AllocsPerRun(10, func() { alts.ValidateTyped(ctx, first) }), "the first alternatives pass")
	}
}

func TestMoreTagsErrors(t *testing.T) {
	tests := []struct {
		name      string
		configure func()
		err       string
	}{
		{"condition without fields", tagsOf[struct {
			A string `validate:"required_with"`
		}], `valigo: tag validate:"required_with" of field A: required_with needs the fields`},
		{"required struct", tagsOf[struct {
			A tagsAddress `validate:"required_with=B"`
			B string
		}], `valigo: tag validate:"required_with=B" of field A: the rule required isn't supported for the field, use the rules in code`},
		{"odd pairs", tagsOf[struct {
			A string `validate:"required_if=B"`
			B string
		}], `valigo: tag validate:"required_if=B" of field A: required_if=B: the pairs of a field and a value expected`},
		{"no condition field", tagsOf[struct {
			A string `validate:"required_with=B"`
		}], `valigo: tag validate:"required_with=B" of field A: required_with=B: no field B`},
		{"no value field", tagsOf[struct {
			A string `validate:"required_if=B 1"`
		}], `valigo: tag validate:"required_if=B 1" of field A: required_if=B 1: no field B`},
		{"struct condition", tagsOf[struct {
			A string `validate:"required_with=B"`
			B tagsAddress
		}], `valigo: tag validate:"required_with=B" of field A: required_with=B: the field B of type tags_test.tagsAddress can't be a condition`},
		{"slice value", tagsOf[struct {
			A string `validate:"required_if=B x"`
			B []string
		}], `valigo: tag validate:"required_if=B x" of field A: required_if=B x: the field B of type []string can't be compared with a value`},
		{"time value", tagsOf[struct {
			A string `validate:"required_if=B x"`
			B time.Time
		}], `valigo: tag validate:"required_if=B x" of field A: required_if=B x: the field B of type time.Time can't be compared with a value`},
		{"bad value", tagsOf[struct {
			A string `validate:"required_if=B x"`
			B int
		}], `valigo: tag validate:"required_if=B x" of field A: required_if=B x: strconv.ParseInt: parsing "x": invalid syntax`},
		{"unknown condition", tagsOf[struct {
			A string `validate:"required_iff=B"`
			B string
		}], `valigo: tag validate:"required_iff=B" of field A: the rule required_iff isn't supported for the field, use the rules in code`},
		{"condition of elements", tagsOf[struct {
			A []string `validate:"dive,required_with=B"`
			B string
		}], `valigo: tag validate:"dive,required_with=B" of field A: required_with isn't supported for the elements`},
		{"keys without endkeys", tagsOf[struct {
			A map[string]string `validate:"dive,keys,alpha"`
		}], `valigo: tag validate:"dive,keys,alpha" of field A: keys without endkeys`},
		{"key rule", tagsOf[struct {
			A map[string]string `validate:"dive,keys,isbn,endkeys"`
		}], `valigo: tag validate:"dive,keys,isbn,endkeys" of field A: the rule isbn isn't supported for the field, use the rules in code`},
		{"unique structs", tagsOf[struct {
			A []tagsAddress `validate:"unique"`
		}], `valigo: tag validate:"unique" of field A: unique of the slices of structs isn't supported`},
		{"rule of struct elements", tagsOf[struct {
			A []tagsAddress `validate:"dive,required"`
		}], `valigo: tag validate:"dive,required" of field A: the rule required isn't supported for the field, use the rules in code`},
		{"alternative", tagsOf[struct {
			A string `validate:"email|isbn"`
		}], `valigo: tag validate:"email|isbn" of field A: the rule isbn isn't supported for the field, use the rules in code`},
		{"time with a parameter", tagsOf[struct {
			A time.Time `validate:"gt=1"`
		}], `valigo: tag validate:"gt=1" of field A: gt=1: the times are compared with the current time, without a parameter`},
		{"bad bool", tagsOf[struct {
			A bool `validate:"eq=yes"`
		}], `valigo: tag validate:"eq=yes" of field A: eq=yes: strconv.ParseBool: parsing "yes": invalid syntax`},
		{"bad length", tagsOf[struct {
			A []string `validate:"gt=x"`
		}], `valigo: tag validate:"gt=x" of field A: gt=x: strconv.Atoi: parsing "x": invalid syntax`},
		{"bad string length", tagsOf[struct {
			A string `validate:"lt=x"`
		}], `valigo: tag validate:"lt=x" of field A: lt=x: strconv.Atoi: parsing "x": invalid syntax`},
		{"element of a struct", tagsOf[struct {
			A []struct {
				B string `validate:"isbn"`
			} `validate:"dive"`
		}], `valigo: tag validate:"isbn" of field B: the rule isbn isn't supported for the field, use the rules in code`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.PanicsWithValue(t, tt.err, tt.configure)
		})
	}
}
