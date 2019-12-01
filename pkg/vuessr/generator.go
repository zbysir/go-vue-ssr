package vuessr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 用来生成模板字符串代码
// 目的是为了解决递归渲染节点造成的性能问题, 不过这是一个难题, 先尝试, 不行就算了.

func genComponentRenderFunc(app *App, pkgName, name string, file string) string {
	ve, err := ParseVue(file)
	if err != nil {
		panic(err)
	}
	code, _ := ve.RenderFunc(app)

	// 处理多余的纯字符串拼接: "a"+"b" => "ab"
	//code = strings.Replace(code, `"+"`, "", -1)

	return fmt.Sprintf("package %s\n\n"+
		"func XComponent_%s(options *Options)string{\n"+
		"%s:= %s\n_ = %s\n"+
		"return %s"+
		"}", pkgName, name, DataKey, PropsKey, DataKey, code)
}

func tuoFeng2SheXing(src []byte) (out []byte) {
	l := len(src)
	out = []byte{}
	for i := 0; i < l; i = i + 1 {
		// 大写变小写
		if 97-32 <= src[i] && src[i] <= 122-32 {
			if i != 0 {
				out = append(out, '-')
			}
			out = append(out, src[i]+32)
		} else {
			out = append(out, src[i])
		}
	}

	return
}

func genRegister(app *App, pkgName string) string {
	m := map[string]string{}
	for k := range app.Components {
		m[k] = fmt.Sprintf(`XComponent_%s`, k)
		k2 := string(tuoFeng2SheXing([]byte(k)))
		if k != k2 {
			m[k2] = fmt.Sprintf(`XComponent_%s`, k)
		}
	}

	return fmt.Sprintf("package %s\n\n"+
		"var components = map[string]ComponentFunc{}\n"+
		"func init(){components = %s}",
		pkgName, mapCodeToGoCode(m, "ComponentFunc"))
}

// 生成并写入文件夹
func GenAllFile(src, desc string) (err error) {
	// 生成文件夹
	err = os.MkdirAll(desc, os.ModePerm)
	if err != nil {
		return
	}

	// 删除老的.vue.go文件

	del, err := walkDir(desc, ".vue.go")
	if err != nil {
		return
	}

	for _, v := range del {
		err = os.Remove(v)
		if err != nil {
			return
		}
	}

	// 生成新的
	vue, err := walkDir(src, ".vue")
	if err != nil {
		return
	}

	var components []string

	app := NewApp()

	for _, v := range vue {
		_, fileName := filepath.Split(v)
		name := strings.Split(fileName, ".")[0]
		app.Component(name)

		components = append(components, name)
	}

	_, pkgName := filepath.Split(desc)

	// 注册vue组件, 用于动态组件
	code := genRegister(app, pkgName)
	err = ioutil.WriteFile(desc+string(os.PathSeparator)+"register.go", []byte(code), 0666)
	if err != nil {
		return
	}

	// 生成vue组件
	for _, v := range vue {
		_, fileName := filepath.Split(v)
		name := strings.Split(fileName, ".")[0]

		code := genComponentRenderFunc(app, pkgName, name, v)
		err = ioutil.WriteFile(desc+string(os.PathSeparator)+fileName+".go", []byte(code), 0666)
		if err != nil {
			return
		}
	}

	// buildin代码
	code = fmt.Sprintf("package %s\n", pkgName) + buildInCode
	err = ioutil.WriteFile(desc+string(os.PathSeparator)+"build.go", []byte(code), 0666)
	if err != nil {
		return
	}

	return
}

func walkDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		//遍历目录
		if err != nil {
			return err
		}
		if fi.IsDir() {
			// 忽略目录
			return nil
		}

		if strings.HasSuffix(filename, suffix) {
			files = append(files, filename)
		}

		return nil
	})

	return
}

