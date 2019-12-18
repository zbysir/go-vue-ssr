package ast_from_api

import (
	"encoding/json"
	"testing"
)

func TestGetAST(t *testing.T) {
	node, err := GetAST(`({[a]: 1})`)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(node, " ", "  ")
	t.Logf("%s", bs)
}
