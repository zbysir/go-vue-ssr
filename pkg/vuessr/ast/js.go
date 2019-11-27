package ast

// 简单的js抽象语法树

// a && b
// a + 1
// !a

// 3种类型的节点: 注释/申明/表达式
// 目前只需要支持表达式
type Node interface {
}

// 表达式
type Expression interface {
	Node
}

// 变量
type Identifier struct {
	Expression
	Name string
}

// 字面值, 不用关心类型
type Literal struct {
	Expression
	Value string
}

// 数组
type Array struct {
	Expression
	Elements []Expression
}

type Object struct {
	Expression
	Properties []Property
}

type Property struct {
	Expression
	Key   string
	Value Expression
}

//  一元运算符
//  "-" | "+" | "!" | "~" | "typeof" | "void" | "delete"
type UnaryOperator struct {
	Expression
	Operator string
	Argument Expression
}

// 二元表达式
// "=" | "+=" | "-=" | "*=" | "/=" | "%="
// | "<<=" | ">>=" | ">>>="
// | "|=" | "^=" | "&="
type Binary struct {
	Expression
	Operator string
	Left     Node
	Right    Node
}

// 逻辑运算
// "||" | "&&"
type LogicalExpression struct {
	Expression
	Operator string
	Left     Node
	Right    Node
}

type Parser struct {
	code []rune
	pos int
}

func (p *Parser) getWord() (word string, l int){
	for {
		if p.pos>=len(p.code){
			break
		}
		curr:=p.code[p.pos]
		// ttodo
	}
	return
}

func parseExpression(code string) (es []Expression, err error) {

	return
}
