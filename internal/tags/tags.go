// Package tags configures the rules of the tags of go-playground/validator by the rules of valigo.
package tags

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/eq"
	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
	"github.com/insei/valigo/timex"
)

// Struct is the configuration of the rules of a struct by its model, see valigo.Configure.
type Struct struct {
	// Storage is the fields of the model, Model is the pointer to the model.
	Storage fmap.Storage
	Model   any
	Helper  shared.Helper
	Fields  shared.Fields
	// Add registers the rules, Describe their constraints, it can be nil.
	Add      func(fn shared.ValidationFn)
	Describe shared.Describe
}

// Apply configures the rules of the tag name of the fields of the struct and of its nested structs,
// see valigo.Configurator.Tags. It panics on an unknown rule.
func Apply(s Struct, name string) {
	for _, path := range s.Storage.GetAllPaths() {
		field := s.Storage.MustFind(path)
		tag := field.GetTag().Get(name)
		if tag == "" || tag == "-" {
			continue
		}
		t := tagContext{field: field, h: s.Helper, fields: s.Fields, add: inStruct(field, s.Add), describe: s.Describe,
			storage: s.Storage, model: s.Model}
		if err := t.apply(parseTag(tag)); err != nil {
			panic(fmt.Sprintf("valigo: tag %s:%q of field %s: %v", name, tag, path, err))
		}
	}
}

// inStruct skips the rules of the field when the nested struct stored by pointer, which the field is in, is nil,
// like go-playground/validator does: the rules of the fields of an optional nested struct apply when it is present.
func inStruct(field fmap.Field, add func(fn shared.ValidationFn)) func(fn shared.ValidationFn) {
	nested := false
	for parent := field.GetParent(); parent != nil; parent = parent.GetParent() {
		nested = nested || parent.GetType().Kind() == reflect.Pointer
	}
	if !nested {
		return add
	}
	present := shared.InStruct(field)
	return func(fn shared.ValidationFn) {
		add(func(ctx context.Context, h shared.Helper, obj any, errs []shared.Error) []shared.Error {
			if !present(obj) {
				return errs
			}
			return fn(ctx, h, obj, errs)
		})
	}
}

// tagRule is a rule of a tag: name=param.
type tagRule struct {
	name, param string
}

func parseTag(tag string) []tagRule {
	parts := strings.Split(tag, ",")
	rules := make([]tagRule, len(parts))
	for i, part := range parts {
		name, param, _ := strings.Cut(part, "=")
		rules[i] = tagRule{name: name, param: param}
	}
	return rules
}

// tagContext is the field of the tag with the arguments of the configurators.
type tagContext struct {
	field    fmap.Field
	h        shared.Helper
	fields   shared.Fields
	add      func(fn shared.ValidationFn)
	describe shared.Describe
	storage  fmap.Storage
	model    any
	// elem is true for the rules of the elements after dive
	elem bool
}

var timeType = reflect.TypeFor[time.Time]()

// apply applies the rules to the field, the rules after dive to its elements.
func (t tagContext) apply(rules []tagRule) error {
	for _, r := range rules {
		if strings.Contains(r.name, "|") {
			return errors.New("the alternatives of | aren't supported, use When or Custom")
		}
	}
	typ := t.field.GetType()
	ptr := typ.Kind() == reflect.Pointer
	if ptr {
		typ = typ.Elem()
	}
	switch {
	case typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map:
		return t.container(typ, ptr, rules)
	case typ == timeType:
		return timeTags.value(t, ptr, rules)
	case typ.Kind() == reflect.Struct:
		return t.structRules(ptr, rules)
	}
	kind, ok := kinds[typ.Kind()]
	if !ok {
		return fmt.Errorf("the fields of type %s aren't supported", t.field.GetType())
	}
	return kind.value(t, ptr, rules)
}

