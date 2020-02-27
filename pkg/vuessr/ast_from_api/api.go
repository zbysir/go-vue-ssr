package ast_from_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zbysir/go-vue-ssr/internal/pkg/log"
	"github.com/zbysir/go-vue-ssr/internal/pkg/obj"
	"reflect"
	"strings"
)

type Node map[string]interface{}
type X interface {
}
type Program struct {
	Body []Node `json:"body"`
}

type ExpressionStatement struct {
	Expression Node `json:"expression"`
}

// 变量
type Identifier struct {
	Name string `json:"name"`
}

type LogicalExpression struct {
	Operator string `json:"operator"` // && || === ==
	Left     Node   `json:"left"`
	Right    Node   `json:"right"`
}

type BinaryExpression struct {
	Operator string `json:"operator"` // +
	Left     Node   `json:"left"`
	Right    Node   `json:"right"`
}

// 一元运算符号 ! -
type UnaryExpression struct {
	Operator string `json:"operator"` // !
	Prefix   bool   `json:"prefix"`
	Argument Node   `json:"argument"`
}

// 字面量, " " , 1, 1.1
type Literal struct {
	Value interface{} `json:"value"`
	Raw   string      `json:"raw"`
}

// 对象
type ObjectExpression struct {
	Properties []Node `json:"properties"`
}

// 对象的成员
type Property struct {
	Key      Node `json:"key"` // 一般都是Identifier
	Value    Node `json:"value"`
	Computed bool `json:"computed"`
}

// a.b.c这样的读取成员变量表达式
type MemberExpression struct {
	Object   Node `json:"object"`
	Property Node `json:"property"`
	Computed bool // Property是否变量
}

func (p Property) GetKey() string {
	if p.Computed {
		panic("Computed不能使用GetKey方法")
	}

	key := ""
	switch t := p.Key.Assert().(type) {
	case Identifier:
		key = t.Name
	case Literal:
		key = t.Value.(string)
	default:
		log.Panicf("%v, %v", t, p)
	}

	return key
}

// 解析js代码: `a.b.c.d[e]`
// 返回 keys表示读取的路径, `["a", "b", "c", "d", interfaceToStr(lookInterface(this, "e"))]`
// root表示读取的对象, 支持变量/字面量/字面对象
func (p MemberExpression) GetKey(dataKey string) (root string, keys []string) {
	currKey := ""
	switch t := p.Property.Assert().(type) {
	case Identifier:
		if p.Computed {
			currKey = fmt.Sprintf(`interfaceToStr(lookInterface(%s, "%s"))`, dataKey, t.Name)
		} else {
			currKey = fmt.Sprintf(`"%s"`, t.Name)
		}
	case Literal:
		currKey = fmt.Sprintf(`"%v"`, t.Value)
	case MemberExpression:
		currKey = fmt.Sprintf("interfaceToStr(%s)", t.GetCode(dataKey))
	default:
		currKey = fmt.Sprintf("interfaceToStr(%s)", genGoCodeByNode(p.Property, dataKey))
		//panic(t)
	}

	root = ""

	switch t := p.Object.Assert().(type) {
	case MemberExpression:
		r, k := t.GetKey(dataKey)
		root = r
		keys = append(k, currKey)
	case Identifier:
		// 变量则读取dataKey
		root = dataKey
		// 只支持父级是一个变量, 如a.b中的a,
		keys = []string{fmt.Sprintf(`"%s"`, t.Name), currKey}
	case Literal:
		root = genGoCodeByNode(p.Object, dataKey)
		keys = []string{currKey}
	default:
		root = genGoCodeByNode(p.Object, dataKey)
		keys = []string{currKey}
	}

	return
}

// 返回读取变量的代码: 如lookInterface(this, "a")
func (p MemberExpression) GetCode(dataKey string) (code string) {
	root, keys := p.GetKey(dataKey)
	return fmt.Sprintf(`lookInterface(%s, %s)`, root, strings.Join(keys, ", "))
}

// a.b.c这样的读取成员变量表达式
type CallExpression struct {
	Arguments []Node `json:"arguments"`
	Callee    Node   `json:"callee"`
}

func (c CallExpression) GetFuncName() string {
	switch t := c.Callee.Assert().(type) {
	case Identifier:
		return t.Name
	}
	return ""
}

type ArrayExpression struct {
	Elements []Node `json:"elements"`
}

// a?b:c
type ConditionalExpression struct {
	Test       Node `json:"test"`       // a
	Consequent Node `json:"consequent"` // b
	Alternate  Node `json:"alternate"`  // c
}

var nodeMap = map[string]interface{}{
	"Program":               Program{},
	"ExpressionStatement":   ExpressionStatement{},
	"BinaryExpression":      BinaryExpression{},
	"LogicalExpression":     LogicalExpression{},
	"Identifier":            Identifier{},
	"UnaryExpression":       UnaryExpression{},
	"Literal":               Literal{},
	"ObjectExpression":      ObjectExpression{},
	"Property":              Property{},
	"MemberExpression":      MemberExpression{},
	"CallExpression":        CallExpression{},
	"ArrayExpression":       ArrayExpression{},
	"ConditionalExpression": ConditionalExpression{},
}

func (n Node) Assert() interface{} {
	typ, ok := n["type"].(string)
	if !ok {
		return nil
	}
	entity, ok := nodeMap[typ]
	if !ok {
		log.Errorf("unhand type:%s, %+v", typ, n)
		return nil
	}

	vPoint := reflect.New(reflect.TypeOf(entity))
	_, err := obj.Copy(n, vPoint.Interface())
	if err != nil {
		log.Error(err)
		return nil
	}

	return vPoint.Elem().Interface()
}

func GetAST(code string) (node Node, err error) {
	bs, _ := json.Marshal(map[string]string{
		"code": code,
	})
	status, res, err := client.Post("", bs, nil, &node)
	if err != nil {
		return
	}

	if status != 200 {
		err = errors.New(string(res))
		return
	}
	return
}
