package ast

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
	"github.com/robertkrimen/otto/token"
	"log"
	"strings"
)

// 生成go代码
// dataKey: 默认为options.data
func Js2Go(code string, dataKey string) (goCode string, err error) {
	// 用括号包裹的原因是让"{x: 1}"这样的语法解析成对象, 而不是label
	code = fmt.Sprintf("(%s)", code)

	p, err := parser.ParseFile(nil, "", code, 0)
	if err != nil {
		err = fmt.Errorf("GetAst err: %w, code:%s", err, code)
		return
	}

	goCode = genGoCodeByNode(p.Body[0], dataKey)
	return
}

func genGoCodeByNode(node ast.Node, dataKey string) (goCode string) {
	switch t := node.(type) {

	case *ast.ExpressionStatement:
		return genGoCodeByNode(t.Expression, dataKey)
	case *ast.Identifier:
		return fmt.Sprintf(`lookInterface(%s, "%s")`, dataKey, t.Name)
	//case ast.MemberExpression:
	//	c:= t.GetCode(dataKey)
	//	return c
	case *ast.DotExpression:
		root, keys := GetKey(t, dataKey)
		return fmt.Sprintf(`lookInterface(%s, %s)`, root, strings.Join(keys, ", "))
	//case *ast.BracketExpression:
	//
	case *ast.StringLiteral:
		// js的字符串可以用'', 但go中必须是"", 所以需要替换
		c := t.Value
		//if strings.HasPrefix(c, "'") {
		//	c = `"` + c[1:len(c)-1] + `"`
		//}
		return `"` + c + `"`
	case *ast.NumberLiteral:
		return fmt.Sprintf("%v", t.Value)
	//case ast.LogicalExpression:
	//	left := genGoCodeByNode(t.Left, dataKey)
	//	right := genGoCodeByNode(t.Right, dataKey)
	//
	//	return fmt.Sprintf(`interfaceToBool(%s) %s interfaceToBool(%s)`, left, t.Operator, right)
	case *ast.BinaryExpression:
		left := genGoCodeByNode(t.Left, dataKey)
		right := genGoCodeByNode(t.Right, dataKey)
		o := t.Operator
		switch o {
		case token.STRICT_EQUAL, token.EQUAL:
			return fmt.Sprintf(`interfaceToStr(%s) == interfaceToStr(%s)`, left, right)
		case token.NOT_EQUAL, token.STRICT_NOT_EQUAL:
			return fmt.Sprintf(`interfaceToStr(%s) != interfaceToStr(%s)`, left, right)
		case token.PLUS:
			return fmt.Sprintf(`interfaceAdd(%s, %s)`, left, right)
		case token.MINUS:
			return fmt.Sprintf(`interfaceToFloat(%s) - interfaceToFloat(%s)`, left, right)
		case token.MULTIPLY:
			return fmt.Sprintf(`interfaceToFloat(%s) * interfaceToFloat(%s)`, left, right)
		case token.SLASH:
			return fmt.Sprintf(`interfaceToFloat(%s) / interfaceToFloat(%s)`, left, right)
		}

		// 可以优化: interfaceToStr("1") 为 "1"
		return fmt.Sprintf(`interfaceToStr(%s) %s interfaceToStr(%s)`, left, o, right)
	case *ast.UnaryExpression:
		arg := genGoCodeByNode(t.Operand, dataKey)
		switch t.Operator {
		case token.NOT:
			return fmt.Sprintf(`%sinterfaceToBool(%s)`, t.Operator, arg)
		case token.MINUS:
			// -1
			return fmt.Sprintf(`%s%s`, t.Operator, arg)
		default:
			panic(fmt.Sprintf("not handle UnaryExpression: %s", t.Operator))
		}
	case *ast.ObjectLiteral:
		if len(t.Value) == 0 {
			return "nil"
		}

		// 对象, 翻译成map[string]interface{}
		var mapCode = "map[string]interface{}"
		mapCode += "{"
		for _, v := range t.Value {
			k := ""
			if v.Kind != "" {
				//k = fmt.Sprintf("interfaceToStr(%s)", genGoCodeByNode(v.Value, dataKey))
			} else {
				k = fmt.Sprintf(`"%s"`, v.Key)
			}
			valueCode := genGoCodeByNode(v.Value, dataKey)
			mapCode += fmt.Sprintf(`%s: %s,`, k, valueCode)
		}
		mapCode += "}"
		return mapCode
	case *ast.CallExpression:
		name := "x"
		args := make([]string, len(t.ArgumentList))
		for i, v := range t.ArgumentList {
			args[i] = genGoCodeByNode(v, dataKey)
		}
		return fmt.Sprintf(`interfaceToFunc(lookInterface(%s,"%s"))(%s)`, dataKey, name, strings.Join(args, ","))
	case *ast.ArrayLiteral:
		args := make([]string, len(t.Value))
		for i, v := range t.Value {
			args[i] = genGoCodeByNode(v, dataKey)
		}
		return fmt.Sprintf(`[]interface{}{%s}`, strings.Join(args, ","))
	case *ast.ConditionalExpression:
		consequent := genGoCodeByNode(t.Consequent, dataKey)
		alternate := genGoCodeByNode(t.Alternate, dataKey)
		test := genGoCodeByNode(t.Test, dataKey)

		return fmt.Sprintf(`func() interface{} {if interfaceToBool(%s){return %s};return %s}()`, test, consequent, alternate)
	default:
		//panic(t)
		bs,_:=json.Marshal(t)
		log.Panicf("%s", bs)
	}
	return
}

func GetKey(p *ast.DotExpression, dataKey string) (root string, keys []string) {

	//switch t:=p.Identifier

	bs, _ := json.MarshalIndent(p, " ", " ")
	print(string(bs))

	return "root", []string{"todo"}
}
