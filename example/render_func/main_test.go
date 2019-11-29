package main

import (
	"github.com/bysir-zl/vue-ssr/internal/vuetpl"
	"testing"
)

func TestX(t *testing.T) {
	html := vuetpl.XComponent_helloworld(&vuetpl.Options{
		Props: map[string]interface{}{
			"name":        "bysir",
			"sex":         "男",
			"age":         "18",
			"list":        []interface{}{"1", map[string]interface{}{"a": 2}},
			"isShow":      true,
			"customClass": "customClass",
		},
	})

	t.Log(html)
}
