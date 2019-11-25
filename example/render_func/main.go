package main

import (
	"github.com/bysir-zl/vue-ssr/genera"
	"go.zhuzi.me/go/log"
)

func main() {
	// run pkg/vuessr/generator_test.go first
	html := genera.XComponent_helloworld(map[string]interface{}{
		"name":   "bysir",
		"sex":    "男",
		"age":    "18",
		"list":   []interface{}{"1", map[string]interface{}{"a": 1}},
		"isShow": "1",
	}, "")

	log.Infof("%v", html)
}
