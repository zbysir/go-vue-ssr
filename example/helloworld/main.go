package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// cd example/helloworld
// exec `go-vue-ssr -src=./vue -to=./ -pkg=main` before run main
func main() {
	r := NewRender()
	// 此指令获取渲染过程中所有v-on指令数据, 用来添加事件.
	r.Directive("v-on-handler", func(b DirectivesBinding, options *Options) {
		options.Slot = map[string]NamedSlotFunc{"default": func(props map[string]interface{}) string {
			bs, _ := json.Marshal(r.VOnBinds)
			return fmt.Sprintf("var vOnBinds = %s;", bs)
		}}
	})
	htmlStr := r.Component_page(&Options{
		Props: map[string]interface{}{
			"title":  "go-vue-ssr",
			"slogan": "Hey vue go",
			"info": map[string]interface{}{
				"author": "bysir",
				"Hey vue go":"Hey vue go",
			},
			"logo":   "https://avatars2.githubusercontent.com/u/13434040?s=88&v=4",
			"height": 100.1,
		},
	})

	// will print like following code(formatted):
	// <html lang="zh">
	// <head>
	//   <meta charset="UTF-8"></meta>
	//   <title>go-vue-ssr</title>
	// </head>
	// <body><h1>go-vue-ssr</h1>
	// <div style="margin-bottom: 10px; padding: 40px; text-align: center;">
	//   <p style="padding: 10px 0; ">Hey vue go</p>
	//   <img alt="todo logo" height="50px" src="https://avatars2.githubusercontent.com/u/13434040?s=88&amp;v=4"></img>
	// </div>
	// </body>
	// </html>
	log.Print(htmlStr)
}
