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
func (e *VueElement) RenderFunc(app *App) (code string) {
	var eleCode = ""

	mySlotCode := ""

	if len(e.Children) != 0 {
		for _, v := range e.Children {
			sr := v.RenderFunc(app)
			if mySlotCode == "" {
				mySlotCode += sr
			} else {
				mySlotCode += "+" + sr
			}
		}
	}

	if mySlotCode == "" {
		mySlotCode = `""`
	}

	// 调用方法
	_, exist := app.Components[e.TagName]
	if exist {
		// 从bind中读取数据, 做为子级数据
		bind := e.Props
		var dataInjectCode = "map[string]interface{}"
		dataInjectCode += "{"
		for k, v := range bind {
			valueCode, err := ast_from_api.JsCode2Go(v)
			if err != nil {
				panic(err)
			}
			dataInjectCode += fmt.Sprintf(`"%s": %s,`, k, valueCode)
		}
		dataInjectCode += "}"

		eleCode = fmt.Sprintf("XComponent_%s(%s, %s)", e.TagName, dataInjectCode, mySlotCode)
	} else if e.TagName == "template" {
		// 使用模板
		eleCode = mySlotCode
	} else if e.TagName == "slot" {
		eleCode = "slot"
	} else if e.TagName == "__string" {
		// 纯字符串节点
		text := strings.Replace(e.Text, "\n", `\n`, -1)
		// 处理变量
		text = injectVal(text)
		eleCode = fmt.Sprintf(`"%s"`, text)
	} else {
		attrs := genAttr(e)
		attrs = injectVal(attrs)
		// 内联元素, slot应该放在标签里
		eleCode = fmt.Sprintf(`"<%s %s>"+%s+"</%s>"`, e.TagName, encodeString(attrs), mySlotCode, e.TagName)
	}

	// 处理指令 如v-for
	eleCode = e.Directives.Exec(e, eleCode)

	return eleCode
}

// 转义字符串中的", 以免和go代码中的"冲突
func encodeString(src string) string {
	return strings.Replace(src, `"`, `\"`, -1)
}

func NewApp() *App {
	return &App{
		Components: map[string]struct{}{},
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

		goCode, err := ast_from_api.JsCode2Go(key)
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf(`"+interfaceToStr(%s)+"`, goCode)
	})

	return src
}
