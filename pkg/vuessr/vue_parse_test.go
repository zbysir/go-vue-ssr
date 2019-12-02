package vuessr

import (
	"encoding/json"
	"testing"
)

func TestParseVue(t *testing.T) {
	e, err := ParseVue(`Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func\vif.vue`)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(e, " ", " ")
	t.Logf("%s", bs)
}
