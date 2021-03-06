package vuessr

import (
	"fmt"
	"github.com/zbysir/go-vue-ssr/internal/pkg/log"
	"github.com/zbysir/go-vue-ssr/pkg/vuessr/ast"
	"github.com/zbysir/go-vue-ssr/pkg/vuessr/parser"
	"regexp"
	"sort"
	"strings"
)

type Compiler struct {
	// 组件的名字, 包含了驼峰/蛇形
	// 如果在编译期间遇到的tag在components中, 就会使用组件方法.
	// key是tag名字, value是驼峰
	Components map[string]string
}

type Prop struct {
	Key, Val string
}

type Props []Prop

func (p Props) Get(key string) (val string, exist bool) {
	for _, v := range p {
		if v.Key == key {
			return v.Val, true
		}
	}
	return
}

func (p *Props) Del(key string) {
	for index, k := range *p {
		if k.Key == key {
			*p = append((*p)[:index], (*p)[index+1:]...)
			break
		}
	}
}

// 用来生成Option代码所需要的数据
type OptionsGen struct {
	Props           Props             // 上级传递的 数据(包含了class和style)
	Attrs           []Attribute       // 上级传递的 静态的attrs (除去class和style), 只会作用在root节点
	Class           []string          // 静态class
	Style           map[string]string // 静态style
	Slot            map[string]string // 插槽节点
	DefaultSlotCode string            // 子节点code, 用于默认的插槽
	NamedSlotCode   map[string]string // 具名插槽
	Directives      []Directive       // 指令代码
}

func sliceStringToGoCode(m []string) string {
	if len(m) == 0 {
		return "nil"
	}
	c := strings.Join(m, `","`)
	c = fmt.Sprintf(`[]string{"%s"}`, c)
	return c
}

