package ast_from_api

import (
	"fmt"
	"strings"
)

// 生成go代码
// dataKey: 默认为options.data
func JsCode2Go(code string, dataKey string) (goCode string, err error) {
	// 用code包裹的原因是让"{x: 1}"这样的语法解析成对象, 而不是label
	code = fmt.Sprintf("(%s)", code)
	node, err := GetAST(code)
	if err != nil {
		return
	}

	goCode = genGoCodeByNode(node, dataKey)
	return
}

func genGoCodeByNode(node Node, dataKey string) (goCode string) {
	switch t := node.Assert().(type) {
	case Program:
		x := ``
		for _, v := range t.Body {
			if x != "" {
				x += `+`
			}
			x += fmt.Sprintf(`%s`, genGoCodeByNode(v, dataKey))
		}

		return x
	case ExpressionStatement:
		return genGoCodeByNode(t.Expression, dataKey)
	case Identifier:
		return fmt.Sprintf(`lookInterface(%s, "%s")`, dataKey, t.Name)
	case MemberExpression:
		return fmt.Sprintf(`lookInterface(%s, "%s")`, dataKey, t.GetKey())
	case Literal:
		// js的字符串可以用'', 但go中必须是"", 所以需要替换
		c := t.Raw
		if strings.HasPrefix(c, "'") {
			c = `"` + c[1:len(c)-1] + `"`
		}
		return c
	case LogicalExpression:
		left := genGoCodeByNode(t.Left, dataKey)
		right := genGoCodeByNode(t.Right, dataKey)
		return fmt.Sprintf(`interfaceToBool(%s) %s interfaceToBool(%s)`, left, t.Operator, right)
	case BinaryExpression:
		left := genGoCodeByNode(t.Left, dataKey)
		right := genGoCodeByNode(t.Right, dataKey)
		return fmt.Sprintf(`interfaceToStr(%s) %s interfaceToStr(%s)`, left, t.Operator, right)
	case UnaryExpression:
		arg := genGoCodeByNode(t.Argument, dataKey)
		return fmt.Sprintf(`%sinterfaceToBool(%s)`, t.Operator, arg)
	case ObjectExpression:
		if len(t.Properties) == 0 {
			return "nil"
		}

		// 对象, 翻译成map[string]interface{}
		var mapCode = "map[string]interface{}"
		mapCode += "{"
		for _, v := range t.Properties {
			p := v.Assert().(Property)
			k := p.GetKey()
			valueCode := genGoCodeByNode(p.Value, dataKey)
			mapCode += fmt.Sprintf(`"%s": %s,`, k, valueCode)
		}
		mapCode += "}"

		return mapCode
	case CallExpression:
		name := t.GetFuncName()
		args := make([]string, len(t.Arguments))
		for i, v := range t.Arguments {
			args[i] = genGoCodeByNode(v, dataKey)
		}
		return fmt.Sprintf("%s(%s)", name, strings.Join(args,","))
	default:
		panic(t)
		//bs,_:=json.Marshal(t)
		//log.Panicf("%v", t)
	}
	return
}
