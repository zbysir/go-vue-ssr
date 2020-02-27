// cd internal/test
// go-vue-ssr -src=./vue -to=./tplgo

package test

import (
	"encoding/json"
	"fmt"
	"github.com/zbysir/go-vue-ssr/internal/test/tplgo"
	"github.com/zbysir/go-vue-ssr/pkg/ssrtool"
	"github.com/zbysir/go-vue-ssr/pkg/ssrtool/rinterface"
	"testing"
)

func TestHelloworld(t *testing.T) {
	r := tplgo.NewRender()
	r.Prototype = map[string]interface{}{"img": func(args ...interface{}) interface{} {
		return fmt.Sprintf("%s?%d", args[0], 10000)
	}}
	str := r.Component_helloworld(&tplgo.Options{
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

	t.Logf("%s", str)
}

func TestVIf(t *testing.T) {
	r := tplgo.NewRender()
	html := r.Component_vif(&tplgo.Options{
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
	r := tplgo.NewRender()
	// 使用闭包访问ctx, 以实现在多个directive中共享数据
	ctx := GetSet{}

	r.Directive("v-animate", func(binding tplgo.DirectivesBinding, options *tplgo.Options) {
		// add class
		c := rinterface.Get(binding.Value, "xclass")
		if c != nil {
			options.Attrs = map[string]string{"data": "2"}
			options.Class = append(options.Class, rinterface.ToStr(c, false))
		}
	})
	r.Directive("v-set", func(binding tplgo.DirectivesBinding, options *tplgo.Options) {
		ctx.Set(
			binding.Arg,
			rinterface.Get(binding.Value, "value"))
	})
	r.Directive("v-get", func(binding tplgo.DirectivesBinding, options *tplgo.Options) {
		options.Slot["default"] = func(props map[string]interface{}) string {
			bs, _ := json.Marshal(ctx.GetAll())
			return string(bs)
		}
	})

	html := r.Component_directive(&tplgo.Options{
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
	r := tplgo.NewRender()
	r.Prototype = map[string]interface{}{
		"img": func(args ...interface{}) interface{} {
			return fmt.Sprintf("%s?100", args[0])
		},
	}

	html := r.Component_xattr(&tplgo.Options{
		Props: map[string]interface{}{
			"imgUrl":      "bysir.jpg",
			"customClass": "customClass",
		},
	})

	t.Log(html)
}

func TestStyle(t *testing.T) {
	r := tplgo.NewRender()

	html := r.Component_xStyle(&tplgo.Options{
		Props: map[string]interface{}{
			"text": "bysir.jpg",
		},
	})

	t.Log(html)
}

func TestVText(t *testing.T) {
	r := tplgo.NewRender()

	html := r.Component_vtext(&tplgo.Options{
		Props: map[string]interface{}{
			"text": "<p color=red>bysir.jpg</p>",
			"html": "<p color=red>bysir.jpg</p>",
		},
	})

	if html != `<div><div>&lt;p color=red&gt;bysir.jpg&lt;/p&gt;</div><div><p color=red>bysir.jpg</p></div>
    "&lt;p color=red&gt;bysir.jpg&lt;/p&gt;
  </div>` {
		t.Fatal(html)
	}

	t.Log(html)
}
