package vuessr

import (
	"github.com/bysir-zl/go-vue-ssr/pkg/vuessr/html"
	"os"
	"strings"
)

type VueElement struct {
	IsRoot    bool // 是否是根节点, 指的是<template>下一级节点, 这个节点会继承父级传递下来的class/style
	TagName   string
	Text      string
	Attrs     map[string]string // 除去指令/props/style/class之外的属性
	AttrsKeys []string          // 属性的key, 用来保证顺序
	// Directives       GenCodeDirectives // genCode指令(如v-if, v-for), 在编译期间运行
	ElseIfConditions []ElseIf          // 将与if指令匹配的elseif/else关联在一起
	Class            []string          // 静态class
	Style            map[string]string // 静态style
	StyleKeys        []string          // 样式的key, 用来保证顺序
	Props            Props             // props, 包括动态的class和style
	Children         []*VueElement     // 子节点
	VIf              *VIf              // 处理v-if需要的数据
	VFor             *VFor
	VSlot            *VSlot
	VElse            bool // 如果是VElse节点则不会生成代码(而是在vif里生成代码)
	VElseIf          bool
}

type ElseIf struct {
	Types      string // else / elseif
	Condition  string // elseif语句的condition表达式
	VueElement *VueElement
}

type VIf struct {
	Condition string // 条件表达式
	ElseIf    []*ElseIf
}

func (p *VIf) AddElseIf(v *ElseIf) {
	p.ElseIf = append(p.ElseIf, v)
}

type VFor struct {
	ArrayKey string
	ItemKey  string
	IndexKey string
}

type VSlot struct {
	SlotName string
	PropsKey string
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
		if _, ok := html[k]; !ok {
			continue
		}

		if !strings.HasPrefix(k, "data-") {
			continue
		}

		a[k] = v
	}
	return a
}


type Element struct {
	Text     string // 只是字
	TagName  string
	Attrs    []html.Attribute
	Children []*Element
}

// parse HTML
func parseHtml(filename string) (*Element, error) {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	decoder := html.NewTokenizer(file)
	var stack []*Element
	var currentElement *Element

	for {
		token := decoder.Token()
		switch token.Type {
		case html.StartTagToken:
			stack = append(stack, &Element{
				"",
				token.Data,
				token.Attr,
				[]*Element{},
			})

			break
		case html.EndTagToken:
			currentNode := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if len(stack) == 0 {
				break
			}

			preNode := stack[len(stack)-1]
			preNode.Children = append(preNode.Children, currentNode)
			currentElement = preNode

			break
		case html.SelfClosingTagToken:
			if len(stack) == 0 {
				break
			}

			preNode := stack[len(stack)-1]
			preNode.Children = append(preNode.Children, &Element{
				"",
				token.Data,
				token.Attr,
				[]*Element{},
			})
			currentElement = preNode
		case html.TextToken:
			if len(stack) == 0 {
				break
			}

			lastNode := stack[len(stack)-1]
			lastNode.Children = append(lastNode.Children, &Element{Text: string(token.Data[:]), TagName: "__string"})
			break
		}

		tp := decoder.Next()
		if tp == html.ErrorToken {
			break
		}
	}

	return currentElement, nil
}

func ParseVue(filename string) (v *VueElement, err error) {
	e, err := parseHtml(filename)
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
	vs := p.parseList([]*Element{e})
	return vs[0]
}

