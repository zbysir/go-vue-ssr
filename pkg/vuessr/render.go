package vuessr

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

type Render interface {
	Render(app *App, data interface{}, slot string) string
}

type Elementer interface {
	Set(attrs []xml.Attr)
}

type App struct {
	Components map[string]Render // name=>node
}

// 组件渲染,
// 如果该组件被components注册, 则使用Element渲染.
func (e *Element) Render(app *App, data interface{}, slot string) string {
	// 节点是slot, 则应该填充传递进来的slot
	if e.TagName == "slot" {
		return slot
	}

	// 如果只是文字, 则直接返回文字
	if e.Text != "" {

		//
		return injectVal(e.Text, data)
	}

	mySlot := ""

	if len(e.Children) != 0 {
		for _, v := range e.Children {
			sr := v.Render(app, data, slot)
			mySlot += sr
		}
	}

	var currTag = ""
	custom, exist := app.Components[e.TagName]
	if exist {
		bind := getBind(e.Attrs)
		// 从bind中读取数据, 做为子级数据
		childData := map[string]interface{}{}

		m, ok := data.(map[string]interface{})
		if ok {
			for k, v := range bind {
				childData[k] = m[v]
			}
		}

		currTag = custom.Render(app, childData, mySlot)
	} else if e.TagName == "template" {
		currTag = mySlot
	} else {
		// 内联元素
		currTag = fmt.Sprintf("<%s>%s</%s>", e.TagName, mySlot, e.TagName)
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
		Components: map[string]Render{
			"text": texTElement,
		},
	}
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

func injectVal(src string, data interface{}) (to string) {
	reg := regexp.MustCompile(`\{\{.+?\}\}`)

	src = reg.ReplaceAllStringFunc(src, func(s string) string {
		key := s[2 : len(s)-2]

		desc, ok := lookInterface(key, data)
		if ok {
			return fmt.Sprintf("%v", desc)
		}
		return ""
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
