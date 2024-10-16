package shared

import (
	"context"
	"reflect"
	"unsafe"

	"github.com/insei/fmap/v3"
)

func UnsafeValigoSliceCast[T any](slice []*any) []*T {
	var val any = slice
	ptr := ((*[2]unsafe.Pointer)(unsafe.Pointer(&val)))[1]
	arr := (*[]*T)(ptr)
	if arr == nil {
		return nil
	}
	return *arr
}

type SliceFieldConfigurator struct {
	field  fmap.Field
	helper Helper
	c      *FieldConfigurator[[]*any]
}

func makeGetValueSliceFn(field fmap.Field) func(value any) ([]*any, bool) {
	fType := field.GetType()
	ptrToSlice := 1 // All values that comes to validator is a pointer
	ptrToSliceElem := 0
	for fType.Kind() == reflect.Ptr {
		ptrToSlice++
		fType = fType.Elem()
	}
	if fType.Kind() != reflect.Slice {
		panic("field value is not a slice")
	}
	sliceElemType := fType.Elem()
	for sliceElemType.Kind() == reflect.Ptr {
		ptrToSliceElem++
		sliceElemType = sliceElemType.Elem()
	}
	if ptrToSlice > 2 || ptrToSliceElem > 1 {
		panic("Are you sure about that? (c)")
	}
	switch {
	case ptrToSlice == 1 && ptrToSliceElem == 0:
		return func(value any) ([]*any, bool) {
			var convertedArr []*any
			ptr := ((*[2]unsafe.Pointer)(unsafe.Pointer(&value)))[1]
			arr := (*[]any)(ptr)
			for i, _ := range *arr {
				convertedArr = append(convertedArr, &(*arr)[i])
			}
			return convertedArr, true
		}
	case ptrToSlice == 2 && ptrToSliceElem == 0:
		return func(value any) ([]*any, bool) {
			var convertedArr []*any
			ptr := ((*[2]unsafe.Pointer)(unsafe.Pointer(&value)))[1]
			arr := (**[]any)(ptr)
			if *arr == nil || **arr == nil {
				return nil, false
			}
			for i, _ := range **arr {
				convertedArr = append(convertedArr, &(**arr)[i])
			}
			return convertedArr, true
		}
	case ptrToSlice == 1 && ptrToSliceElem == 1:
		return func(value any) ([]*any, bool) {
			var convertedArr []*any
			ptr := ((*[2]unsafe.Pointer)(unsafe.Pointer(&value)))[1]
			arr := (*[]*any)(ptr)
			for _, v := range *arr {
				convertedArr = append(convertedArr, &(*v))
			}
			return convertedArr, true
		}
	case ptrToSlice == 2 && ptrToSliceElem == 1:
		return func(value any) ([]*any, bool) {
			var convertedArr []*any
			ptr := ((*[2]unsafe.Pointer)(unsafe.Pointer(&value)))[1]
			arr := (**[]*any)(ptr)
			if *arr == nil || **arr == nil {
				return nil, false
			}
			for _, v := range **arr {
				convertedArr = append(convertedArr, &(*v))
			}
			return convertedArr, true
		}
	}
	//ptr := unsafe.Pointer(&convertedArr)
	//ee := *((*[]*string)(ptr))
	//for _, v := range ee {
	//	*v = "222"
	//}
	//fmt.Println(ee)
	return nil
}

type SliceFieldConfiguratorParams struct {
	Field    fmap.Field
	Helper   Helper
	AppendFn func(fn FieldValidationFn)
}

func NewSliceFieldConfigurator(p SliceFieldConfiguratorParams) *SliceFieldConfigurator {
	getValueFn := makeGetValueSliceFn(p.Field)
	mk := NewSimpleFieldFnMaker(SimpleFieldFnMakerParams[[]*any]{
		GetValue: getValueFn,
		Field:    p.Field,
		Helper:   p.Helper,
	})
	return &SliceFieldConfigurator{
		field:  p.Field,
		helper: p.Helper,
		c: NewFieldConfigurator(FieldConfiguratorParams[[]*any]{
			Maker:    mk,
			AppendFn: p.AppendFn,
		}),
	}
}

func (s *SliceFieldConfigurator) MaxLen(maxLen int) *SliceFieldConfigurator {
	s.c.Append(func(v []*any) bool {
		if len(v) > maxLen {
			return false
		}
		return true
	}, "max len error")
	return s
}

func (s *SliceFieldConfigurator) MinLen(MinLen int) *SliceFieldConfigurator {
	s.c.Append(func(v []*any) bool {
		if len(v) < MinLen {
			return false
		}
		return true
	}, "min len error")
	return s
}

func (s *SliceFieldConfigurator) Required() *SliceFieldConfigurator {
	s.c.Append(func(v []*any) bool {
		if v == nil {
			return false
		}
		return true
	}, "required error")
	return s
}

// Custom allows for custom validation logic.
func (s *SliceFieldConfigurator) Custom(f func(ctx context.Context, h *FieldCustomHelper, value []*any) []Error) *SliceFieldConfigurator {
	customHelper := NewFieldCustomHelper(s.field, s.helper)
	s.c.appendFn(s.c.mk.CustomMake(func(ctx context.Context, h Helper, value any) []Error {
		return f(ctx, customHelper, value.([]*any))
	}))
	return s
}

// When allows for conditional validation based on a given condition.
func (s *SliceFieldConfigurator) When(whenFn func(ctx context.Context, value []*any) bool) *SliceFieldConfigurator {
	if whenFn == nil {
		return s
	}
	s.c.appendFn = func(fn FieldValidationFn) {
		fnWithEnabler := func(ctx context.Context, h Helper, v any) []Error {
			if !whenFn(ctx, v.([]*any)) {
				return nil
			}
			return fn(ctx, h, v)
		}
		s.c.appendFn(fnWithEnabler)
	}
	return s
}
