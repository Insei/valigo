package valigo

import (
	"context"
	"reflect"
	"unsafe"

	"github.com/insei/fmap/v3"

	"github.com/insei/valigo/shared"
)

// Validator holds validation rules of the configured types.
type Validator struct {
	validators     map[reflect.Type][]shared.ValidationFn
	helper         *helper
	transformError func(errs []shared.Error) []error
}

// ValidateTyped validates obj using the rules configured for its type.
// obj is a pointer to the configured struct, a struct value is validated by a copy,
// so modifications like Trim are not visible to the caller.
// It returns nil if there are no errors, no rules for the obj type or obj is a nil pointer.
func (v *Validator) ValidateTyped(ctx context.Context, obj any) []shared.Error {
	validators, ok := v.validators[reflect.TypeOf(obj)]
	if !ok {
		validators, obj = v.byValue(obj)
	} else if dataPtr(obj) == nil {
		return nil
	}
	var errs []shared.Error
	for _, validator := range validators {
		errs = append(errs, validator(ctx, v.helper, obj)...)
	}
	return errs
}

// byValue returns validators and a pointer to the copy of obj if obj is a value of the configured struct.
func (v *Validator) byValue(obj any) ([]shared.ValidationFn, any) {
	val := reflect.ValueOf(obj)
	if !val.IsValid() || val.Kind() == reflect.Pointer {
		return nil, obj
	}
	validators, ok := v.validators[reflect.PointerTo(val.Type())]
	if !ok {
		return nil, obj
	}
	ptr := reflect.New(val.Type())
	ptr.Elem().Set(val)
	return validators, ptr.Interface()
}

// Validate is similar to ValidateTyped, but returns errors transformed by the errors transformer.
// It returns nil if there are no errors.
func (v *Validator) Validate(ctx context.Context, obj any) []error {
	errsTyped := v.ValidateTyped(ctx, obj)
	if len(errsTyped) == 0 {
		return nil
	}
	return v.transformError(errsTyped)
}

// dataPtr returns the data pointer of the interface value, i.e. the pointer itself for pointer values.
func dataPtr(v any) unsafe.Pointer {
	return (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[1]
}

// GetHelper returns the helper object associated with the Validator instance.
func (v *Validator) GetHelper() shared.Helper {
	return v.helper
}

// New creates a new Validator instance with default values for storage and helper.
// It also applies any options passed to the function to the new instance.
func New(opts ...Option) *Validator {
	v := &Validator{
		validators: make(map[reflect.Type][]shared.ValidationFn),
		helper:     newHelper(),
		transformError: func(errs []shared.Error) []error {
			errsNew := make([]error, len(errs))
			for i, err := range errs {
				errsNew[i] = err
			}
			return errsNew
		},
	}
	for _, opt := range opts {
		opt.apply(v)
	}
	return v
}

// Configure configures validation rules for the struct type T, replacing the rules of previous Configure calls for T.
// fn receives a model with all nested pointer structs allocated, rules are configured
// by pointers to the model fields. Configure panics on invalid configuration,
// call it on application start and not concurrently with Validate.
func Configure[T any](v *Validator, fn func(builder Configurator[T], m *T)) {
	delete(v.validators, reflect.TypeFor[*T]())
	model := new(T)
	// allocate all pointer fields values recursively
	mustZero(model)
	fields, err := fmap.GetFrom(model)
	if err != nil {
		panic(err)
	}
	fn(Configurator[T]{b: &builder[T]{v: v, obj: model, fields: fields}}, model)
}
