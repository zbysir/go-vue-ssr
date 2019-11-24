package vuessr

import "testing"

func TestGenAllFile(t *testing.T) {
	err := genAllFile(`Z:\golang\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func`, `Z:\golang\go_path\src\github.com\bysir-zl\vue-ssr\genera`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("OK")
}

func TestGenComponentRenderFunc(t *testing.T) {
	app := NewApp()

	code := genComponentRenderFunc(app, "gebera", "xx", `Z:\golang\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func\v-for.vue.html`)
	t.Log(code)
}
