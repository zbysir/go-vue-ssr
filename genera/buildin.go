package genera

import (
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

	return fmt.Sprintf("%v", m)
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
