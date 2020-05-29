package vuessr

import (
	"encoding/json"
	"testing"
)

func TestParseVue(t *testing.T) {
	e, err := ParseVue(`Z:\go_project\go-vue-ssr\internal\test\vue\page.vue`)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(e, " ", " ")
	t.Logf("%s", bs)
}
