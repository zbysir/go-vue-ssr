package test

import (
	"encoding/json"
	"fmt"
	"github.com/bysir-zl/go-vue-ssr/internal/vuetpl"
	"testing"
)

func TestHelloworld(t *testing.T) {
	r := vuetpl.NewRender()
	r.Prototype = map[string]interface{}{"img": func(args ...interface{}) interface{} {
		return fmt.Sprintf("%s?%d", args[0], 10000)
	}}
	html := r.Component_helloworld(&vuetpl.Options{
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

	t.Log(html)
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
	r.Ctx = GetSet{}

	r.Directive("v-animate", func(value interface{}, r *vuetpl.Render, options *vuetpl.Options) {
		// add class
		c := vuetpl.LookInterface(value, "xclass")
		if c != nil {
			options.Attrs = map[string]string{"data": "2"}
			options.Class = append(options.Class, vuetpl.InterfaceToStr(c))
		}
	})
	r.Directive("v-set", func(value interface{}, r *vuetpl.Render, options *vuetpl.Options) {
		r.Ctx.Set(vuetpl.InterfaceToStr(vuetpl.LookInterface(value, "key")), vuetpl.LookInterface(value, "value"))
	})
	r.Directive("v-get", func(value interface{}, r *vuetpl.Render, options *vuetpl.Options) {
		options.Slot["default"] = func(props map[string]interface{}) string {
			bs, _ := json.Marshal(r.Ctx.GetAll())
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

	html := r.Component_xstyle(&vuetpl.Options{
		Props: map[string]interface{}{
			"text": "bysir.jpg",
		},
	})

	t.Log(html)
}
