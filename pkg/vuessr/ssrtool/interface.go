package ssrtool

import (
	"encoding/json"
	"fmt"
)

func InterfaceToStr(s interface{}) (d string) {
	switch a := s.(type) {
	case int, int64, string, float64:
		return fmt.Sprintf("%v", a)
	default:
		bs, _ := json.Marshal(a)
		return string(bs)
	}
}

func InterfaceToInt(s interface{}) (d int64) {
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

func InterfaceToSliceInt(s interface{}) (d []int64) {
	ss := Interface2Slice(s)
	d = make([]int64, len(ss))
	for i, v := range ss {
		d[i] = InterfaceToInt(v)
	}
	return
}

func Interface2Slice(s interface{}) (d []interface{}) {
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
func InterfaceToBool(s interface{}) (d bool) {
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
