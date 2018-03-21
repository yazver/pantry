package pantry

import (
	"reflect"
	"testing"
)

func Test_assignStringToValue(t *testing.T) {
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
		if err := assignStringToValue(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
			t.Errorf("%q. assignStringToValue() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_assignValue(t *testing.T) {
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
		if err := assignValue(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
			t.Errorf("%q. assignValue() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_traverseStruct(t *testing.T) {
	type args struct {
		v       interface{}
		process processField
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := traverseStruct(tt.args.v, tt.args.process); (err != nil) != tt.wantErr {
			t.Errorf("%q. traverseStruct() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_traverseStructValue(t *testing.T) {
	type args struct {
		v       reflect.Value
		process processField
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := traverseStructValue(tt.args.v, tt.args.process); (err != nil) != tt.wantErr {
			t.Errorf("%q. traverseStructValue() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
