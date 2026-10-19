package shared

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/insei/fmap/v3"
)

// ValidationFn validates an object, obj is a pointer to the configured struct.
type ValidationFn func(ctx context.Context, h Helper, obj any, errs []Error) []Error

// checkFn checks a single value: a struct field or a slice element, obj is the validated object.
// v is nil when the value is absent, i.e. a nil pointer field or a field of a nil nested struct.
type checkFn[T any] func(ctx context.Context, h Helper, obj any, v *T, errs []Error) []Error

// Fields resolves a pointer to a field of the model of the configured struct to the field, it panics for other pointers.
type Fields func(fieldPtr any) fmap.Field

// Field is a typed rules set for a single value: a struct field or every element of a slice.
// Its rules are the primitives of the typed configurators, convert it to one with As, e.g. f.As(str.Of).
type Field[T any] struct {
	field  fmap.Field
	h      Helper
	fields Fields
	add    func(fn checkFn[T])
	// root registers validation functions of the struct field itself, nil for the fields derived by When,
	// conversions and for slice elements
	root func(fn ValidationFn)
	// rootDescribe describes the constraints of the struct field itself for SlicePtr
	rootDescribe Describe
	schema       constraints
}

// NewValueField creates a Field for the struct field of type T or *T,
// named types with the same underlying layout are allowed.
// fields resolves the pointers to other fields of the model, add registers a validation function of the configured struct,
// describe registers the constraints of its rules, it can be nil.
func NewValueField[T any](field fmap.Field, h Helper, fields Fields, add func(fn ValidationFn), describe Describe) *Field[T] {
	return &Field[T]{field: field, h: h, fields: fields, add: fieldSource(add, valueGetter[T](field)), root: add,
		rootDescribe: describe, schema: newConstraints(describe, field)}
}

// newConstraints returns the constraints of the struct field, describe can be nil.
func newConstraints(describe Describe, field fmap.Field) constraints {
	if describe == nil {
		return constraints{}
	}
	return constraints{describe: describe, field: field, modified: new(bool)}
}

// Schema describes the rules of the value by the keyword of JSON Schema with the value, e.g. Schema("maxLength", 64),
// see Constraint. The conversions call it for their rules that JSON Schema can express.
func (f *Field[T]) Schema(keyword string, value any) *Field[T] {
	f.schema.add(keyword, value)
	return f
}

// Check adds a rule that fails when ok returns false. Absent values are skipped.
// The code of the error is localeKey, use CheckMsg for a stable code and named parameters.
func (f *Field[T]) Check(ok func(v T) bool, localeKey string, args ...any) *Field[T] {
	return f.CheckMsg(ok, key(localeKey, args))
}

// CheckMsg adds a rule that fails with msg when ok returns false. Absent values are skipped.
func (f *Field[T]) CheckMsg(ok func(v T) bool, msg Msg) *Field[T] {
	field := f.field
	f.add(func(ctx context.Context, h Helper, obj any, v *T, errs []Error) []Error {
		if v == nil || ok(*v) {
			return errs
		}
		return append(errs, msg.Error(ctx, h, field, *v))
	})
	return f
}

// Required adds a rule that fails when the value is absent or ok returns false.
// The code of the error is localeKey, use RequiredMsg for a stable code.
func (f *Field[T]) Required(ok func(v T) bool, localeKey string) *Field[T] {
	return f.RequiredMsg(ok, key(localeKey, nil))
}

// RequiredMsg adds a rule that fails with msg when the value is absent or ok returns false.
func (f *Field[T]) RequiredMsg(ok func(v T) bool, msg Msg) *Field[T] {
	field := f.field
	f.add(func(ctx context.Context, h Helper, obj any, v *T, errs []Error) []Error {
		if v != nil && ok(*v) {
			return errs
		}
		var val any
		if v != nil {
			val = *v
		}
		return append(errs, msg.Error(ctx, h, field, val))
	})
	return f
}

