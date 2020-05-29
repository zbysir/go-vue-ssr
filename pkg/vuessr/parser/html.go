package parser

import (
	"github.com/zbysir/go-vue-ssr/internal/pkg/html"
)

type HtmlParser interface {
	Parse(html string) (es []*Element, err error)
}

type Element struct {
	NodeType NodeType
	TagName  string // 节点类型: html基础节点如div/span/input, 也可能是自定义组件
	Text     string // 字节点的值
	DocType  string // 特殊的docType值
	Attrs    []html.Attribute
	Children []*Element
}

type NodeType int

const (
	TextNode NodeType = iota + 1
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

// 是否只有一个子节点
// 和正常html不同的是, vue中串联的v-if/v-else只算一个节点
func (e *Element) hasOnlyOneChildren() bool {
	c := 0
	for _, v := range e.Children {
		hasIf := false
		hasElse := false
		hasElseIf := false
		for _, a := range v.Attrs {
			if a.Key == "v-if" {
				hasIf = true
				break
			} else if a.Key == "v-else" {
				hasElse = true
				break
			} else if a.Key == "v-else-if" {
				hasElseIf = true
				break
			}
		}

		if hasIf {
			c++
		} else if hasElse || hasElseIf {
		} else {
			c++
		}
	}

	return c == 1
}

