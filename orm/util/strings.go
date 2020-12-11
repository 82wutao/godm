package util

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"dm.net/datamine/orm/errors"
)

// StringValue 将一个对象字符串化
func StringValue(val interface{}) string {
	if val == nil {
		return "NULL"
	}
	reflector := func(v interface{}, vt reflect.Type, vk reflect.Kind, vv reflect.Value) string {
		switch vk {
		case reflect.String:
			return fmt.Sprintf("'%s'", v.(string))
		case reflect.Bool:
			return fmt.Sprintf("%t", vv.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%d", vv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fmt.Sprintf("%d", vv.Uint())
		case reflect.Uintptr:
			return fmt.Sprintf("%d", vv.Uint())
		case reflect.Float32, reflect.Float64:
			return fmt.Sprintf("%f", vv.Float())
		case reflect.Complex64, reflect.Complex128:
			return fmt.Sprintf("(%f + %fi)", real(vv.Complex()), imag(vv.Complex()))
		case reflect.Array:
			length := vt.Len()
			results := make([]string, length)
			for i := 0; i < length; i++ {
				s := stringValue(vv.Index(i).Interface())
				results[i] = s
			}
			return strings.Join(results, ",")
		case reflect.Slice:
			if vv.IsNil() {
				return "NULL"
			}
			length := vv.Len()
			results := make([]string, length)
			for i := 0; i < length; i++ {
				s := stringValue(vv.Index(i).Interface())
				results[i] = s
			}
			return strings.Join(results, ",")
		case reflect.Ptr:
			if vv.IsNil() {
				return "NULL"
			}
			return stringValue(vv.Elem().Interface())
		case reflect.Func, reflect.Interface:
			if vv.IsNil() {
				return "NULL"
			}

			return fmt.Sprintf("'%s'", runtime.FuncForPC(vv.Pointer()).Name())
		case reflect.Struct:
			mth, existed := vt.MethodByName("String")
			if !existed {
				panic(errors.ERR_ORM_STRUCT_MUST_OVERWRITE_STRING)
			}
			function := mth.Func
			rets := function.Call([]reflect.Value{vv})
			return fmt.Sprintf("'%s'", rets[0].String())
		}
		// case reflect.Chan, reflect.Map, reflect.UnsafePointer:
		panic(errors.ERR_REFLECT_DATAKIND_CANOT_BE_REFLECT)
	}
	vt := reflect.TypeOf(val)
	vk := reflect.TypeOf(val).Kind()
	vv := reflect.ValueOf(val)

	return reflector(val, vt, vk, vv)
}

// StringValues 将一个对象切片字符串化
func StringValues(values []interface{}) string {
	length := len(values)
	results := make([]string, length)
	for i, v := range values {
		s := StringValue(v)
		results[i] = s
	}

	return strings.Join(results, ",")
}
