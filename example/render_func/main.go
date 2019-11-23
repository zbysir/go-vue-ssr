package main

import (
	"fmt"
	"github.com/bysir-zl/vue-ssr/pkg/vuessr"
	"go.zhuzi.me/go/log"
)

func main() {
	e, err := vuessr.H(`Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func\helloworld.vue`)
	if err != nil {
		panic(err)
	}

	app := vuessr.NewApp()
	app.ComponentFile("text", `Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func\text.vue`)

	str := e.RenderFunc(app, nil, "")

	log.Infof("%v", fmt.Sprintf(`function XComponent_%s(data map[string]interface{}, slot string)string{return %s}`, "main", str))
	//log.Infof("%+v", e)
}