const buildInCode = `
import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// 内置组件
func XComponent_slot(options *Options) string {
	name := options.Attrs["name"]
	if name == "" {
		name = "default"
	}
	props := options.Props
	injectSlotFunc := options.P.Slot[name]

	// 如果没有传递slot 则使用默认的code
	if injectSlotFunc == nil {
		return options.Slot["default"](nil)
	}

	return injectSlotFunc(props)
}

func XComponent_component(options *Options) string {
	is, ok := options.Props["is"].(string)
	if !ok {
		return ""
	}
	if c, ok := components[is]; ok {
		return c(options)
	}

	return fmt.Sprintf("not register com: %s", is)
}

// 渲染组件需要的结构
type Options struct {
	Props     map[string]interface{}   // 上级传递的 数据(包含了class和style)
	Attrs     map[string]string        // 上级传递的 静态的attrs (除去class和style), 只会作用在root节点
	Class     []string                 // 静态class, 只会作用在root节点
	Style     map[string]string        // 静态style, 只会作用在root节点
	StyleKeys []string                 // 样式的key, 用来保证顺序, 只会作用在root节点
	Slot      map[string]namedSlotFunc // 插槽代码, 支持多个不同名字的插槽, 如果没有名字则是"default"
	P         *Options                 // 父级options, 在渲染插槽会用到. (根据name取到父级的slot)
}

type Props map[string]interface{}

func (p Props) CanBeAttr() Props {
	html := map[string]struct{}{
		"id":  {},
		"src": {},
	}

	a := Props{}
	for k, v := range p {
		if _, ok := html[k]; ok {
			a[k] = v
			continue
		}

		if strings.HasPrefix(k, "data-") {
			a[k] = v
			continue
		}
	}
	return a
}

// 组件的render函数
type ComponentFunc func(options *Options) string

// 用来生成slot的方法
// 由于slot具有自己的作用域, 所以只能使用闭包实现(而不是字符串).
type namedSlotFunc func(props map[string]interface{}) string

// 混合动态和静态的标签, 主要是style/class需要混合
// todo) 如果style/class没有冲突, 则还可以优化
// tip: 纯静态的class应该在编译时期就生成字符串, 而不应调用这个
// classProps: 支持 obj, array, string
func mixinClass(options *Options, staticClass []string, classProps interface{}) (str string) {
	var class []string
	// 静态
	for _, c := range staticClass {
		if c != "" {
			class = append(class, c)
		}
	}

	// 本身的props
	for _, c := range getClassFromProps(classProps) {
		class = append(class, c)
	}

	if options != nil {
		// 上层传递的props
		if options.Props != nil {
			prop, ok := options.Props["class"]
			if ok {
				for _, c := range getClassFromProps(prop) {
					class = append(class, c)
				}
			}
		}

		// 上层传递的静态class
		if len(options.Class) != 0 {
			for _, c := range options.Class {
				if c != "" {
					class = append(class, c)
				}
			}
		}
	}

	if len(class) != 0 {
		str = fmt.Sprintf(" class=\"%s\"", strings.Join(class, " "))
	}

	return
}

// 构建style, 生成如style="color: red"的代码, 如果style代码为空 则只会返回空字符串
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
		str = fmt.Sprintf(" style=\"%s\"", styleCode)
	}

	return
}

// 生成除了style和class的attr
func mixinAttr(options *Options, staticAttr map[string]string, propsAttr map[string]interface{}) string {
	attrs := map[string]string{}

	// 静态
	for k, v := range staticAttr {
		attrs[k] = v
	}

	// 当前props
	ps := getStyleFromProps(propsAttr)
	for k, v := range ps {
		attrs[k] = v
	}

	if options != nil {
		// 上层传递的props
		if options.Props != nil {
			for k, v := range (Props(options.Props)).CanBeAttr() {
				attrs[k] = fmt.Sprintf("%v", v)
			}
		}

		// 上层传递的静态style
		for k, v := range options.Attrs {
			attrs[k] = v
		}
	}

	c := genAttr(attrs)
	if c == "" {
		return ""
	}

	return fmt.Sprintf(" %s", c)
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

func genAttr(attr map[string]string) string {
	sortedKeys := getSortedKey(attr)

	st := ""
	for _, k := range sortedKeys {
		v := attr[k]
		st += fmt.Sprintf("%s=\"%s\" ", k, v)
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
func getClassFromProps(classProps interface{}) []string {
	if classProps == nil {
		return nil
	}
	switch t := classProps.(type) {
	case []string:
		return t
	case string:
		return []string{t}
	case map[string]interface{}:
		var c []string
		for k, v := range t {
			if interfaceToBool(v) {
				c = append(c, k)
			}
		}

		return c
	}

	return nil
}

func lookInterface(data interface{}, key string) (desc interface{}) {
	m, ok := shouldLookInterface(data, key)
	if !ok {
		return ""
	}

	return m
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
	default:
		return true
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
`
