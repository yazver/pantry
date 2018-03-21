package reflect

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var (
	errValueIsNil         = errors.New("Invalid value, should be not nil")
	errValueIsNotPointer  = errors.New("Invalid value, should be pointer")
	errValueIsNotStruct   = errors.New("Invalid value, should be struct")
	errValueIsNotEditable = errors.New("Invalid value, should be editable")
)

func getMarshaler(v reflect.Value) encoding.TextMarshaler {
	if v.CanInterface() {
		if m, ok := v.Interface().(encoding.TextMarshaler); ok {
			return m
		}
	}
	return nil
}

func findMarshaler(v reflect.Value) encoding.TextMarshaler {
	if m := getMarshaler(v); m != nil {
		return m
	}
	if v.CanAddr() {
		if m := getMarshaler(v.Addr()); m != nil {
			return m
		}
	}
	if m := getMarshaler(reflect.Indirect(v)); m != nil {
		return m
	}
	return nil
}

func getUnmarshaler(v reflect.Value) encoding.TextUnmarshaler {
	if v.CanInterface() {
		if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
			return u
		}
	}
	return nil
}

func findUnmarshaler(v reflect.Value) encoding.TextUnmarshaler {
	if u := getUnmarshaler(v); u != nil {
		return u
	}
	if v.CanAddr() {
		if u := getUnmarshaler(v.Addr()); u != nil {
			return u
		}
	}
	if u := getUnmarshaler(reflect.Indirect(v)); u != nil {
		return u
	}
	return nil
}

func AssignStringToValue(dst reflect.Value, src string) (err error) {
	if u := findUnmarshaler(dst); u != nil {
		return u.UnmarshalText([]byte(src))
	}
	dst = reflect.Indirect(dst)
	if !dst.CanSet() {
		return errors.New("Value is not assignable")
	}
	if dst.CanInterface() {
		if _, ok := dst.Interface().(time.Duration); ok {
			duration, err := time.ParseDuration(src)
			if err == nil {
				dst.Set(reflect.ValueOf(duration))
			}
			return err
		}
	}
	switch dst.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(src, 0, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui, err := strconv.ParseUint(src, 0, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetUint(ui)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(src, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(src)
		if err != nil {
			return err
		}
		dst.SetBool(b)
	case reflect.String:
		dst.SetString(src)
	default:
		err = fmt.Errorf("Unable to convert string \"%s\" to type \"%s\"", src, dst.Type().Name())
	}
	return
}

func AssignValue(dst, src reflect.Value) (err error) {
	dst = reflect.Indirect(dst)
	if !dst.CanSet() {
		return errValueIsNotEditable
	}
	src = reflect.Indirect(src)

	defer func() {
		if e := recover(); e != nil {
			switch e := e.(type) {
			case error:
				err = e
			case string:
				err = errors.New(e)
			default:
				panic(e)
			}
		}
	}()
	if src.Kind() == reflect.String {
		return AssignStringToValue(dst, src.String())
	}
	value := src.Convert(dst.Type())
	dst.Set(value)
	return nil
}

type ProcessValue func(value reflect.Value, path string, level uint, field *reflect.StructField) error

func Traverse(v interface{}, process ProcessValue) error {
	return traverseValue(reflect.ValueOf(v), "", 0, nil, process)
}

func TraverseValue(v reflect.Value, process ProcessValue) error {
	return traverseValue(v, "", 0, nil, process)
}

func traverseValue(v reflect.Value, path string, depth uint, field *reflect.StructField, process ProcessValue) error {
	v = reflect.Indirect(v)
	if err := process(v, path, depth, field); err != nil {
		return err
	}
	depth++

	switch v.Kind() {
	case reflect.Struct:
		structType := v.Type()
		for i := 0; i < structType.NumField(); i++ {
			structField := structType.Field(i)
			fieldValue := v.Field(i)
			if err := traverseValue(fieldValue, path+"."+structField.Name, depth, &structField, process); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		length := v.Len()
		for i := 0; i < length; i++ {
			if err := traverseValue(v.Index(i), path+"["+strconv.Itoa(i)+"]", depth, nil, process); err != nil {
				return err
			}
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			keyStr := "[]"
			if key.CanInterface() {
				keyStr = fmt.Sprintf("[%v]", key.Interface())
			}
			if err := traverseValue(v.MapIndex(key), path+keyStr, depth, nil, process); err != nil {
				return err
			}
		}
	}

	return nil
}

func TraverseFields(v interface{}, processField ProcessValue) error {
	return TraverseValueFields(reflect.ValueOf(v), processField)
}

func TraverseValueFields(v reflect.Value, processField ProcessValue) error {
	process := func(value reflect.Value, path string, level uint, field *reflect.StructField) error {
		if field != nil {
			return processField(value, path, level, field)
		}
		return nil
	}
	return traverseValue(v, "", 0, nil, process)
}

func Clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}
