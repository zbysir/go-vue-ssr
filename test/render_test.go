package test

import (
	"encoding/json"
	"fmt"
	"github.com/bysir-zl/go-vue-ssr/internal/vuetpl"
	"github.com/bysir-zl/go-vue-ssr/pkg/ssrtool"
	"testing"
)

func TestHelloworld(t *testing.T) {
	r := vuetpl.NewRender()
	r.Prototype = map[string]interface{}{"img": func(args ...interface{}) interface{} {
		return fmt.Sprintf("%s?%d", args[0], 10000)
	}}
	str := r.Component_helloworld(&vuetpl.Options{
		Props: map[string]interface{}{
			"name":        "bysir",
			"sex":         "男",
			"age":         "18",
			"list":        []interface{}{"1", map[string]interface{}{"a": 2}},
			"isShow":      true,
			"customClass": "customClass",
			"imgUrl":      "https://s3.cn-north-1.amazonaws.com.cn/lcavatar/00b5aeb3-e45b-4aa1-a530-ff21c4d5835c_80x80.png",
		},
	})

	t.Logf("%s", ssrtool.FormatHtml(str, 2))
}

func TestVIf(t *testing.T) {
	r := vuetpl.NewRender()
	html := r.Component_vif(&vuetpl.Options{
		Props: map[string]interface{}{
			"name": "bysir",
			//"name2": "b2",
		},
	})

	t.Log(html)
}

type GetSet map[string]interface{}

func (g GetSet) Get(key string) interface{} {
	return g[key]
}

func (g GetSet) Set(key string, val interface{}) {
	g[key] = val
}

func (g GetSet) GetAll() map[string]interface{} {
	return g
}

func TestVDirective(t *testing.T) {
	r := vuetpl.NewRender()
	// 使用闭包访问ctx, 以实现在多个directive中共享数据
	ctx := GetSet{}

	r.Directive("v-animate", func(binding vuetpl.DirectivesBinding, options *vuetpl.Options) {
		// add class
		c := vuetpl.LookInterface(binding.Value, "xclass")
		if c != nil {
			options.Attrs = map[string]string{"data": "2"}
			options.Class = append(options.Class, vuetpl.InterfaceToStr(c))
		}
	})
	r.Directive("v-set", func(binding vuetpl.DirectivesBinding, options *vuetpl.Options) {
		ctx.Set(
			binding.Arg,
			vuetpl.LookInterface(binding.Value, "value"))
	})
	r.Directive("v-get", func(binding vuetpl.DirectivesBinding, options *vuetpl.Options) {
		options.Slot["default"] = func(props map[string]interface{}) string {
			bs, _ := json.Marshal(ctx.GetAll())
			return string(bs)
		}
	})

	html := r.Component_directive(&vuetpl.Options{
		Props: map[string]interface{}{
			"name":   "bysir",
			"xclass": "v-animate",
			"speed":  "5s",
			"id":     "_123",
			"show":   1,
		},
	})

	html = ssrtool.FormatHtml(html, 2)

	t.Log(html)
}

func TestAttr(t *testing.T) {
	r := vuetpl.NewRender()
	r.Prototype = map[string]interface{}{
		"img": func(args ...interface{}) interface{} {
			return fmt.Sprintf("%s?100", args[0], )
		},
	}

	html := r.Component_xattr(&vuetpl.Options{
		Props: map[string]interface{}{
			"imgUrl":      "bysir.jpg",
			"customClass": "customClass",
		},
	})

	t.Log(html)
}

func TestStyle(t *testing.T) {
	r := vuetpl.NewRender()

	html := r.Component_xStyle(&vuetpl.Options{
		Props: map[string]interface{}{
			"text": "bysir.jpg",
		},
	})

	t.Log(html)
}
