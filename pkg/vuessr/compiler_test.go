package vuessr

import "testing"

func TestParseVueVif(t *testing.T) {
	e, err := ParseVue(`Z:\go_path\src\github.com\bysir-zl\go-vue-ssr\internal\test\vue\page.vue`)
	if err != nil {
		t.Fatal(err)
	}
	c := NewCompiler()
	code, _ := c.GenEleCode(e)
	t.Log(code)
	return
}

func TestQuote(t *testing.T) {
	t.Log(quote(` '"a""{{b + "" +c }} "c" {{d}}`))
}
