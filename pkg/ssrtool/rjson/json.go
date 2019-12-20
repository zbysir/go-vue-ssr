package rjson

import (
	"bytes"
	"encoding/json"
	"github.com/buger/jsonparser"
	"strings"
)

// 更高效的处理json字符串(使用Token, 而非反射)
// 在json字符串大但需要读取的数据少的时候使用.

func Get(bs []byte, key string) (desc interface{}) {
	v, t, _, err := jsonparser.Get(bs, strings.Split(key, ".")...)
	if err != nil {
		return
	}

	switch t {
	case jsonparser.Null, jsonparser.NotExist, jsonparser.Unknown:
	case jsonparser.Boolean:
		// fast path
		desc, _ = jsonparser.ParseBoolean(v)
	case jsonparser.String:
		desc, _ = jsonparser.ParseString(v)
	case jsonparser.Number:
		desc, _ = jsonparser.ParseFloat(v)
	case jsonparser.Object, jsonparser.Array:
		// low path
		_ = json.Unmarshal(v, &desc)
	}

	return
}

func GetBool(bs []byte, key string) (desc bool) {
	v, t, _, err := jsonparser.Get(bs, strings.Split(key, ".")...)
	if err != nil {
		return
	}
	switch t {
	case jsonparser.Boolean:
		return bytes.Equal(v, []byte("true"))
	case jsonparser.Number:
		f, err := jsonparser.ParseFloat(v)
		if err != nil {
			return false
		}
		return f != 0
	case jsonparser.String:
		a, err := jsonparser.ParseString(v)
		if err != nil {
			return false
		}

		return a != "" && a != "false" && a != "0"
	case jsonparser.Null:
		return false
	default:
		return true
	}
}

func GetStr(bs []byte, key string) (desc string) {
	v, err := jsonparser.GetString(bs, strings.Split(key, ".")...)
	if err != nil {
		return
	}
	return v
}

// 模糊获取字符串
// Number/Boolean等类似会返回原始json
func GetStrObscure(bs []byte, key string) (desc string) {
	v, t, _, err := jsonparser.Get(bs, strings.Split(key, ".")...)
	if err != nil {
		return
	}
	switch t {
	case jsonparser.Null, jsonparser.NotExist, jsonparser.Unknown:
		return
	case jsonparser.String:
		desc, _ = jsonparser.ParseString(v)
		return desc
	case jsonparser.Number, jsonparser.Object, jsonparser.Array, jsonparser.Boolean:
		return string(v)
	}

	return
}

func GetNumber(bs []byte, key string) (desc float64) {
	v, err := jsonparser.GetFloat(bs, strings.Split(key, ".")...)
	if err != nil {
		return
	}
	return v
}
