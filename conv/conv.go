package conv

import (
	"fmt"
	"reflect"
)

func Serialize(obj any) (map[string]any, error) {
	result := make(map[string]any)

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // 解引用指针
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct but got %s", val.Kind())
	}

	vType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		var key string
		f := vType.Field(i)
		if key = f.Tag.Get("json"); key == "" || key == "-" {
			continue
		}

		switch f := val.Field(i); f.Kind() {
		case reflect.Struct: // 递归处理嵌套结构体
			m, err := Serialize(f.Interface())
			if err != nil {
				return nil, err
			}
			result[key] = m
		default:
			result[key] = f.Interface()
		}
	}

	return result, nil
}