// structRules applies the rules of a nested struct field: required of a pointer, the fields have their own tags.
func (t tagContext) structRules(ptr bool, rules []tagRule) error {
	for _, r := range rules {
		switch {
		case r.name == "omitempty":
		case r.name == "required" && ptr:
			f := shared.NewValueField[*struct{}](t.field, t.h, t.fields, t.add, t.describe)
			f.Schema(shared.KeywordRequired, true)
			f.RequiredMsg(func(p *struct{}) bool { return p != nil }, shared.M("value.required", valueRequiredKey))
		default:
			return unsupported(r)
		}
	}
	return nil
}

const (
	valueRequiredKey = "validation:value:Should be fulfilled"
	sliceRequiredKey = "validation:slice:Should be fulfilled"
	sliceMaxLenKey   = "validation:slice:Cannot contain more than %d elements"
	sliceMinLenKey   = "validation:slice:Cannot contain less than %d elements"
	mapRequiredKey   = "validation:map:Should be fulfilled"
	mapMaxLenKey     = "validation:map:Cannot contain more than %d elements"
	mapMinLenKey     = "validation:map:Cannot contain less than %d elements"
)

// container applies the rules of a slice or a map field, the rules after dive to its elements or values.
func (t tagContext) container(typ reflect.Type, ptr bool, rules []tagRule) error {
	own, elems := rules, []tagRule(nil)
	for i, r := range rules {
		if r.name == "dive" {
			own, elems = rules[:i], rules[i+1:]
			break
		}
	}
	isMap := typ.Kind() == reflect.Map
	if err := t.lenRules(isMap, ptr, own); err != nil {
		return err
	}
	unique := false
	for _, r := range own {
		unique = unique || r.name == "unique"
	}
	if elems == nil && !unique {
		return nil
	}
	elem := typ.Elem()
	elemPtr := elem.Kind() == reflect.Pointer
	if elemPtr {
		elem = elem.Elem()
	}
	kind, ok := kinds[elem.Kind()]
	if elem == timeType {
		kind, ok = timeTags, true
	}
	switch {
	case !ok:
		return fmt.Errorf("the elements of type %s aren't supported, configure the rules of their type", typ.Elem())
	case isMap && (typ.Key().Kind() != reflect.String || elemPtr):
		return errors.New("dive is supported for the maps of strings to values, map[string]V")
	case isMap:
		return kind.values(t, elems)
	case unique && elemPtr:
		return errors.New("unique of the slices of pointers isn't supported")
	case unique:
		kind.unique(t)
	}
	if elems == nil {
		return nil
	}
	return kind.elems(t, elemPtr, elems)
}

// lenRules applies the rules of the length of a slice or a map field.
func (t tagContext) lenRules(isMap, ptr bool, rules []tagRule) error {
	// the length is read by the header of any slice or map, the type of the elements doesn't matter;
	// the values of these types never get to the errors, only the lengths
	if isMap {
		return lengthRules(t, ptr, rules, func(m map[struct{}]struct{}) int { return len(m) }, lengthMsgs{
			code: "map", required: mapRequiredKey, minKey: mapMinLenKey, maxKey: mapMaxLenKey,
			minWord: "minProperties", maxWord: "maxProperties",
		})
	}
	return lengthRules(t, ptr, rules, func(s []struct{}) int { return len(s) }, lengthMsgs{
		code: "slice", required: sliceRequiredKey, minKey: sliceMinLenKey, maxKey: sliceMaxLenKey,
		minWord: "minItems", maxWord: "maxItems", unique: true,
	})
}

// lengthMsgs are the codes, the messages and the keywords of JSON Schema of the rules of the length.
type lengthMsgs struct {
	code, required, minKey, maxKey, minWord, maxWord string
	// unique is true for slices, unique is checked by the elements
	unique bool
}