// 递归处理同级节点
func (p VueElementParser) parseList(es []*Element) []*VueElement {
	vs := make([]*VueElement, len(es))

	var ifVueEle *VueElement
	for i, e := range es {
		props := map[string]string{}
		//ds := GenCodeDirectives{}
		var class []string
		style := map[string]string{}
		var styleKeys []string
		attrs := map[string]string{}
		var attrsKeys []string
		var vIf *VIf
		var vFor *VFor
		var vSlot *VSlot

		// 标记节点是不是if
		var vElse *ElseIf
		var vElseIf *ElseIf

		for _, attr := range e.Attrs {
			oriKey := attr.Key
			ss := strings.Split(oriKey, ":")
			nameSpace := "-"
			key := oriKey
			if len(ss) == 2 {
				key = ss[1]
				nameSpace = ss[0]
			}

			if nameSpace == "v-bind" || nameSpace == "" {
				// v-bind & shorthands
				props[key] = attr.Val
			} else if strings.HasPrefix(oriKey, "v-") {
				// 指令
				// v-if=""
				// v-slot:name=""
				// v-show=""
				switch {
				case key == "v-for":
					val := attr.Val

					ss := strings.Split(val, " in ")
					arrayKey := strings.Trim(ss[1], " ")

					left := strings.Trim(ss[0], " ")
					var itemKey string
					var indexKey string
					// (item, index) in list
					if strings.Contains(left, ",") {
						left = strings.Trim(left, "()")
						ss := strings.Split(left, ",")
						itemKey = strings.Trim(ss[0], " ")
						indexKey = strings.Trim(ss[1], " ")
					} else {
						itemKey = left
						indexKey = "$index"
					}

					vFor = &VFor{
						ArrayKey: arrayKey,
						ItemKey:  itemKey,
						IndexKey: indexKey,
					}
				case key == "v-if":
					vIf = &VIf{
						Condition: strings.Trim(attr.Val, " "),
						ElseIf:    nil,
					}
				case nameSpace == "v-slot":
					oriKey := attr.Key
					key := oriKey
					ss := strings.Split(oriKey, ":")
					if len(ss) == 2 {
						key = ss[1]
					}
					slotName := key
					propsKey := attr.Val
					// 不应该为空, 否则可能会导致生成的go代码有误
					if propsKey == "" {
						propsKey = "slotProps"
					}
					vSlot = &VSlot{
						SlotName: slotName,
						PropsKey: propsKey,
					}
				case key == "v-else-if":
					vElseIf = &ElseIf{
						Types:     "elseif",
						Condition: strings.Trim(attr.Val, " "),
					}
				case key == "v-else":
					vElse = &ElseIf{
						Types:     "else",
						Condition: strings.Trim(attr.Val, " "),
					}
				}
			} else if attr.Key == "class" {
				ss := strings.Split(attr.Val, " ")
				for _, v := range ss {
					if v != "" {
						class = append(class, v)
					}
				}
			} else if attr.Key == "style" {
				ss := strings.Split(attr.Val, ";")
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
				key := attr.Key
				if attr.Namespace != "" {
					key = attr.Namespace + ":" + attr.Key
				}
				attrs[key] = attr.Val
				attrsKeys = append(attrsKeys, key)
			}
		}

		ch := p.parseList(e.Children)

		isRoot := false
		if !p.hasRoot {
			// 最外层的template下的节点是根节点
			if e.TagName == "template" {
				isRoot = true
				p.hasRoot = true
			}
		}

		for _, v := range ch {
			v.IsRoot = isRoot
		}

		v := &VueElement{
			IsRoot:           false,
			TagName:          e.TagName,
			Text:             e.Text,
			Attrs:            attrs,
			AttrsKeys:        attrsKeys,
			ElseIfConditions: []ElseIf{},
			//Directives:       ds,
			Class:     class,
			Style:     style,
			StyleKeys: styleKeys,
			Props:     props,
			Children:  ch,
			VIf:       vIf,
			VFor:      vFor,
			VSlot:     vSlot,
			VElse:     vElse != nil,
			VElseIf:   vElseIf != nil,
		}

		// 记录vif, 接下来的elseif将与这个节点关联
		if vIf != nil {
			ifVueEle = v
		} else if e.TagName != "__string" && vElse == nil && vElseIf == nil {
			ifVueEle = nil
		}

		if vElseIf != nil {
			if ifVueEle == nil {
				panic("v-else-if must below v-if")
			}
			vElseIf.VueElement = v
			ifVueEle.VIf.AddElseIf(vElseIf)
		}
		if vElse != nil {
			if ifVueEle == nil {
				panic("v-else must below v-if")
			}
			vElse.VueElement = v
			ifVueEle.VIf.AddElseIf(vElse)
			ifVueEle = nil
		}

		vs[i] = v
	}

	return vs
}
