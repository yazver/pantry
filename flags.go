package pantry

import (
	"encoding"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"

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

type FlagsUsing int

const (
	FlagsDontUse FlagsUsing = iota
	FlagsUsePrediffenned
	FlagsUseAll
)

type Flags struct {
	Using   FlagsUsing
	FlagSet *flag.FlagSet
	parsed  bool
	values  map[string]interface{}
}

func (f *Flags) Add(v reflect.Value, name, usage string) error {
	if f.FlagSet == nil {
		f.FlagSet = flag.CommandLine
	}

	if name == "" {
		return errors.New("Undefined flag name")
	}

	if f.Using == FlagsUseAll && f.FlagSet.Lookup(name) == nil {
		f.parsed = false

		if v.CanInterface() {
			switch value := v.Interface().(type) {
			case time.Duration:
				f.FlagSet.Duration(name, value, usage)
				return nil
			case encoding.TextUnmarshaler:
				if m, ok := value.(encoding.TextMarshaler); ok {
					b, err := m.MarshalText()
					if err != nil {
						return fmt.Errorf("Unable to assign default value (%#v) to flag (%s): %s", value, name, err)
					}
					f.FlagSet.String(name, string(b), usage)
					return nil
				}
			}
		}
		switch v.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.FlagSet.Int64(name, v.Int(), usage)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			f.FlagSet.Uint64(name, v.Uint(), usage)
		case reflect.Int:
			f.FlagSet.Int(name, int(v.Int()), usage)
		case reflect.Uint:
			f.FlagSet.Uint(name, uint(v.Uint()), usage)
		case reflect.Float32, reflect.Float64:
			f.FlagSet.Float64(name, v.Float(), usage)
		case reflect.Bool:
			f.FlagSet.Bool(name, v.Bool(), usage)
		case reflect.String:
			f.FlagSet.String(name, v.String(), usage)
		default:
			return UnsupportedFlagTypeError{name, v.Type().Name()}
		}
	}
	//fo.flagSet.
	return nil
}

func (f *Flags) Get(dst reflect.Value, name string) (err error) {
	if f.Using == FlagsDontUse {
		return nil
	}
	if !f.parsed {
		if err = f.parse(); err != nil {
			return
		}
	}
	if value, ok := f.values[name]; ok {
		src := reflect.ValueOf(value)
		if err = ref.AssignValue(dst, src); err != nil {
			return fmt.Errorf("Unsupported type of field \"%s\": %s; Error: %s", name, dst.Type().Name(), err)
		}
	}
	return nil
}

func (f *Flags) parse() error {
	if !f.parsed {
		f.values = map[string]interface{}{}
		if err := f.FlagSet.Parse(os.Args[1:]); err != nil {
			return err
		}
		f.FlagSet.Visit(func(fl *flag.Flag) {
			if getter, ok := fl.Value.(flag.Getter); ok {
				f.values[fl.Name] = getter.Get()
			}
		})
	}
	return nil
}
