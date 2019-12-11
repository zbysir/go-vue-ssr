package rinterface

import (
	"strings"
)

// 在map[string]interface{}中找到多级key的value
// kye: e.g. info.name
func Get(data interface{}, key string) (desc interface{}) {
	m, isObj := data.(map[string]interface{})

	kk := strings.Split(key, ".")
	currKey := kk[0]

	// 如果是对象, 则继续查找下一级
	if len(kk) != 1 && isObj {
		c, ok := m[currKey]
		if !ok {
			return
		}
		return Get(c, strings.Join(kk[1:], "."))
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
					return len(ToSlice(t))
				}
			}
		}
	} else {
		// key不只有一个, 但是data不是对象, 说明出现了undefined的问题, 直接return
		return
	}

	return
}

func GetStr(data interface{}, key string, escaped ...bool) string {
	return ToStr(Get(data, key), escaped...)
}

func GetInt(data interface{}, key string) int64 {
	return ToInt(Get(data, key))
}

func GetSlice(data interface{}, key string) []interface{} {
	return ToSlice(Get(data, key))
}

func GetSliceInt(data interface{}, key string) []int64 {
	return ToSliceInt(Get(data, key))
}
