package ast_from_api

import "testing"

func TestGenGoCode(t *testing.T) {
	node, err := GetAST(`(data.r.st.pc['custom-class'-1].name)`)
	if err != nil {
		t.Fatal(err)
	}
	goCode := genGoCodeByNode(node, "data")
	t.Log(goCode)
}
