package pantry

import (
	"errors"
	"reflect"
	"strings"
)

var (
	errConfigIsNotStruct = errors.New("Invalid config, should be struct")

	decodersTag = []string{"json", "toml", "yaml"}
)

func parseTagSetting(tags reflect.StructTag) map[string]string {
	setting := map[string]string{}
	for _, str := range []string{tags.Get("plantry"), tags.Get("store")} {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if len(v) >= 2 {
				setting[k] = strings.Join(v[1:], ":")
			} else {
				setting[k] = k
			}
		}
	}
	return setting
}

func traverseStruct(v interface{}, processField func(value interface{}, settings map[string]string) error) error {
	structValue := reflect.Indirect(reflect.ValueOf(v))
	if structValue.Kind() != reflect.Struct {
		return errConfigIsNotStruct
	}

	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		value := structValue.Field(i)
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.Struct:
			if err := traverseStruct(value.Addr().Interface(), processField); err != nil {
				return err
			}
		case reflect.Slice:
			length := value.Len()
			for i := 0; i < length; i++ {
				if item := value.Index(i); reflect.Indirect(item).Kind() == reflect.Struct {
					if err := traverseStruct(item.Addr().Interface(), processField); err != nil {
						return err
					}
				}
			}
		case reflect.Map:
			for _, key := range value.MapKeys() {
				if item := value.MapIndex(key); reflect.Indirect(item).Kind() == reflect.Struct {
					if err := traverseStruct(item.Addr().Interface(), processField); err != nil {
						return err
					}
				}
			}
		default:
			settings := parseTagSetting(structField.Tag)

			if err := processField(value, settings); err != nil {
				return err
			}

		}
	}
	return nil

}
