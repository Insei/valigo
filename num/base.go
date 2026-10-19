package num

import (
	"context"
	"slices"

	"github.com/insei/valigo/shared"
)

const (
	minLocaleKey          = "validation:num:Cannot be less than %v"
	maxLocaleKey          = "validation:num:Cannot be greater than %v"
	requiredLocaleKey     = "validation:num:Should be fulfilled"
	anyOfLocaleKey        = "validation:num:Only %v values is allowed"
	anyOfIntervalLocalKey = "validation:num:Only interval[%v - %v] is allowed"
)

// Number is a constraint for numeric types supported by BaseConfigurator.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// BaseConfigurator configures rules for a numeric value.
type BaseConfigurator[N Number] struct {
	f        *shared.Field[N]
	optional bool
}

// NewBaseConfigurator creates a new BaseConfigurator.
// optional is true for pointer fields: Required then checks only that the value is present, zero is allowed.
func NewBaseConfigurator[N Number](f *shared.Field[N], optional bool) *BaseConfigurator[N] {
	return &BaseConfigurator[N]{f: f, optional: optional}
}

// Of converts a numeric field to the number configurator: b.Field(&m.Age).As(num.Of).
func Of[N Number](f *shared.Field[N]) *BaseConfigurator[N] {
	return NewBaseConfigurator(f, false)
}

// Ptr converts a pointer to numeric field to the number configurator: b.Field(&m.Discount).As(num.Ptr).
// A nil pointer is absent, Required checks only that the pointer is not nil, zero is allowed.
func Ptr[N Number](f *shared.Field[*N]) *BaseConfigurator[N] {
	return NewBaseConfigurator(shared.Deref(f), true)
}

// Required checks if the number is not zero, for pointer fields checks if the pointer is not nil.
func (i *BaseConfigurator[N]) Required() *BaseConfigurator[N] {
	optional := i.optional
	i.f.Required(func(v N) bool { return optional || v != 0 }, requiredLocaleKey)
	return i
}

// Max checks if the number is not greater than maxNum.
func (i *BaseConfigurator[N]) Max(maxNum N) *BaseConfigurator[N] {
	i.f.Check(func(v N) bool { return v <= maxNum }, maxLocaleKey, maxNum)
	return i
}

// Min checks if the number is not less than minNum.
func (i *BaseConfigurator[N]) Min(minNum N) *BaseConfigurator[N] {
	i.f.Check(func(v N) bool { return v >= minNum }, minLocaleKey, minNum)
	return i
}

// MinField checks if the number is not less than the value of another numeric field of the object:
// b.Field(&m.Price).As(num.Of).MinField(&m.MinPrice).
func (i *BaseConfigurator[N]) MinField(other *N) *BaseConfigurator[N] {
	i.f.CheckWith(other, func(v, o N) bool { return v >= o }, minLocaleKey)
	return i
}

// MaxField checks if the number is not greater than the value of another numeric field of the object:
// b.Field(&m.Used).As(num.Of).MaxField(&m.Limit).
func (i *BaseConfigurator[N]) MaxField(other *N) *BaseConfigurator[N] {
	i.f.CheckWith(other, func(v, o N) bool { return v <= o }, maxLocaleKey)
	return i
}

// AnyOf checks if the number is one of the allowed values.
func (i *BaseConfigurator[N]) AnyOf(allowed ...N) *BaseConfigurator[N] {
	i.f.Check(func(v N) bool { return slices.Contains(allowed, v) }, anyOfLocaleKey, allowed)
	return i
}

// AnyOfInterval checks if the number is within [begin, end] interval, both ends inclusive.
func (i *BaseConfigurator[N]) AnyOfInterval(begin, end N) *BaseConfigurator[N] {
	i.f.Check(func(v N) bool { return v >= begin && v <= end }, anyOfIntervalLocalKey, begin, end)
	return i
}

// Custom allows for custom validation logic to be applied to the number.
func (i *BaseConfigurator[N]) Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value N) []shared.Error) *BaseConfigurator[N] {
	i.f.Custom(f)
	return i
}

// When allows for conditional validation logic to be applied to the number.
func (i *BaseConfigurator[N]) When(whenFn func(ctx context.Context, value N) bool) *BaseConfigurator[N] {
	return NewBaseConfigurator(i.f.When(whenFn), i.optional)
}
