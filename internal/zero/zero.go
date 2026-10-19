// Package zero allocates the pointer fields of the models of Configure, so the rules are configured by the pointers
// to the fields of the nested structs.
package zero

import (
	"fmt"
	"reflect"
	"slices"
)

// Allocate allocates an input structure such that all pointer fields
// are fully allocated, i.e. rather than having a nil value,
// the pointer contains a pointer to an initialized value,
// e.g. an *int field will be a pointer to 0 instead of a nil pointer.
//
// Allocate does not allocate private fields.
func Allocate(obj interface{}) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Pointer || val.IsNil() {
		return fmt.Errorf("input interface is not addressable (can't Set the memory address): %#v",
			obj)
	}
	if val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("allocate.Zero currently only works with [pointers to] structs, not type %v",
			val.Elem().Kind())
	}
	zeroStruct(val.Elem(), nil)
	return nil
}

// zeroStruct allocates exported pointer and map fields of the struct recursively.
// path holds the types of the structs being allocated, pointers to them are left nil,
// so recursive types like struct{ Next *Node } are allocated once.
func zeroStruct(val reflect.Value, path []reflect.Type) {
	path = append(path, val.Type())
	for i := 0; i < val.NumField(); i++ {
		if !val.Type().Field(i).IsExported() {
			continue
		}
		field := val.Field(i)

		// pre-allocate pointer fields
		if field.Kind() == reflect.Pointer && field.IsNil() {
			if slices.Contains(path, field.Type().Elem()) {
				continue
			}
			field.Set(reflect.New(field.Type().Elem()))
		}

		indirectField := reflect.Indirect(field)
		switch indirectField.Kind() {
		case reflect.Map:
			indirectField.Set(reflect.MakeMap(indirectField.Type()))
		case reflect.Struct:
			// recursively allocate each of the structs embedded fields
			zeroStruct(indirectField, path)
		}
	}
}

// Must is Allocate that panics on an error.
func Must(obj interface{}) {
	err := Allocate(obj)
	if err != nil {
		panic(err)
	}
}
