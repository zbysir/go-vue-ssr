package main

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {
	p := NewScope(nil)
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

	p.WriteString("1")
	p.WriteString("2")
	{
		s := NewChanSpan()
		go func() {
			time.Sleep(4 * time.Second)
			s.Done("3")
		}()
		p.WriteSpan(s)
	}

	{
		s := NewChanSpan()
		go func() {
			time.Sleep(2 * time.Second)
			s.Done("4")
		}()
		p.WriteSpan(s)
	}

	t.Log("3", p.Result())

	p3 := NewListSpans()
	p3.WriteString("5")

	//t.Log("5", p.Join())

	p2 := NewListSpans()
	p2.WriteSpan(p3)

	p.WriteSpan(p2)

	t.Log("5", p.Result())

	p.WriteString("6")

	// for cur := p; cur != nil; cur = cur.Next {
	// 	t.Log(cur.Note)
	// }

	want := "123456"
	// 由于并发计算的特性, 只应该执行4s
	r := p.Result()

	if r != want {
		t.Fatalf("want:%s but:%s", want, r)
	}
	t.Log("ok")

}

func TestBufferSpans(t *testing.T) {
	var p = NewBufferSpans()

	p.WriteString("1")
	p.WriteString("2")
	{
		s := NewChanSpan()
		go func() {
			time.Sleep(4 * time.Second)
			s.Done("3")
		}()
		p.WriteSpan(s)
	}
	{
		s := NewChanSpan()
		go func() {
			time.Sleep(2 * time.Second)
			s.Done("4")
		}()
		p.WriteSpan(s)
	}

	t.Log("3", p.Result())

	p3 := NewListSpans()
	p3.WriteString("5")

	//t.Log("5", p.Join())

	p2 := NewListSpans()
	p2.WriteSpan(p3)

	p.WriteSpan(p2)

	t.Log("5", p.Result())

	p.WriteString("6")

	want := "123456"
	// 会执行6s
	r := p.Result()

	if r != want {
		t.Fatalf("want:%s but:%s", want, r)
	}
	t.Log("ok")

}

func Test_getAttrFromProps(t *testing.T) {
	as := getAttrFromProps(NewProps(map[string]interface{}{
		"autoplay":  false,
		"id":        1,
		"autofocus": true,
	}))
	t.Logf("%+v", as)
}
