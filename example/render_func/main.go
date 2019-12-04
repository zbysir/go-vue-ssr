package main

import (
	"github.com/bysir-zl/go-vue-ssr/internal/vuetpl"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":10000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// run pkg/vuessr/generator_test.go first

		render := vuetpl.NewRender()
		html := render.Component_helloworld(&vuetpl.Options{
			Props: map[string]interface{}{
				"name":   "bysir",
				"sex":    "男",
				"age":    "18",
				"list":   []interface{}{"1", map[string]interface{}{"a": 1}},
				"isShow": true,
			},
		})
		w.Write([]byte(html))

		return
	}))

	if err != nil {
		panic(err)
	}
}
