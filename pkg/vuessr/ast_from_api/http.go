package ast_from_api

import "github.com/zbysir/go-vue-ssr/internal/pkg/http"

var client *http.Client

const apiHost = "http://test.zhuzi.me:23000/"

//const apiHost = "http://localhost:3000/"

func init() {
	var err error
	client, err = http.NewClient(apiHost)
	if err != nil {
		panic(err)
	}
}
