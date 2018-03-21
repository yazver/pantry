package reflect

import (
	"math/rand"
	"net"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

func Test_assignStringToValue(t *testing.T) {
	//var i int
	//var i8 int8
	//var i16 int16
	//var i32 int32
	//var i64 int64
	//var ui uint
	//var ui8 uint8
	//var ui16 uint16
	//var ui32 uint32
	//var ui64 uint64
	//var f32 float32
	//var f64 float64
	//var b bool
	//var s string
	//var ti time.Time
	//var d time.Duration
	//var ip net.IP

	//type args struct {
	//	dst reflect.Value
	//	src string
	//}
	//
	//tests := []struct {
	//	args    args
	//	wantErr bool
	//}{
	//	{args{reflect.ValueOf(&i), "8546778"}, false},
	//	{args{reflect.ValueOf(&i8), "16"}, false},
	//	{args{reflect.ValueOf(&i16), "-30000"}, false},
	//	{args{reflect.ValueOf(&i32), "0x7fffffff"}, false},
	//	{args{reflect.ValueOf(&i64), "-0xffffffff"}, false},
	//	{args{reflect.ValueOf(&i), "100.1"}, true},
	//	{args{reflect.ValueOf(&i8), "255"}, true},
	//	{args{reflect.ValueOf(&i16), "0xffff"}, true},
	//	{args{reflect.ValueOf(&i32), "0xffffffff"}, true},
	//	{args{reflect.ValueOf(&i64), "true"}, true},
	//	{args{reflect.ValueOf(&ui), "8546778"}, false},
	//	{args{reflect.ValueOf(&ui8), "16"}, false},
	//	{args{reflect.ValueOf(&ui16), "64000"}, false},
	//	{args{reflect.ValueOf(&ui32), "0xffffffff"}, false},
	//	{args{reflect.ValueOf(&ui64), "0xffffffffffffffff"}, false},
	//	{args{reflect.ValueOf(&ui), "100.1"}, true},
	//	{args{reflect.ValueOf(&ui8), "-255"}, true},
	//	{args{reflect.ValueOf(&ui16), "0xffffffff"}, true},
	//	{args{reflect.ValueOf(&ui32), "0x8ffffffff"}, true},
	//	{args{reflect.ValueOf(&ui64), "fish"}, true},
	//	{args{reflect.ValueOf(&f32), "10.1"}, false},
	//	{args{reflect.ValueOf(&f64), "-5.12345678e42"}, false},
	//	{args{reflect.ValueOf(&f32), "rabbit"}, true},
	//	{args{reflect.ValueOf(&f64), "5.1234.5678"}, true},
	//	{args{reflect.ValueOf(&b), "true"}, false},
	//	{args{reflect.ValueOf(&b), "0"}, false},
	//	{args{reflect.ValueOf(&b), "FaLsE"}, true},
	//	{args{reflect.ValueOf(&b), "10"}, true},
	//	{args{reflect.ValueOf(&s), "I DOWN THE RABBIT HOLE"}, false},
	//	{args{reflect.ValueOf(&ti), "1832-01-27T01:02:03+05:00"}, false},
	//	{args{reflect.ValueOf(&ti), "1832-01-27T01-02-03"}, true},
	//	{args{reflect.ValueOf(&d), "22h49m22s0ms"}, false},
	//	{args{reflect.ValueOf(&d), "25r"}, true},
	//	{args{reflect.ValueOf(&ip), "127.0.0.1"}, false},
	//	{args{reflect.ValueOf(&ip), "2001:db8::"}, false},
	//	{args{reflect.ValueOf(&ip), "1000.0.0.1"}, true},
	//}
	//
	//for _, tt := range tests {
	//	value := reflect.Indirect(tt.args.dst)
	//	value.Set(reflect.Zero(value.Type()))
	//	if err := assignStringToValue(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
	//		t.Errorf("assignStringToValue(%v(%v), %v), error = %v, wantErr %v",
	//			value.Type().Name(), value.Interface(), tt.args.src, err, tt.wantErr)
	//	}
	//}

	//types := []reflect.Type{
	//	reflect.TypeOf(int(0)),
	//	reflect.TypeOf(&i),
	//	reflect.TypeOf(&i8),
	//	reflect.TypeOf(&i16),
	//	reflect.TypeOf(&i32),
	//	reflect.TypeOf(&i64),
	//	reflect.TypeOf(&ui),
	//	reflect.TypeOf(&ui8),
	//	reflect.TypeOf(&ui16),
	//	reflect.TypeOf(&ui32),
	//	reflect.TypeOf(&ui64),
	//	reflect.TypeOf(&f32),
	//	reflect.TypeOf(&f64),
	//	reflect.TypeOf(&b),
	//	reflect.TypeOf(&s),
	//	reflect.TypeOf(&ti),
	//	reflect.TypeOf(&d),
	//	reflect.TypeOf(&ip),
	//}
	types := []reflect.Type{
		reflect.TypeOf(int(0)),
		reflect.TypeOf(int(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(int16(0)),
		reflect.TypeOf(int32(0)),
		reflect.TypeOf(int64(0)),
		reflect.TypeOf(uint(0)),
		reflect.TypeOf(uint8(0)),
		reflect.TypeOf(uint16(0)),
		reflect.TypeOf(uint32(0)),
		reflect.TypeOf(uint64(0)),
		reflect.TypeOf(float32(0)),
		reflect.TypeOf(float64(0)),
		reflect.TypeOf(bool(false)),
		//reflect.TypeOf(""),
		//reflect.TypeOf(time.Now()),
		//reflect.TypeOf(time.Second),
		//reflect.TypeOf(net.IPv4(127, 0, 0, 1)),
	}

	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for _, tt := range types {

		dst := reflect.New(tt).Elem()
		if rand.Intn(2) != 0 {
			dst = dst.Addr()
		}
		srcValue, _ := quick.Value(tt, r)
		src := srcValue.String()
		//if rand.Intn(2) != 0 {
		//	src = src.Addr()
		//}
		wantErr := false
		value := reflect.Indirect(dst)
		if err := AssignStringToValue(dst, src); (err != nil) != wantErr {
			t.Errorf("assignStringToValue(%v(%v), %v), error = %v, wantErr %v",
				value.Type().Name(), value.Interface(), src, err, wantErr)
		} else {
			t.Logf("assignStringToValue(%v(%v), %v)",
				value.Type().Name(), value.Interface(), src)
		}
	}

}

func Test_assignValue(t *testing.T) {
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var ui uint
	var ui8 uint8
	var ui16 uint16
	var ui32 uint32
	var ui64 uint64
	var f32 float32
	var f64 float64
	var b bool
	var s string
	var ti time.Time
	var d time.Duration
	var ip net.IP

	b2 := true

	type args struct {
		dst reflect.Value
		src reflect.Value
	}
	tests := []struct {
		args    args
		wantErr bool
	}{
		{args{reflect.ValueOf(&i), reflect.ValueOf("8546778")}, false},
		{args{reflect.ValueOf(&i8), reflect.ValueOf("16")}, false},
		{args{reflect.ValueOf(&i16), reflect.ValueOf("-30000")}, false},
		{args{reflect.ValueOf(&i32), reflect.ValueOf("0x7fffffff")}, false},
		{args{reflect.ValueOf(&i64), reflect.ValueOf("-0xffffffff")}, false},
		{args{reflect.ValueOf(&i), reflect.ValueOf(int(100))}, false},
		{args{reflect.ValueOf(&i8), reflect.ValueOf(&i)}, false},
		{args{reflect.ValueOf(&i16), reflect.ValueOf(i8)}, false},
		{args{reflect.ValueOf(&i32), reflect.ValueOf(int64(1000000000))}, false},
		{args{reflect.ValueOf(&i64), reflect.ValueOf(&i32)}, false},
		{args{reflect.ValueOf(&i), reflect.ValueOf("100.1")}, true},
		{args{reflect.ValueOf(&i8), reflect.ValueOf("255")}, true},
		{args{reflect.ValueOf(&i16), reflect.ValueOf("0xffff")}, true},
		{args{reflect.ValueOf(&i32), reflect.ValueOf("0xffffffff")}, true},
		{args{reflect.ValueOf(&i64), reflect.ValueOf("true")}, true},
		{args{reflect.ValueOf(&i), reflect.ValueOf(1 - 0.707i)}, true},
		{args{reflect.ValueOf(&i8), reflect.ValueOf(ip)}, true},
		{args{reflect.ValueOf(&i16), reflect.ValueOf(&ti)}, true},
		{args{reflect.ValueOf(&i32), reflect.ValueOf(int64(1000000000000))}, true},
		{args{reflect.ValueOf(&i64), reflect.ValueOf(&b)}, true},

		{args{reflect.ValueOf(&ui), reflect.ValueOf("8546778")}, false},
		{args{reflect.ValueOf(&ui8), reflect.ValueOf("16")}, false},
		{args{reflect.ValueOf(&ui16), reflect.ValueOf("64000")}, false},
		{args{reflect.ValueOf(&ui32), reflect.ValueOf("0xffffffff")}, false},
		{args{reflect.ValueOf(&ui64), reflect.ValueOf("0xffffffffffffffff")}, false},
		{args{reflect.ValueOf(&ui), reflect.ValueOf(uint(100))}, false},
		{args{reflect.ValueOf(&ui8), reflect.ValueOf(ui)}, false},
		{args{reflect.ValueOf(&ui16), reflect.ValueOf(&ui8)}, false},
		{args{reflect.ValueOf(&ui32), reflect.ValueOf(uint16(64000))}, false},
		{args{reflect.ValueOf(&ui64), reflect.ValueOf(&ui32)}, false},
		{args{reflect.ValueOf(&ui), reflect.ValueOf("100.1")}, true},
		{args{reflect.ValueOf(&ui8), reflect.ValueOf("-255")}, true},
		{args{reflect.ValueOf(&ui16), reflect.ValueOf("0xffffffff")}, true},
		{args{reflect.ValueOf(&ui32), reflect.ValueOf("0x8ffffffff")}, true},
		{args{reflect.ValueOf(&ui64), reflect.ValueOf("fish")}, true},
		{args{reflect.ValueOf(&ui), reflect.ValueOf(int(-100))}, true},
		{args{reflect.ValueOf(&ui8), reflect.ValueOf(256)}, true},
		{args{reflect.ValueOf(&ui16), reflect.ValueOf(true)}, true},
		{args{reflect.ValueOf(&ui32), reflect.ValueOf(&ti)}, true},
		{args{reflect.ValueOf(&ui64), reflect.ValueOf(complex64(1.1))}, true},

		{args{reflect.ValueOf(&f32), reflect.ValueOf("10.1")}, false},
		{args{reflect.ValueOf(&f64), reflect.ValueOf("-5.12345678e42")}, false},
		{args{reflect.ValueOf(&f32), reflect.ValueOf(int(100))}, false},
		{args{reflect.ValueOf(&f64), reflect.ValueOf(&f32)}, false},
		{args{reflect.ValueOf(&f32), reflect.ValueOf("rabbit")}, true},
		{args{reflect.ValueOf(&f64), reflect.ValueOf("5.1234.5678")}, true},
		{args{reflect.ValueOf(&f32), reflect.ValueOf(b)}, true},
		{args{reflect.ValueOf(&f64), reflect.ValueOf(&ip)}, true},

		{args{reflect.ValueOf(&b), reflect.ValueOf("true")}, false},
		{args{reflect.ValueOf(&b), reflect.ValueOf("0")}, false},
		{args{reflect.ValueOf(&b), reflect.ValueOf(false)}, false},
		{args{reflect.ValueOf(&b), reflect.ValueOf(&b2)}, false},
		{args{reflect.ValueOf(&b), reflect.ValueOf("FaLsE")}, true},
		{args{reflect.ValueOf(&b), reflect.ValueOf("10")}, true},
		{args{reflect.ValueOf(&b), reflect.ValueOf(uint(20))}, true},
		{args{reflect.ValueOf(&b), reflect.ValueOf(&f64)}, true},

		{args{reflect.ValueOf(&s), reflect.ValueOf("I DOWN THE RABBIT HOLE")}, false},
		{args{reflect.ValueOf(&s), reflect.ValueOf(&i8)}, false},

		{args{reflect.ValueOf(&ti), reflect.ValueOf("1832-01-27T01:02:03+05:00")}, false},
		{args{reflect.ValueOf(&ti), reflect.ValueOf(time.Now())}, false},
		{args{reflect.ValueOf(&ti), reflect.ValueOf("1832-01-27T01-02-03")}, true},
		{args{reflect.ValueOf(&ti), reflect.ValueOf(&b)}, true},

		{args{reflect.ValueOf(&d), reflect.ValueOf("22h49m22s0ms")}, false},
		{args{reflect.ValueOf(&d), reflect.ValueOf(time.Hour)}, false},
		{args{reflect.ValueOf(&d), reflect.ValueOf("25r")}, true},
		{args{reflect.ValueOf(&d), reflect.ValueOf(&ti)}, true},

		{args{reflect.ValueOf(&ip), reflect.ValueOf("127.0.0.1")}, false},
		{args{reflect.ValueOf(&ip), reflect.ValueOf("2001:db8::")}, false},
		{args{reflect.ValueOf(&ip), reflect.ValueOf(net.IPv4(4, 31, 198, 44))}, false},
		{args{reflect.ValueOf(&ip), reflect.ValueOf("1000.0.0.1")}, true},
		{args{reflect.ValueOf(&ip), reflect.ValueOf(d)}, true},
	}
	for _, tt := range tests {
		value1 := reflect.Indirect(tt.args.dst)
		value1.Set(reflect.Zero(value1.Type()))
		value2 := reflect.Indirect(tt.args.src)
		if err := AssignValue(tt.args.dst, tt.args.src); (err != nil) != tt.wantErr {
			t.Errorf("assignValue(%v(%v), %v(%v)), error = %v, wantErr %v",
				value1.Type().Name(), value1.Interface(), value2.Type().Name(), value2.Interface(),
				err, tt.wantErr)
		}
	}
}

func Test_traverseStruct(t *testing.T) {
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
		if err := TraverseStruct(tt.args.v, tt.args.process); (err != nil) != tt.wantErr {
			t.Errorf("%q. traverseStruct() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func Test_traverseStructValue(t *testing.T) {
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
		if err := TraverseStructValue(tt.args.v, tt.args.process); (err != nil) != tt.wantErr {
			t.Errorf("%q. traverseStructValue() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
