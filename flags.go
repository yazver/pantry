package pantry

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	ref "github.com/yazver/golibs/reflect"
)

type UnsupportedFlagTypeError struct {
	name     string
	flagType string
}

// Returns the formatted configuration error.
func (e UnsupportedFlagTypeError) Error() string {
	return fmt.Sprintf("Unsupported type of flag \"%s\": %s", e.name, e.flagType)
}

type Flag struct {
	Name     string // name as it appears on command line
	Short    string // optional short name
	Usage    string // help message
	DefValue string // default value (as text); for usage message
	//Type     reflect.Type
	Value reflect.Value
	//actual   bool
}

type FlagsUsing int

const (
	FlagsDontUse FlagsUsing = iota
	FlagsUsePredefined
	FlagsUseAll
)

type Flags struct {
	Using          FlagsUsing
	FlagSet        *flag.FlagSet
	Args           []string
	ConfigPathFlag string
	Hierarchically bool
	parsed         bool
	values         []*Flag
}

func (flags *Flags) Init(flagSet *flag.FlagSet, arguments []string) {
	if flagSet == nil {
		flags.FlagSet = flag.CommandLine
	} else {
		flags.FlagSet = flagSet
	}
	if arguments == nil {
		flags.Args = os.Args[1:]
	} else {
		flags.Args = arguments
	}
}

func (flags *Flags) GetConfigPath() string {
	if flags.ConfigPathFlag == "" {
		return ""
	}
	fs := flag.NewFlagSet("", flag.ContinueOnError) // Temporary FlagSet
	fs.Usage = func() {}                            // Prevent showing usage message
	if flag := flags.FlagSet.Lookup(flags.ConfigPathFlag); flag == nil {
		flags.FlagSet.String(flags.ConfigPathFlag, "", "Path to config file")
	}
	configPath := fs.String(flags.ConfigPathFlag, "", "")
	_ = fs.Parse(flags.Args) // Ingnore error
	return *configPath
}

func (flags *Flags) Add(f *Flag) error {
	// if flags.values == nil {
	// 	flags.flags = []*Flag{}
	// }
	if flags.FlagSet == nil {
		flags.FlagSet = flag.CommandLine
	}

	if f.Name == "" {
		return errors.New("Undefined flag name")
	}

	if flags.Using == FlagsUseAll && flags.FlagSet.Lookup(f.Name) == nil {
		flags.parsed = false

		t := f.Value.Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		v := reflect.Value{}
		switch t.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
			reflect.Float32, reflect.Float64, reflect.Bool:
			v = reflect.New(t).Elem()
			defValue := strings.TrimSpace(f.DefValue)
			if defValue != "" {
				if err := ref.AssignStringToValue(v, defValue); err != nil {
					return fmt.Errorf("Can't convert default value \"%s\" to type \"%s\": %s", defValue, t.Name, err)
				}
			}
		}	

		switch t.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			flags.FlagSet.Int64(f.Name, v.Int(), f.Usage)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			flags.FlagSet.Uint64(f.Name, v.Uint(), f.Usage)
		case reflect.Int:
			flags.FlagSet.Int(f.Name, int(v.Int()), f.Usage)
		case reflect.Uint:
			flags.FlagSet.Uint(f.Name, uint(v.Uint()), f.Usage)
		case reflect.Float32, reflect.Float64:
			flags.FlagSet.Float64(f.Name, v.Float(), f.Usage)
		case reflect.Bool:
			flags.FlagSet.Bool(f.Name, v.Bool(), f.Usage)
		default:
			flags.FlagSet.String(f.Name, f.DefValue, f.Usage)
		}
	}
	flags.values = append(flags.values, f)
	return nil
}

// func (flags *Flags) Get(dst reflect.Value, name string) (err error) {
// 	if flags.Using == FlagsDontUse {
// 		return nil
// 	}
// 	if !flags.parsed {
// 		if err = flags.parse(); err != nil {
// 			return
// 		}
// 	}
// 	if f, ok := flags.flags[name]; ok {
// 		if f.actual {
// 			if err = ref.AssignValue(dst, f.value); err != nil {
// 			return fmt.Errorf("Unsupported type of field \"%s\": %s; Error: %s", name, dst.Type().Name(), err)
// 		}
// 	}
// 	}
// 	return nil
// }

func (flags *Flags) Process() error {
	if !flags.parsed {
		if err := flags.FlagSet.Parse(flags.Args); err != nil {
			return err
		}
		values := map[string]reflect.Value{}
		flags.FlagSet.Visit(func(fl *flag.Flag) {
			if getter, ok := fl.Value.(flag.Getter); ok {
				values[fl.Name] = reflect.ValueOf(getter.Get())
			}
		})
		for _, value := range flags.values {
			if v, ok := values[value.Name]; ok {
				if err := ref.AssignValue(value.Value, v); err != nil {
					return errors.New("Can't assign flag to value: " + err.Error())
				}
			}
		}

	}
	return nil
}
