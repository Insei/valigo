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
	"unsafe"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/internal/zero"
	"github.com/insei/valigo/rules/eq"
	"github.com/insei/valigo/rules/num"
	"github.com/insei/valigo/rules/str"
	"github.com/insei/valigo/rules/timex"
	"github.com/insei/valigo/shared"
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
	b := &builder{name: name, helper: s.Helper, structs: map[reflect.Type]*[]shared.ValidationFn{}}
	b.apply(s)
}

// builder configures the rules of the tags of a struct and of the structs of its elements.
type builder struct {
	name   string
	helper shared.Helper
	// structs are the rules of the struct types of the elements, built once: a recursive type refers to its rules
	structs map[reflect.Type]*[]shared.ValidationFn
}

func (b *builder) apply(s Struct) {
	for _, path := range s.Storage.GetAllPaths() {
		field := s.Storage.MustFind(path)
		tag := field.GetTag().Get(b.name)
		if tag == "" || tag == "-" {
			continue
		}
		t := tagContext{field: field, h: s.Helper, fields: s.Fields, add: inStruct(field, s.Add), describe: s.Describe,
			storage: s.Storage, model: s.Model, b: b}
		if err := t.apply(parseTag(tag)); err != nil {
			panic(fmt.Sprintf("valigo: tag %s:%q of field %s: %v", b.name, tag, path, err))
		}
	}
}

// rulesOf returns the rules of the tags of the struct type of the elements of a slice.
func (b *builder) rulesOf(typ reflect.Type) *[]shared.ValidationFn {
	if fns, ok := b.structs[typ]; ok {
		return fns
	}
	fns := new([]shared.ValidationFn)
	b.structs[typ] = fns
	model := reflect.New(typ).Interface()
	zero.Must(model)
	storage, _ := fmap.GetFrom(model) // a pointer to a struct
	b.apply(Struct{Storage: storage, Model: model, Helper: b.helper,
		Fields: func(fieldPtr any) fmap.Field {
			field, _ := storage.GetFieldByPtr(model, fieldPtr) // the pointers of otherField
			return field
		},
		Add: func(fn shared.ValidationFn) { *fns = append(*fns, fn) },
	})
	return fns
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
	return when(add, shared.InStruct(field))
}

