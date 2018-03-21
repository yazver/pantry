package reflect

import (
	"encoding"
	"reflect"
	"testing"
)

func Test_getMarshaler(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want encoding.TextMarshaler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := getMarshaler(tt.args.v); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getMarshaler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_findMarshaler(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want encoding.TextMarshaler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := findMarshaler(tt.args.v); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. findMarshaler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getUnmarshaler(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want encoding.TextUnmarshaler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := getUnmarshaler(tt.args.v); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getUnmarshaler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_findUnmarshaler(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want encoding.TextUnmarshaler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := findUnmarshaler(tt.args.v); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. findUnmarshaler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAssignStringToValue(t *testing.T) {
	type args struct {
		dst reflect.Value
		src string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := AssignStringToValue(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
			t.Errorf("%q. AssignStringToValue() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestAssignValue(t *testing.T) {
	type args struct {
		dst reflect.Value
		src reflect.Value
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := AssignValue(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
			t.Errorf("%q. AssignValue() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTraverse(t *testing.T) {
	type args struct {
		v       interface{}
		process ProcessValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := Traverse(tt.args.v, tt.args.process); (err != nil) != tt.wantErr {
			t.Errorf("%q. Traverse() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTraverseValue(t *testing.T) {
	type args struct {
		v       reflect.Value
		process ProcessValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := TraverseValue(tt.args.v, tt.args.process); (err != nil) != tt.wantErr {
			t.Errorf("%q. TraverseValue() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_traverseValue(t *testing.T) {
	type args struct {
		v       reflect.Value
		path    string
		depth   uint
		field   *reflect.StructField
		process ProcessValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := traverseValue(tt.args.v, tt.args.path, tt.args.depth, tt.args.field, tt.args.process); (err != nil) != tt.wantErr {
			t.Errorf("%q. traverseValue() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTraverseFields(t *testing.T) {
	type args struct {
		v            interface{}
		processField ProcessValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := TraverseFields(tt.args.v, tt.args.processField); (err != nil) != tt.wantErr {
			t.Errorf("%q. TraverseFields() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestTraverseValueFields(t *testing.T) {
	type args struct {
		v            reflect.Value
		processField ProcessValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := TraverseValueFields(tt.args.v, tt.args.processField); (err != nil) != tt.wantErr {
			t.Errorf("%q. TraverseValueFields() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestClear(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		Clear(tt.args.v)
	}
}
