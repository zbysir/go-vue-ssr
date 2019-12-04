package main

import (
	"fmt"
	"github.com/bysir-zl/go-vue-ssr/internal/vuetpl"
	"testing"
)

func TestX(t *testing.T) {
	r := vuetpl.NewRender()
	r.Prototype = map[string]interface{}{"img": func(args ...interface{}) interface{} {
		return fmt.Sprintf("%s?%d", args[0], 10000)
	}}
	html := r.Component_helloworld(&vuetpl.Options{
		Props: map[string]interface{}{
			"name":        "bysir",
			"sex":         "男",
			"age":         "18",
			"list":        []interface{}{"1", map[string]interface{}{"a": 2}},
			"isShow":      true,
			"customClass": "customClass",
			"imgUrl":      "https://s3.cn-north-1.amazonaws.com.cn/lcavatar/00b5aeb3-e45b-4aa1-a530-ff21c4d5835c_80x80.png",
		},
	})

	t.Log(html)
}

func TestVIf(t *testing.T) {
	r := vuetpl.NewRender()
	html := r.Component_vif(&vuetpl.Options{
		Props: map[string]interface{}{
			"name": "bysir",
			//"name2": "b2",
		},
	})

	t.Log(html)
}
