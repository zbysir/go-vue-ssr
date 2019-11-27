package genera

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

// 混合动态和静态的标签, 主要是style/class需要混合
// todo) 如果style/class没有冲突, 则还可以优化
// tip: 纯静态的attr应该在编译时期就生成字符串, 而不应调用这个
func mixinAttr(attr []xml.Attr, data interface{}) (str string) {

	return
}

func lookInterface(data interface{}, key string) (desc interface{}) {
	m, ok := shouldLookInterface(data, key)
	if !ok {
		return ""
	}

	return m
}

func lookInterfaceToStr(data interface{}, key string) (desc string) {
	m, ok := shouldLookInterface(data, key)
	if !ok {
		return ""
	}

	return interfaceToStr(m)
}

func lookInterfaceToBool(data interface{}, key string, re ...bool) (desc bool) {
	m, ok := shouldLookInterface(data, key)
	if !ok {
		desc = false
	} else {
		desc = interfaceToBool(m)
	}

	if len(re) != 0 && re[0] {
		desc = !desc
	}
	return
}

// 扩展map, 实现作用域
func extendMap(src map[string]interface{}, ext map[string]interface{}) (desc map[string]interface{}) {
	//desc = make(map[string]interface{}, len(src))
	for i, v := range ext {
		src[i] = v
	}
	return src
}

func lookInterfaceToSlice(data interface{}, key string) (desc []interface{}) {
	m, ok := shouldLookInterface(data, key)
	if !ok {
		return nil
	}

	return interface2Slice(m)
}

func interfaceToStr(s interface{}) (d string) {
	switch a := s.(type) {
	case map[string]interface{}:
		bs, _ := json.Marshal(a)
		return string(bs)
	case int, string, float64:
		return fmt.Sprintf("%v", a)
	}

	return
}

// 字符串false,0 会被认定为false
func interfaceToBool(s interface{}) (d bool) {
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
	}

	return
}

func interface2Slice(s interface{}) (d []interface{}) {
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

func shouldLookInterface(data interface{}, key string) (desc interface{}, exist bool) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, false
	}

	kk := strings.Split(key, ".")
	c, ok := m[kk[0]]
	if len(kk) == 1 {
		if !ok {
			return nil, false
		}

		return c, true
	}

	return shouldLookInterface(c, strings.Join(kk[1:], "."))
}

func injectVal(src string, data interface{}) (to string) {
	reg := regexp.MustCompile(`\{\{.+?\}\}`)

	src = reg.ReplaceAllStringFunc(src, func(s string) string {
		key := s[2 : len(s)-2]

		desc, ok := shouldLookInterface(data, key)
		if ok {
			return fmt.Sprintf("%v", desc)
		}
		return ""
	})
	return src
}

// 判断, 支持简单的表达式:
// && || ! (并或非)
// isShow && !isHide
func condition(data map[string]interface{}, exp string) bool {
	return lookInterfaceToBool(data, exp)
}
