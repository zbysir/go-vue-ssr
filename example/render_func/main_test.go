package main

import (
	"github.com/bysir-zl/vue-ssr/genera"
	"testing"
)

func TestX(t *testing.T) {
	html := genera.XComponent_helloworld(&genera.Options{
		Props: map[string]interface{}{
			"name":        "bysir",
			"sex":         "男",
			"age":         "18",
			"list":        []interface{}{"1", map[string]interface{}{"a": 1}},
			"isShow":      true,
			"customClass": "customClass",
		},
	})

	t.Log(html)
}
