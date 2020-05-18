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

func TestHtml(t *testing.T) {
	e, err := parseHtml(`Z:\go_path\src\github.com\zbysir\go-vue-ssr\internal\test\vue\select.vue`)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(e, " ", " ")
	t.Logf("%s", bs)
}

func TestVif(t *testing.T) {
	e, err := parseHtml(`Z:\go_path\src\github.com\zbysir\go-vue-ssr\test\base\vif.vue`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(e[0].hasOnlyOneChildren())
}
