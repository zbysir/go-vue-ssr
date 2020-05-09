package vuessr

import "testing"

func TestParseVueVif(t *testing.T) {
	e, err := ParseVue(`Z:\go_project\go-vue-ssr\internal\test\async_true\bench.vue`)
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
