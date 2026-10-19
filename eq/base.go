// Package eq configures the rules of comparable values: enums, bools, named types.
// b.Field(&m.Status).As(eq.Of).OneOf(Active, Blocked), b.Field(&m.Agreed).As(eq.Of).Equal(true).
package eq

import (
	"context"
	"slices"

	"github.com/insei/valigo/shared"
)

const (
	requiredLocaleKey = "validation:value:Should be fulfilled"
	oneOfLocaleKey    = "validation:value:Only %v values is allowed"
	notOneOfLocaleKey = "validation:value:Values %v are not allowed"
	equalLocaleKey    = "validation:value:Should be %v"
	notEqualLocaleKey = "validation:value:Cannot be %v"
	equalFieldKey     = "validation:value:Should be equal to %s"
)

// BaseConfigurator configures rules for a comparable value.
type BaseConfigurator[T comparable] struct {
	f        *shared.Field[T]
	optional bool
}

// NewBaseConfigurator creates a new BaseConfigurator.
// optional is true for pointer fields: Required then checks only that the value is present, zero is allowed.
func NewBaseConfigurator[T comparable](f *shared.Field[T], optional bool) *BaseConfigurator[T] {
	return &BaseConfigurator[T]{f: f, optional: optional}
}

// Of converts a field of a comparable type to the configurator: b.Field(&m.Status).As(eq.Of).
// For a pointer field use Ptr, Of compares the pointers themselves.
func Of[T comparable](f *shared.Field[T]) *BaseConfigurator[T] {
	return NewBaseConfigurator(f, false)
}

// Ptr converts a pointer to comparable field to the configurator of the pointed value: b.Field(&m.Role).As(eq.Ptr).
// A nil pointer is absent, Required checks only that the pointer is not nil, zero is allowed.
func Ptr[T comparable](f *shared.Field[*T]) *BaseConfigurator[T] {
	return NewBaseConfigurator(shared.Deref(f), true)
}

// Required checks if the value is not zero, for pointer fields checks if the pointer is not nil.
func (i *BaseConfigurator[T]) Required() *BaseConfigurator[T] {
	optional := i.optional
	var zero T
	i.f.Required(func(v T) bool { return optional || v != zero }, requiredLocaleKey)
	return i
}

// OneOf checks if the value is one of the allowed values.
func (i *BaseConfigurator[T]) OneOf(allowed ...T) *BaseConfigurator[T] {
	i.f.Check(func(v T) bool { return slices.Contains(allowed, v) }, oneOfLocaleKey, allowed)
	return i
}

// NotOneOf checks if the value is none of the forbidden values.
func (i *BaseConfigurator[T]) NotOneOf(forbidden ...T) *BaseConfigurator[T] {
	i.f.Check(func(v T) bool { return !slices.Contains(forbidden, v) }, notOneOfLocaleKey, forbidden)
	return i
}

// Equal checks if the value is equal to expected, e.g. Equal(true) for an accepted agreement.
func (i *BaseConfigurator[T]) Equal(expected T) *BaseConfigurator[T] {
	i.f.Check(func(v T) bool { return v == expected }, equalLocaleKey, expected)
	return i
}

// EqualField checks if the value is equal to the value of another field of the object:
// b.Field(&m.Confirmation).As(eq.Of).EqualField(&m.Password).
func (i *BaseConfigurator[T]) EqualField(other *T) *BaseConfigurator[T] {
	i.f.CheckWith(other, func(v, o T) bool { return v == o }, equalFieldKey)
	return i
}

// NotEqual checks if the value is not equal to forbidden.
func (i *BaseConfigurator[T]) NotEqual(forbidden T) *BaseConfigurator[T] {
	i.f.Check(func(v T) bool { return v != forbidden }, notEqualLocaleKey, forbidden)
	return i
}

// Custom allows for custom validation logic to be applied to the value.
func (i *BaseConfigurator[T]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value T) []shared.Error) *BaseConfigurator[T] {
	i.f.Custom(f)
	return i
}

// When allows for conditional validation logic to be applied to the value.
func (i *BaseConfigurator[T]) When(whenFn func(ctx context.Context, value T) bool) *BaseConfigurator[T] {
	return NewBaseConfigurator(i.f.When(whenFn), i.optional)
}
