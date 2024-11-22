package mapstruct

import (
	"fmt"
	"reflect"
)

func defaultDecoderConfig(config ...*DecoderConfig) *DecoderConfig {
	cfg := append(config, &DecoderConfig{})[0]
	if cfg.TagName == "" {
		cfg.TagName = "json"
	}
	return cfg
}

type DecoderConfig struct {
	DecodeHook func(reflect.Type, any) (any, error)
	TagName    string
}

func Decode(obj any, config ...*DecoderConfig) (map[string]any, error) {
	cfg := defaultDecoderConfig(config...)
	result := make(map[string]any)

	sv := reflect.ValueOf(obj)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem() // 解引用指针
	}

	if sv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct but got %s", sv.Kind())
	}

	for i, st := 0, sv.Type(); i < sv.NumField(); i++ {
		key := st.Field(i).Tag.Get(cfg.TagName)
		if key == "" || key == "-" {
			continue
		}

		fv := sv.Field(i) // 获得字段的值
		if hook := cfg.DecodeHook; hook != nil {
			value, err := hook(fv.Type(), fv.Interface())
			if err != nil {
				return nil, err
			}
			result[key] = value
			continue
		}

		switch fv.Kind() {
		case reflect.Struct: // 递归处理嵌套结构体
			m, err := Decode(fv.Interface())
			if err != nil {
				return nil, err
			}
			result[key] = m
		default:
			result[key] = fv.Interface()
		}
	}

	return result, nil
}
