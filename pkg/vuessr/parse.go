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
	Children []Element
}

func H(filename string) (*Element, error) {
	file, err := os.Open(filename)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

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
				[]Element{},
			})

			break
		case xml.EndElement:
			currentNode := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if len(stack) == 0 {
				break
			}

			preNode := stack[len(stack)-1]
			preNode.Children = append(preNode.Children, *currentNode)
			currentElement = preNode

			break
		case xml.CharData:
			if len(stack) == 0 {
				break
			}

			lastNode := stack[len(stack)-1]
			lastNode.Children = append(lastNode.Children, Element{Text: string(token[:]), TagName: "__string"})
			break
		}
	}

	return currentElement, nil
}