func lengthRules[C any](t tagContext, ptr bool, rules []tagRule, length func(c C) int, m lengthMsgs) error {
	f := valueField[C](t, ptr)
	check := func(ok func(n int) bool, msg shared.Msg) {
		f.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, c C) []shared.Error {
			if n := length(c); !ok(n) {
				return []shared.Error{h.ErrorMsg(ctx, n, msg)}
			}
			return nil
		})
	}
	for _, r := range rules {
		switch r.name {
		case "omitempty":
			f = f.When(func(_ context.Context, c C) bool { return length(c) > 0 })
		case "unique":
			if !m.unique {
				return unsupported(r)
			}
		case "required":
			msg := shared.M(m.code+".required", m.required)
			f.Schema(shared.KeywordRequired, true).Schema(m.minWord, 1)
			// an absent value, the value isn't in the error
			f.RequiredMsg(func(C) bool { return true }, msg)
			check(func(n int) bool { return n > 0 }, msg)
		case "min", "max", "len":
			n, err := strconv.Atoi(r.param)
			if err != nil {
				return fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			if r.name != "max" {
				f.Schema(m.minWord, n)
				check(func(v int) bool { return v >= n }, shared.M(m.code+".min_len", m.minKey, shared.P("min", n)))
			}
			if r.name != "min" {
				f.Schema(m.maxWord, n)
				check(func(v int) bool { return v <= n }, shared.M(m.code+".max_len", m.maxKey, shared.P("max", n)))
			}
		default:
			return unsupported(r)
		}
	}
	return nil
}

// valueField returns the rules of the value of the field of type V or *V, a nil pointer is absent.
func valueField[V any](t tagContext, ptr bool) *shared.Field[V] {
	if ptr {
		return shared.Deref(shared.NewValueField[*V](t.field, t.h, t.fields, t.add, t.describe))
	}
	return shared.NewValueField[V](t.field, t.h, t.fields, t.add, t.describe)
}

// kindTags applies the rules of the values of one type.
type kindTags interface {
	// value applies the rules to the field of the type or of a pointer to it
	value(t tagContext, ptr bool, rules []tagRule) error
	// elems applies the rules to the elements of the slice field
	elems(t tagContext, ptr bool, rules []tagRule) error
	// values applies the rules to the values of the map[string]V field
	values(t tagContext, rules []tagRule) error
	// unique checks that the slice field has no duplicates
	unique(t tagContext)
}

// typedTags is kindTags of the values of type V by the rules of apply.
type typedTags[V comparable] struct {
	apply func(t tagContext, f *shared.Field[V], optional bool, rules []tagRule) error
}

func (k typedTags[V]) value(t tagContext, ptr bool, rules []tagRule) error {
	return k.applyValue(t, valueField[V](t, ptr), ptr, rules)
}

func (k typedTags[V]) elems(t tagContext, ptr bool, rules []tagRule) error {
	t.elem = true
	s := shared.NewSliceFieldConfigurator[V](t.field, t.h, t.fields, t.add, t.describe)
	return k.applyValue(t, s.EachField(), ptr, rules)
}

func (k typedTags[V]) values(t tagContext, rules []tagRule) error {
	t.elem = true
	var err error
	shared.NewMapFieldConfigurator[string, V](t.field, t.h, t.fields, t.add, t.describe).EachValue(func(f *shared.Field[V]) {
		err = k.applyValue(t, f, false, rules)
	})
	return err
}

func (k typedTags[V]) unique(t tagContext) {
	shared.NewSliceFieldConfigurator[V](t.field, t.h, t.fields, t.add, t.describe).As(eq.Slice).Unique()
}

// applyValue applies the rules, omitempty skips the next rules for the zero value.
func (k typedTags[V]) applyValue(t tagContext, f *shared.Field[V], optional bool, rules []tagRule) error {
	for i, r := range rules {
		if r.name == "omitempty" {
			var zero V
			f = f.When(func(_ context.Context, v V) bool { return v != zero })
			rules = append(rules[:i:i], rules[i+1:]...)
			break
		}
	}
	return k.apply(t, f, optional, rules)
}