// when applies the rules registered by the returned function only when cond returns true for the object.
func when(add func(fn shared.ValidationFn), cond func(obj any) bool) func(fn shared.ValidationFn) {
	return func(fn shared.ValidationFn) {
		add(func(ctx context.Context, h shared.Helper, obj any, errs []shared.Error) []shared.Error {
			if !cond(obj) {
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

// parseTag splits the tag to the rules, 0x2C in the parameters is a comma and 0x7C is a pipe, like in
// go-playground/validator.
func parseTag(tag string) []tagRule {
	parts := strings.Split(tag, ",")
	rules := make([]tagRule, len(parts))
	for i, part := range parts {
		if strings.Contains(part, "|") {
			// the alternatives, see alternatives
			rules[i] = tagRule{name: part}
			continue
		}
		name, param, _ := strings.Cut(part, "=")
		rules[i] = tagRule{name: name, param: unescape(param)}
	}
	return rules
}

func unescape(param string) string {
	return strings.NewReplacer("0x2C", ",", "0x7C", "|").Replace(param)
}

// alternatives returns the rules of the alternatives of a rule with |, e.g. rgb|rgba, or false.
func (r tagRule) alternatives() ([]tagRule, bool) {
	if !strings.Contains(r.name, "|") {
		return nil, false
	}
	var alts []tagRule
	for _, part := range strings.Split(r.name, "|") {
		name, param, _ := strings.Cut(part, "=")
		alts = append(alts, tagRule{name: name, param: unescape(param)})
	}
	return alts, true
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
	b        *builder
	// elem is true for the rules of the elements after dive
	elem bool
}

var timeType = reflect.TypeFor[time.Time]()

// apply applies the rules to the field, the rules after dive to its elements.
func (t tagContext) apply(rules []tagRule) error {
	rules, err := t.conditions(rules)
	if err != nil {
		return err
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

// conditions applies the conditional rules required_if, required_unless, required_with, required_with_all,
// required_without and required_without_all, and returns the other rules.
func (t tagContext) conditions(rules []tagRule) ([]tagRule, error) {
	var rest []tagRule
	for i, r := range rules {
		if r.name == "dive" {
			// the rules of the elements, the conditions of the elements are checked by the elements
			return append(rest, rules[i:]...), nil
		}
		if !strings.HasPrefix(r.name, "required_") {
			rest = append(rest, r)
			continue
		}
		cond, err := t.condition(r)
		if err != nil {
			return nil, err
		}
		// the conditional rules aren't described in JSON Schema
		conditional := t
		conditional.add, conditional.describe = when(t.add, cond), nil
		if err := conditional.apply([]tagRule{{name: "required"}}); err != nil {
			return nil, err
		}
	}
	return rest, nil
}

// condition returns the condition of a conditional rule by the values of the fields of its parameter.
func (t tagContext) condition(r tagRule) (func(obj any) bool, error) {
	params := strings.Fields(r.param)
	if len(params) == 0 {
		return nil, fmt.Errorf("%s needs the fields", r.name)
	}
	switch r.name {
	case "required_if", "required_unless":
		if len(params)%2 != 0 {
			return nil, fmt.Errorf("%s=%s: the pairs of a field and a value expected", r.name, r.param)
		}
		var equals []func(obj any) bool
		for i := 0; i < len(params); i += 2 {
			other, err := t.other(r, params[i])
			if err != nil {
				return nil, err
			}
			eq, err := equalsOf(other, params[i+1])
			if err != nil {
				return nil, fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			equals = append(equals, eq)
		}
		all := func(obj any) bool {
			for _, eq := range equals {
				if !eq(obj) {
					return false
				}
			}
			return true
		}
		if r.name == "required_if" {
			return all, nil
		}
		return func(obj any) bool { return !all(obj) }, nil
	case "required_with", "required_with_all", "required_without", "required_without_all":
		var presents []func(obj any) bool
		for _, p := range params {
			other, err := t.other(r, p)
			if err != nil {
				return nil, err
			}
			present, err := presentOf(other)
			if err != nil {
				return nil, fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			presents = append(presents, present)
		}
		// with: any is present, with_all: all are present, without: any is absent, without_all: all are absent
		want, all := !strings.HasPrefix(r.name, "required_without"), strings.HasSuffix(r.name, "_all")
		return func(obj any) bool {
			for _, present := range presents {
				if present(obj) != want && all {
					return false
				}
				if present(obj) == want && !all {
					return true
				}
			}
			return all
		}, nil
	}
	return nil, unsupported(r)
}

// presentOf returns whether the value of the field isn't zero: not nil, not empty, not the zero value.
func presentOf(field fmap.Field) (func(obj any) bool, error) {
	typ := field.GetType()
	switch typ.Kind() {
	case reflect.Pointer:
		get := shared.Getter[*struct{}](field)
		return func(obj any) bool { p := get(obj); return p != nil && *p != nil }, nil
	case reflect.Slice:
		get := shared.Getter[[]struct{}](field)
		return func(obj any) bool { p := get(obj); return p != nil && len(*p) > 0 }, nil
	case reflect.Map:
		// only the length is read, the type of the elements doesn't matter
		get := shared.Getter[map[struct{}]struct{}](field)
		return func(obj any) bool { p := get(obj); return p != nil && len(*p) > 0 }, nil
	}
	kind, ok := kinds[typ.Kind()]
	if typ == timeType {
		kind, ok = timeTags, true
	}
	if !ok {
		return nil, fmt.Errorf("the field %s of type %s can't be a condition", field.GetStructPath(), typ)
	}
	return kind.present(field), nil
}

// equalsOf returns whether the value of the field is equal to the parameter.
func equalsOf(field fmap.Field, param string) (func(obj any) bool, error) {
	typ := field.GetType()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	kind, ok := kinds[typ.Kind()]
	if !ok {
		return nil, fmt.Errorf("the field %s of type %s can't be compared with a value", field.GetStructPath(), field.GetType())
	}
	return kind.equals(field, param)
}

// structRules applies the rules of a nested struct field: required of a pointer, the fields have their own tags.
func (t tagContext) structRules(ptr bool, rules []tagRule) error {
	for _, r := range rules {
		switch {
		case r.name == "omitempty" || r.name == "omitnil":
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

// container applies the rules of a slice or a map field, the rules after dive to its elements or values,
// the rules between keys and endkeys after dive to the keys of a map.
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
	if !isMap && elem.Kind() == reflect.Struct && elem != timeType {
		if unique {
			return errors.New("unique of the slices of structs isn't supported")
		}
		return t.structElems(ptr, elem, elemPtr, elems)
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
		keys, values, err := splitKeys(elems)
		if err != nil {
			return err
		}
		return kind.mapRules(t, keys, values)
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

// splitKeys returns the rules of the keys between keys and endkeys and the rules of the values after them.
func splitKeys(rules []tagRule) (keys, values []tagRule, err error) {
	if len(rules) == 0 || rules[0].name != "keys" {
		return nil, rules, nil
	}
	for i, r := range rules {
		if r.name == "endkeys" {
			return rules[1:i], rules[i+1:], nil
		}
	}
	return nil, nil, errors.New("keys without endkeys")
}

// structElems applies the rules of the tags of the struct type of the elements of the slice field to every element,
// the location of an error is Items[1].City. The rules after dive are required and omitempty of the elements.
func (t tagContext) structElems(ptr bool, elem reflect.Type, elemPtr bool, rules []tagRule) error {
	required := false
	for _, r := range rules {
		switch {
		case r.name == "omitempty" || r.name == "omitnil":
		case r.name == "required" && elemPtr:
			required = true
		default:
			return unsupported(r)
		}
	}
	fns, size, h, field := t.b.rulesOf(elem), elem.Size(), t.h, t.field
	msg := shared.M("value.required", valueRequiredKey)
	check := func(ctx context.Context, i int, p unsafe.Pointer, errs []shared.Error) []shared.Error {
		n := len(errs)
		if p == nil {
			if !required {
				return errs
			}
			errs = append(errs, msg.Error(ctx, h, field, nil))
		} else {
			obj := reflect.NewAt(elem, p).Interface()
			for _, fn := range *fns {
				errs = fn(ctx, h, obj, errs)
			}
		}
		if len(errs) == n {
			return errs
		}
		location := h.FieldLocation(field) + "[" + strconv.Itoa(i) + "]"
		for j := n; j < len(errs); j++ {
			if p == nil {
				errs[j].Location = location
			} else {
				errs[j].Location = location + "." + errs[j].Location
			}
		}
		return errs
	}
	if elemPtr {
		valueField[[]*struct{}](t, ptr).Custom(func(ctx context.Context, _ *shared.FieldCustomHelper, s []*struct{}) []shared.Error {
			var errs []shared.Error
			for i, p := range s {
				errs = check(ctx, i, unsafe.Pointer(p), errs)
			}
			return errs
		})
		return nil
	}
	// the elements are read by the data of the slice and the size of the element type
	valueField[[]struct{}](t, ptr).Custom(func(ctx context.Context, _ *shared.FieldCustomHelper, s []struct{}) []shared.Error {
		var errs []shared.Error
		data := unsafe.Pointer(unsafe.SliceData(s))
		for i := range s {
			errs = check(ctx, i, unsafe.Add(data, uintptr(i)*size), errs)
		}
		return errs
	})
	return nil
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
	minLen := func(n int) {
		f.Schema(m.minWord, n)
		check(func(v int) bool { return v >= n }, shared.M(m.code+".min_len", m.minKey, shared.P("min", n)))
	}
	maxLen := func(n int) {
		f.Schema(m.maxWord, n)
		check(func(v int) bool { return v <= n }, shared.M(m.code+".max_len", m.maxKey, shared.P("max", n)))
	}
	for _, r := range rules {
		switch r.name {
		case "omitempty", "omitnil":
			if r.name == "omitempty" {
				f = f.When(func(_ context.Context, c C) bool { return length(c) > 0 })
			}
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
		case "min", "max", "len", "eq", "gt", "gte", "lt", "lte":
			n, err := strconv.Atoi(r.param)
			if err != nil {
				return fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			switch r.name {
			case "min", "gte":
				minLen(n)
			case "gt":
				minLen(n + 1)
			case "max", "lte":
				maxLen(n)
			case "lt":
				maxLen(n - 1)
			default:
				minLen(n)
				maxLen(n)
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
	// mapRules applies the rules to the keys and the values of the map[string]V field
	mapRules(t tagContext, keys, values []tagRule) error
	// unique checks that the slice field has no duplicates
	unique(t tagContext)
	// present returns whether the value of the field isn't zero
	present(field fmap.Field) func(obj any) bool
	// equals returns whether the value of the field is equal to the parameter
	equals(field fmap.Field, param string) (func(obj any) bool, error)
}

// typedTags is kindTags of the values of type V by the rules of apply, parse parses a parameter of the conditions,
// nil for time.Time.
type typedTags[V comparable] struct {
	apply func(t tagContext, f *shared.Field[V], optional bool, rules []tagRule) error
	parse func(s string) (V, error)
}

func (k typedTags[V]) value(t tagContext, ptr bool, rules []tagRule) error {
	return k.applyValue(t, valueField[V](t, ptr), ptr, rules)
}

func (k typedTags[V]) elems(t tagContext, ptr bool, rules []tagRule) error {
	t.elem = true
	s := shared.NewSliceFieldConfigurator[V](t.field, t.h, t.fields, t.add, t.describe)
	return k.applyValue(t, s.EachField(), ptr, rules)
}

func (k typedTags[V]) mapRules(t tagContext, keys, values []tagRule) error {
	t.elem = true
	m := shared.NewMapFieldConfigurator[string, V](t.field, t.h, t.fields, t.add, t.describe)
	var err error
	if len(keys) > 0 {
		m.EachKey(func(f *shared.Field[string]) {
			err = stringKind.applyValue(t, f, false, keys)
		})
	}
	if err == nil && len(values) > 0 {
		m.EachValue(func(f *shared.Field[V]) {
			err = k.applyValue(t, f, false, values)
		})
	}
	return err
}

func (k typedTags[V]) unique(t tagContext) {
	shared.NewSliceFieldConfigurator[V](t.field, t.h, t.fields, t.add, t.describe).As(eq.Slice).Unique()
}

func (k typedTags[V]) present(field fmap.Field) func(obj any) bool {
	get := shared.Getter[V](field)
	var zero V
	return func(obj any) bool {
		p := get(obj)
		return p != nil && *p != zero
	}
}

// equals is called for the kinds of the parameters, they have parse.
func (k typedTags[V]) equals(field fmap.Field, param string) (func(obj any) bool, error) {
	value, err := k.parse(param)
	if err != nil {
		return nil, err
	}
	get := shared.Getter[V](field)
	return func(obj any) bool {
		p := get(obj)
		return p != nil && *p == value
	}, nil
}

// applyValue applies the rules: omitempty skips the next rules for the zero value, omitnil is the absent value
// of a nil pointer, the alternatives of | pass when one of them passes.
func (k typedTags[V]) applyValue(t tagContext, f *shared.Field[V], optional bool, rules []tagRule) error {
	for _, r := range rules {
		switch alts, ok := r.alternatives(); {
		case r.name == "omitempty":
			var zero V
			f = f.When(func(_ context.Context, v V) bool { return v != zero })
		case r.name == "omitnil":
		case strings.HasPrefix(r.name, "required_") && t.elem:
			return fmt.Errorf("%s isn't supported for the elements", r.name)
		case ok:
			var err error
			configure := make([]func(alt *shared.Field[V]), len(alts))
			for i, alt := range alts {
				configure[i] = func(f *shared.Field[V]) {
					if e := k.apply(t, f, optional, []tagRule{alt}); e != nil && err == nil {
						err = e
					}
				}
			}
			f.Or(configure...)
			if err != nil {
				return err
			}
		default:
			if err := k.apply(t, f, optional, []tagRule{r}); err != nil {
				return err
			}
		}
	}
	return nil
}

var stringKind = typedTags[string]{apply: stringTags, parse: func(s string) (string, error) { return s, nil }}

var kinds = map[reflect.Kind]kindTags{
	reflect.String:  stringKind,
	reflect.Bool:    typedTags[bool]{apply: boolTags, parse: strconv.ParseBool},
	reflect.Int:     numberKind[int](),
	reflect.Int8:    numberKind[int8](),
	reflect.Int16:   numberKind[int16](),
	reflect.Int32:   numberKind[int32](),
	reflect.Int64:   numberKind[int64](),
	reflect.Uint:    numberKind[uint](),
	reflect.Uint8:   numberKind[uint8](),
	reflect.Uint16:  numberKind[uint16](),
	reflect.Uint32:  numberKind[uint32](),
	reflect.Uint64:  numberKind[uint64](),
	reflect.Float32: numberKind[float32](),
	reflect.Float64: numberKind[float64](),
}

func numberKind[N num.Number]() typedTags[N] {
	return typedTags[N]{apply: numberTags[N], parse: parseNumber[N]}
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
	"alpha":            (*str.BaseConfigurator).Alpha,
	"alphanum":         (*str.BaseConfigurator).Alphanumeric,
	"number":           (*str.BaseConfigurator).Digits,
	"numeric":          (*str.BaseConfigurator).Numeric,
	"hexadecimal":      (*str.BaseConfigurator).Hexadecimal,
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

// stringParamRules are the string rules with a parameter.
var stringParamRules = map[string]func(c *str.BaseConfigurator, param string) *str.BaseConfigurator{
	"eq":             func(c *str.BaseConfigurator, v string) *str.BaseConfigurator { return c.AnyOf(v) },
	"startswith":     (*str.BaseConfigurator).Prefix,
	"endswith":       (*str.BaseConfigurator).Suffix,
	"contains":       (*str.BaseConfigurator).Contains,
	"startsnotwith":  (*str.BaseConfigurator).NotPrefix,
	"endsnotwith":    (*str.BaseConfigurator).NotSuffix,
	"excludes":       (*str.BaseConfigurator).Excludes,
	"containsany":    (*str.BaseConfigurator).ContainsAny,
	"excludesall":    (*str.BaseConfigurator).ExcludesAny,
	"eq_ignore_case": (*str.BaseConfigurator).EqualFold,
	"datetime":       (*str.BaseConfigurator).Datetime,
}

func stringTags(t tagContext, f *shared.Field[string], _ bool, rules []tagRule) error {
	c := str.NewBaseConfigurator(f)
	for _, r := range rules {
		if rule, ok := stringRules[r.name]; ok {
			rule(c)
			continue
		}
		if rule, ok := stringParamRules[r.name]; ok {
			rule(c, r.param)
			continue
		}
		switch r.name {
		case "min", "max", "len", "gt", "gte", "lt", "lte":
			n, err := strconv.Atoi(r.param)
			if err != nil {
				return fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			switch r.name {
			case "min", "gte":
				c.MinRunes(n)
			case "gt":
				c.MinRunes(n + 1)
			case "max", "lte":
				c.MaxRunes(n)
			case "lt":
				c.MaxRunes(n - 1)
			default:
				c.MinRunes(n).MaxRunes(n)
			}
		case "oneof":
			c.AnyOf(strings.Fields(r.param)...)
		case "ne":
			eq.NewBaseConfigurator(f, false).NotEqual(r.param)
		case "eqfield", "nefield":
			other, err := otherField[string](t, r)
			if err != nil {
				return err
			}
			if r.name == "eqfield" {
				c.EqualField(other)
			} else {
				c.NotEqualField(other)
			}
		default:
			return unsupported(r)
		}
	}
	return nil
}

func boolTags(_ tagContext, f *shared.Field[bool], optional bool, rules []tagRule) error {
	c := eq.NewBaseConfigurator(f, optional)
	for _, r := range rules {
		switch r.name {
		case "required":
			c.Required()
		case "eq", "ne":
			b, err := strconv.ParseBool(r.param)
			if err != nil {
				return fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			if r.name == "eq" {
				c.Equal(b)
			} else {
				c.NotEqual(b)
			}
		default:
			return unsupported(r)
		}
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
		case "min", "gte", "max", "lte", "gt", "lt", "eq", "ne":
			n, err := parseNumber[N](r.param)
			if err != nil {
				return fmt.Errorf("%s=%s: %w", r.name, r.param, err)
			}
			switch r.name {
			case "min", "gte":
				c.Min(n)
			case "max", "lte":
				c.Max(n)
			case "gt":
				c.GreaterThan(n)
			case "lt":
				c.LessThan(n)
			case "ne":
				eq.NewBaseConfigurator(f, optional).NotEqual(n)
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
		case "gtefield", "ltefield", "gtfield", "ltfield", "eqfield", "nefield":
			other, err := otherField[N](t, r)
			if err != nil {
				return err
			}
			switch r.name {
			case "gtefield":
				c.MinField(other)
			case "ltefield":
				c.MaxField(other)
			case "gtfield":
				c.GreaterField(other)
			case "ltfield":
				c.LessField(other)
			case "eqfield":
				eq.NewBaseConfigurator(f, optional).EqualField(other)
			default:
				eq.NewBaseConfigurator(f, optional).NotEqualField(other)
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
		case "gt", "gte", "lt", "lte":
			// without a parameter, the comparison with the current time
			if r.param != "" {
				return fmt.Errorf("%s=%s: the times are compared with the current time, without a parameter", r.name, r.param)
			}
			if strings.HasPrefix(r.name, "gt") {
				c.InFuture()
			} else {
				c.InPast()
			}
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

// other returns the field of the model by the name of the parameter of the rule: a sibling field of the field
// of the tag or its nested field, Inner.Field.
func (t tagContext) other(r tagRule, name string) (fmap.Field, error) {
	path := name
	if parent := t.field.GetParent(); parent != nil {
		path = parent.GetStructPath() + "." + path
	}
	other, ok := t.storage.Find(path)
	if !ok {
		return nil, fmt.Errorf("%s=%s: no field %s", r.name, r.param, path)
	}
	return other, nil
}

// otherField returns the pointer to the field of the model named by the parameter of the rule, see other.
// Its type must be V or a named type of V.
func otherField[V any](t tagContext, r tagRule) (*V, error) {
	if t.elem {
		return nil, fmt.Errorf("%s isn't supported for the elements", r.name)
	}
	other, err := t.other(r, r.param)
	if err != nil {
		return nil, err
	}
	want := reflect.TypeFor[V]()
	if typ := other.GetType(); typ.Kind() != want.Kind() || typ.Size() != want.Size() {
		return nil, fmt.Errorf("%s=%s: the field %s has type %s, %s expected", r.name, r.param, other.GetStructPath(), typ, want)
	}
	return (*V)(reflect.ValueOf(other.GetPtr(t.model)).UnsafePointer()), nil
}

func unsupported(r tagRule) error {
	return fmt.Errorf("the rule %s isn't supported for the field, use the rules in code", r.name)
}
