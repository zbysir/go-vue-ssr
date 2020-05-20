package ast

import (
	"fmt"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
	"github.com/robertkrimen/otto/token"
	"strings"
)

// 生成go代码
// dataKey: 默认为options.data
func Js2Go(code string, scopeKey string) (goCode string, err error) {
	// 用括号包裹的原因是让"{x: 1}"这样的语法解析成对象, 而不是label
	code = fmt.Sprintf("(%s)", code)

	p, err := parser.ParseFile(nil, "", code, 0)
	if err != nil {
		err = fmt.Errorf("GetAst err: %w, code:%s", err, code)
		return
	}

	goCode = genGoCodeByNode(p.Body[0], scopeKey)
	return
}

func genGoCodeByNode(node ast.Node, scopeKey string) (goCode string) {
	switch t := node.(type) {

	case *ast.ExpressionStatement:
		return genGoCodeByNode(t.Expression, scopeKey)
	case *ast.Identifier:
		return fmt.Sprintf(`%s.Get("%s")`, scopeKey, t.Name)
	case *ast.DotExpression:
		root, keys := lookExpress(t, scopeKey)
		return fmt.Sprintf(`%s.Get(%s)`, root, strings.Join(keys, ", "))
	case *ast.BracketExpression:
		// a[b]
		root, keys := lookExpress(t, scopeKey)
		return fmt.Sprintf(`%s.Get(%s)`, root, strings.Join(keys, ", "))
	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, t.Value)
	case *ast.NumberLiteral:
		return fmt.Sprintf("%v", t.Value)
	case *ast.BooleanLiteral:
		return fmt.Sprintf("%v", t.Value)
	case *ast.NullLiteral:
		return fmt.Sprintf("%v", "nil")
	case *ast.BinaryExpression:
		left := genGoCodeByNode(t.Left, scopeKey)
		right := genGoCodeByNode(t.Right, scopeKey)
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
		case token.LOGICAL_AND, token.LOGICAL_OR:
			return fmt.Sprintf(`interfaceToBool(%s) %s interfaceToBool(%s)`, left, t.Operator, right)
		case token.LESS:
			return fmt.Sprintf(`interfaceLess(%s, %s)`, left, right)
		case token.GREATER:
			return fmt.Sprintf(`interfaceGreater(%s, %s)`, left, right)

		default:
			panic(fmt.Sprintf("bad Operator for BinaryExpression: %s", o))
		}

	case *ast.UnaryExpression:
		arg := genGoCodeByNode(t.Operand, scopeKey)
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

			switch v.Kind {
			case "value":
				k = fmt.Sprintf(`"%s"`, v.Key)
			default:
				panic(fmt.Sprintf("bad Value kind of ObjectLiteral: %v", v.Kind))
			}

			valueCode := genGoCodeByNode(v.Value, scopeKey)
			mapCode += fmt.Sprintf(`%s: %s,`, k, valueCode)
		}
		mapCode += "}"
		return mapCode
	case *ast.CallExpression:
		funcName := genGoCodeByNode(t.Callee, scopeKey)

		args := make([]string, len(t.ArgumentList))
		for i, v := range t.ArgumentList {
			args[i] = genGoCodeByNode(v, scopeKey)
		}
		return fmt.Sprintf(`interfaceToFunc(%s)(options, %s)`, funcName, strings.Join(args, ","))
	case *ast.ArrayLiteral:
		args := make([]string, len(t.Value))
		for i, v := range t.Value {
			args[i] = genGoCodeByNode(v, scopeKey)
		}
		return fmt.Sprintf(`[]interface{}{%s}`, strings.Join(args, ","))
	case *ast.ConditionalExpression:
		// 三元运算
		consequent := genGoCodeByNode(t.Consequent, scopeKey)
		alternate := genGoCodeByNode(t.Alternate, scopeKey)
		test := genGoCodeByNode(t.Test, scopeKey)

		return fmt.Sprintf(`func() interface{} {if interfaceToBool(%s){return %s};return %s}()`, test, consequent, alternate)

	default:
		panic(fmt.Sprintf("bad type %T for genGoCodeByNode", t))
		//bs, _ := json.Marshal(t)
		//log.Panicf("%s", bs)
	}
	return
}

// 读取值
// 将a.b.c解析成 root 和keys
// 如a.b.c, root: this, keys: [a ,b ,c]
// 如"a".length, root: "a", keys: [length]
func lookExpress(e ast.Expression, scopeKey string) (root string, keys []string) {
	switch r := e.(type) {
	case *ast.DotExpression:
		// a.b 中的b
		currKey := fmt.Sprintf(`"%s"`, r.Identifier.Name)
		root, keys = lookExpress(r.Left, scopeKey)
		keys = append(keys, currKey)
	case *ast.Identifier:
		// a.b 中的a
		// 使用dataKey读取变量
		root = scopeKey
		keys = []string{fmt.Sprintf(`"%s"`, r.Name)}
	case *ast.ObjectLiteral:
		root = genGoCodeByNode(r, scopeKey)
	case *ast.BinaryExpression:
		root = genGoCodeByNode(r, scopeKey)
	case *ast.BracketExpression:
		var currKey string
		switch m := r.Member.(type) {
		case *ast.StringLiteral:
			// a['b']
			// 也可以走default语句, 但这是fastPath, 可以少调用interfaceToStr函数
			currKey = fmt.Sprintf(`"%s"`, m.Value)
		default:
			// a[b]
			// a[a+1]
			// ... 各种表达式
			currKey = fmt.Sprintf(`interfaceToStr(%s)`, genGoCodeByNode(r.Member, scopeKey))
		}

		root, keys = lookExpress(r.Left, scopeKey)
		keys = append(keys, currKey)
	default:
		panic(fmt.Sprintf("bad type for lookExpress: %T, %s", r, r))
	}

	return
}
