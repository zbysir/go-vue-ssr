package ssrtool

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"golang.org/x/net/html"
	"strings"
)

func InterfaceToStr(s interface{}, escaped ...bool) (d string) {
	switch a := s.(type) {
	case int, string, float64:
		d = fmt.Sprintf("%v", a)
	default:
		bs, _ := json.Marshal(a)
		d = string(bs)
	}
	if len(escaped) == 1 && escaped[0] {
		d = escape(d)
	}
	return
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
	ss := InterfaceToSlice(s)
	d = make([]int64, len(ss))
	for i, v := range ss {
		d[i] = InterfaceToInt(v)
	}
	return
}

func InterfaceToSlice(s interface{}) (d []interface{}) {
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

func LookJson(bs []byte, key string) (desc interface{}) {
	v, _, _, err := jsonparser.Get(bs, strings.Split(key, ".")...)
	if err != nil {
		return
	}
	_ = json.Unmarshal(v, &desc)
	return
}

// 在map[string]interface{}中找到多级key的value
// kye: e.g. info.name
func LookInterface(data interface{}, key string) (desc interface{}) {
	m, isObj := data.(map[string]interface{})

	kk := strings.Split(key, ".")
	currKey := kk[0]

	// 如果是对象, 则继续查找下一级
	if len(kk) != 1 && isObj {
		c, ok := m[currKey]
		if !ok {
			return
		}
		return LookInterface(c, strings.Join(kk[1:], "."))
	}

	if len(kk) == 1 {
		if isObj {
			c, ok := m[currKey]
			if !ok {
				return
			}
			return c
		} else {
			switch currKey {
			case "length":
				switch t := data.(type) {
				// string
				case string:
					return len(t)
				default:
					// slice
					return len(InterfaceToSlice(t))
				}
			}
		}
	} else {
		// key不只有一个, 但是data不是对象, 说明出现了undefined的问题, 直接return
		return
	}

	return
}

func LookStr(data interface{}, key string, escaped ...bool) string {
	return InterfaceToStr(LookInterface(data, key), escaped...)
}

func LookInt(data interface{}, key string) int64 {
	return InterfaceToInt(LookInterface(data, key))
}

func LookSlice(data interface{}, key string) []interface{} {
	return InterfaceToSlice(LookInterface(data, key))
}

func LookSliceInt(data interface{}, key string) []int64 {
	return InterfaceToSliceInt(LookInterface(data, key))
}

func escape(src string) string {
	return html.EscapeString(src)
}
