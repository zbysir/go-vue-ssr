package vuessr

import (
	"strings"
)

type VueElement struct {
	TagName    string
	Text       string
	Attrs      map[string]string // 除去指令/props/style/class之外的属性
	Directives Directives
	Class      []string
	Style      map[string]string
	Props      map[string]string
	Children   []*VueElement
}

func ParseVue(filename string) (v *VueElement, err error) {
	e, err := H(filename)
	if err != nil {
		return
	}
	v = toVueElement(e)
	return
}

func toVueElement(e *Element) *VueElement {
	props := map[string]string{}
	ds := Directives{}
	var class []string
	var style map[string]string
	attrs := map[string]string{}

	for _, v := range e.Attrs {
		if v.Name.Space == "v-bind" {
			props[v.Name.Local] = v.Value
		} else if strings.HasPrefix(v.Name.Local, "v-") {
			name := v.Name.Local
			switch name {
			case "v-for":
				ds[name] = getVForDirective(v.Value)
			case "v-if":
				ds[name] = getVIfDirective(v.Value)
			}
		} else if v.Name.Local == "class" {
			ss := strings.Split(v.Value, " ")
			for _, v := range ss {
				if v != "" {
					class = append(class, v)
				}
			}
		} else if v.Name.Local == "style" {
			// todo parse style
		} else {
			attrs[v.Name.Space+":"+v.Name.Local] = v.Value
		}
	}

	ch := make([]*VueElement, len(e.Children))
	for i, v := range e.Children {
		ch[i] = toVueElement(v)
	}
	v := &VueElement{
		TagName:    e.TagName,
		Text:       e.Text,
		Attrs:      attrs,
		Directives: ds,
		Class:      class,
		Style:      style,
		Props:      props,
		Children:   ch,
	}
	return v
}
