package num

import (
	"reflect"
	"unsafe"

	"github.com/insei/fmap/v3"
)

type Struct struct {
	Slice       []string
	SlicePtr    []*string
	PtrSlice    *[]string
	PtrSlicePtr *[]*string
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
		panic("Are you sure about that?")
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

var str = "999"
var str1 = "9999"
var testStruct = &Struct{
	Slice:       []string{"999"},
	SlicePtr:    []*string{&str},
	PtrSlice:    &[]string{"999"},
	PtrSlicePtr: &[]*string{&str1},
}

//
//import (
//	"context"
//	"fmt"
//	"testing"
//
//	"github.com/insei/fmap/v3"
//
//	"github.com/insei/valigo/shared"
//)
//
//type helperIntImpl struct{}
//
//func (h *helperIntImpl) ErrorT(ctx context.Context, field fmap.Field, value any, localeKey string, args ...any) shared.Error {
//	return shared.Error{
//		Message: fmt.Sprintf(localeKey, value),
//	}
//}
//
//type user struct {
//	Age       int
//	Height    int
//	AgePtr    *int
//	HeightPtr *int
//}
//
//func TestIntConfiguratorMax(t *testing.T) {
//	testUser := user{
//		Age:    40,
//		Height: 185,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		max           int
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "Age maxT check",
//			fieldName:     "Age",
//			max:           35,
//			value:         &testUser.Age,
//			expectedError: 1,
//		},
//		{
//			name:          "Height maxT check",
//			fieldName:     "Height",
//			max:           190,
//			value:         &testUser.Height,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewValueConfigurator[int](field, appendFn, &helper)
//			configurator.Max(tc.max)
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntPtrConfiguratorMax(t *testing.T) {
//	age := 40
//	height := 185
//	testUser := user{
//		AgePtr:    &age,
//		HeightPtr: &height,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		max           int
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "AgePtr maxT check",
//			fieldName:     "AgePtr",
//			max:           35,
//			value:         &testUser.AgePtr,
//			expectedError: 1,
//		},
//		{
//			name:          "HeightPtr maxT check",
//			fieldName:     "HeightPtr",
//			max:           190,
//			value:         &testUser.HeightPtr,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewPtrConfigurator[int](field, appendFn, &helper)
//			configurator.Max(tc.max)
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntConfiguratorMin(t *testing.T) {
//	testUser := user{
//		Age:    20,
//		Height: 185,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		min           int
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "Age min check",
//			fieldName:     "Age",
//			min:           25,
//			value:         &testUser.Age,
//			expectedError: 1,
//		},
//		{
//			name:          "Height min check",
//			fieldName:     "Height",
//			min:           165,
//			value:         &testUser.Height,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewValueConfigurator[int](field, appendFn, &helper)
//			configurator.Min(tc.min)
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntPtrConfiguratorMin(t *testing.T) {
//	age := 20
//	height := 185
//	testUser := user{
//		AgePtr:    &age,
//		HeightPtr: &height,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		min           int
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "AgePtr min check",
//			fieldName:     "AgePtr",
//			min:           25,
//			value:         &testUser.AgePtr,
//			expectedError: 1,
//		},
//		{
//			name:          "HeightPtr min check",
//			fieldName:     "HeightPtr",
//			min:           165,
//			value:         &testUser.HeightPtr,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewPtrConfigurator[int](field, appendFn, &helper)
//			configurator.Min(tc.min)
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntConfiguratorRequired(t *testing.T) {
//	testUser := user{
//		Age: 20,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "Height required check",
//			fieldName:     "Height",
//			value:         &testUser.Height,
//			expectedError: 0,
//		},
//		{
//			name:          "Age required check",
//			fieldName:     "Age",
//			value:         &testUser.Age,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewValueConfigurator[int](field, appendFn, &helper)
//			configurator.Required()
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntPtrConfiguratorRequired(t *testing.T) {
//	age := 20
//	testUser := user{
//		AgePtr: &age,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "HeightPtr required check",
//			fieldName:     "HeightPtr",
//			value:         &testUser.HeightPtr,
//			expectedError: 1,
//		},
//		{
//			name:          "AgePtr required check",
//			fieldName:     "AgePtr",
//			value:         &testUser.AgePtr,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewPtrConfigurator[int](field, appendFn, &helper)
//			configurator.Required()
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntConfiguratorAnyOf(t *testing.T) {
//	testUser := user{
//		Age:    18,
//		Height: 185,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		allowed       []int
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "Age any of check",
//			fieldName:     "Age",
//			allowed:       []int{20, 30, 40},
//			value:         &testUser.Age,
//			expectedError: 1,
//		},
//		{
//			name:          "Height any of check",
//			fieldName:     "Height",
//			allowed:       []int{180, 185, 190},
//			value:         &testUser.Height,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewValueConfigurator[int](field, appendFn, &helper)
//			configurator.AnyOf(tc.allowed...)
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntPtrConfiguratorAnyOf(t *testing.T) {
//	age := 18
//	height := 185
//	testUser := user{
//		AgePtr:    &age,
//		HeightPtr: &height,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		allowed       []int
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "AgePtr any of check",
//			fieldName:     "AgePtr",
//			allowed:       []int{20, 30, 40},
//			value:         &testUser.AgePtr,
//			expectedError: 1,
//		},
//		{
//			name:          "HeightPtr any of check",
//			fieldName:     "HeightPtr",
//			allowed:       []int{180, 185, 190},
//			value:         &testUser.HeightPtr,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewPtrConfigurator[int](field, appendFn, &helper)
//			configurator.AnyOf(tc.allowed...)
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntConfiguratorAnyOfInterval(t *testing.T) {
//	testUser := user{
//		Age:    18,
//		Height: 185,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		begin         int
//		end           int
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "Age any of interval check",
//			fieldName:     "Age",
//			begin:         20,
//			end:           50,
//			value:         &testUser.Age,
//			expectedError: 1,
//		},
//		{
//			name:          "Height any of interval  check",
//			fieldName:     "Height",
//			begin:         160,
//			end:           190,
//			value:         &testUser.Height,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewValueConfigurator[int](field, appendFn, &helper)
//			configurator.AnyOfInterval(tc.begin, tc.end)
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntPtrConfiguratorAnyOfInterval(t *testing.T) {
//	age := 18
//	height := 185
//	testUser := user{
//		AgePtr:    &age,
//		HeightPtr: &height,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		begin         int
//		end           int
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "AgePtr any of interval check",
//			fieldName:     "AgePtr",
//			begin:         20,
//			end:           50,
//			value:         &testUser.AgePtr,
//			expectedError: 1,
//		},
//		{
//			name:          "HeightPtr any of interval  check",
//			fieldName:     "HeightPtr",
//			begin:         160,
//			end:           190,
//			value:         &testUser.HeightPtr,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewPtrConfigurator[int](field, appendFn, &helper)
//			configurator.AnyOfInterval(tc.begin, tc.end)
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntConfiguratorCustom(t *testing.T) {
//	testUser := user{
//		Age: 18,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "Height custom function check",
//			fieldName:     "Height",
//			value:         &testUser.Height,
//			expectedError: 1,
//		},
//		{
//			name:          "Age custom function check",
//			fieldName:     "Age",
//			value:         &testUser.Age,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewValueConfigurator[int](field, appendFn, &helper)
//			configurator.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value *int) []shared.Error {
//				if value == nil || *value == 0 {
//					return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
//				}
//				return nil
//			})
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
//
//func TestIntPtrConfiguratorCustom(t *testing.T) {
//	age := 18
//	testUser := user{
//		AgePtr: &age,
//	}
//	storage, _ := fmap.GetFrom(testUser)
//	helper := helperIntImpl{}
//
//	testCases := []struct {
//		name          string
//		fieldName     string
//		value         any
//		expectedError int
//	}{
//		{
//			name:          "HeightPtr custom function check",
//			fieldName:     "HeightPtr",
//			value:         &testUser.HeightPtr,
//			expectedError: 1,
//		},
//		{
//			name:          "AgePtr custom function check",
//			fieldName:     "AgePtr",
//			value:         &testUser.AgePtr,
//			expectedError: 0,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			var errs []shared.Error
//			field := storage.MustFind(tc.fieldName)
//			appendFn := func(field fmap.Field, fn shared.FieldValidationFn) {
//				errs = fn(context.Background(), &helper, tc.value)
//			}
//			configurator := NewPtrConfigurator[int](field, appendFn, &helper)
//			configurator.Custom(func(ctx context.Context, h *shared.FieldCustomHelper, value **int) []shared.Error {
//				if value == nil || *value == nil || **value == 0 {
//					return []shared.Error{h.ErrorT(ctx, *value, requiredLocaleKey)}
//				}
//				return nil
//			})
//			if len(errs) != tc.expectedError {
//				t.Errorf("expected %v, got %v", tc.expectedError, len(errs))
//			}
//		})
//	}
//}
