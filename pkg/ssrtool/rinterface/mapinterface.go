package rinterface

import (
	"strconv"
	"strings"
)

// 在map[string]interface{}中找到多级key的value
// kye: e.g. info.name
func Get(data interface{}, key ...string) (desc interface{}) {
	keys := key
	// 兼容老写法
	if len(key) == 1 {
		keys = strings.Split(key[0], ".")
	}
	desc, _ = shouldLookInterface(data, keys...)
	return
}

// shouldLookInterface会返回interface(map[string]interface{})中指定的keys路径的值
func shouldLookInterface(data interface{}, keys ...string) (desc interface{}, exist bool) {
	if len(keys) == 0 {
		return data, true
	}

	currKey := keys[0]

	switch data := data.(type) {
	case map[string]interface{}:
		c, ok := data[currKey]
		if !ok {
			return
		}

		return shouldLookInterface(c, keys[1:]...)
	case []interface{}:
		switch currKey {
		case "length":
			// length
			return len(data), true
		default:
			// index
			index, ok := strconv.ParseInt(currKey, 10, 64)
			if ok != nil {
				return
			}
			if int(index) >= len(data) {
				return
			}
			return shouldLookInterface(data[index], keys[1:]...)
		}
	case string:
		switch currKey {
		case "length":
			// length
			return len(data), true
		default:
		}
	}

	return
}

func GetStr(data interface{}, key ...string) string {
	return ToStr(Get(data, key...), false)
}

func GetBool(data interface{}, key ...string) bool {
	return ToBool(Get(data, key...))
}

func GetInt(data interface{}, key ...string) int64 {
	return ToInt(Get(data, key...))
}

func GetSlice(data interface{}, key ...string) []interface{} {
	return ToSlice(Get(data, key...))
}

func GetSliceInt(data interface{}, key ...string) []int64 {
	return ToSliceInt(Get(data, key...))
}
