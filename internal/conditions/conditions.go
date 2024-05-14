package conditions

import "reflect"

func DeepCopy(i interface{}) interface{} {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		result := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			result.Field(i).Set(reflect.ValueOf(DeepCopy(v.Field(i).Interface())))
		}
		return result.Interface()
	default:
		return v.Interface()
	}
}
