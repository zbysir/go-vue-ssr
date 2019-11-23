package vuessr

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type Render interface {
	Render(app *App, data interface{}, slot string) string
	RenderFunc(app *App, slot string) (function string)
}

type Elementer interface {
	Set(attrs []xml.Attr)
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
// dataInject:
func (e *Element) RenderFunc(app *App, slot string) (code string) {
	mySlotCode := ""

	if len(e.Children) != 0 {
		for _, v := range e.Children {
			sr := v.RenderFunc(app, slot)
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

	var currTag = ""

	// 调用方法
	_, exist := app.Components[e.TagName]
	if exist {
		// 从bind中读取数据, 做为子级数据
		bind := getBind(e.Attrs)
		var dataInjectCode = "map[string]interface{}"
		dataInjectCode += "{"
		for k, v := range bind {
			dataInjectCode += fmt.Sprintf(`"%s": lookInterface(data, "%s"),`, k, v)
		}
		dataInjectCode += "}"

		currTag = fmt.Sprintf("XComponent_%s(%s, %s)", e.TagName, dataInjectCode, mySlotCode)
	} else if e.TagName == "template" {
		// 使用模板
		currTag = mySlotCode
	} else if e.TagName == "slot" {
		if slot == "" {
			slot = `""`
		}
		currTag = "slot"
	} else if e.TagName == "__string" {
		text := strings.Replace(e.Text, "\n", `\n`, -1)
		// 处理变量
		text = injectVal(text)
		currTag = fmt.Sprintf(`"%s"`, text)
	} else {
		// 内联元素, slot应该放在标签里
		currTag = fmt.Sprintf(`"<%s>"+%s+"</%s>"`, e.TagName, mySlotCode, e.TagName)
	}

	return currTag
}

func getBind(as []xml.Attr) (bind map[string]string) {
	bind = map[string]string{}
	for _, v := range as {
		if v.Name.Space == "v-bind" {
			bind[v.Name.Local] = v.Value
		}
	}
	return
}

type Document struct {
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

func lookInterface(key string, data interface{}) (desc interface{}, exist bool) {
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

	return lookInterface(strings.Join(kk[1:], "."), c)
}

func injectVal(src string) (to string) {
	reg := regexp.MustCompile(`\{\{.+?\}\}`)

	src = reg.ReplaceAllStringFunc(src, func(s string) string {
		key := s[2 : len(s)-2]

		return fmt.Sprintf(`"+lookInterfaceToStr(data, "%s")+"`, key)
	})
	return src
}

//type Node struct {
//	Components map[string]interface{} // name=>node
//	Ctx interface{}
//}
//
//func (r *Render) renderNode(node interface{}, ctx interface{}) (str string) {
//	switch n := node.(type) {
//	case *Element:
//		ch := ""
//		if len(n.Children) != 0 {
//			for _, v := range n.Children {
//				sr := r.renderNode(v, ctx)
//				ch += sr
//			}
//		}
//
//		currNode := r.renderTag(n.TagName, ctx)
//
//		str = fmt.Sprintf(currNode, ch)
//	case string:
//		str = n
//	default:
//		panic(n)
//	}
//
//	return
//}
//
//func (r *Render) renderTag(tagName string, ctx interface{}) (h string) {
//	node, exist := r.Components[tagName]
//	if exist {
//		h = r.renderNode(node, ctx)
//	} else {
//		h = fmt.Sprintf("<%s>%%s</%s>", tagName, tagName)
//	}
//	return
//}
