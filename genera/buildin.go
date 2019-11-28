package genera

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// 渲染组件需要的结构
type Options struct {
	Props     map[string]interface{}   // 上级传递的 数据(包含了class和style)
	Attrs     map[string]string        // 上级传递的 静态的attrs (除去class和style), 只会作用在root节点
	Class     []string                 // 静态class, 只会作用在root节点
	Style     map[string]string        // 静态style, 只会作用在root节点
	StyleKeys []string                 // 样式的key, 用来保证顺序, 只会作用在root节点
	Slot      map[string]namedSlotFunc // 插槽代码, 支持多个不同名字的插槽, 如果没有名字则是"default"
}

// 混合动态和静态的标签, 主要是style/class需要混合
// todo) 如果style/class没有冲突, 则还可以优化
// tip: 纯静态的class应该在编译时期就生成字符串, 而不应调用这个
// classProps: 支持 obj, array, string
func mixinClass(options *Options, staticClass []string, classProps interface{}) (str string) {
	// 静态
	str = strings.Join(staticClass, " ")

	if str != "" {
		str += " "
	}
	// 本身的props
	str += genClassFromProps(classProps)

	if options != nil {
		// 上层传递的props
		if options.Props != nil {
			prop, ok := options.Props["class"]
			if ok {
				if str != "" {
					str += " "
				}
				str += genClassFromProps(prop)
			}
		}

		// 上层传递的静态class
		if len(options.Class) != 0 {
			if str != "" {
				str += " "
			}
			str += strings.Join(options.Class, " ")
		}
	}

	if str != "" {
		str = fmt.Sprintf(` class="%s"`, str)
	}

	return
}

// 构建style, 生成如`style="color: red"`的代码, 如果style代码为空 则只会返回空字符串
func mixinStyle(options *Options, staticStyle map[string]string, styleProps interface{}) (str string) {
	style := map[string]string{}

	// 静态
	for k, v := range staticStyle {
		style[k] = v
	}

	// 当前props
	ps := getStyleFromProps(styleProps)
	for k, v := range ps {
		style[k] = v
	}

	if options != nil {
		// 上层传递的props
		if options.Props != nil {
			prop, ok := options.Props["style"]
			if ok {
				ps := getStyleFromProps(prop)
				for k, v := range ps {
					style[k] = v
				}
			}
		}

		// 上层传递的静态style
		for k, v := range options.Style {
			style[k] = v
		}
	}

	styleCode := genStyle(style)
	if styleCode != "" {
		str = fmt.Sprintf(` style="%s"`, styleCode)
	}

	return
}

func getSortedKey(m map[string]string) (keys []string) {
	keys = make([]string, len(m))
	index := 0
	for k := range m {
		keys[index] = k
		index++
	}
	if len(m) < 2 {
		return keys
	}

	sort.Strings(keys)

	return
}

func genStyle(style map[string]string) string {
	sortedKeys := getSortedKey(style)

	st := ""
	for _, k := range sortedKeys {
		v := style[k]
		st += fmt.Sprintf("%s: %s; ", k, v)
	}
	return st
}

func getStyleFromProps(styleProps interface{}) map[string]string {
	pm, ok := styleProps.(map[string]interface{})
	if !ok {
		return nil
	}
	st := map[string]string{}
	for k, v := range pm {
		st[k] = fmt.Sprintf("%v", v)
	}
	return st
}

// classProps: 支持 obj, array, string
func genClassFromProps(classProps interface{}) string {
	if classProps == nil {
		return ""
	}
	switch t := classProps.(type) {
	case []string:
		return strings.Join(t, " ")
	case string:
		return t
	case map[string]interface{}:
		c := ""
		for k, v := range t {
			if interfaceToBool(v) {
				if c != "" {
					c += " "
				}
				c += k
			}
		}

		return c
	}

	return ""
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
	case int, string, float64:
		return fmt.Sprintf("%v", a)
	default:
		bs, _ := json.Marshal(a)
		return string(bs)
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

// 用来生成slot的方法
type namedSlotFunc func(props map[string]interface{}) string

// 执行slot返回代码
func xSlot(injectSlotFunc namedSlotFunc, props map[string]interface{}, defaultCode string) string {
	// 如果没有传递slot 则使用默认的code
	if injectSlotFunc == nil {
		return defaultCode
	}

	return injectSlotFunc(props)
}
