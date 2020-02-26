package ast

import "testing"

func TestObject(t *testing.T) {
	gocode, err := Js2Go(`{a+1: 1}[c]`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)
}

func TestBracket(t *testing.T) {
	gocode, err := Js2Go(`ab.c[c]`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)
}

func TestMulti(t *testing.T) {
	gocode, err := Js2Go(`data.r.st.pc['custom-class'-1].name`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)
}

func TestArray(t *testing.T) {
	gocode, err := Js2Go(`[data.i,data.r.st.pc['custom-class'].name]`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)
}

func TestLogical(t *testing.T) {
	gocode, err := Js2Go(`a || b`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)
}

func TestUnary(t *testing.T) {
	gocode, err := Js2Go(`!b`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)
}

func TestFunc(t *testing.T) {
	gocode, err := Js2Go(`a(b)`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)
}

func TestAll(t *testing.T) {
	gocode, err := Js2Go(`ab.c[c]`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)

	gocode, err = Js2Go(`1`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)

	gocode, err = Js2Go(`"1"`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)

	gocode, err = Js2Go(`1+1`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)

	gocode, err = Js2Go(`1+"1"`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)

	gocode, err = Js2Go(`1-1`, "this")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", gocode)

}
