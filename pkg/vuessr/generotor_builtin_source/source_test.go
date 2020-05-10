package main

import "testing"

func TestName(t *testing.T) {
	p := NewScope()
	p.Set("a", 2)
	scope := extendScope(p, map[string]interface{}{"a": 1})
	if scope.Get("a") != 1 {
		t.Fatal(scope.Get("a"))
	}

	scope = extendScope(p, map[string]interface{}{"b": 1})
	if scope.Get("a") != 2 {
		t.Fatal(scope.Get("a"))
	}
}
