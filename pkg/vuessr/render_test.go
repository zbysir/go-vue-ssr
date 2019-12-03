package vuessr

import "testing"

func TestParseVueVif(t *testing.T) {
	e, err := ParseVue(`Z:\go_path\src\github.com\bysir-zl\go-vue-ssr\example\render_func\vif.vue`)
	if err != nil {
		t.Fatal(err)
	}
	app := NewApp()
	code, _ := e.GenCode(app)
	t.Log(code)
	return
}
