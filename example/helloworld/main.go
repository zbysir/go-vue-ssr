package main

import "log"

// exec `go-vue-ssr -src=./vue -to=./ -pkg=main` before run main func
func main() {
	r := NewRender()
	htmlStr := r.Component_page(&Options{
		Props: map[string]interface{}{
			"title":  "go-vue-ssr",
			"slogan": "Hey vue go",
			"logo":   "https://avatars2.githubusercontent.com/u/13434040?s=88&v=4",
		},
	})

	// will print like below code(formatted):
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
