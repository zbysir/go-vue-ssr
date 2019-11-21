package main

import (
	"github.com/bysir-zl/vue-ssr/pkg/vuessr"
	"go.zhuzi.me/go/log"
)

func main() {
	e, err := vuessr.H(`Z:\golang\go_path\src\github.com\bysir-zl\vue-ssr\example\helloword\helloworld.vue`)
	if err != nil {
		panic(err)
	}

	app := vuessr.NewApp()
	str := e.Render(app, map[string]interface{}{
		"name": "bysir",
		"sex":  "男",
		"age":  "18",
	}, "")

	log.Infof("%v", str)
	//log.Infof("%+v", e)
}
