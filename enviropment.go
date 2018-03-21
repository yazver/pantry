package pantry

import (
	"errors"
	"os"
	"reflect"

	ref "github.com/yazver/golibs/reflect"
)

type Enviropment struct {
	Use    bool
	Prefix string
}

func (e *Enviropment) Get(v reflect.Value, name string) (err error) {
	if e.Use {
		defer func() {
			if v := recover(); v != nil {
				switch errorValue := v.(type) {
				case error:
					err = errorValue
				case string:
					err = errors.New(errorValue)
				default:
					panic(errorValue) // err := errors.New("Unable to get enviropment variable")
				}
			}

		}()

		if env, ok := os.LookupEnv(e.Prefix + name); ok {
			err = ref.AssignStringToValue(v, env)
		}
	}
	return
}
