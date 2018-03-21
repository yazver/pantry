package pantry

import (
	"fmt"
	"reflect"
	"strings"

	ref "github.com/yazver/golibs/reflect"
)

func parseTagSettings(tag reflect.StructTag, configTagName string) map[string]string {
	settings := map[string]string{}
	if str, ok := tag.Lookup(configTagName); ok {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToLower(v[0]))
			if len(v) >= 2 {
				settings[k] = strings.Join(v[1:], ":")
			} else {
				settings[k] = ""
			}
		}
	}
	for _, key := range []string{"flag", "env"} {
		if str, ok := tag.Lookup(configTagName + "." + key); ok {
			settings[key] = str
		}
	}
	return settings
}

func parseFlagSettings(str string) (name, usage string) {
	parts := strings.Split(str, "|")
	name = strings.TrimSpace(parts[0])
	if len(parts) >= 2 {
		usage = strings.Join(parts[1:], "|")
	}
	return
}

func processDefaultValues(v interface{}, options *Options) error {
	if !options.Tags.Default.Use {
		return nil
	}
	initValues := func(value reflect.Value, name string, tag reflect.StructTag) error {
		if defaultValue, ok := tag.Lookup(options.Tags.Default.Name); ok {
			if err := ref.AssignStringToValue(value, defaultValue); err != nil {
				return fmt.Errorf("Unable to assign default value \"%s\" to field \"%s\": %s", defaultValue, name, err)
			}
		}
		return nil
	}
	return traverseStruct(v, initValues)
}

func processTags(v interface{}, options *Options) error {
	flags := &(options.Flags)
	env := &(options.Enviropment)
	if options.Tags.Config.Use && (flags.Using != FlagsDontUse || options.Enviropment.Use) {
		initFlags := func(value reflect.Value, name string, tag reflect.StructTag) error {
			settings := parseTagSettings(tag, options.Tags.Config.Name)
			if flagSettings, ok := settings["flag"]; ok {
				flag, usage := parseFlagSettings(flagSettings)
				if options.Tags.Description.Use {
					if description, ok := tag.Lookup(options.Tags.Description.Name); ok {
						usage = description
					}
				}
				if err := flags.Add(value, flag, usage); err != nil {
					return err
				}
			}
			return nil
		}
		if flags.Using == FlagsUseAll {
			if err := traverseStruct(v, initFlags); err != nil {
				return err
			}
		}

		initValues := func(value reflect.Value, name string, tag reflect.StructTag) error {
			settings := parseTagSettings(tag, options.Tags.Config.Name)
			if envVarName, ok := settings["env"]; ok {
				if err := env.Get(value, envVarName); err != nil {
					return err
				}
			}
			if flagSettings, ok := settings["flag"]; ok {
				flag, _ := parseFlagSettings(flagSettings)
				if err := flags.Get(value, flag); err != nil {
					return err
				}
			}
			return nil
		}
		return traverseStruct(v, initValues)
	}
	return nil
}