var kinds = map[reflect.Kind]kindTags{
	reflect.String:  typedTags[string]{apply: stringTags},
	reflect.Bool:    typedTags[bool]{apply: boolTags},
	reflect.Int:     typedTags[int]{apply: numberTags[int]},
	reflect.Int8:    typedTags[int8]{apply: numberTags[int8]},
	reflect.Int16:   typedTags[int16]{apply: numberTags[int16]},
	reflect.Int32:   typedTags[int32]{apply: numberTags[int32]},
	reflect.Int64:   typedTags[int64]{apply: numberTags[int64]},
	reflect.Uint:    typedTags[uint]{apply: numberTags[uint]},
	reflect.Uint8:   typedTags[uint8]{apply: numberTags[uint8]},
	reflect.Uint16:  typedTags[uint16]{apply: numberTags[uint16]},
	reflect.Uint32:  typedTags[uint32]{apply: numberTags[uint32]},
	reflect.Uint64:  typedTags[uint64]{apply: numberTags[uint64]},
	reflect.Float32: typedTags[float32]{apply: numberTags[float32]},
	reflect.Float64: typedTags[float64]{apply: numberTags[float64]},
}

var timeTags = typedTags[time.Time]{apply: timeRules}

// stringRules are the string rules without a parameter.
var stringRules = map[string]func(c *str.BaseConfigurator) *str.BaseConfigurator{
	"required":         (*str.BaseConfigurator).Required,
	"email":            (*str.BaseConfigurator).Email,
	"url":              (*str.BaseConfigurator).URL,
	"http_url":         (*str.BaseConfigurator).URL,
	"uuid":             (*str.BaseConfigurator).UUID,
	"ip":               (*str.BaseConfigurator).IP,
	"ipv4":             (*str.BaseConfigurator).IPv4,
	"ipv6":             (*str.BaseConfigurator).IPv6,
	"hostname":         (*str.BaseConfigurator).Hostname,
	"hostname_rfc1123": (*str.BaseConfigurator).Hostname,
	"alphanum":         (*str.BaseConfigurator).Alphanumeric,
	"ascii":            (*str.BaseConfigurator).ASCII,
	"lowercase":        (*str.BaseConfigurator).Lowercase,
	"uppercase":        (*str.BaseConfigurator).Uppercase,
	"e164":             (*str.BaseConfigurator).E164,
	"base64":           (*str.BaseConfigurator).Base64,
	"base64url":        (*str.BaseConfigurator).Base64URL,
	"json":             (*str.BaseConfigurator).JSON,
	"semver":           (*str.BaseConfigurator).Semver,
	"cidr":             (*str.BaseConfigurator).CIDR,
	"mac":              (*str.BaseConfigurator).MAC,
	"hexcolor":         (*str.BaseConfigurator).HexColor,
	"iso3166_1_alpha2": (*str.BaseConfigurator).CountryCode,
	"iso4217":          (*str.BaseConfigurator).CurrencyCode,
	"credit_card":      (*str.BaseConfigurator).CreditCard,
	"luhn_checksum":    (*str.BaseConfigurator).Luhn,
}

