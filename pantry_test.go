package pantry

import (
	"testing"

	"github.com/yazver/golibs/reflect"
)

func TestPantry_Load(t *testing.T) {
	type args struct {
		path string
		v    interface{}
		opt  []func(*LoadOptions)
	}
	tests := []struct {
		name    string
		p       *Pantry
		args    args
		want    Box
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Load(tt.args.path, tt.args.v, tt.args.opt...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pantry.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pantry.Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
