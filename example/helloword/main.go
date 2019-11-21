package main

import (
	"github.com/bysir-zl/vue-ssr/pkg/vuessr"
	"go.zhuzi.me/go/log"
)

func main() {
	e, err := vuessr.H(`Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\helloword\helloworld.vue`)
	if err != nil {
		panic(err)
	}
	str, err := vuessr.RenderNode(e)
	if err != nil {
		panic(err)
	}

	log.Infof("%v", str)
	//log.Infof("%+v", e)
}
