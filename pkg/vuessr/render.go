package vuessr

import (
	"fmt"
)

func RenderNode(node interface{}) (str string, err error) {
	switch n := node.(type) {
	case *Element:
		ch := ""
		if len(n.Children) != 0 {
			for _, v := range n.Children {
				sr, e := RenderNode(v)
				if e != nil {
					err = e
					return
				}
				ch += sr
			}
		}
		str = fmt.Sprintf("<%s>%s</%s>", n.TagName, ch, n.TagName)
	case string:
		str = n
	default:
		panic(n)
	}

	return
}
