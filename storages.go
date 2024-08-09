package valigo

import (
	"reflect"

	"github.com/insei/fmap/v3"
)

type validatorFn func(h Helper, obj, v any) []error

var storage = &storageData{data: make(map[reflect.Type]map[fmap.Field][]validatorFn)}

type storageData struct {
	data map[reflect.Type]map[fmap.Field][]validatorFn
}

func (s *storageData) new(temp any, enabler func(any) bool) func(field fmap.Field, fn func(h Helper, v any) []error) {
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
