package main

import (
	"testing"
	"time"
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

func TestListSpans(t *testing.T) {
	var p = NewListSpans()

	p.AppendString("1")
	p.AppendString("2")
	p.AppendString("3")
	{
		s := NewChanSpan()
		go func() {
			time.Sleep(2 * time.Second)
			s.Done("4")
		}()
		p.AppendSpan(s)
	}

	t.Log("3", p.Result())

	p3 := NewListSpans()
	p3.AppendString("5")

	//t.Log("5", p.Join())

	p2 := NewListSpans()
	p2.AppendSpans(p3)

	p.AppendSpan(p2)

	t.Log("5", p.Result())

	p.AppendString("6")

	// for cur := p; cur != nil; cur = cur.Next {
	// 	t.Log(cur.Note)
	// }

	want := "123456"
	r := p.Result()

	if r != want {
		t.Fatalf("want:%s but:%s", want, r)
	}
	t.Log("ok")

}

func TestBufferSpans(t *testing.T) {
	var p = NewBufferSpans()

	p.AppendString("1")
	p.AppendString("2")
	p.AppendString("3")
	{
		s := NewChanSpan()
		go func() {
			time.Sleep(2 * time.Second)
			s.Done("4")
		}()
		p.AppendSpan(s)
	}

	t.Log("3", p.Result())

	p3 := NewListSpans()
	p3.AppendString("5")

	//t.Log("5", p.Join())

	p2 := NewListSpans()
	p2.AppendSpans(p3)

	p.AppendSpan(p2)

	t.Log("5", p.Result())

	p.AppendString("6")

	// for cur := p; cur != nil; cur = cur.Next {
	// 	t.Log(cur.Note)
	// }

	want := "123456"
	r := p.Result()

	if r != want {
		t.Fatalf("want:%s but:%s", want, r)
	}
	t.Log("ok")

}
