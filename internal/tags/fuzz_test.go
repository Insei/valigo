package tags_test

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/internal/helper"
	"github.com/insei/valigo/internal/tags"
	"github.com/insei/valigo/internal/zero"
	"github.com/insei/valigo/shared"
)

type fuzzItem struct {
	Name string `validate:"required"`
}

// fuzzFields are the fields of the struct of the fuzzed tag, every field but X has the tag.
var fuzzFields = []reflect.StructField{
	{Name: "S", Type: reflect.TypeFor[string]()},
	{Name: "P", Type: reflect.TypeFor[*string]()},
	{Name: "N", Type: reflect.TypeFor[int]()},
	{Name: "U", Type: reflect.TypeFor[*uint8]()},
	{Name: "F", Type: reflect.TypeFor[float64]()},
	{Name: "B", Type: reflect.TypeFor[bool]()},
	{Name: "T", Type: reflect.TypeFor[time.Time]()},
	{Name: "L", Type: reflect.TypeFor[[]string]()},
	{Name: "Q", Type: reflect.TypeFor[[]*int]()},
	{Name: "M", Type: reflect.TypeFor[map[string]int]()},
	{Name: "I", Type: reflect.TypeFor[[]fuzzItem]()},
	{Name: "R", Type: reflect.TypeFor[*fuzzItem]()},
}

// FuzzTags checks that a tag either configures the rules or panics with the error of the tag in Configure,
// and that the rules validate any value without panics.
func FuzzTags(f *testing.F) {
	for _, tag := range []string{
		"required", "omitempty,min=1,max=3", "dive,required", "email|url", "required_with=X", "required_if=X a",
		"gt=1,lt=5", "eqfield=X", "nefield=X", "dive,keys,alpha,endkeys,min=1", "len=2", "unique", "oneof=a b",
		"datetime=2006-01-02", "containsany=0x2C", "required_without_all=X N", "keys", "dive,dive", "|", "=", ",,",
	} {
		f.Add(tag, "value", int64(3))
	}
	f.Fuzz(func(t *testing.T, tag, s string, n int64) {
		for _, field := range fuzzFields {
			field.Tag = reflect.StructTag(`validate:"` + strings.NewReplacer(`"`, "", `\`, "").Replace(tag) + `"`)
			typ := reflect.StructOf([]reflect.StructField{field, {Name: "X", Type: reflect.TypeFor[string]()}})
			fns, ok := configure(t, typ)
			if !ok {
				continue
			}
			for _, obj := range []reflect.Value{reflect.New(typ), filled(typ, s, n)} {
				for _, fn := range fns {
					fn(context.Background(), helper.New(), obj.Interface(), nil)
				}
			}
		}
	})
}

// configure applies the tags of the struct type, false if they panic with the error of a tag.
func configure(t *testing.T, typ reflect.Type) (fns []shared.ValidationFn, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			if msg, isString := r.(string); !isString || !strings.HasPrefix(msg, "valigo: tag ") {
				t.Fatalf("%s: a panic of Configure isn't the error of a tag: %v", typ, r)
			}
			ok = false
		}
	}()
	model := reflect.New(typ).Interface()
	zero.Must(model)
	storage, err := fmap.GetFrom(model)
	if err != nil {
		t.Fatal(err)
	}
	tags.Apply(tags.Struct{Storage: storage, Model: model, Helper: helper.New(),
		Fields: func(ptr any) fmap.Field {
			field, err := storage.GetFieldByPtr(model, ptr)
			if err != nil {
				t.Fatalf("%s: %v", typ, err)
			}
			return field
		},
		Add: func(fn shared.ValidationFn) { fns = append(fns, fn) },
	}, "validate")
	return fns, true
}

// filled returns a value of the struct type with the fields made of s and n.
func filled(typ reflect.Type, s string, n int64) reflect.Value {
	v := reflect.New(typ)
	e := v.Elem()
	for i := range e.NumField() {
		f := e.Field(i)
		switch f.Interface().(type) {
		case string:
			f.SetString(s)
		case *string:
			f.Set(reflect.ValueOf(&s))
		case int:
			f.SetInt(n)
		case *uint8:
			u := uint8(n)
			f.Set(reflect.ValueOf(&u))
		case float64:
			f.SetFloat(float64(n) / 3)
		case bool:
			f.SetBool(n%2 == 0)
		case time.Time:
			f.Set(reflect.ValueOf(time.Unix(n, 0)))
		case []string:
			f.Set(reflect.ValueOf([]string{s, "", s}))
		case []*int:
			k := int(n)
			f.Set(reflect.ValueOf([]*int{&k, nil}))
		case map[string]int:
			f.Set(reflect.ValueOf(map[string]int{s: int(n), "": 0}))
		case []fuzzItem:
			f.Set(reflect.ValueOf([]fuzzItem{{Name: s}, {}}))
		case *fuzzItem:
			f.Set(reflect.ValueOf(&fuzzItem{Name: s}))
		}
	}
	return v
}
