package vuessr

import (
	"encoding/xml"
	"go.zhuzi.me/go/log"
)

type TexTElement struct {
	Attrs []xml.Attr
}

func (e *TexTElement) Set(attrs []xml.Attr) {
	e.Attrs = attrs
}

// parentData: 上一层的数据
//func (e *TexTElement) Render(app *App, parentData interface{}) string {
//	bind := getBind(e.Attrs)
//	// 从bind中读取数据, 做为自己的数据
//	data := map[string]interface{}{}
//
//	m, ok := parentData.(map[string]interface{})
//	if ok {
//		for k, v := range bind {
//			data[k] = m[v]
//		}
//	}
//
//	log.Infof("%+v", data)
//
//	return fmt.Sprintf("<p>%v%%s</p>", "1")
//}

var texTElement Render

func init() {
	e, err := H(`Z:\golang\go_path\src\github.com\bysir-zl\vue-ssr\example\helloword\text.vue`)
	if err != nil {
		panic(err)
	}

	texTElement = e

	log.Infof("%+v", e)

}
