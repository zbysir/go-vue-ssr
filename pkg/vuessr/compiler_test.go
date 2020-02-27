package vuessr

import "testing"

func TestParseVueVif(t *testing.T) {
	e, err := ParseVue(`Z:\golang\go_path\src\github.com\zbysir\go-vue-ssr\internal\test\vue\svg.vue`)
	if err != nil {
		t.Fatal(err)
	}
	c := NewCompiler()
	code, _ := c.GenEleCode(e)
	t.Log(code)
	return
}

func TestQuote(t *testing.T) {
	t.Log(safeStringCode(` '"a""{{b + "" +c }} "c" {{d}}`))
}
