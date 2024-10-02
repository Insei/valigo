package num

import (
	"context"

	"github.com/insei/valigo/shared"
)

type numbers interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float64 | float32
}

type ptrNumbers interface {
	*int | *int8 | *int16 | *int32 | *int64 | *uint | *uint8 | *uint16 | *uint32 | *uint64 | *float64 | *float32
}

type BaseConfigurator interface {
	// Required checks if the integer is not empty.
	Required() BaseConfigurator

	// AnyOf checks if the integer is one of the allowed values.
	AnyOf(allowed ...any) BaseConfigurator

	// AnyOfInterval checks if the integer is one of the allowed values intervals.
	AnyOfInterval(begin, end any) BaseConfigurator

	// Max checks if the integer is not greater than the given maximum number.
	Max(any) BaseConfigurator

	// Min checks if the integer is not less than the given minimum number.
	Min(any) BaseConfigurator

	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value any) []shared.Error) BaseConfigurator
	When(whenFn func(ctx context.Context, value any) bool) BaseConfigurator
}

//// ValueConfigurator is a builder interface for number fields.
//// It provides methods for adding validation rules to an number field.
//type ValueConfigurator[T numbers] interface {
//	BaseConfigurator[T, ValueConfigurator[T]]
//
//	//Custom allows for custom validation logic.
//	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) ValueConfigurator[T]
//
//	// When allows for conditional validation based on a given condition.
//	When(f func(ctx context.Context, value *T) bool) ValueConfigurator[T]
//}

//// PtrConfigurator is a builder interface for ptr number fields.
//// It provides methods for adding validation rules to a number field.
//type PtrConfigurator[T numbers] interface {
//	BaseConfigurator[T, PtrConfigurator[T]]
//
//	//Custom allows for custom validation logic.
//	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value **T) []shared.Error) PtrConfigurator[T]
//
//	// When allows for conditional validation based on a given condition.
//	When(f func(ctx context.Context, value **T) bool) PtrConfigurator[T]
//}

// IntSliceBuilder is a builder interface for integer slice fields.
// It provides methods for adding validation rules to an integer slice field.
type IntSliceBuilder[T []int | *[]int | []int8 | *[]int8 | []int16 | *[]int16 | []int32 | *[]int32 | []int64 | *[]int64 |
	[]uint | *[]uint | []uint8 | *[]uint8 | []uint16 | *[]uint16 | []uint32 | *[]uint32 | []uint64 | *[]uint64] interface {
	// Required checks if the slice is not empty.
	Required() IntSliceBuilder[T]

	// Custom allows for custom validation logic.
	Custom(f func(ctx context.Context, h *shared.FieldCustomHelper, value *T) []shared.Error) IntSliceBuilder[T]

	// Max checks if the of each integer in the slice is not greater than the given maximum number.
	Max(int) IntSliceBuilder[T]

	// Min checks if the of each integer in the slice is not less than the given minimum number.
	Min(int) IntSliceBuilder[T]

	// When allows for conditional validation based on a given condition.
	When(f func(ctx context.Context, value *T) bool) IntSliceBuilder[T]
}

// BundleConfigurator is a builder interface for a bundle of uint fields.
// It provides methods for adding integer fields to the bundle.
type BundleConfigurator interface {
	Number(fieldPtr any) BaseConfigurator
	//NumberSlice(fieldPtr any) *shared.SliceFieldConfigurator
	//// Int returns an int field validation configurator.
	//Int(field *int) ValueConfigurator[int]
	//
	//// IntPtr returns a pointer to an int field validation configurator.
	//IntPtr(field **int) PtrConfigurator[int]
	//
	//// Int8 returns an int8 field validation configurator.
	//Int8(field *int8) ValueConfigurator[int8]
	//
	//// Int8Ptr returns a pointer to an int8 field validation configurator.
	//Int8Ptr(field **int8) PtrConfigurator[int8]

	//// IntSlice adds an integer slice field to the bundle.
	//IntSlice(field *[]int) IntSliceBuilder[[]int]
	//
	//// IntSlicePtr adds a pointer to an integer slice field to the bundle.
	//IntSlicePtr(field **[]int) IntSliceBuilder[*[]int]
}
