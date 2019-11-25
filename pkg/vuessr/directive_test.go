package vuessr

import "testing"

func TestVFor(t *testing.T) {
	d := getVForDirective("item in list")
	t.Log(d.Exec(nil, `"<div></div>"`))
}

func TestVIf(t *testing.T) {
	d := getVIfDirective("isShow")
	t.Log(d.Exec(nil, `"<div></div>"`))
}
