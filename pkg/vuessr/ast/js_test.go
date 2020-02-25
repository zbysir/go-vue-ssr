package ast

import (
	"encoding/json"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
	"testing"
)

func TestIdentifier(t *testing.T) {
	code := `a`

	t.Log(code)
}

func TestParse(t *testing.T) {
	src:="a+b"
	p,err:=parser.ParseFile(nil, "",src, 0)
	if err != nil {
		t.Fatal(err)
	}

	//p.Body
	bs,_:=json.MarshalIndent(p.Body, " "," ")
	t.Logf("%s %+v", bs, p.Body)
}

func TestVM(t *testing.T) {
	vm := otto.New()
	v,err:=vm.Run(`
    abc = 2 + 2;
    console.log("The value of abc is " + abc); // 4
`)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", v)
}
