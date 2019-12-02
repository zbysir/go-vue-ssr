package vuessr

import (
	"golang.org/x/net/html"
	"os"
)

type Element struct {
	Text     string // 只是字
	TagName  string
	Attrs    []html.Attribute
	Children []*Element
}

// parse HTML
func H(filename string) (*Element, error) {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	// 处理@/:缩写, 缩写不能通过xml解析
	// ps: 这个正则有点难写, 先不做

	//bs,err:=ioutil.ReadAll(file)
	//if err != nil {
	//	panic(err)
	//}
	//reg:=regexp.MustCompile(`<[\s\S]+? (:.*)=`)
	//bs = reg.ReplaceAllFunc(bs, func(src []byte) []byte {
	//	log.Infof("%s", src)
	//	return src
	//})
	//
	//var buff bytes.Buffer
	//buff.Write(bs)

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
