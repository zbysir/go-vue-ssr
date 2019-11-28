package vuessr

import (
	"fmt"
	"github.com/bysir-zl/vue-ssr/pkg/vuessr/ast_from_api"
	"regexp"
	"strings"
)

type Render interface {
	Render(app *App, data interface{}, slot string) string
	RenderFunc(app *App, slot string) (function string)
}

type App struct {
	Components map[string]struct{} // name=>node
}

// 用来生成Option代码所需要的数据
type OptionsGen struct {
	Props           map[string]string // 上级传递的 数据(包含了class和style)
	Attrs           map[string]string // 上级传递的 静态的attrs (除去class和style), 只会作用在root节点
	Class           []string          // 静态class
	Style           map[string]string // 静态style
	StyleKeys       []string          // 样式的key, 用来保证顺序
	Slot            map[string]string // 插槽节点
	DefaultSlotCode string            // 子节点code, 用于默认的插槽
	NamedSlotCode   map[string]string // 具名插槽
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
	for k, v := range m {
		c += fmt.Sprintf(`"%s": "%s",`, k, v)
	}
	c += "}"

	return c
}

func mapCodeToGoCode(m map[string]string, valueType string) string {
	c := "map[string]" + valueType
	c += "{"
	for k, v := range m {
		c += fmt.Sprintf(`"%s": %s,`, k, v)
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

func (o *OptionsGen) ToGoCode() string {
	c := "&Options{"
	if len(o.Props) != 0 {
		props := "Props: "
		props += "map[string]interface{}"
		props += "{"
		for k, v := range o.Props {
			valueCode, err := ast_from_api.JsCode2Go(v, DataKey)
			if err != nil {
				panic(err)
			}
			props += fmt.Sprintf(`"%s": %s,`, k, valueCode)
		}
		props += "}"
		c += props + ",\n"
	}
	if len(o.Attrs) != 0 {
		c += fmt.Sprintf("Attrs: %s,\n", mapStringToGoCode(o.Attrs))
	}
	if len(o.Class) != 0 {
		c += fmt.Sprintf("Class: %s,\n", sliceToGoCode(o.Class))
	}
	if len(o.Style) != 0 {
		c += fmt.Sprintf("Style: %s,\n", mapStringToGoCode(o.Style))
	}
	if len(o.StyleKeys) != 0 {
		c += fmt.Sprintf("StyleKeys: %s,\n", sliceToGoCode(o.StyleKeys))
	}
	slot := map[string]string{}

	children := o.DefaultSlotCode
	if children == "" {
		children = `""`
	}
	slot["default"] = fmt.Sprintf(`func (props map[string]interface{})string{return %s}`, children)

	for k, v := range o.NamedSlotCode {
		slot[k] = v
	}
	c += fmt.Sprintf("Slot: %s,\n", mapCodeToGoCode(slot, "namedSlotFunc"))
	c += fmt.Sprintf("P: options,\n")

	c += "}"
	return c
}

// 生成代码中的key
const (
	DataKey  = "data" // 变量作用域的key, 相当于js的this.
	PropsKey = "options.Props"
	SlotKey  = "options.Slot"
)

// 组件渲染,
// 如果该组件被components注册, 则使用Element渲染.
// todo 如果将slot改为map[string]string应该就可以实现多个slot.
//
// 用节点直接渲染可能会有的性能问题:
// - 处理文字时会使用正则来匹配{{变量, 会消耗过多性能
// - 很多没有变量的节点可以被预先处理成字符串, 就不会走递归流程
//

// 每个组件都是一个func或者是一个字符串
// slot: 子级代码
func (e *VueElement) RenderFunc(app *App) (code string, namedSlotCode map[string]string) {
	var eleCode = ""

	defaultSlotCode := ""

	namedSlotCode = map[string]string{}
	if len(e.Children) != 0 {
		for _, v := range e.Children {
			childCode, childNamedSlotCode := v.RenderFunc(app)
			if defaultSlotCode == "" {
				defaultSlotCode += childCode
			} else {
				defaultSlotCode += "+" + childCode
			}

			for k, v := range childNamedSlotCode {
				namedSlotCode[k] = v
			}
		}
	}

	if defaultSlotCode == "" {
		defaultSlotCode = `""`
	}

	// 调用组件
	_, exist := app.Components[e.TagName]
	if exist {
		options := OptionsGen{
			StyleKeys:       e.StyleKeys,
			Class:           e.Class,
			Attrs:           e.Attrs,
			Props:           e.Props,
			Style:           e.Style,
			DefaultSlotCode: defaultSlotCode,
			NamedSlotCode:   namedSlotCode,
		}
		optionsCode := options.ToGoCode()
		eleCode = fmt.Sprintf("XComponent_%s(%s)", e.TagName, optionsCode)
	} else if e.TagName == "template" {
		// 使用子级
		eleCode = defaultSlotCode
	} else if e.TagName == "__string" {
		// 纯字符串节点
		text := strings.Replace(e.Text, "\n", `\n`, -1)
		// 处理变量
		text = injectVal(text)
		eleCode = fmt.Sprintf(`"%s"`, text)
	} else {
		attrs := genAttr(e)
		attrs = injectVal(attrs)
		// attr: 如果是root元素, 则还需要处理上层传递而来的style/class
		// 内联元素, slot应该放在标签里
		eleCode = fmt.Sprintf(`"<%s"+%s+">"+%s+"</%s>"`, e.TagName, attrs, defaultSlotCode, e.TagName)
	}

	// 处理指令 如v-for

	eleCode, namedSlotCode2 := e.Directives.Exec(e, eleCode)
	for i, v := range namedSlotCode2 {
		namedSlotCode[i] = v
	}

	return eleCode, namedSlotCode
}

// 转义字符串中的", 以免和go代码中的"冲突
func encodeString(src string) string {
	return strings.Replace(src, `"`, `\"`, -1)
}

func NewApp() *App {
	return &App{
		Components: map[string]struct{}{
			"component": {},
			"slot": {},
		},
	}
}

func (a *App) Component(name string) {
	a.Components[name] = struct {
	}{}
}

// 处理 {{}} 变量
func injectVal(src string) (to string) {
	reg := regexp.MustCompile(`\{\{.+?\}\}`)

	src = reg.ReplaceAllStringFunc(src, func(s string) string {
		key := s[2 : len(s)-2]

		goCode, err := ast_from_api.JsCode2Go(key, DataKey)
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf(`"+interfaceToStr(%s)+"`, goCode)
	})

	return src
}