func mapStringToGoCode(m map[string]string) string {
	if len(m) == 0 {
		return "nil"
	}
	c := "map[string]string"
	c += "{"

	for _, k := range getSortedKey(m) {
		v := m[k]
		c += fmt.Sprintf(`"%s": "%s",`, k, v)
	}
	c += "}"

	return c
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

func mapGoCodeToCode(m map[string]string, valueType string, newLine bool) string {
	c := "map[string]" + valueType
	c += "{"
	if newLine {
		c += "\n"
	}

	for _, k := range getSortedKey(m) {
		v := m[k]
		c += fmt.Sprintf(`"%s": %s,`, k, v)
		if newLine {
			c += "\n"
		}
	}
	c += "}"

	return c
}

func sliceToGoCode(m []string) string {
	c := "[]string"
	c += "{"
	for _, v := range m {
		c += fmt.Sprintf(`"%s", `, v)
	}
	c += "}"

	return c
}

// 根据js代码生成go代码(基于js AST)
func mapJsCodeToCode(m map[string]string) string {
	if len(m) == 0 {
		return "nil"
	}
	props := "map[string]interface{}"
	props += "{"
	for _, k := range getSortedKey(m) {
		v := m[k]
		valueCode, err := ast.Js2Go(v, ScopeKey)
		if err != nil {
			log.Panicf("%v, %s", err, v)
		}
		props += fmt.Sprintf(`"%s": %s,`, k, valueCode)
	}
	props += "}"

	return props
}

// 生成Options代码
func (o *OptionsGen) ToGoCode() string {
	c := "&Options{\n"

	if len(o.Props) != 0 {
		// class
		classJs, ok := o.Props.Get("class")
		if ok {
			o.Props.Del("class")
			cCode := genPropsClassCode(classJs)
			c += fmt.Sprintf("PropsClass: %s, \n", cCode)
		}
		styleJs, ok := o.Props.Get("style")
		// style
		if ok {
			o.Props.Del("style")
			cStyle := genPropsStyleCode(styleJs)
			c += fmt.Sprintf("PropsStyle: %s, \n", cStyle)
		}

		// 除了class/style的props
		if len(o.Props) != 0 {
			c += fmt.Sprintf("Props: %s, \n", genProps(o.Props))
		}
	}

	if len(o.Attrs) != 0 {
		c += fmt.Sprintf("Attrs: %s,\n", genAttrsCode(o.Attrs))
	}
	if len(o.Class) != 0 {
		c += fmt.Sprintf("Class: %s,\n", sliceToGoCode(o.Class))
	}
	if len(o.Style) != 0 {
		c += fmt.Sprintf("Style: %s,\n", mapStringToGoCode(o.Style))
	}

	// slot
	slot := map[string]string{}

	children := o.DefaultSlotCode
	if children != `""` {
		slot["default"] = fmt.Sprintf(`func(w Writer, props Props){
%s
}`, children)
	}

	for k, v := range o.NamedSlotCode {
		slot[k] = v
	}
	c += fmt.Sprintf("Slots: %s,\n", mapGoCodeToCode(slot, "NamedSlotFunc", false))

	// p 父级option
	c += fmt.Sprintf("P: options,\n")

	// directive
	if len(o.Directives) != 0 {
		// 数组
		dir := "[]directive{\n"
		for _, v := range o.Directives {
			valueCode := "nil"
			if v.Value != "" {
				var err error
				valueCode, err = ast.Js2Go(v.Value, ScopeKey)
				if err != nil {
					panic(err)
				}
			}
			dir += fmt.Sprintf("{Name: \"%s\", Value: %s, Arg: \"%s\"},\n", v.Name, valueCode, v.Arg)
		}
		dir += "}"

		c += fmt.Sprintf("Directives: %s,\n", dir)
	}

	// Scope
	c += fmt.Sprintf("Scope: %s,\n", ScopeKey)

	c += "}"
	return c
}

// 生成组件根节点所需要的Option代码
// 和ToGoCode不同的是: ToGoCodeForRoot会合并上层Options数据用于渲染当前节点,
//   实现<component :id=1>这样的写法会在组件下的根节点上生成id.
// 会合并以下数据:
//   Props: 用于生成attr
//   Attr: 用于生成attr
//   Class: 用于生成Class
//   Style: 用于生成Class
//   Directives: 在组件外层的指令并没有实用价值(无法操作Dom), 放在根节点上运行更实用.
func (o *OptionsGen) ToGoCodeForRoot() string {
	c := "&Options{\n"

	if len(o.Props) != 0 {
		// class
		classJs, ok := o.Props.Get("class")
		if ok {
			o.Props.Del("class")
			cCode := genPropsClassCode(classJs)
			c += fmt.Sprintf("PropsClass: %s, \n", cCode)
		}
		styleJs, ok := o.Props.Get("style")
		// style
		if ok {
			o.Props.Del("style")
			cStyle := genPropsStyleCode(styleJs)
			c += fmt.Sprintf("PropsStyle: %s, \n", cStyle)
		}

		// 除了class/style的props
		if len(o.Props) != 0 {
			c += fmt.Sprintf("Props: %s, \n", genProps(o.Props))
		}
	}

	if len(o.Attrs) != 0 {
		c += fmt.Sprintf("Attrs: %s,\n", genAttrsCode(o.Attrs))
	}
	if len(o.Class) != 0 {
		c += fmt.Sprintf("Class: %s,\n", sliceToGoCode(o.Class))
	}
	if len(o.Style) != 0 {
		c += fmt.Sprintf("Style: %s,\n", mapStringToGoCode(o.Style))
	}

	// slot
	slot := map[string]string{}

	children := o.DefaultSlotCode
	if children != `""` {
		slot["default"] = fmt.Sprintf(`func(w Writer, props Props){
%s
}`, children)
	}

	for k, v := range o.NamedSlotCode {
		slot[k] = v
	}
	c += fmt.Sprintf("Slots: %s,\n", mapGoCodeToCode(slot, "NamedSlotFunc", false))

	// p 父级option
	c += fmt.Sprintf("P: options,\n")

	// directive
	{
		// 数组
		dir := "append(options.Directives,\n"
		for _, v := range o.Directives {
			valueCode := "nil"
			if v.Value != "" {
				var err error
				valueCode, err = ast.Js2Go(v.Value, ScopeKey)
				if err != nil {
					panic(err)
				}
			}
			dir += fmt.Sprintf("directive{Name: \"%s\", Value: %s, Arg: \"%s\"},\n", v.Name, valueCode, v.Arg)
		}
		dir += ")"

		c += fmt.Sprintf("Directives: %s,\n", dir)
	}

	// Scope
	c += fmt.Sprintf("Scope: %s,\n", ScopeKey)

	c += "}"
	return c
}

type Code struct {
	Src  string
	Type string // string 纯字符串 / async 异步(PromiseGroup)
}

// 生成代码中的key
const (
	ScopeKey = "scope" // 变量作用域的key, 模拟js作用域.
)

// voidElements 没有子元素, 会渲染成 <br/> 这样的格式
var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"keygen": true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

// 组件渲染,
// 如果该组件被components注册, 则使用Element渲染.
//
// 用节点直接渲染可能会有的性能问题:
// - 处理文字时会使用正则来匹配{{变量, 会消耗过多性能
// - 很多没有变量的节点可以被预先处理成字符串, 就不会走递归流程
//

// 每个组件都是一个func或者是一个字符串
// slot: 子级代码
// 返回的code 是一行代码,
func (c *Compiler) GenEleCode(e *VueElement) (code string, namedSlotCode map[string]string) {
	var eleCode = ""

	defaultSlotCode := ""

	namedSlotCode = map[string]string{}
	if len(e.Children) != 0 {
		for _, v := range e.Children {
			// 跳过生成else节点的代码, 真正生成else节点的代码在if节点中
			if v.VElse || v.VElseIf {
				continue
			}
			childCode, childNamedSlotCode := c.GenEleCode(v)
			for k, v := range childNamedSlotCode {
				namedSlotCode[k] = v
			}

			if childCode == "" {
				continue
			}
			defaultSlotCode += childCode + "\n"
		}
	}
	defaultSlotCode = strings.TrimSuffix(defaultSlotCode, "\n")

	switch e.NodeType {
	case parser.TextNode:
		// 纯字符串节点
		// 将文本处理成go代码的字符串写法: "xxx"
		// 注意{{表达式中的"不应该被处理, 因为这是js代码, 需要解析成为JS AST.
		text := safeStringCode(e.Text)
		// 处理变量
		text = injectVal(text)
		eleCode = fmt.Sprintf(`w.WriteString(%s)`, text)
	case parser.DocumentNode:
		log.Infof("DocumentNode %+v", e)
	case parser.ElementNode:
		// 判断是否是自定义组件
		componentName, exist := c.Components[e.TagName]
		if exist {
			options := OptionsGen{
				Class:           e.Class,
				Attrs:           e.Attrs,
				Props:           e.Props,
				Style:           e.Style,
				DefaultSlotCode: defaultSlotCode,
				NamedSlotCode:   namedSlotCode,
				Directives:      e.Directives,
			}
			optionsCode := options.ToGoCode()
			eleCode = fmt.Sprintf("xx_%s(r, w, %s)", componentName, optionsCode)
		} else if e.TagName == "component" || e.TagName == "slot" || e.TagName == "async" {
			// 自带组件
			options := OptionsGen{
				Class:           e.Class,
				Attrs:           e.Attrs,
				Props:           e.Props,
				Style:           e.Style,
				DefaultSlotCode: defaultSlotCode,
				NamedSlotCode:   namedSlotCode,
				Directives:      e.Directives,
			}
			optionsCode := options.ToGoCode()
			eleCode = fmt.Sprintf("_%s(r, w, %s)", e.TagName, optionsCode)
		} else if e.TagName == "template" {
			// template和其他自带组件不一样: 它可以包含额外多个功能: 使用v-html/v-text
			children := defaultSlotCode
			if e.VHtml != "" {
				children = genVHtml(e.VHtml)
			} else if e.VText != "" {
				children = genVText(e.VText)
			}

			// 如果没有指令, 则直接输出子级
			if len(e.Directives) == 0 {
				eleCode = children
			} else {
				options := OptionsGen{
					Class:           nil, // dom相关都不需要处理
					Attrs:           nil, // dom相关都不需要处理
					Props:           e.Props,
					Style:           nil, // dom相关都不需要处理
					DefaultSlotCode: children,
					NamedSlotCode:   namedSlotCode,
					Directives:      e.Directives,
				}
				optionsCode := options.ToGoCode()
				eleCode = fmt.Sprintf("_%s(r, w, %s)", e.TagName, optionsCode)
			}

		} else {
			// 基础html标签

			// 判断节点是否是动态节点, 动态则使用r.Tag渲染节点, 否则使用字符串拼接
			// 动态节点
			// - 自定义指令: 在指令中会修改任何一个属性(class/style/attr...), 所以是动态的
			// - 组件的root节点: root节点会继承上层传递的(class/style/attr)

			// 动态节点
			if e.IsRoot || len(e.Directives) != 0 {
				children := defaultSlotCode
				if e.VHtml != "" {
					children = genVHtml(e.VHtml)
				} else if e.VText != "" {
					children = genVText(e.VText)
				}

				options := OptionsGen{
					Props:           e.Props,
					Attrs:           e.Attrs,
					Class:           e.Class,
					Style:           e.Style,
					Slot:            nil,
					DefaultSlotCode: children,
					NamedSlotCode:   namedSlotCode,
					Directives:      e.Directives,
				}

				if e.IsRoot {
					optionsCode := options.ToGoCodeForRoot()
					eleCode = fmt.Sprintf(`_tag(r, w, "%s", true, %s)`, e.TagName, optionsCode)
				} else {
					optionsCode := options.ToGoCode()
					eleCode = fmt.Sprintf(`_tag(r, w, "%s", false, %s)`, e.TagName, optionsCode)
				}

			} else {
				// 静态节点
				attrs := genAllAttrCode(e)
				children := defaultSlotCode
				if e.VHtml != "" {
					children = genVHtml(e.VHtml)
				} else if e.VText != "" {
					children = genVText(e.VText)
				}

				if children != "" {
					eleCode = fmt.Sprintf("w.WriteString(\"<%s\"+%s+\">\")\n%s\nw.WriteString(\"</%s>\")", e.TagName, attrs, children, e.TagName)
				} else {
					if voidElements[e.TagName] {
						eleCode = fmt.Sprintf("w.WriteString(\"<%s\"+%s+\"/>\")", e.TagName, attrs)
					} else {
						eleCode = fmt.Sprintf("w.WriteString(\"<%s\"+%s+\"></%s>\")", e.TagName, attrs, e.TagName)
					}
				}
			}
		}

	case parser.CommentNode:
	case parser.DoctypeNode:
		eleCode = fmt.Sprintf(`w.WriteString("<!doctype %s>")`, e.DocType)
	default:
		panic(fmt.Sprintf("bad nodeType, %+v", e))
	}

	// 优先级 vSlot > vFor > vIf, 所以先处理VIf(后处理的可覆盖前处理的)

	if e.VIf != nil {
		var namedSlotCodeElseIf map[string]string
		eleCode, namedSlotCodeElseIf = genVIf(e.VIf, eleCode, c)
		for i, v := range namedSlotCodeElseIf {
			namedSlotCode[i] = v
		}
	}
	if e.VFor != nil {
		eleCode = genVFor(e.VFor, eleCode)
	}
	if e.VSlot != nil {
		var namedSlotCode2 map[string]string
		eleCode, namedSlotCode2 = genVSlot(e.VSlot, eleCode)
		for i, v := range namedSlotCode2 {
			namedSlotCode[i] = v
		}
	}

	return eleCode, namedSlotCode
}

// vIf处理if节点与elseif/else节点, 会返回elseif节点的namedSlotCode
func genVIf(e *VIf, srcCode string, c *Compiler) (code string, namedSlotCode map[string]string) {
	// 自己的conditions
	condition, err := ast.Js2Go(e.Condition, ScopeKey)
	if err != nil {
		panic(err)
	}
	namedSlotCode = map[string]string{}

	// open if
	code = fmt.Sprintf(`
if interfaceToBool(%s) { %s`, condition, srcCode)
	// 继续处理else节点
	for _, v := range e.ElseIf {
		eleCode, namedSlotCode2 := c.GenEleCode(v.VueElement)
		for k, v := range namedSlotCode2 {
			namedSlotCode[k] = v
		}
		switch v.Types {
		case "else":
			code += fmt.Sprintf(`} else { %s`, eleCode)
		case "elseif":
			condition, err := ast.Js2Go(v.Condition, ScopeKey)
			if err != nil {
				panic(err)
			}
			code += fmt.Sprintf(`} else if interfaceToBool(%s) { %s`, condition, eleCode)
		}
	}

	// close if
	code += `
}`
	return
}

func genVSlot(e *VSlot, srcCode string) (code string, namedSlotCode map[string]string) {
	namedSlotCode = map[string]string{
		e.SlotName: fmt.Sprintf(`func(w Writer, props Props){
	%s := extendScope(%s, map[string]interface{}{"%s": props})
_ = %s
%s
}`, ScopeKey, ScopeKey, e.PropsKey, ScopeKey, srcCode),
	}

	// 插槽会将原来的子代码去掉, 并将代码放在namedSlot里.
	code = `""`
	return
}

func genVFor(e *VFor, srcCode string) (code string) {
	vfArray := e.ArrayKey
	vfItem := e.ItemKey
	vfIndex := e.IndexKey
	vfArrayCode, err := ast.Js2Go(vfArray, ScopeKey)
	if err != nil {
		panic(err)
	}

	// 将自己for, 将子代码的data字段覆盖, 实现作用域的修改
	return fmt.Sprintf(`
  for index, item := range interface2Slice(%s) {
    func(xscope *Scope){
        %s := extendScope(xscope, map[string]interface{}{
          "%s": index,
          "%s": item,
        })
		_ = %s
		%s
    }(%s)
  }
`, vfArrayCode, ScopeKey, vfIndex, vfItem, ScopeKey, srcCode, ScopeKey)
}

func genVHtml(value string) (code string) {
	goCode, err := ast.Js2Go(value, ScopeKey)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`w.WriteString(interfaceToStr(%s))`, goCode)
}

