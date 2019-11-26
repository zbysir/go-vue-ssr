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
	StyleKeys  []string // 样式的key, 用来保证顺序
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
	style := map[string]string{}
	var styleKeys []string
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
			ss := strings.Split(v.Value, ";")
			for _, v := range ss {
				v = strings.Trim(v, " ")
				ss := strings.Split(v, ":")
				if len(ss) != 2 {
					continue
				}
				key := strings.Trim(ss[0], " ")
				val := strings.Trim(ss[1], " ")
				style[key] = val
				styleKeys = append(styleKeys, key)
			}
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
		StyleKeys:  styleKeys,
		Props:      props,
		Children:   ch,
	}
	return v
}
