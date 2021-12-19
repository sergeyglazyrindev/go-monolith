package utils

import (
	"fmt"
	"reflect"
)

// GetString returns string representation on an instance of
// a model
func GetString(a interface{}) string {
	str, ok := a.(fmt.Stringer)
	if ok {
		return str.String()
	}
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	if t.Kind() != reflect.Ptr && t.Kind() == reflect.Struct {
		v = reflect.Indirect(reflect.New(t))
		v.Set(reflect.ValueOf(a))

		sp := v.Addr().Interface()
		str, ok := sp.(fmt.Stringer)
		if ok {
			return str.String()
		}
		if _, ok := t.FieldByName("Name"); ok {
			return v.FieldByName("Name").String()
		}
	} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		// Check if nil
		if v.IsNil() {
			return ""
		}

		if _, ok := t.Elem().FieldByName("Name"); ok {
			return v.Elem().FieldByName("Name").String()
		}
	} else if t.Kind() == reflect.Int && t.Name() != "int" {
		val := v.Int()
		// This is a static list type
		for i := 0; i < v.NumMethod(); i++ {
			ret := v.Method(i).Call([]reflect.Value{})
			if len(ret) > 0 {
				if ret[0].Int() == val {
					return t.Method(i).Name
				}
			}
		}
	}
	return fmt.Sprint(a)
}
