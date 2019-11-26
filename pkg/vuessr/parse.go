package vuessr

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Element struct {
	Text     string // 只是字
	TagName  string
	Attrs    []xml.Attr
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

	decoder := xml.NewDecoder(file)
	var stack []*Element
	var currentElement *Element

	for {
		token, err := decoder.Token()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}

		switch token := token.(type) {
		case xml.StartElement:
			stack = append(stack, &Element{
				"",
				token.Name.Local,
				token.Attr,
				[]*Element{},
			})

			break
		case xml.EndElement:
			currentNode := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if len(stack) == 0 {
				break
			}

			preNode := stack[len(stack)-1]
			preNode.Children = append(preNode.Children, currentNode)
			currentElement = preNode

			break
		case xml.CharData:
			if len(stack) == 0 {
				break
			}

			lastNode := stack[len(stack)-1]
			lastNode.Children = append(lastNode.Children, &Element{Text: string(token[:]), TagName: "__string"})
			break
		}
	}

	return currentElement, nil
}