// Modify adds a step that changes the value in place, e.g. trims a string. Absent values are skipped.
func (f *Field[T]) Modify(fn func(v *T)) *Field[T] {
	if f.schema.modified != nil {
		*f.schema.modified = true
	}
	f.add(func(_ context.Context, _ Helper, _ any, v *T, errs []Error) []Error {
		if v != nil {
			fn(v)
		}
		return errs
	})
	return f
}

// Custom adds a user defined rule. Absent values are skipped.
func (f *Field[T]) Custom(fn func(ctx context.Context, h *FieldCustomHelper, value T) []Error) *Field[T] {
	customHelper := NewFieldCustomHelper(f.field, f.h)
	f.add(func(ctx context.Context, _ Helper, _ any, v *T, errs []Error) []Error {
		if v == nil {
			return errs
		}
		return append(errs, fn(ctx, customHelper, *v)...)
	})
	return f
}

// CheckWith adds a rule that compares the value with the value of another field of the validated object,
// it fails when ok(value, other) returns false. other is a pointer to a field of the model, it panics in Configure
// for other pointers. The location of the other field is appended to args of the message.
// The rule is skipped when the value or the value of the other field is absent.
// The code of the error is localeKey, use CheckWithMsg for a stable code and named parameters.
func (f *Field[T]) CheckWith[U any](other *U, ok func(v T, other U) bool, localeKey string, args ...any) *Field[T] {
	return f.CheckWithMsg(other, ok, key(localeKey, args))
}

// CheckWithMsg is CheckWith with msg, the location of the other field is its last format argument
// and its "field" parameter.
func (f *Field[T]) CheckWithMsg[U any](other *U, ok func(v T, other U) bool, msg Msg) *Field[T] {
	field, otherField := f.field, f.fields(other)
	get := valueGetter[U](otherField)
	f.add(func(ctx context.Context, h Helper, obj any, v *T, errs []Error) []Error {
		if v == nil {
			return errs
		}
		if o := get(obj); o == nil || ok(*v, *o) {
			return errs
		}
		return append(errs, msg.with("field", h.FieldLocation(otherField)).Error(ctx, h, field, *v))
	})
	return f
}

// When returns a Field whose rules are applied only if cond returns true.
// When the value is absent, cond is not called and the rules, including Required, are skipped.
func (f *Field[T]) When(cond func(ctx context.Context, value T) bool) *Field[T] {
	if cond == nil {
		return f
	}
	add := f.add
	return &Field[T]{field: f.field, h: f.h, fields: f.fields, add: func(fn checkFn[T]) {
		add(func(ctx context.Context, h Helper, obj any, v *T, errs []Error) []Error {
			if v == nil || !cond(ctx, *v) {
				return errs
			}
			return fn(ctx, h, obj, v, errs)
		})
	}}
}

// Or adds a rule that passes when the rules of any of the alternatives pass, every alternative configures its rules
// on its own field: f.Or(func(a *Field[T]) { a.As(str.Of).Email() }, func(a *Field[T]) { a.As(str.Of).URL() }).
// When no alternative passes, the errors of the first one are reported. The alternatives are checked in order
// and the next ones are skipped after one passes. A value that passes the first alternative is checked without
// allocations; the errors of the failed alternatives are built and dropped, so the next alternatives allocate them.
func (f *Field[T]) Or(alternatives ...func(alt *Field[T])) *Field[T] {
	alts := make([][]checkFn[T], len(alternatives))
	for i, configure := range alternatives {
		configure(&Field[T]{field: f.field, h: f.h, fields: f.fields, add: func(fn checkFn[T]) {
			alts[i] = append(alts[i], fn)
		}})
	}
	f.add(func(ctx context.Context, h Helper, obj any, v *T, errs []Error) []Error {
		n, first := len(errs), len(errs)
		for i, checks := range alts {
			start := len(errs)
			for _, check := range checks {
				errs = check(ctx, h, obj, v, errs)
			}
			if len(errs) == start {
				clear(errs[n:start])
				return errs[:n]
			}
			if i == 0 {
				first = len(errs)
				continue
			}
			clear(errs[start:])
			errs = errs[:start]
		}
		return errs[:first]
	})
	return f
}

