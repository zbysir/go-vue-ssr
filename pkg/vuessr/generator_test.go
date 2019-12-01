package vuessr

import (
	"strings"
	"testing"
)

func TestGenAllFile(t *testing.T) {
	//err := GenAllFile(`/Users/bysir/go/src/github.com/bysir-zl/vue-ssr/example/render_func`, `/Users/bysir/go/src/github.com/bysir-zl/vue-ssr/genera`)
	err := GenAllFile(`Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func`, `Z:\go_path\src\github.com\bysir-zl\vue-ssr\genera`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("OK")
}

func TestGenComponentRenderFunc(t *testing.T) {
	app := NewApp()

	code := genComponentRenderFunc(app, "gebera", "xx", `Z:\go_path\src\github.com\bysir-zl\vue-ssr\example\render_func\v-for.vue`)
	t.Log(code)
}



func shouldLookInterface(data interface{}, key string) (desc interface{}, exist bool) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return nil, false
	}


	kk := strings.Split(key, ".")

	key:=kk[0]
	if len(kk)==1{
		switch t:=data.(type){
		case "string":
			if p,ok:=properties[key];ok{
				return p()
			}
			return
		}
	}
	c, ok := m[key]
	if len(kk) == 1 {
		if !ok {
			return nil, false
		}

		return c, true
	}

	return shouldLookInterface(c, strings.Join(kk[1:], "."))
}