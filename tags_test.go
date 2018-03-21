package pantry

import (
	"reflect"
	"testing"
)

func Test_tagLookupValue(t *testing.T) {
	type args struct {
		tag reflect.StructTag
		key string
	}
	tests := []struct {
		name      string
		args      args
		wantValue string
		wantOk    bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotValue, gotOk := tagLookupValue(tt.args.tag, tt.args.key)
		if gotValue != tt.wantValue {
			t.Errorf("%q. tagLookupValue() gotValue = %v, want %v", tt.name, gotValue, tt.wantValue)
		}
		if gotOk != tt.wantOk {
			t.Errorf("%q. tagLookupValue() gotOk = %v, want %v", tt.name, gotOk, tt.wantOk)
		}
	}
}

func Test_parseTagSettings(t *testing.T) {
	type args struct {
		tag reflect.StructTag
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := parseTagSettings(tt.args.tag); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. parseTagSettings() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_parseFlagSettings(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name      string
		args      args
		wantName  string
		wantUsage string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		gotName, gotUsage := parseFlagSettings(tt.args.str)
		if gotName != tt.wantName {
			t.Errorf("%q. parseFlagSettings() gotName = %v, want %v", tt.name, gotName, tt.wantName)
		}
		if gotUsage != tt.wantUsage {
			t.Errorf("%q. parseFlagSettings() gotUsage = %v, want %v", tt.name, gotUsage, tt.wantUsage)
		}
	}
}

func Test_processTags(t *testing.T) {
	type args struct {
		v     interface{}
		flags *Flags
		env   *Enviropment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := processTags(tt.args.v, tt.args.flags, tt.args.env); (err != nil) != tt.wantErr {
			t.Errorf("%q. processTags() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
