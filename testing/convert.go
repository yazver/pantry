package testing

// type ConvertOptions struct {
// 	// IntBase
// 	// Default is 10
// 	IntBase int
// 	// FloatFormat is one of
// 	// 'b' (-ddddp±ddd, a binary exponent),
// 	// 'e' (-d.dddde±dd, a decimal exponent),
// 	// 'E' (-d.ddddE±dd, a decimal exponent),
// 	// 'f' (-ddd.dddd, no exponent),
// 	// 'g' ('e' for large exponents, 'f' otherwise), or
// 	// 'G' ('E' for large exponents, 'f' otherwise).
// 	// Default is 'g'
// 	FloatFormat byte
// }

// var defaultConvertOptions = &ConvertOptions{IntBase: 10, FloatFormat: 'g'}

// func ConvertValueToString(v reflect.Value, opt *ConvertOptions) (string, bool) {
// 	if opt == nil {
// 		opt = defaultConvertOptions
// 	}

// 	v = reflect.Indirect(v)
// 	switch v.Kind() {
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		return strconv.FormatInt(v.Int(), 10), true
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
// 		return strconv.FormatUint(v.Uint(), 10), true
// 	case reflect.Float32:
// 		return strconv.FormatFloat(v.Float(), 'g', -1, 32), true
// 	case reflect.Float64:
// 		return strconv.FormatFloat(v.Float(), 'g', -1, 64), true
// 	case reflect.Bool:
// 		return strconv.FormatBool(v.Bool()), true
// 	case reflect.String:
// 		return v.String(), true
// 	default:
// 		if v.CanInterface() {
// 			switch i := v.Interface().(type) {
// 			case encoding.TextMarshaler:
// 				if b, err := i.MarshalText(); err == nil {
// 					return string(b), true
// 				}
// 				return "", false
// 			case time.Duration:
// 				return i.String(), true
// 			}
// 		}
// 		if v.CanAddr() {
// 			v = v.Addr()
// 			if v.CanInterface() {
// 				if m, ok := v.Interface().(encoding.TextMarshaler); ok {
// 					if b, err := m.MarshalText(); err == nil {
// 						return string(b), true
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return "", false
// }

// func Value(t reflect.Type, rand *rand.Rand) (value reflect.Value, ok bool) {

// 	return quick.Value(t, rand)
// }
// func RandomValue(t reflect.Type, rand *rand.Rand) (value reflect.Value, ok bool) {
// 	if reflect.TypeOf(time.Time{}).AssignableTo(t) {

// 	} else if reflect.TypeOf(time.Duration(0)).AssignableTo(t) {

// 	} else if reflect.TypeOf(net.IP{}).AssignableTo(t) {

// 	} else if reflect.TypeOf().AssignableTo(t) {

// 	}
// 	quick.Value()
// }
