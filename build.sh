go generate Z:/go_project/go-vue-ssr/pkg/vuessr/generator_builtin.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o go-vue-ssr main.go