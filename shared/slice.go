package shared

import (
	"reflect"

	"github.com/insei/fmap/v3"
)

var sliceConvertersCache = map[reflect.Type]func(value any) (any, bool){}

func PtrSlicePtrElem[T any](value any) ([]*T, bool) {
	v, ok := value.(*[]*T)
	if !ok {
		return nil, false
	}
	return *v, true
}

func SlicePtrElem[T any](value any) ([]*T, bool) {
	v, ok := value.([]*T)
	if !ok {
		return nil, false
	}
	return v, true
}

func SliceElem[T any](value any) ([]*T, bool) {
	v, ok := value.([]T)
	if !ok {
		return nil, false
	}
	var rv []*T
	for i, _ := range v {
		rv = append(rv, &v[i])
	}
	return rv, true
}

func PtrSliceElem[T any](value any) ([]*T, bool) {
	v, ok := value.(*[]T)
	if !ok {
		return nil, false
	}
	var rv []*T
	for i, _ := range *v {
		rv = append(rv, &(*v)[i])
	}
	return rv, true
}

func addSliceMutation[T any](field fmap.Field) bool {
	var isPtrToSlice, isPtrToValue bool

	sliceType := field.GetType()
	typeOf := sliceType
	if typeOf.Kind() == reflect.Ptr {
		isPtrToSlice = true
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Slice {
		return false
	}
	typeOf = typeOf.Elem()
	if typeOf.Kind() == reflect.Ptr {
		isPtrToValue = true
	}
	switch {
	case isPtrToSlice && isPtrToValue:
		convFn := func(value any) (any, bool) {
			return PtrSlicePtrElem[T](value)
		}
		sliceConvertersCache[sliceType] = convFn
		return true
	case !isPtrToSlice && isPtrToValue:
		convFn := func(value any) (any, bool) {
			return SlicePtrElem[T](value)
		}
		sliceConvertersCache[sliceType] = convFn
		return true
	case isPtrToSlice && !isPtrToValue:
		convFn := func(value any) (any, bool) {
			return PtrSliceElem[T](value)
		}
		sliceConvertersCache[sliceType] = convFn
		return true
	case !isPtrToSlice && !isPtrToValue:
		convFn := func(value any) (any, bool) {
			return SliceElem[T](value)
		}
		sliceConvertersCache[sliceType] = convFn
		return true
	}
	return false
}

func ConvertSliceValue[T any](slice any) ([]*T, bool) {
	fn, ok := sliceConvertersCache[reflect.TypeOf(slice)]
	if !ok {
		return nil, false
	}
	s, ok := fn(slice)
	if !ok {
		return nil, false
	}
	sc, ok := s.([]*T)
	return sc, ok
}

type SliceFieldConfigurator struct {
	*FieldConfigurator[[]*any]
}

func validateSliceField(field fmap.Field) {
	typeOf := field.GetType()
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Slice {
	}
}

func NewSliceFieldConfigurator(params FieldConfiguratorParams[[]*any]) *SliceFieldConfigurator {
	fType := params.Field.GetType()
	if fType.Kind() != reflect.Slice {
		panic("wrong configuration field is not a slice")
	}
	elemType := fType.Elem()
	for elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	switch elemType.Kind() {
	case reflect.Int8:
	case reflect.Int16:

	}
	//addSliceMutation(f)
	return &SliceFieldConfigurator{
		FieldConfigurator: NewFieldConfigurator(params),
	}
}

func (s *SliceFieldConfigurator) MaxLen(maxLen int) *SliceFieldConfigurator {
	s.Append(func(v []*any) bool {
		if len(v) > maxLen {
			return false
		}
		return true
	}, "max len error")
	return s
}

func (s *SliceFieldConfigurator) MinLen(MinLen int) *SliceFieldConfigurator {
	s.Append(func(v []*any) bool {
		if len(v) < MinLen {
			return false
		}
		return true
	}, "min len error")
	return s
}

func (s *SliceFieldConfigurator) Required() *SliceFieldConfigurator {
	s.Append(func(v []*any) bool {
		if len(v) < 1 {
			return false
		}
		return true
	}, "required error")
	return s
}
