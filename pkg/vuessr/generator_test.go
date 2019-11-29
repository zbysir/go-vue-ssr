package vuessr

import "testing"

func TestGenAllFile(t *testing.T) {
	//err := GenAllFile(`/Users/bysir/go/src/github.com/bysir-zl/vue-ssr/example/render_func`, `/Users/bysir/go/src/github.com/bysir-zl/vue-ssr/genera`)
	err := GenAllFile(`Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func`, `Z:\go_path\src\github.com\bysir-zl\vue-ssr\genera`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("OK")
}

func TestGenComponentRenderFunc(t *testing.T) {
	app := NewApp()

	code := genComponentRenderFunc(app, "gebera", "xx", `Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func\v-for.vue`)
	t.Log(code)
}
