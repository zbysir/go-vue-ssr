package vuessr

import (
	"golang.org/x/net/html"
	"testing"
)

func TestVFor(t *testing.T) {
	d := getVForDirective(html.Attribute{Val: "item in list"})
	t.Log(d.Exec(nil, `"<div></div>"`))
}

func TestVIf(t *testing.T) {
	d := getVIfDirective(html.Attribute{Val: "isShow"})
	t.Log(d.Exec(nil, `"<div></div>"`))
}
