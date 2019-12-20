package errors

import (
	"errors"
	"go.zhuzi.me/go/log"
	"testing"
)

func TestConcat(t *testing.T) {
	e := Concat(NewCoder(400, "e1"), NewCoder(500, "e2"))
	if e.Code() != 400 {
		t.Fatal("code is't 400")
	}
	log.Errorf2("%+v; %d", e.Where(), e, e.Code())
}

func TestNewError(t *testing.T) {
	err := NewCoder(500, "plan")
	t.Log(err.Error())
}

func TestWarp(t *testing.T) {
	err := NewCoder(500, "打开文件", "z:/1.txt not found")
	err2 := NewCoder(400, err, "文件微服务")
	t.Log(err2.Error())
}

func TestWarp2(t *testing.T) {
	err := NewCoder(400, "页面未发布")
	err2 := NewCoder(err, "调用站点微服务")
	t.Log(err2.Error())
}

func TestNil(t *testing.T) {
	var err error
	err = NewCodere(400, errors.New("No content found to be updated"), "xx")
	//t.Log(err.Error())
	t.Log(err == nil)
}
