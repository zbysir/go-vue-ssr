package rjson

import "testing"

func TestLookJson(t *testing.T) {
	t.Logf("%v", Get([]byte(`{"a": 1}`), "a") == 1.0)
}

func TestGetBool(t *testing.T) {
	a := `{"a" :{"b": true}}`
	if !GetBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": 1}}`
	if !GetBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": ""}}`
	if GetBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": ''}}`
	if GetBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": "1231232"}}`
	if !GetBool([]byte(a), "a.b") {
		t.Fatal(a)
	}
	a = `{"a" :{"b": {}}}`
	if !GetBool([]byte(a), "a.b") {
		t.Fatal(a)
	}

	t.Log("OK")
}
