package vuessr

import (
	"encoding/xml"
	"testing"
)

func TestVFor(t *testing.T) {
	d := getVForDirective(xml.Attr{Value:"item in list"})
	t.Log(d.Exec(nil, `"<div></div>"`))
}

func TestVIf(t *testing.T) {
	d := getVIfDirective(xml.Attr{Value:"isShow"})
	t.Log(d.Exec(nil, `"<div></div>"`))
}
