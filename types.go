package valigo

import (
	"context"
)

type Translator interface {
	ErrorT(ctx context.Context, format string, args ...any) error
	T(ctx context.Context, format string, args ...any) string
}
type StringBuilder[T string | *string] interface {
	Trim() StringBuilder[T]
	Required() StringBuilder[T]
	AnyOf(vals ...string) StringBuilder[T]
	Custom(f func(ctx context.Context, h *Helper, value *T) []error) StringBuilder[T]
	When(f func(ctx context.Context, value *T) bool) StringBuilder[T]
}

type NumberBuilder[T ~int | int8 | int16 | int32 | int64 | *int | *int8 | *int16 | *int32 | *int64 |
	uint | uint8 | uint16 | uint32 | uint64 | *uint | *uint8 | *uint16 | *uint32 | *uint64 |
	float64 | float32 | *float64 | *float32] interface {
	Max(T) NumberBuilder[T]
	Min(T) NumberBuilder[T]
	Custom(func(h *Helper, value *T) []error) NumberBuilder[T]
	When(func(value *T) bool) NumberBuilder[T]
}

type NumbersBundleBuilder interface {
	Int(field *int) NumberBuilder[int]
	Int8(field *int8) NumberBuilder[int8]
	Int16(field *int16) NumberBuilder[int16]
	Int32(field *int32) NumberBuilder[int32]
	Int64(field *int64) NumberBuilder[int64]
	IntPtr(field **int) NumberBuilder[*int]
	Int8Ptr(field **int8) NumberBuilder[*int8]
	Int16Ptr(field **int16) NumberBuilder[*int16]
	Int32Ptr(field **int32) NumberBuilder[*int32]
	Int64Ptr(field **int64) NumberBuilder[*int64]
	Uint(field *uint) NumberBuilder[uint]
	Uint8(field *uint8) NumberBuilder[uint8]
	Uint16(field *uint16) NumberBuilder[uint16]
	Uint32(field *uint32) NumberBuilder[uint32]
	Uint64(field *uint64) NumberBuilder[uint64]
	UintPtr(field **uint) NumberBuilder[*uint]
	Uint8Ptr(field **uint8) NumberBuilder[*uint8]
	Uint16Ptr(field **uint16) NumberBuilder[*uint16]
	Uint32Ptr(field **uint32) NumberBuilder[*uint32]
	Uint64Ptr(field **uint64) NumberBuilder[*uint64]
}

type StringsBundleBuilder interface {
	String(field *string) StringBuilder[string]
	StringPtr(field **string) StringBuilder[*string]
}

type StringSliceBuilder[T string | *string] interface {
	Trim() StringSliceBuilder[T]
	Max(uint) StringSliceBuilder[T]
	Min(uint) StringSliceBuilder[T]
	Unique() StringSliceBuilder[T]
}

type SlicesBundleBuilder interface {
	SliceStrings(value *[]string) StringSliceBuilder[string]
	SlicePtrStrings(value **[]string) StringSliceBuilder[string]
	SliceStringsPtr(value *[]*string) StringSliceBuilder[*string]
	SlicePtrStringsPtr(value **[]*string) StringSliceBuilder[*string]
}

type Builder[T any] interface {
	//NumbersBundleBuilder
	StringsBundleBuilder
	When(func(obj *T) bool) Builder[T]
	//Custom(func(obj *T) []error) Builder[T]
	//SlicesBundleBuilder
}
