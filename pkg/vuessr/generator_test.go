package vuessr

import "testing"

func TestGenAllFile(t *testing.T) {
	err := genAllFile(`Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func`, `Z:\go_path\src\github.com\bysir-zl\vue-ssr\generat`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("OK")
}
