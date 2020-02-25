package ast

import "testing"

func TestBase(t *testing.T) {
	gocode, err := Js2Go("{a: 1}", "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)
}
