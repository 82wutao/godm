package orm

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func stringValue(val interface{}) (string, error) {
	if val == nil {
		return "NULL", nil
	}
	reflector := func(v interface{}, vt reflect.Type, vk reflect.Kind, vv reflect.Value) (string, error) {
		switch vk {
		case reflect.String:
			return fmt.Sprintf("'%s'", v.(string)), nil
		case reflect.Bool:
			return fmt.Sprintf("%t", vv.Bool()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%d", vv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fmt.Sprintf("%d", vv.Uint()), nil
		case reflect.Uintptr:
			return fmt.Sprintf("%d", vv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return fmt.Sprintf("%f", vv.Float()), nil
		case reflect.Complex64, reflect.Complex128:
			return fmt.Sprintf("(%f + %fi)", real(vv.Complex()), imag(vv.Complex())), nil
		case reflect.Array:
			length := vt.Len()
			results := make([]string, length)
			for i := 0; i < length; i++ {
				s, e := stringValue(vv.Index(i).Interface())
				if e != nil {
					return "", e
				}
				results[i] = s
			}
			return strings.Join(results, ","), nil
		case reflect.Slice:
			if vv.IsNil() {
				return "NULL", nil
			}
			length := vv.Len()
			results := make([]string, length)
			for i := 0; i < length; i++ {
				s, e := stringValue(vv.Index(i).Interface())
				if e != nil {
					return "", e
				}
				results[i] = s
			}
			return strings.Join(results, ","), nil
		case reflect.Ptr:
			if vv.IsNil() {
				return "NULL", nil
			}
			return stringValue(vv.Elem().Interface())
		case reflect.Func, reflect.Interface:
			if vv.IsNil() {
				return "NULL", nil
			}

			return fmt.Sprintf("'%s'", runtime.FuncForPC(vv.Pointer()).Name()), nil
		case reflect.Struct:
			mth, existed := vt.MethodByName("String")
			if !existed {
				return "", ERR_ORM_STRUCT_MUST_OVERWRITE_STRING
			}
			function := mth.Func
			rets := function.Call([]reflect.Value{vv})
			return fmt.Sprintf("'%s'", rets[0].String()), nil
		}
		// case reflect.Chan, reflect.Map, reflect.UnsafePointer:
		return "", ERR_REFLECT_DATAKIND_CANOT_BE_REFLECT
	}
	vt := reflect.TypeOf(val)
	vk := reflect.TypeOf(val).Kind()
	vv := reflect.ValueOf(val)

	return reflector(val, vt, vk, vv)
}

/**

 */
func stringValues(values []interface{}) (string, error) {
	length := len(values)
	results := make([]string, length)
	for i, v := range values {
		s, err := stringValue(v)
		if err != nil {
			return "", err
		}
		results[i] = s
	}

	return strings.Join(results, ","), nil
}
