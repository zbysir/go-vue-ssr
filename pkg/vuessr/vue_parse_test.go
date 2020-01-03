package vuessr

import (
	"encoding/json"
	"testing"
)

func TestParseVue(t *testing.T) {
	e, err := ParseVue(`Z:\go_path\src\github.com\bysir-zl\go-vue-ssr\example\base\directive.vue`)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(e, " ", " ")
	t.Logf("%s", bs)
}

func TestHtml(t *testing.T) {
	e, err := parseHtml(`Z:\go_path\src\github.com\bysir-zl\go-vue-ssr\example\base\text.vue`)
	if err != nil {
		t.Fatal(err)
	}

	bs, _ := json.MarshalIndent(e, " ", " ")
	t.Logf("%s", bs)
}

func TestVif(t *testing.T) {
	e, err := parseHtml(`Z:\go_path\src\github.com\bysir-zl\go-vue-ssr\test\base\vif.vue`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(e[0].hasOnlyOneChildren())
}
