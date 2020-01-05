package rinterface

import (
	"encoding/json"
	"fmt"
	"html"
)

// 用于处理{{}}语法, 当传递的是对象时, 应该和vue处理逻辑一致: 序列化为json字符串.
func ToStr(s interface{}, escaped bool) (d string) {
	switch a := s.(type) {
	case int, string, float64, int64, int32, float32:
		d = fmt.Sprintf("%v", a)
	default:
		bs, _ := json.Marshal(a)
		d = string(bs)
	}
	if  escaped{
		d = html.EscapeString(d)
	}
	return
}

func ToInt(s interface{}) (d int64) {
	switch a := s.(type) {
	case int:
		return int64(a)
	case int8:
		return int64(a)
	case int32:
		return int64(a)
	case int64:
		return a
	case float64:
		return int64(a)
	case float32:
		return int64(a)
	default:
		return 0
	}
}

func ToSliceInt(s interface{}) (d []int64) {
	ss := ToSlice(s)
	d = make([]int64, len(ss))
	for i, v := range ss {
		d[i] = ToInt(v)
	}
	return
}

func ToSlice(s interface{}) (d []interface{}) {
	switch a := s.(type) {
	case []interface{}:
		return a
	case []map[string]interface{}:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int64:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int32:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []string:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []float64:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	}
	return
}

// 字符串false,0 会被认定为false
func ToBool(s interface{}) (d bool) {
	if s == nil {
		return false
	}
	switch a := s.(type) {
	case bool:
		return a
	case int, float64, float32, int8, int64, int32, int16:
		return a != 0
	case string:
		return a != "" && a != "false" && a != "0"
	default:
		return true
	}
}
