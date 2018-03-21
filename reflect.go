package pantry

import (
	"errors"
	"reflect"

	ref "github.com/yazver/golibs/reflect"
)

var (
	errValueIsNil         = errors.New("Invalid value, should be not nil")
	errValueIsNotPointer  = errors.New("Invalid value, should be pointer")
	errValueIsNotStruct   = errors.New("Invalid value, should be struct")
	errValueIsNotEditable = errors.New("Invalid value, should be editable")
)

type processField func(value reflect.Value, name string, tag reflect.StructTag) error

func traverseStruct(v interface{}, process processField) error {
	value := reflect.ValueOf(v)
	if value.IsNil() {
		return errValueIsNil
	}
	if value.Kind() != reflect.Ptr {
		return errValueIsNotPointer
	}
	value = reflect.Indirect(value)
	if value.Kind() != reflect.Struct {
		return errValueIsNotStruct
	}
	if !value.CanSet() {
		return errValueIsNotEditable
	}
	return traverseStructValue(value, process)
}

func traverseStructValue(v reflect.Value, process processField) error {
	return ref.TraverseValueFields(v, func(value reflect.Value, path string, level uint, field *reflect.StructField) error {
		return process(value, field.Name, field.Tag)
	})
}