func stringTags(t tagContext, f *shared.Field[string], _ bool, rules []tagRule) error {
	c := str.NewBaseConfigurator(f)
	for _, r := range rules {
		if rule, ok := stringRules[r.name]; ok {
			rule(c)
			continue
		}
		switch r.name {
		case "min", "max", "len":
			n, err := strconv.Atoi(r.param)
			if err != nil {
				return fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			if r.name != "max" {
				c.MinRunes(n)
			}
			if r.name != "min" {
				c.MaxRunes(n)
			}
		case "oneof":
			c.AnyOf(strings.Fields(r.param)...)
		case "eq":
			c.AnyOf(r.param)
		case "startswith":
			c.Prefix(r.param)
		case "endswith":
			c.Suffix(r.param)
		case "contains":
			c.Contains(r.param)
		case "datetime":
			c.Datetime(r.param)
		case "eqfield":
			other, err := otherField[string](t, r)
			if err != nil {
				return err
			}
			c.EqualField(other)
		default:
			return unsupported(r)
		}
	}
	return nil
}

func boolTags(_ tagContext, f *shared.Field[bool], optional bool, rules []tagRule) error {
	c := eq.NewBaseConfigurator(f, optional)
	for _, r := range rules {
		if r.name != "required" {
			return unsupported(r)
		}
		c.Required()
	}
	return nil
}

func numberTags[N num.Number](t tagContext, f *shared.Field[N], optional bool, rules []tagRule) error {
	c := num.NewBaseConfigurator(f, optional)
	for _, r := range rules {
		switch r.name {
		case "required":
			c.Required()
		case "latitude":
			c.Latitude()
		case "longitude":
			c.Longitude()
		case "min", "gte", "max", "lte", "eq":
			n, err := parseNumber[N](r.param)
			if err != nil {
				return fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			switch r.name {
			case "min", "gte":
				c.Min(n)
			case "max", "lte":
				c.Max(n)
			default:
				c.AnyOf(n)
			}
		case "oneof":
			var allowed []N
			for _, p := range strings.Fields(r.param) {
				n, err := parseNumber[N](p)
				if err != nil {
					return fmt.Errorf("oneof=%s: %w", r.param, err)
				}
				allowed = append(allowed, n)
			}
			c.AnyOf(allowed...)
		case "gtefield", "ltefield":
			other, err := otherField[N](t, r)
			if err != nil {
				return err
			}
			if r.name == "gtefield" {
				c.MinField(other)
			} else {
				c.MaxField(other)
			}
		default:
			return unsupported(r)
		}
	}
	return nil
}

// parseNumber parses the number of the type N: an integer for the integer types.
func parseNumber[N num.Number](s string) (N, error) {
	switch reflect.TypeFor[N]().Kind() {
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		return N(f), err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, int(reflect.TypeFor[N]().Bits()))
		return N(u), err
	default:
		i, err := strconv.ParseInt(s, 10, int(reflect.TypeFor[N]().Bits()))
		return N(i), err
	}
}

func timeRules(t tagContext, f *shared.Field[time.Time], _ bool, rules []tagRule) error {
	c := timex.NewBaseConfigurator(f)
	for _, r := range rules {
		switch r.name {
		case "required":
			c.Required()
		case "gtfield", "ltfield":
			other, err := otherField[time.Time](t, r)
			if err != nil {
				return err
			}
			if r.name == "gtfield" {
				c.AfterField(other)
			} else {
				c.BeforeField(other)
			}
		default:
			return unsupported(r)
		}
	}
	return nil
}

// otherField returns the pointer to the field of the model named by the parameter of the rule, a sibling field
// of the field of the tag or its nested field, Inner.Field. Its type must be V or a named type of V.
func otherField[V any](t tagContext, r tagRule) (*V, error) {
	if t.elem {
		return nil, fmt.Errorf("%s isn't supported for the elements", r.name)
	}
	path := r.param
	if parent := t.field.GetParent(); parent != nil {
		path = parent.GetStructPath() + "." + path
	}
	other, ok := t.storage.Find(path)
	if !ok {
		return nil, fmt.Errorf("%s=%s: no field %s", r.name, r.param, path)
	}
	want := reflect.TypeFor[V]()
	if typ := other.GetType(); typ.Kind() != want.Kind() || typ.Size() != want.Size() {
		return nil, fmt.Errorf("%s=%s: the field %s has type %s, %s expected", r.name, r.param, path, typ, want)
	}
	return (*V)(reflect.ValueOf(other.GetPtr(t.model)).UnsafePointer()), nil
}

func unsupported(r tagRule) error {
	return fmt.Errorf("the rule %s isn't supported for the field, use the rules in code", r.name)
}
