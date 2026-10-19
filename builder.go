package valigo

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	guuid "github.com/google/uuid"
	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
	"github.com/insei/valigo/uuid"
)

// Configurator configures validation rules for the struct type T.
// Field pointers passed to its methods must point to the fields of the model passed to the Configure function.
type Configurator[T any] struct {
	b *builder[T]
}

// builder holds the configuration state shared by Configurator copies.
type builder[T any] struct {
	v      *Validator
	obj    *T
	fields fmap.Storage
	when   func(ctx context.Context, obj *T) bool
}

// field returns the model field by the pointer to it.
func (c Configurator[T]) field(fieldPtr any) fmap.Field {
	field, err := c.b.fields.GetFieldByPtr(c.b.obj, fieldPtr)
	if err != nil {
		panic(fmt.Sprintf("valigo: %T is not a pointer to a field of %T: %v", fieldPtr, c.b.obj, err))
	}
	return field
}

// add registers a validation function, the struct level condition is applied if set.
func (c Configurator[T]) add(fn shared.ValidationFn) {
	if when := c.b.when; when != nil {
		inner := fn
		fn = func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
			if !when(ctx, obj.(*T)) {
				return nil
			}
			return inner(ctx, h, obj)
		}
	}
	t := reflect.TypeFor[*T]()
	c.b.v.validators[t] = append(c.b.v.validators[t], fn)
}

// When returns a Configurator whose rules are applied only if fn returns true.
// Conditions of chained When calls are combined with AND.
func (c Configurator[T]) When(fn func(ctx context.Context, obj *T) bool) Configurator[T] {
	if fn == nil {
		return c
	}
	when := fn
	if prev := c.b.when; prev != nil {
		when = func(ctx context.Context, obj *T) bool {
			return prev(ctx, obj) && fn(ctx, obj)
		}
	}
	b := *c.b
	b.when = when
	return Configurator[T]{b: &b}
}

// structHelper implements shared.StructCustomHelper for a single struct Custom call.
type structHelper struct {
	h      shared.Helper
	fields fmap.Storage
	obj    any
}

// ErrorT returns a shared.Error for the field of the validated object by the pointer to it.
func (s *structHelper) ErrorT(ctx context.Context, ptrToField, fieldValue any, localeKey string, args ...any) shared.Error {
	field, err := s.fields.GetFieldByPtr(s.obj, ptrToField)
	if err != nil {
		panic(fmt.Sprintf("valigo: %T is not a pointer to a field of %T: %v", ptrToField, s.obj, err))
	}
	return s.h.ErrorT(ctx, field, fieldValue, localeKey, args...)
}

// structHelpers reuses helpers, so struct Custom rules don't allocate on valid objects.
var structHelpers = sync.Pool{New: func() any { return new(structHelper) }}

// Custom adds a custom validation function on the whole struct.
// The helper passed to fn is valid only during the fn call.
func (c Configurator[T]) Custom(fn func(ctx context.Context, h shared.StructCustomHelper, obj *T) []shared.Error) {
	fields := c.b.fields
	c.add(func(ctx context.Context, h shared.Helper, obj any) []shared.Error {
		sh := structHelpers.Get().(*structHelper)
		sh.h, sh.fields, sh.obj = h, fields, obj
		errs := fn(ctx, sh, obj.(*T))
		*sh = structHelper{}
		structHelpers.Put(sh)
		return errs
	})
}

// String returns a configurator for a string or *string field.
func (c Configurator[T]) String[F ~string | *string](fieldPtr *F) *str.BaseConfigurator {
	field := c.field(fieldPtr)
	return str.NewBaseConfigurator(shared.NewValueField[string](field, c.b.v.helper, c.add))
}

// Number returns a configurator for a numeric field.
func (c Configurator[T]) Number[N num.Number](fieldPtr *N) *num.BaseConfigurator[N] {
	field := c.field(fieldPtr)
	return num.NewBaseConfigurator(shared.NewValueField[N](field, c.b.v.helper, c.add), false)
}

// NumberPtr returns a configurator for a pointer to numeric field, nil values are skipped by all rules except Required.
func (c Configurator[T]) NumberPtr[N num.Number](fieldPtr **N) *num.BaseConfigurator[N] {
	field := c.field(fieldPtr)
	return num.NewBaseConfigurator(shared.NewValueField[N](field, c.b.v.helper, c.add), true)
}

// UUID returns a configurator for an uuid.UUID or *uuid.UUID field.
func (c Configurator[T]) UUID[F guuid.UUID | *guuid.UUID](fieldPtr *F) *uuid.BaseConfigurator {
	field := c.field(fieldPtr)
	return uuid.NewBaseConfigurator(shared.NewValueField[guuid.UUID](field, c.b.v.helper, c.add))
}

// Slice returns a configurator for a slice field.
func (c Configurator[T]) Slice[E any](sliceFieldPtr *[]E) *shared.SliceFieldConfigurator[E] {
	field := c.field(sliceFieldPtr)
	return shared.NewSliceFieldConfigurator[E](field, c.b.v.helper, c.add)
}

// StringSlice returns a configurator for a slice of strings field.
func (c Configurator[T]) StringSlice[F []string | *[]string | []*string | *[]*string](sliceFieldPtr *F) *str.StringSliceFieldConfigurator {
	field := c.field(sliceFieldPtr)
	return str.NewStringSliceFieldConfigurator(shared.NewSliceFieldConfigurator[string](field, c.b.v.helper, c.add))
}

// UUIDSlice returns a configurator for a slice of uuids field.
func (c Configurator[T]) UUIDSlice[F []guuid.UUID | *[]guuid.UUID | []*guuid.UUID | *[]*guuid.UUID](sliceFieldPtr *F) *uuid.UUIDSliceFieldConfigurator {
	field := c.field(sliceFieldPtr)
	return uuid.NewUUIDSliceFieldConfigurator(shared.NewSliceFieldConfigurator[guuid.UUID](field, c.b.v.helper, c.add))
}
