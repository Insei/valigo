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
	latitudeLocaleKey     = "validation:num:Should be latitude"
	longitudeLocaleKey    = "validation:num:Should be longitude"
	greaterLocaleKey      = "validation:num:Should be greater than %v"
	lessLocaleKey         = "validation:num:Should be less than %v"
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
	i.f.Schema(shared.KeywordRequired, true)
	optional := i.optional
	i.f.RequiredMsg(func(v N) bool { return optional || v != 0 }, shared.M("number.required", requiredLocaleKey))
	return i
}

// Max checks if the number is not greater than maxNum.
func (i *BaseConfigurator[N]) Max(maxNum N) *BaseConfigurator[N] {
	i.f.Schema("maximum", maxNum)
	i.f.CheckMsg(func(v N) bool { return v <= maxNum }, shared.M("number.max", maxLocaleKey, shared.P("max", maxNum)))
	return i
}

// Min checks if the number is not less than minNum.
func (i *BaseConfigurator[N]) Min(minNum N) *BaseConfigurator[N] {
	i.f.Schema("minimum", minNum)
	i.f.CheckMsg(func(v N) bool { return v >= minNum }, shared.M("number.min", minLocaleKey, shared.P("min", minNum)))
	return i
}

// GreaterThan checks if the number is greater than n.
func (i *BaseConfigurator[N]) GreaterThan(n N) *BaseConfigurator[N] {
	i.f.Schema("exclusiveMinimum", n)
	i.f.CheckMsg(func(v N) bool { return v > n }, shared.M("number.greater_than", greaterLocaleKey, shared.P("min", n)))
	return i
}

// LessThan checks if the number is less than n.
func (i *BaseConfigurator[N]) LessThan(n N) *BaseConfigurator[N] {
	i.f.Schema("exclusiveMaximum", n)
	i.f.CheckMsg(func(v N) bool { return v < n }, shared.M("number.less_than", lessLocaleKey, shared.P("max", n)))
	return i
}

// GreaterField checks if the number is greater than the value of another numeric field of the object:
// b.Field(&m.Max).As(num.Of).GreaterField(&m.Min).
func (i *BaseConfigurator[N]) GreaterField(other *N) *BaseConfigurator[N] {
	i.f.CheckWithMsg(other, func(v, o N) bool { return v > o }, shared.M("number.greater_field", greaterLocaleKey))
	return i
}

// LessField checks if the number is less than the value of another numeric field of the object:
// b.Field(&m.Min).As(num.Of).LessField(&m.Max).
func (i *BaseConfigurator[N]) LessField(other *N) *BaseConfigurator[N] {
	i.f.CheckWithMsg(other, func(v, o N) bool { return v < o }, shared.M("number.less_field", lessLocaleKey))
	return i
}

// MinField checks if the number is not less than the value of another numeric field of the object:
// b.Field(&m.Price).As(num.Of).MinField(&m.MinPrice).
func (i *BaseConfigurator[N]) MinField(other *N) *BaseConfigurator[N] {
	i.f.CheckWithMsg(other, func(v, o N) bool { return v >= o }, shared.M("number.min_field", minLocaleKey))
	return i
}

// MaxField checks if the number is not greater than the value of another numeric field of the object:
// b.Field(&m.Used).As(num.Of).MaxField(&m.Limit).
func (i *BaseConfigurator[N]) MaxField(other *N) *BaseConfigurator[N] {
	i.f.CheckWithMsg(other, func(v, o N) bool { return v <= o }, shared.M("number.max_field", maxLocaleKey))
	return i
}

// AnyOf checks if the number is one of the allowed values.
func (i *BaseConfigurator[N]) AnyOf(allowed ...N) *BaseConfigurator[N] {
	i.f.Schema("enum", shared.Values(allowed))
	i.f.CheckMsg(func(v N) bool { return slices.Contains(allowed, v) }, shared.M("number.any_of", anyOfLocaleKey, shared.P("values", allowed)))
	return i
}

// AnyOfInterval checks if the number is within [begin, end] interval, both ends inclusive.
func (i *BaseConfigurator[N]) AnyOfInterval(begin, end N) *BaseConfigurator[N] {
	i.f.Schema("minimum", begin).Schema("maximum", end)
	i.f.CheckMsg(func(v N) bool { return v >= begin && v <= end }, shared.M("number.interval", anyOfIntervalLocalKey, shared.P("min", begin), shared.P("max", end)))
	return i
}

// Latitude checks if the number is a latitude in degrees, within [-90, 90].
func (i *BaseConfigurator[N]) Latitude() *BaseConfigurator[N] {
	i.f.Schema("minimum", -90).Schema("maximum", 90)
	i.f.CheckMsg(func(v N) bool { return float64(v) >= -90 && float64(v) <= 90 }, shared.M("number.latitude", latitudeLocaleKey))
	return i
}

// Longitude checks if the number is a longitude in degrees, within [-180, 180].
func (i *BaseConfigurator[N]) Longitude() *BaseConfigurator[N] {
	i.f.Schema("minimum", -180).Schema("maximum", 180)
	i.f.CheckMsg(func(v N) bool { return float64(v) >= -180 && float64(v) <= 180 }, shared.M("number.longitude", longitudeLocaleKey))
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
