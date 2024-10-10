package valigo

import (
	"context"

	"github.com/insei/valigo/num"
	"github.com/insei/valigo/shared"
	"github.com/insei/valigo/str"
	"github.com/insei/valigo/uuid"
)

// NumberBuilder is an interface that defines methods for building validators for numeric types.
type NumberBuilder[T ~int | int8 | int16 | int32 | int64 | *int | *int8 | *int16 | *int32 | *int64 |
	uint | uint8 | uint16 | uint32 | uint64 | *uint | *uint8 | *uint16 | *uint32 | *uint64 |
	float64 | float32 | *float64 | *float32] interface {
	// Max sets a maximum value for the validator.
	Max(T) NumberBuilder[T]
	// Min sets a minimum value for the validator.
	Min(T) NumberBuilder[T]
	// Custom adds a custom validation function to the validator.
	Custom(func(ctx context.Context, h *shared.Helper, value *T) []error) NumberBuilder[T]
	// When sets a condition for when the validator should be applied.
	When(func(ctx context.Context, value *T) bool) NumberBuilder[T]
}

// NumbersBundleBuilder is an interface that defines methods for building validators for numeric types.
type NumbersBundleBuilder interface {
	// Int returns a NumberBuilder for validating an int field.
	Int(field *int) NumberBuilder[int]
	// Int8 returns a NumberBuilder for validating an int8 field.
	Int8(field *int8) NumberBuilder[int8]
	// Int16 returns a NumberBuilder for validating an int16 field.
	Int16(field *int16) NumberBuilder[int16]
	// Int32 returns a NumberBuilder for validating an int32 field.
	Int32(field *int32) NumberBuilder[int32]
	// Int64 returns a NumberBuilder for validating an int64 field.
	Int64(field *int64) NumberBuilder[int64]
	// IntPtr returns a NumberBuilder for validating a pointer to an int field.
	IntPtr(field **int) NumberBuilder[*int]
	// Int8Ptr returns a NumberBuilder for validating a pointer to an int8 field.
	Int8Ptr(field **int8) NumberBuilder[*int8]
	// Int16Ptr returns a NumberBuilder for validating a pointer to an int16 field.
	Int16Ptr(field **int16) NumberBuilder[*int16]
	// Int32Ptr returns a NumberBuilder for validating a pointer to an int32 field.
	Int32Ptr(field **int32) NumberBuilder[*int32]
	// Int64Ptr returns a NumberBuilder for validating a pointer to an int64 field.
	Int64Ptr(field **int64) NumberBuilder[*int64]
	// Uint returns a NumberBuilder for validating an uint field.
	Uint(field *uint) NumberBuilder[uint]
	// Uint8 returns a NumberBuilder for validating an uint8 field.
	Uint8(field *uint8) NumberBuilder[uint8]
	// Uint16 returns a NumberBuilder for validating an uint16 field.
	Uint16(field *uint16) NumberBuilder[uint16]
	// Uint32 returns a NumberBuilder for validating an uint32 field.
	Uint32(field *uint32) NumberBuilder[uint32]
	// Uint64 returns a NumberBuilder for validating an uint64 field.
	Uint64(field *uint64) NumberBuilder[uint64]
	// UintPtr returns a NumberBuilder for validating a pointer to an uint field.
	UintPtr(field **uint) NumberBuilder[*uint]
	// Uint8Ptr returns a NumberBuilder for validating a pointer to an uint8 field.
	Uint8Ptr(field **uint8) NumberBuilder[*uint8]
	// Uint16Ptr returns a NumberBuilder for validating a pointer to an uint16 field.
	Uint16Ptr(field **uint16) NumberBuilder[*uint16]
	// Uint32Ptr returns a NumberBuilder for validating a pointer to an uint32 field.
	Uint32Ptr(field **uint32) NumberBuilder[*uint32]
	// Uint64Ptr returns a NumberBuilder for validating a pointer to an uint64 field.
	Uint64Ptr(field **uint64) NumberBuilder[*uint64]
}

// StringSliceBuilder is an interface that defines methods for building validators for string slices.
type StringSliceBuilder[T string | *string] interface {
	// Trim removes leading and trailing whitespace from each string in the slice.
	Trim() StringSliceBuilder[T]
	// Max sets a maximum length for each string in the slice.
	Max(uint) StringSliceBuilder[T]
	// Min sets a minimum length for each string in the slice.
	Min(uint) StringSliceBuilder[T]
	// Unique ensures that all strings in the slice are unique.
	Unique() StringSliceBuilder[T]
}

// SlicesBundleBuilder is an interface that defines methods for building validators for string slice types.
type SlicesBundleBuilder interface {
	// SliceStrings returns a StringSliceBuilder for validating a string slice.
	SliceStrings(value *[]string) StringSliceBuilder[string]
	// SlicePtrStrings returns a StringSliceBuilder for validating a pointer to a string slice.
	SlicePtrStrings(value **[]string) StringSliceBuilder[string]
	// SliceStringsPtr returns a StringSliceBuilder for validating a string slice of pointers.
	SliceStringsPtr(value *[]*string) StringSliceBuilder[*string]
	// SlicePtrStringsPtr returns a StringSliceBuilder for validating a pointer to a string slice of pointers.
	SlicePtrStringsPtr(value **[]*string) StringSliceBuilder[*string]
}

// Configurator is an interface that defines methods for building validators for any type T.
type Configurator[T any] interface {
	str.StringBundleConfigurator
	num.NumberBundleConfigurator
	uuid.UUIDBundleConfigurator
	Slice(sliceFieldPtr any) *shared.SliceFieldConfigurator
	// When sets a condition for when the validator should be applied.
	When(func(ctx context.Context, obj *T) bool) Configurator[T]
	// Custom adds a custom validation function to the validator.
	Custom(fn func(ctx context.Context, h shared.StructCustomHelper, obj *T) []shared.Error)
	//Custom(func(obj *T) []error) FieldConfigurator[T]
	//SlicesBundleBuilder
}