// Getter returns a function that reads a pointer to the value of the field of type T or *T from obj,
// a pointer to the configured struct, nil when the value is absent. Named types with the layout of T are allowed,
// it panics for other types. It reads the value without allocations, e.g. for the conditions of the rules.
func Getter[T any](field fmap.Field) func(obj any) *T {
	return valueGetter[T](field)
}

// fieldSource applies checks to the struct field value read by get.
// All checks of the field are registered as a single validation function, so the value is read once.
func fieldSource[T any](add func(fn ValidationFn), get func(obj any) *T) func(fn checkFn[T]) {
	var checks []checkFn[T]
	return func(fn checkFn[T]) {
		checks = append(checks, fn)
		if len(checks) > 1 {
			return
		}
		add(func(ctx context.Context, h Helper, obj any, errs []Error) []Error {
			v := get(obj)
			for _, check := range checks {
				errs = check(ctx, h, obj, v, errs)
			}
			return errs
		})
	}
}

// valueGetter returns a function that reads a pointer to the T value of the field from obj.
// The field type must be T or *T, named types with the same underlying layout are allowed.
// The returned function returns nil when the value is absent.
func valueGetter[T any](field fmap.Field) func(obj any) *T {
	isPtr := field.GetType().Kind() == reflect.Pointer && reflect.TypeFor[T]().Kind() != reflect.Pointer
	t := field.GetType()
	if isPtr {
		t = t.Elem()
	}
	mustLayout[T](field, t)
	addr := fieldAddr(field)
	if isPtr {
		return func(obj any) *T {
			p := (**T)(addr(dataPtr(obj)))
			if p == nil {
				return nil
			}
			return *p
		}
	}
	return func(obj any) *T {
		return (*T)(addr(dataPtr(obj)))
	}
}

// InStruct returns a function that reports whether the nested structs stored by pointer, which the field is in,
// are allocated in obj, a pointer to the configured struct: false for the fields of a nil nested struct.
func InStruct(field fmap.Field) func(obj any) bool {
	addr := fieldAddr(field)
	return func(obj any) bool {
		return addr(dataPtr(obj)) != nil
	}
}

// fieldAddr returns a function that computes the address of the field in the struct at root by the offsets of fmap,
// without reflect: the offset of a field is from the root struct or from the struct of its nearest ancestor stored
// by pointer, so the address follows the pointers of these ancestors. It returns nil if one of them is nil.
func fieldAddr(field fmap.Field) func(root unsafe.Pointer) unsafe.Pointer {
	// the offsets of the ancestors stored by pointer, from the root down to the field
	var owners []uintptr
	for parent := field.GetParent(); parent != nil; parent = parent.GetParent() {
		if t := parent.GetType(); t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
			owners = append([]uintptr{parent.GetOffset()}, owners...)
		}
	}
	offset := field.GetOffset()
	if len(owners) == 0 {
		return func(root unsafe.Pointer) unsafe.Pointer {
			return unsafe.Add(root, offset)
		}
	}
	return func(root unsafe.Pointer) unsafe.Pointer {
		base := root
		for _, owner := range owners {
			base = *(*unsafe.Pointer)(unsafe.Add(base, owner))
			if base == nil {
				return nil
			}
		}
		return unsafe.Add(base, offset)
	}
}

// dataPtr returns the data pointer of the interface value, i.e. the pointer itself for pointer values.
func dataPtr(v any) unsafe.Pointer {
	return (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[1]
}

// mustLayout panics if values of type t can't be read as T.
func mustLayout[T any](field fmap.Field, t reflect.Type) {
	want := reflect.TypeFor[T]()
	if t.Kind() != want.Kind() || t.Size() != want.Size() {
		panic(fmt.Sprintf("valigo: field %s has type %s, %s expected", field.GetStructPath(), field.GetType(), want))
	}
}
