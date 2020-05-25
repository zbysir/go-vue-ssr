package main

import (
	"log"
)

// cd example/helloworld
// exec `go-vue-ssr -src=./vue -to=./ -pkg=main -watch` before run main
func main() {
	c := NewRenderCreator()

	r := c.NewRender()

	w := r.NewWriter()
	r.Render("page", w, &Options{
		Props: map[string]interface{}{
			"title":  "go-vue-ssr",
			"slogan": "Hey vue go",
			"info": map[string]interface{}{
				"author":     "bysir",
				"Hey vue go": "Hey vue go",
			},
			"logo":   "https://avatars2.githubusercontent.com/u/13434040?s=88&v=4",
			"height": 100.1,
		},
	})

	log.Print(w.Result())

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
}
