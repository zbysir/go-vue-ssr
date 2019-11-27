package ast_from_api

import "testing"

func TestGenGoCode(t *testing.T) {
	node,err:=GetAST(`({a: 1})`)
	if err != nil {
		t.Fatal(err)
	}
	goCode:= genGoCodeByNode(node)
	t.Log(goCode)
}
