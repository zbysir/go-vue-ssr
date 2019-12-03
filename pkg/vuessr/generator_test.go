package vuessr

import (
	"strings"
	"testing"
)

func TestGenAllFile(t *testing.T) {
	//err := GenAllFile(`/Users/bysir/go/src/github.com/bysir-zl/go-vue-ssr/example/render_func`, `/Users/bysir/go/src/github.com/bysir-zl/go-vue-ssr/genera`)
	err := GenAllFile(`Z:\go_path\src\github.com\bysir-zl\go-vue-ssr\example\render_func`, `Z:\go_path\src\github.com\bysir-zl\go-vue-ssr\internal\vuetpl`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("OK")
}

func TestGenComponentRenderFunc(t *testing.T) {
	app := NewApp()

	code := genComponentRenderFunc(app, "gebera", "xx", `Z:\go_path\src\github.com\bysir-zl\go-vue-ssr\example\render_func\v-for.vue`)
	t.Log(code)
}

func TestShouldLookInterface(t *testing.T) {
	d, exist := shouldLookInterface([]int64{1, 3}, "length")
	t.Log(exist, d)
}

func interface2Slice(s interface{}) (d []interface{}) {
	switch a := s.(type) {
	case []interface{}:
		return a
	case []map[string]interface{}:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int64:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []int32:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []string:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	case []float64:
		d = make([]interface{}, len(a))
		for i, v := range a {
			d[i] = v
		}
	}
	return
}

func shouldLookInterface(data interface{}, key string) (desc interface{}, exist bool) {
	m, isObj := data.(map[string]interface{})

	kk := strings.Split(key, ".")
	currKey := kk[0]

	// 如果是对象, 则继续查找下一级
	if len(kk) != 1 && isObj {
		c, ok := m[currKey]
		if !ok {
			return
		}
		return shouldLookInterface(c, strings.Join(kk[1:], "."))
	}

	if len(kk) == 1 {
		if isObj {
			c, ok := m[currKey]
			if !ok {
				return
			}
			return c, true
		} else {
			switch currKey {
			case "length":
				switch t := data.(type) {
				// string
				case string:
					return len(t), true
				default:
					// slice
					return len(interface2Slice(t)), true
				}
			}
		}
	} else {
		// key不只有一个, 但是data不是对象, 说明出现了undefined的问题, 直接return
		return
	}

	return
}
