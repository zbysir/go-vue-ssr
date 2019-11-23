package main

import (
	"github.com/bysir-zl/vue-ssr/generat"
	"go.zhuzi.me/go/log"
)

func main() {
	// run pkg/vuessr/generator_test.go first
	html := generat.XComponent_helloworld(map[string]interface{}{
		"name": "bysir",
		"sex":  "男",
		"age":  "18",
	}, "")

	log.Infof("%v", html)
}
