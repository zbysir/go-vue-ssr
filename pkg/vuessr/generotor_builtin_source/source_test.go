package main

import (
	"testing"
)

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

func TestPromise(t *testing.T) {
	var p = &PromiseGroup{
	}

	p.Append("1")
	p.Append("2")
	p.Append(PromiseString("3"))
	//{
	//	s := make(PromiseChan, 1)
	//	go func() {
	//		time.Sleep(2 * time.Second)
	//		s <- "4"
	//	}()
	//	p.Append(s)
	//}

	t.Log("3", p.Join())

	p3 := &PromiseGroup{
		Note: "p3",
	}
	p3.AppendString("5")

	//t.Log("5", p.Join())

	p2 := &PromiseGroup{
		Note: "p2",
	}
	p2.AppendGroup(p3)

	p.AppendGroup(p2)

	t.Log("5", p.Join())

	p.Append("6")

	for cur := p; cur != nil; cur = cur.Next {
		t.Log(cur.Note)

	}

	t.Log(p.Join())
}
