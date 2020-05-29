package parser

import (
	"fmt"
	"github.com/zbysir/go-vue-ssr/internal/pkg/html"
	"github.com/zbysir/go-vue-ssr/internal/pkg/html/atom"
	"os"
	"strings"
)

// GoHtml 使用go原生html库解析html字符串
// 目前有的问题:
// - 不支持不规则的html, 如在<select>里嵌套<slot>, 如在<head>里嵌套<div>/<template>
// 还在寻求另一个解决方案.
type GoHtml struct {
}

func (g GoHtml) Parse(html string) (es []*Element, err error) {
	return parseHtml(html)
}

// parse HTML
func parseHtml(filename string) (es []*Element, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	var nodes []*html.Node

	// 两个情况: 一种是<template>开头的 则是标准的vue组件, 一种vue组件如html页面. 但为了简化流程, html页面也可以被当为vue组件来渲染.
	peekWant := "<template"
	peek := make([]byte, len(peekWant))
	_, err = file.Read(peek)
	if err != nil {
		return
	}
	_, _ = file.Seek(0, 0)

	if string(peek) == peekWant {
		root := &html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Div,
			Data:     atom.Div.String(),
		}
		nodes, err = html.ParseFragment(file, root)
		if err != nil {
			return
		}
	} else {
		var node *html.Node
		node, err = html.Parse(file)
		if err != nil {
			return
		}
		if node.Type == html.DocumentNode {
			c := node.FirstChild
			for c != nil {
				nodes = append(nodes, c)
				c = c.NextSibling
			}
		} else {
			err = fmt.Errorf("bad nodeType: %d, want DocumentNode", node.Type)
			return
		}
	}

	es = hNodeToElement(nodes)
	return
}

func hNodeToElement(nodes []*html.Node) []*Element {
	var es []*Element
	for _, node := range nodes {
		var e Element
		omitNode := false
		switch node.Type {
		case html.TextNode:
			// html中多个空格和\n都应该被替换为空格
			// 注意 <script> 中的节点不应该别替换
			// 注意下面的实现方式有bug, 没有处理在<script>中的情况
			//text = strings.ReplaceAll(node.Data, "\n", " ")
			//reg := regexp.MustCompile(`\s+`)
			//text = reg.ReplaceAllString(text, " ")

			// 忽略空节点
			if strings.Trim(node.Data, "\n ") == "" {
				omitNode = true
				break
			}
			e = Element{
				NodeType: TextNode,
				Text:     node.Data,
			}
		case html.DocumentNode:
			e = Element{
				NodeType: DocumentNode,
				TagName:  "document",
			}
		case html.ElementNode:
			e = Element{
				NodeType: ElementNode,
				TagName:  node.Data,
			}
		case html.CommentNode:
			omitNode = true
		case html.DoctypeNode:
			e = Element{
				NodeType: DoctypeNode,
				DocType:  node.Data,
			}
		default:
			panic(uint32(node.Type))
		}

		if omitNode {
			continue
		}

		var children []*Element
		if node.FirstChild != nil {
			c := node.FirstChild
			var allC []*html.Node
			for c != nil {
				allC = append(allC, c)
				c = c.NextSibling
			}

			children = hNodeToElement(allC)
		}

		e.Children = children
		e.Attrs = node.Attr

		es = append(es, &e)
	}
	return es
}
