package valigo

import (
	"reflect"

	"github.com/insei/fmap/v3"
)

type validatorFn func(h Helper, obj, v any) []error

type cached struct {
	field fmap.Field
	fn    validatorFn
}

type storage struct {
	data  map[reflect.Type]map[fmap.Field][]validatorFn
	cache map[reflect.Type][]cached
}

func (s *storage) toCache(objTypeOf reflect.Type) {
	for field, validators := range s.get(objTypeOf) {
		for _, validator := range validators {
			s.cache[objTypeOf] = append(s.cache[objTypeOf], cached{field: field, fn: validator})
		}
	}
	s.data = make(map[reflect.Type]map[fmap.Field][]validatorFn)
}

func (s *storage) getCache(objTypeOf reflect.Type) ([]cached, bool) {
	cache, ok := s.cache[objTypeOf]
	return cache, ok
}

func (s *storage) new(temp any, enabler func(any) bool) func(field fmap.Field, fn func(h Helper, v any) []error) {
	t := reflect.TypeOf(temp)
	return func(field fmap.Field, fn func(h Helper, v any) []error) {
		fnNew := func(h Helper, obj, v any) []error {
			if enabler != nil && !enabler(obj) {
				return nil
			}
			return fn(h, v)
		}
		_, ok := s.data[t]
		if !ok {
			s.data[t] = make(map[fmap.Field][]validatorFn)
		}
		s.data[t][field] = append(s.data[t][field], fnNew)
	}
}

func (s *storage) get(objType reflect.Type) map[fmap.Field][]validatorFn {
	return s.data[objType]
}

func newStorage() *storage {
	return &storage{data: make(map[reflect.Type]map[fmap.Field][]validatorFn), cache: make(map[reflect.Type][]cached)}
}
