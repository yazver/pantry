package pantry

import (
	"os"
	"reflect"

	ref "github.com/yazver/golibs/reflect"
)

type Enviropment struct {
	Use            bool
	Prefix         string
	Hierarchically bool
}

func (e *Enviropment) Get(v reflect.Value, name string) (err error) {
	if e.Use {
		if env, ok := os.LookupEnv(e.Prefix + name); ok {
			err = ref.AssignStringToValue(v, env)
		}
	}
	return
}
