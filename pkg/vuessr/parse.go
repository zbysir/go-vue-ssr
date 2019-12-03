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