func genVText(value string) (code string) {
	goCode, err := ast.Js2Go(value, ScopeKey)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`w.WriteString(interfaceToStr(%s, true))`, goCode)
}

func NewCompiler() *Compiler {
	return &Compiler{
		Components: map[string]string{},
	}
}

func (a *Compiler) AddComponent(name string) {
	// 蛇形
	tagName := tuoFeng2SheXing(name)
	// 驼峰
	compName := sheXing2TuoFeng(name)
	a.Components[tagName] = compName
	a.Components[compName] = compName
}

// 处理 Mustache {{}} 插值
// 生成代码（字符串类型）, .e.g: "123" + interfaceToStr(scope.Get("total"),true)
func injectVal(src string) (to string) {
	reg := regexp.MustCompile(`{{.+?}}`)

	src = reg.ReplaceAllStringFunc(src, func(s string) string {
		key := s[2 : len(s)-2]

		goCode, err := ast.Js2Go(key, ScopeKey)
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf(`"+interfaceToStr(%s, true)+"`, goCode)
	})

	src = strings.TrimPrefix(src, `""+`)
	src = strings.TrimSuffix(src, `+""`)
	return src
}

// 包裹字符串
// 需要处理如: 将"变为 \"
// 跳过处理{{表达式中的字符串.
func safeStringCode(s string) (to string) {
	var t strings.Builder
	for _, v := range strings.Split(s, "{{") {
		sp := strings.Split(v, "}}")
		if len(sp) == 2 {
			// 跳过处理{{表达式中的字符串.
			t.WriteString("{{")
			t.WriteString(sp[0])
			t.WriteString("}}")
			t.WriteString(strings.ReplaceAll(sp[1], `"`, `\"`))
		} else {
			t.WriteString(strings.ReplaceAll(sp[0], `"`, `\"`))
		}
	}

	to = `"` + strings.Replace(t.String(), "\n", `\n`, -1) + `"`
	return
}
