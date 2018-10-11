package pantry

import (
	"fmt"
	"reflect"
	"strings"

	ref "github.com/yazver/golibs/reflect"
)

func addIfNotSpace(m map[string]string, k, v string) {
	k = strings.TrimSpace(strings.ToLower(k))
	v = strings.TrimSpace(v)
	if k != "" && v != "" {
		m[k] = v
	}
}

func parseTagSettings(tag reflect.StructTag, opt *Options) map[string]string {
	configTagName := opt.Tags.Config.Name
	settings := map[string]string{}
	if str, ok := tag.Lookup(configTagName); ok {
		tags := strings.Split(str, ";")
		for _, value := range tags {
			v := strings.SplitN(value, ":", 2)
			if len(v) >= 2 {
				addIfNotSpace(settings, v[0], v[1])
			}
		}
	}

	if configTagName != "" {
		configTagName = configTagName + "."
	}
	for _, key := range []string{opt.Tags.Flag.Name, opt.Tags.Env.Name} {
		if str, ok := tag.Lookup(configTagName + key); ok {
			addIfNotSpace(settings, key, str)
		}
	}
	return settings
}

func parseFlagSettings(str string) (flag *Flag) {
	flag = &Flag{Name: "-"}
	if strings.HasPrefix(str, "-") {
		return
	}
	splitNTo(str, "|", &(flag.Name), &(flag.Usage))
	fieldsTo(strings.TrimSpace(flag.Name), &(flag.Name), &(flag.Short))
	return
}

func processDefaultValues(v interface{}, opt *Options) error {
	if !opt.Tags.Default.Use {
		return nil
	}
	initValues := func(value reflect.Value, name string, tag reflect.StructTag, state *ref.State) error {
		if defaultValue, ok := tag.Lookup(opt.Tags.Default.Name); ok {
			if err := ref.AssignStringToValue(value, defaultValue); err != nil {
				return fmt.Errorf("Unable to assign default value \"%s\" to field \"%s\": %s", defaultValue, name, err)
			}
		}
		return nil
	}
	return traverseStruct(v, initValues)
}

type traverseState struct {
	Flag               string
	FlagHierarchically bool
	Env                string
	EnvHierarchically  bool
}

func processTags(v interface{}, opt *Options) error {
	flags := &(opt.Flags)
	enviropment := &(opt.Enviropment)
	if opt.Tags.Config.Use && (flags.Using != FlagsDontUse || opt.Enviropment.Use) {
		initFlags := func(value reflect.Value, name string, tag reflect.StructTag, s *ref.State) error {
			state, _ := s.Value.(traverseState)

			settings := parseTagSettings(tag, opt)
			if flagSettings, ok := settings[opt.Tags.Flag.Name]; opt.Tags.Flag.Use && state.Flag != "-" && ok {
				flag := parseFlagSettings(flagSettings)
				if flag.Name != "" && flag.Name != "-" {
					if state.FlagHierarchically {
						if state.Flag != "" {
							flag.Name = state.Flag + "-" + flag.Name
						}
						state.Flag = flag.Name
					}
					if opt.Tags.Description.Use {
						if description, ok := tag.Lookup(opt.Tags.Description.Name); ok {
							flag.Usage = description
						}
					}
					if opt.Tags.Default.Use {
						if defValue, ok := tag.Lookup(opt.Tags.Default.Name); ok {
							flag.DefValue = defValue
						}
					}
					flag.Value = value
					if err := flags.Add(flag); err != nil {
						return err
					}
				} else {
					state.Flag = "-"
					state.FlagHierarchically = false
				}
			}

			if env, ok := settings[opt.Tags.Env.Name]; opt.Tags.Env.Use && state.Env != "-" && ok {
				if env != "-" {
					if state.EnvHierarchically {
						if state.Env != "" {
							env = state.Flag + "_" + env
						}
						state.Env = env
					}
					if err := enviropment.Get(value, env); err != nil {
						return err
					}
				} else {
					state.Env = "-"
					state.EnvHierarchically = false
				}
			}
			return nil
		}
		if flags.Using != FlagsDontUse || enviropment.Use {
			if err := traverseStruct(v, initFlags); err != nil {
				return err
			}
		}

		return flags.Process()
		
		// initValues := func(value reflect.Value, name string, tag reflect.StructTag, state *ref.State) error {
		// 	settings := parseTagSettings(tag, opt)
		// 	if envVarName, ok := settings["env"]; ok {
		// 		if err := env.Get(value, envVarName); err != nil {
		// 			return err
		// 		}
		// 	}
		// 	if flagSettings, ok := settings["flag"]; ok {
		// 		flag := parseFlagSettings(flagSettings)
		// 		if err := flags.Get(value, flag.Name); err != nil {
		// 			return err
		// 		}
		// 	}
		// 	return nil
		// }
		// return traverseStruct(v, initValues)
		//return nil
	}
	return nil
}
