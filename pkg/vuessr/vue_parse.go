package vuessr

import (
	"strings"
)

type VueElement struct {
	IsRoot     bool // 是否是根节点, 指的是<template>下一级节点, 这个节点会继承父级传递下来的class/style
	TagName    string
	Text       string
	Attrs      map[string]string // 除去指令/props/style/class之外的属性
	AttrsKeys  []string          // 属性的key, 用来保证顺序
	Directives Directives
	Class      []string          // 静态class
	Style      map[string]string // 静态style
	StyleKeys  []string          // 样式的key, 用来保证顺序
	Props      Props             // props, 包括动态的class和style
	Children   []*VueElement
}

type Props map[string]string

func (p Props) Get(key string) string {
	if p == nil {
		return ""
	}
	return p[key]
}

func (p Props) Omit(key ...string) Props {
	kMap := map[string]struct{}{}
	for _, k := range key {
		kMap[k] = struct{}{}
	}

	a := Props{}
	for k, v := range p {
		if _, ok := kMap[k]; ok {
			continue
		}
		a[k] = v
	}
	return a
}

func (p Props) Only(key ...string) Props {
	kMap := map[string]struct{}{}
	for _, k := range key {
		kMap[k] = struct{}{}
	}

	a := Props{}
	for k, v := range p {
		if _, ok := kMap[k]; !ok {
			continue
		}

		a[k] = v
	}
	return a
}

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

func ParseVue(filename string) (v *VueElement, err error) {
	e, err := H(filename)
	if err != nil {
		return
	}
	p := VueElementParser{}
	v = p.Parse(e)
	return
}

type VueElementParser struct {
	hasRoot bool // 是否已经有了root节点, 如果有了就不会再查找root节点了.
}

func (p VueElementParser) Parse(e *Element) *VueElement {
	props := map[string]string{}
	ds := Directives{}
	var class []string
	style := map[string]string{}
	var styleKeys []string
	attrs := map[string]string{}
	var attrsKeys []string

	for _, v := range e.Attrs {
		oriKey := v.Key
		ss := strings.Split(oriKey, ":")
		nameSpace := "-"
		key := oriKey
		if len(ss) == 2 {
			key = ss[1]
			nameSpace = ss[0]
		}

		if nameSpace == "v-bind" || nameSpace == "" {
			// v-bind & shorthands
			props[key] = v.Val
		} else if strings.HasPrefix(oriKey, "v-") {
			// 指令
			// v-if=""
			// v-slot:name=""
			// v-show=""
			switch {
			case key == "v-for":
				ds["v-for"] = getVForDirective(v)
			case key == "v-if":
				ds["v-if"] = getVIfDirective(v)
			case key == "v-else-if":
				ds["v-else-if"] = getVElseIfDirective(v)
			case key == "v-else":
				ds["v-else"] = getVElseDirective(v)
			case nameSpace == "v-slot":
				ds["v-slot"] = getVSlotDirective(v)
			}
		} else if v.Key == "class" {
			ss := strings.Split(v.Val, " ")
			for _, v := range ss {
				if v != "" {
					class = append(class, v)
				}
			}
		} else if v.Key == "style" {
			ss := strings.Split(v.Val, ";")
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
			key := v.Key
			if v.Namespace != "" {
				key = v.Namespace + ":" + v.Key
			}
			attrs[key] = v.Val
			attrsKeys = append(attrsKeys, key)
		}
	}

	ch := make([]*VueElement, len(e.Children))

	isRoot := false
	if !p.hasRoot {
		if e.TagName == "template" {
			isRoot = true
			p.hasRoot = true
		}
	}

	for i, v := range e.Children {
		ch[i] = p.Parse(v)
		ch[i].IsRoot = isRoot
	}

	v := &VueElement{
		TagName:    e.TagName,
		Text:       e.Text,
		Attrs:      attrs,
		AttrsKeys:  attrsKeys,
		Directives: ds,
		Class:      class,
		Style:      style,
		StyleKeys:  styleKeys,
		Props:      props,
		Children:   ch,
	}
	return v
}
