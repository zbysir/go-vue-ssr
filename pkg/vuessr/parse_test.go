package vuessr

import (
	"encoding/json"
	"testing"
)

func TestH(t *testing.T) {
	e,err:=H(`Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func\vif.vue`)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(e, " ", " ")
	t.Logf("%s", bs)
}