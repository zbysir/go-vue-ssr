package rinterface

import "testing"

func TestLookJson(t *testing.T) {
	t.Logf("%v", LookJson([]byte(`{"a": 1}`), "a") == 1.0)
}

func TestLookJsonBool(t *testing.T) {
	a := `{"a" :{"b": true}}`
	if !LookJsonBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": 1}}`
	if !LookJsonBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": ""}}`
	if LookJsonBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": ''}}`
	if LookJsonBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": "1231232"}}`
	if !LookJsonBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": {}}}`
	if !LookJsonBool([]byte(a), "a.b") {
		t.Fatal(a)
	}

	t.Log("OK")
}
