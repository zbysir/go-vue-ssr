package ast_from_api

import (
	"fmt"
)

// 生成go代码
func JsCode2Go(code string) (goCode string, err error) {
	// 用code包裹的原因是让"{x: 1}"这样的语法解析成对象, 而不是label
	code = fmt.Sprintf("(%s)", code)
	node, err := GetAST(code)
	if err != nil {
		return
	}

	goCode = genGoCodeByNode(node)
	return
}

func genGoCodeByNode(node Node) (goCode string) {
	switch t := node.Assert().(type) {
	case Program:
		x := ``
		for _, v := range t.Body {
			if x != "" {
				x += `+`
			}
			x += fmt.Sprintf(`%s`, genGoCodeByNode(v))
		}

		return x
	case ExpressionStatement:
		return genGoCodeByNode(t.Expression)
	case Identifier:
		return fmt.Sprintf(`lookInterface(data, "%s")`, t.Name)
	case Literal:
		return fmt.Sprintf(`%s`, t.Raw)
	case LogicalExpression:
		left := genGoCodeByNode(t.Left)
		right := genGoCodeByNode(t.Right)
		return fmt.Sprintf(`interfaceToBool(%s) %s interfaceToBool(%s)`, left, t.Operator, right)
	case BinaryExpression:
		left := genGoCodeByNode(t.Left)
		right := genGoCodeByNode(t.Right)
		return fmt.Sprintf(`interfaceToStr(%s) %s interfaceToStr(%s)`, left, t.Operator, right)
	case UnaryExpression:
		arg := genGoCodeByNode(t.Argument)
		return fmt.Sprintf(`%sinterfaceToBool(%s)`, t.Operator, arg)
	case ObjectExpression:
		// 对象, 翻译成map[string]interface{}
		var mapCode = "map[string]interface{}"
		mapCode += "{"
		for _, v := range t.Properties {
			p := v.Assert().(Property)
			k := p.GetKey()
			valueCode := genGoCodeByNode(p.Value)
			mapCode += fmt.Sprintf(`"%s": %s,`, k, valueCode)
		}
		mapCode += "}"

		return mapCode
	default:
		panic(t)
		//bs,_:=json.Marshal(t)
		//log.Panicf("%v", t)
	}
	return
}
